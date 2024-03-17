// Transfig is a library that aims at providing a reactiful state, similar to what is found in
// React for the JS world.

package transfig

import "github.com/google/uuid"

// Key represents a key pointing to a value in the state
type Key interface {
	ExtractFrom(map[Key]interface{}) interface{}
}

// StringKey is a string representing a specific key in the state. When used in
// a subscription, it will add a `key: value` pair to the subscription's arguments,
// where `key` is the string and `value` is the value in the state.
type KeyString string

// ExtractFrom for a KeyString looks for the key in the state and returns its value.
// If the key is not found, it returns nil. If the value is a nested state, it will
// return a map with all values in the nested state.
func (k KeyString) ExtractFrom(state map[Key]interface{}) interface{} {
	value, exists := state[k]
	if !exists {
		return nil
	}
	if stateValue, ok := value.(*State); ok {
		return Wildcard.ExtractFrom(stateValue.values)
	}
	return value
}

// keyWildcard is a special key that represents all keys in the state. When used in
// a subscription, it will add a `key: value` pair to the subscription's arguments,
// where `key` is wildcard itself and `value` is a map with all values in the state.
// Nested states will be represented as maps as well.
type keyWildcard struct{}

// ExtractFrom for a keyWildcard returns the entire state.
func (k keyWildcard) ExtractFrom(state map[Key]interface{}) interface{} {
	out := make(map[Key]interface{})
	for k := range state {
		out[k] = k.ExtractFrom(state)
	}
	return out
}

var Wildcard = keyWildcard{}

type SubscriptionCallback func(map[Key]interface{})

// Subscription represents a func that will be called when the state changes
// for specific keys of the state.
type Subscription struct {
	name      string
	keys      []Key
	callbacks []SubscriptionCallback
}

// With adds keys to the subscription
func (s *Subscription) With(keys ...Key) *Subscription {
	s.keys = append(s.keys, keys...)
	return s
}

// Calls adds a function to be called when the state changes
func (s *Subscription) Calls(callback SubscriptionCallback) *Subscription {
	s.callbacks = append(s.callbacks, callback)
	return s
}

func (s *Subscription) subscribedTo(key Key) bool {
	for _, k := range s.keys {
		if k == Wildcard {
			return true
		}
		if k == key {
			return true
		}
	}
	return false
}

func (s *Subscription) notify(values map[Key]interface{}) {
	args := make(map[Key]interface{})
	for _, k := range s.keys {
		value := k.ExtractFrom(values)
		args[k] = value
	}
	for _, callback := range s.callbacks {
		callback(args)
	}
}

func NewSubscription(name string) *Subscription {
	return &Subscription{
		name:      name,
		keys:      make([]Key, 0),
		callbacks: make([]SubscriptionCallback, 0),
	}
}

type State struct {
	subscriptions                 map[string]*Subscription
	values                        map[Key]interface{}
	nestedStatesSubscriptionNames map[Key]string
}

func (s *State) Set(key Key, value interface{}) {
	if value == s.values[key] {
		return
	}
	if oldVal, ok := s.values[key]; ok {
		if oldStateVal, ok := oldVal.(*State); ok {
			if subName, ok := s.nestedStatesSubscriptionNames[key]; ok {
				oldStateVal.Unsubscribe(subName)
			}
		}
	}
	s.values[key] = value
	if stateValue, ok := value.(*State); ok {
		s.subscribeToNestedState(key, stateValue)
	}
	s.notifySubscribers(key)
}

func (s *State) Get(key Key) interface{} {
	return key.ExtractFrom(s.values)
}

func (s *State) Subscribe(subscription *Subscription) {
	s.subscriptions[subscription.name] = subscription
}

func (s *State) Unsubscribe(subscriptionName string) {
	delete(s.subscriptions, subscriptionName)
}

func (s *State) subscribeToNestedState(key Key, state *State) {
	subName := uuid.New().String()
	sub := NewSubscription(subName).With(Wildcard).Calls(func(_ map[Key]interface{}) {
		s.notifySubscribers(key)
	})
	s.nestedStatesSubscriptionNames[key] = subName
	state.Subscribe(sub)
}

func (s *State) notifySubscribers(key Key) {
	for _, sub := range s.subscriptions {
		if sub.subscribedTo(key) {
			sub.notify(s.values)
		}
	}
}

// NewState creates a new state
func NewState() *State {
	return &State{
		subscriptions:                 make(map[string]*Subscription),
		values:                        make(map[Key]interface{}),
		nestedStatesSubscriptionNames: make(map[Key]string),
	}
}
