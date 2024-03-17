package transfig_test

import (
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/transfig"
	"testing"
)

var (
	Name         Key = KeyString("name")
	Age          Key = KeyString("age")
	Job          Key = KeyString("job")
	Title        Key = KeyString("title")
	Compensation Key = KeyString("compensation")
	Ammount      Key = KeyString("ammount")
)

func DefaultState() *State {
	x := NewState()
	x.Set(Name, "John")
	x.Set(Age, 30)
	return x
}

func Test_Get(t *testing.T) {
	state := DefaultState()
	name := state.Get(Name)
	assert.Equal(t, "John", name.(string))
}

func Test_Get_NonExistentKey(t *testing.T) {
	state := DefaultState()
	value := state.Get(KeyString("nonExistent"))
	assert.Nil(t, value)
}

func Test_Get_Wildcard(t *testing.T) {
	state := DefaultState()
	value := state.Get(Wildcard)
	assert.Equal(t, map[Key]interface{}{Name: "John", Age: 30}, value)
}

func Test_Get_NestedState(t *testing.T) {
	state := DefaultState()
	subState := NewState()
	state.Set(Job, subState)
	subState.Set(Title, "Developer")
	value := state.Get(Job)
	assert.Equal(t, map[Key]interface{}{Title: "Developer"}, value)
}

func Test_Get_WildcardWithNestedState(t *testing.T) {
	state := DefaultState()
	subState := NewState()
	state.Set(Job, subState)
	subState.Set(Title, "Developer")
	value := state.Get(Wildcard)
	assert.Equal(t, map[Key]interface{}{
		Name: "John",
		Age:  30,
		Job: map[Key]interface{}{
			Title: "Developer",
		}},
		value,
	)
}

func Test_Subscribe_OneVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) {
		callbackArgs = args
	}
	sub := NewSubscription("subName").With(Name).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[Key]interface{}{Name: "Mike"}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_TwoVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) {
		callbackArgs = args
	}
	sub := NewSubscription("subName").With(Name, Age).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[Key]interface{}{Name: "Mike", Age: 30}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_MissingVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) {
		callbackArgs = args
	}
	missingKey := KeyString("NON_EXISTANT")
	sub := NewSubscription("subName").With(Name, missingKey).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[Key]interface{}{Name: "Mike", missingKey: nil}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_Wildcard(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) { callbackArgs = args }
	sub := NewSubscription("subName").With(Wildcard).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[Key]interface{}{Wildcard: map[Key]interface{}{Name: "Mike", Age: 30}}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_TwoCallbacks(t *testing.T) {
	state := DefaultState()
	callbackArgs1 := make(map[Key]interface{})
	callback1 := func(args map[Key]interface{}) { callbackArgs1 = args }
	callbackArgs2 := make(map[Key]interface{})
	callback2 := func(args map[Key]interface{}) { callbackArgs2 = args }
	sub := NewSubscription("subName").With(Name).Calls(callback1).Calls(callback2)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[Key]interface{}{Name: "Mike"}
	assert.Equal(t, expCallArgs, callbackArgs1)
	assert.Equal(t, expCallArgs, callbackArgs2)
}

func Test_Subscribe_NotCalledIfNotSubscribed(t *testing.T) {
	state := DefaultState()
	callbackCalled := false
	callback := func(args map[Key]interface{}) { callbackCalled = true }
	sub := NewSubscription("subName").With(Name).Calls(callback)
	state.Subscribe(sub)
	anotherKey := KeyString("NOT_NAME")
	state.Set(anotherKey, "John")
	assert.False(t, callbackCalled)
}

func Test_Subscribe_NestedState(t *testing.T) {
	state := DefaultState()
	subState := NewState()
	state.Set(Job, subState)

	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) { callbackArgs = args }

	sub := NewSubscription("subName").With(Job).Calls(callback)
	state.Subscribe(sub)

	subState.Set(Title, "Developer")

	expCallArg := map[Key]interface{}{Job: map[Key]interface{}{Title: "Developer"}}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_TwoNestedState(t *testing.T) {
	state := DefaultState()
	jobState := NewState()
	compensationState := NewState()
	state.Set(Job, jobState)
	jobState.Set(Compensation, compensationState)

	callbackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) { callbackArgs = args }

	sub := NewSubscription("subName").With(Job).Calls(callback)
	state.Subscribe(sub)

	compensationState.Set(Ammount, 1000)

	expCallArg := map[Key]interface{}{Job: map[Key]interface{}{Compensation: map[Key]interface{}{Ammount: 1000}}}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_NestedStateWithWildcard(t *testing.T) {
	state := DefaultState()
	subState := NewState()
	state.Set(Job, subState)
	callCount := 0
	callback := func(args map[Key]interface{}) { callCount++ }

	sub := NewSubscription("subName").With(Wildcard).Calls(callback)
	state.Subscribe(sub)

	subState.Set(Title, "Developer")

	assert.Equal(t, 1, callCount)
}

func Test_Set_SetSubstateTwice(t *testing.T) {
	state := DefaultState()
	subState := NewState()
	state.Set(Job, subState)
	state.Set(Job, subState)

	callCount := 0
	callback := func(args map[Key]interface{}) { callCount++ }

	sub := NewSubscription("subName").With(Wildcard).Calls(callback)
	state.Subscribe(sub)

	subState.Set(Title, "Developer")

	assert.Equal(t, 1, callCount)
}

func Test_Set_SetSubstateTwiceWithDifferentValues(t *testing.T) {
	state := DefaultState()
	subState1 := NewState()
	state.Set(Job, subState1)
	subState2 := NewState()
	state.Set(Job, subState2)

	callbcackArgs := make(map[Key]interface{})
	callback := func(args map[Key]interface{}) { callbcackArgs = args }

	sub := NewSubscription("subName").With(Job).Calls(callback)
	state.Subscribe(sub)

	subState1.Set(Title, "Developer")
	subState2.Set(Title, "Manager")

	expCallArgs := map[Key]interface{}{Job: map[Key]interface{}{Title: "Manager"}}
	assert.Equal(t, expCallArgs, callbcackArgs)
}
