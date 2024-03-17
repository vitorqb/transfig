// Transfig is a library that aims at providing a reactiful state, similar to what is found in
// React for the JS world.

package transfig

// Key represents a key pointing to a value in the state
type Key string

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
		if k == key {
			return true
		}
	}
	return false
}

func (s *Subscription) notify(values map[Key]interface{}) {
	args := make(map[Key]interface{})
	for _, k := range s.keys {
		if value, exists := values[k]; exists {
			args[k] = value
		} else {
			args[k] = nil
		}
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
	subscriptions map[string]*Subscription
	values        map[Key]interface{}
}

func (s *State) Set(key Key, value interface{}) {
	s.values[key] = value
	s.notifySubscribers(key)
}

func (s *State) Get(key Key) (value interface{}, exists bool) {
	if value, exists = s.values[key]; exists {
		return value, true
	}
	return nil, false
}

func (s *State) Subscribe(subscription *Subscription) {
	s.subscriptions[subscription.name] = subscription
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
		subscriptions: make(map[string]*Subscription),
		values:        make(map[Key]interface{}),
	}
}
