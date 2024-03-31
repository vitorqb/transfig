// Transfig is a library that aims at providing a reactiful state, similar to what is found in
// React for the JS world.
package transfig

import (
	"reflect"
)

// KeyValIter is an iterator for (key, value) pairs.
type KeyValIter func() (key KeyString, value interface{}, finished bool)

// Path is a sequence of KeyString to access a nested value in a state
type Path = []KeyString

// Selector is an interface to select (key, value) pairs from the state.
type Selector interface {
	// Select returns an iterator for the selected (key, value) pairs.
	Select(map[KeyString]interface{}) KeyValIter
	// Contains returns true if the selector contains a path of KeyStrings.
	Contains(Path) bool
}

// NestedSelector is a Selector for a nested value in the state
type NestedSelector struct {
	keys Path
}

func (s NestedSelector) Select(m map[KeyString]interface{}) KeyValIter {
	called := false
	return func() (key KeyString, value interface{}, finished bool) {
		if called || len(s.keys) == 0 {
			return key, nil, true
		}
		called = true
		if len(s.keys) == 1 {
			return s.keys[0], m[s.keys[0]], false
		}
		topKey, restKeys := s.keys[0], s.keys[1:]
		topValue := m[topKey]
		topValueAsMap, ok := topValue.(map[KeyString]interface{})
		if !ok {
			topValueAsMap = make(map[KeyString]interface{})
		}
		restIt := NestedSelector{restKeys}.Select(topValueAsMap)
		nestedKey, nestedValue, _ := restIt()
		return s.keys[0], map[KeyString]interface{}{nestedKey: nestedValue}, false
	}
}

func (s NestedSelector) Contains(p Path) bool {
	for i, k := range p {
		if s.keys[i] != k {
			return false
		}
	}
	return true
}

// StringKey is a string representing a specific key in the state. When used in
// a subscription, it will add a `key: value` pair to the subscription's arguments,
// where `key` is the string and `value` is the value in the state.
type KeyString string

// KeyString implements Selector for a single key
func (k KeyString) Select(m map[KeyString]interface{}) KeyValIter {
	called := false
	return func() (key KeyString, value interface{}, finished bool) {
		if called {
			return key, nil, true
		}
		called = true
		return k, m[k], false
	}
}

func (k KeyString) Contains(p Path) bool {
	if len(p) == 0 {
		return false
	}
	return p[0] == k
}

// Wildcard is a special Selector that selects all keys
type Wildcard struct{}

func (Wildcard) Select(m map[KeyString]interface{}) KeyValIter {
	keys := Path{}
	for k := range m {
		keys = append(keys, k)
	}
	return func() (key KeyString, value interface{}, finished bool) {
		if len(keys) == 0 {
			return "", nil, true
		}
		firstKey := keys[0]
		keys = keys[1:]
		return firstKey, m[firstKey], false
	}
}

func (Wildcard) Contains(Path) bool { return true }

type SubscriptionCallback func(map[KeyString]interface{})

// Subscription represents a func that will be called when the state changes
// for specific keys of the state.
type Subscription struct {
	name      string
	selectors []Selector
	callbacks []SubscriptionCallback
}

// With add keys selectors to the subscription
func (s *Subscription) With(selector Selector) *Subscription {
	s.selectors = append(s.selectors, selector)
	return s
}

// WithNested allows subscribing to nested values in the state
func (s *Subscription) WithNested(keys ...KeyString) *Subscription {
	s.selectors = append(s.selectors, NestedSelector{keys})
	return s
}

// Calls adds a function to be called when the state changes
func (s *Subscription) Calls(callback SubscriptionCallback) *Subscription {
	s.callbacks = append(s.callbacks, callback)
	return s
}

func (s *Subscription) subscribedTo(p Path) bool {
	for _, s := range s.selectors {
		if s.Contains(p) {
			return true
		}
	}
	return false
}

// notify calls the subscription's callbacks with the subscribed values
func (s *Subscription) notify(values map[KeyString]interface{}) {
	args := make(map[KeyString]interface{})
	for _, selector := range s.selectors {
		it := selector.Select(values)
		for {
			key, value, finished := it()
			if finished {
				break
			}
			args[key] = value
		}
	}
	for _, callback := range s.callbacks {
		callback(args)
	}
}

// NewSubscription creates a new subscription
func NewSubscription(name string) *Subscription {
	return &Subscription{
		name:      name,
		callbacks: make([]SubscriptionCallback, 0),
	}
}

// State represents a potentially nested key -> value state that can
// be subscribed to and updated.
type State struct {
	subscriptions map[string]*Subscription
	values        map[KeyString]interface{}
}

// Set updates the state with a new value for a specific key
func (s *State) Set(key KeyString, value interface{}) {
	s.SetNested(Path{key}, value)
}

// SetNested updates the state with a new value for a nested key
func (s *State) SetNested(path Path, value interface{}) {
	oldValue, found := MapGetNested(s.values, path)
	if found && reflect.DeepEqual(oldValue, value) {
		return
	}
	MapSetNested(s.values, path, value)
	notifiedSubs := make(map[string]bool)
	for _, sub := range s.subscriptions {
		if sub.subscribedTo(path) {
			if _, ok := notifiedSubs[sub.name]; !ok {
				notifiedSubs[sub.name] = true
				sub.notify(s.values)
			}
		}
	}
}

// Get returns the value for a specific key
func (s *State) Get(key KeyString) (value interface{}, found bool) {
	value, found = s.values[key]
	return value, found
}

// GetNested returns the value for a nested key
func (s *State) GetNested(keys ...KeyString) (value interface{}, found bool) {
	if len(keys) == 0 {
		return nil, false
	}
	topKey, restKeys := keys[0], keys[1:]
	value, found = s.values[topKey]
	for _, k := range restKeys {
		if valueAsMap, ok := value.(map[KeyString]interface{}); ok {
			value = valueAsMap[k]
		} else {
			return nil, false
		}
	}
	return value, found
}

// Subscribe adds a subscription to the state
func (s *State) Subscribe(subscription *Subscription) {
	s.subscriptions[subscription.name] = subscription
}

// Unsubscribe removes a subscription from the state
func (s *State) Unsubscribe(subscriptionName string) {
	delete(s.subscriptions, subscriptionName)
}

// AsMap returns a copy of the state as a map
func (s *State) AsMap() map[KeyString]interface{} {
	return MapDeepCopy(s.values)
}

// NewState creates a new state
func NewState() *State {
	return &State{
		subscriptions: make(map[string]*Subscription),
		values:        make(map[KeyString]interface{}),
	}
}
