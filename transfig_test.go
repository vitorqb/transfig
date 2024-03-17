package transfig_test

import (
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/transfig"
	"testing"
)

var (
	Name Key = "name"
	Age  Key = "age"
)

func DefaultState() *State {
	x := NewState()
	x.Set(Name, "John")
	x.Set(Age, 30)
	return x
}

func Test_SimpleGet(t *testing.T) {
	state := DefaultState()
	name, exists := state.Get(Name)
	assert.True(t, exists)
	assert.Equal(t, "John", name.(string))
}

func Test_GetNonExistentKey(t *testing.T) {
	state := DefaultState()
	_, exists := state.Get("nonExistent")
	assert.False(t, exists)
}

func TestSubscribe_OneVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) {
		callbackArgs = args
	}
	sub := NewSubscription("subName").With(Name).Calls(callback)
	state.Subscribe(sub)
	state.Set("name", "John")
	expCallArgs := map[Key]interface{}{Name: "John"}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func TestSubscribe_TwoVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) {
		callbackArgs = args
	}
	sub := NewSubscription("subName").With("name", "age").Calls(callback)
	state.Subscribe(sub)
	state.Set("name", "Mike")
	expCallArgs := map[Key]interface{}{"name": "Mike", "age": 30}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func TestSubscribe_MissingVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) {
		callbackArgs = args
	}
	sub := NewSubscription("subName").With("name", "NOT_A_KEY").Calls(callback)
	state.Subscribe(sub)
	state.Set("name", "Mike")
	expCallArgs := map[Key]interface{}{"name": "Mike", "NOT_A_KEY": nil}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func TestSubscribe_TwoCallbacks(t *testing.T) {
	state := DefaultState()
	callbackArgs1 := make(map[Key]interface{})
	callback1 := func(args map[Key]interface{}) { callbackArgs1 = args }
	callbackArgs2 := make(map[Key]interface{})
	callback2 := func(args map[Key]interface{}) { callbackArgs2 = args }
	sub := NewSubscription("subName").With("name").Calls(callback1).Calls(callback2)
	state.Subscribe(sub)
	state.Set("name", "Mike")
	expCallArgs := map[Key]interface{}{"name": "Mike"}
	assert.Equal(t, expCallArgs, callbackArgs1)
	assert.Equal(t, expCallArgs, callbackArgs2)
}

func TestSubscribe_NotCalledIfNotSubscribed(t *testing.T) {
	state := DefaultState()
	callbackCalled := false
	callback := func(args map[Key]interface{}) { callbackCalled = true }
	sub := NewSubscription("subName").With("name").Calls(callback)
	state.Subscribe(sub)
	state.Set("NOT_NAME", "John")
	assert.False(t, callbackCalled)
}
