package transfig_test

import (
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/transfig"
	"testing"
)

var (
	Name         = KeyString("name")
	Age          = KeyString("age")
	Job          = KeyString("job")
	Title        = KeyString("title")
	Compensation = KeyString("compensation")
	Ammount      = KeyString("ammount")
	MissingKey   = KeyString("missingKey")
)

func DefaultState() *State {
	x := NewState()
	x.Set(Name, "John")
	x.Set(Age, 30)
	return x
}

func Test_Subscribe_OneVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) {
		callbackArgs = args
	}
	sub := NewSubscription("subName").With(Name).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[KeyString]interface{}{Name: "Mike"}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_TwoVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) {
		callbackArgs = args
	}
	sub := NewSubscription("subName").With(Name).With(Age).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[KeyString]interface{}{Name: "Mike", Age: 30}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_MissingVariable(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) {
		callbackArgs = args
	}
	missingKey := KeyString("NON_EXISTANT")
	sub := NewSubscription("subName").With(Name).With(missingKey).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[KeyString]interface{}{Name: "Mike", missingKey: nil}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_Wildcard(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbackArgs = args }
	sub := NewSubscription("subName").With(Wildcard{}).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[KeyString]interface{}{Name: "Mike", Age: 30}
	assert.Equal(t, expCallArgs, callbackArgs)
}

func Test_Subscribe_TwoCallbacks(t *testing.T) {
	state := DefaultState()
	callbackArgs1 := make(map[KeyString]interface{})
	callback1 := func(args CallbackArgs) { callbackArgs1 = args }
	callbackArgs2 := make(map[KeyString]interface{})
	callback2 := func(args CallbackArgs) { callbackArgs2 = args }
	sub := NewSubscription("subName").With(Name).Calls(callback1).Calls(callback2)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
	expCallArgs := map[KeyString]interface{}{Name: "Mike"}
	assert.Equal(t, expCallArgs, callbackArgs1)
	assert.Equal(t, expCallArgs, callbackArgs2)
}

func Test_Subscribe_NotCalledIfNotSubscribed(t *testing.T) {
	state := DefaultState()
	callbackCalled := false
	callback := func(args CallbackArgs) { callbackCalled = true }
	sub := NewSubscription("subName").With(Name).Calls(callback)
	state.Subscribe(sub)
	anotherKey := KeyString("NOT_NAME")
	state.Set(anotherKey, "John")
	assert.False(t, callbackCalled)
}

func Test_Subscribe_NestedValue(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Developer")

	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbackArgs = args }

	sub := NewSubscription("subName").With(Job).Calls(callback)
	state.Subscribe(sub)
	state.SetNested(Path{Job, Title}, "Manager")

	expCallArg := map[KeyString]interface{}{Job: map[KeyString]interface{}{Title: "Manager"}}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_TwoNestedState(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Dev")

	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbackArgs = args }

	sub := NewSubscription("subName").With(Job).Calls(callback)
	state.Subscribe(sub)

	state.SetNested(Path{Job, Compensation, Ammount}, 1000)

	expCallArg := map[KeyString]interface{}{
		Job: map[KeyString]interface{}{
			Title: "Dev",
			Compensation: map[KeyString]interface{}{
				Ammount: 1000,
			},
		},
	}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_NestedStateWithWildcard(t *testing.T) {
	state := DefaultState()
	jobMap := map[KeyString]interface{}{Title: "Developer"}
	state.Set(Job, jobMap)
	callCount := 0
	callback := func(args CallbackArgs) { callCount++ }

	sub := NewSubscription("subName").With(Wildcard{}).Calls(callback)
	state.Subscribe(sub)

	state.SetNested(Path{Title}, "Developer")

	assert.Equal(t, 1, callCount)
}

func Test_Subscribe_WithNested(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Developer")

	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbackArgs = args }

	sub := NewSubscription("subName").WithNested(Job, Title).Calls(callback)
	state.Subscribe(sub)
	state.SetNested(Path{Job, Title}, "Manager")

	expCallArg := map[KeyString]interface{}{Job: map[KeyString]interface{}{Title: "Manager"}}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_WithNestedThreeLong(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbackArgs = args }

	sub := NewSubscription("subName").WithNested(Job, Compensation, Ammount).Calls(callback)
	state.Subscribe(sub)
	state.SetNested(Path{Job, Compensation, Ammount}, 1000)

	expCallArg := map[KeyString]interface{}{
		Job: map[KeyString]interface{}{
			Compensation: map[KeyString]interface{}{
				Ammount: 1000,
			},
		},
	}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_WithNestedModfyingTop(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbackArgs = args }

	sub := NewSubscription("subName").WithNested(Job, Compensation, Ammount).Calls(callback)
	state.Subscribe(sub)

	state.Set(Job, map[KeyString]interface{}{})

	expCallArg := map[KeyString]interface{}{
		Job: map[KeyString]interface{}{
			Compensation: map[KeyString]interface{}{
				Ammount: nil,
			},
		},
	}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_WithNestedRepeatedKeys(t *testing.T) {
	state := DefaultState()
	callbackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbackArgs = args }

	KeyFoo := KeyString("foo")
	KeyBar := KeyString("bar")
	sub := NewSubscription("subName").WithNested(KeyFoo, KeyBar).Calls(callback)
	state.Subscribe(sub)

	state.Set(KeyBar, "Bar")

	expCallArg := map[KeyString]interface{}{}
	assert.Equal(t, expCallArg, callbackArgs)
}

func Test_Subscribe_WithRepeatedKeysNestedMap(t *testing.T) {
	state := DefaultState()
	callbackCount := 0
	callback := func(args CallbackArgs) { callbackCount++ }

	KeyFoo := KeyString("foo")
	KeyBar := KeyString("bar")
	KeyBaz := KeyString("baz")
	sub := NewSubscription("subName").WithNested(KeyFoo, KeyBar).Calls(callback)
	state.Subscribe(sub)

	state.SetNested(Path{KeyFoo, KeyBar}, "1")
	assert.Equal(t, 1, callbackCount)

	state.Set(KeyBar, "2")
	assert.Equal(t, 1, callbackCount)

	state.SetNested(Path{KeyFoo, KeyBaz}, "3")
	assert.Equal(t, 1, callbackCount)
}

func Test_Set_SetSubstateTwice(t *testing.T) {
	state := DefaultState()
	subState := map[KeyString]interface{}{Title: "Developer"}
	state.Set(Job, subState)

	callCount := 0
	callback := func(args CallbackArgs) { callCount++ }

	sub := NewSubscription("subName").With(Wildcard{}).Calls(callback)
	state.Subscribe(sub)

	state.SetNested(Path{Job, Title}, "Senior Developer")

	assert.Equal(t, 1, callCount)
}

func Test_Set_SetSubstateTwiceWithDifferentValues(t *testing.T) {
	state := DefaultState()
	subState1 := map[KeyString]interface{}{Title: "Developer"}
	state.Set(Job, subState1)
	subState2 := map[KeyString]interface{}{Title: "Manager"}
	state.Set(Job, subState2)

	callbcackArgs := make(map[KeyString]interface{})
	callback := func(args CallbackArgs) { callbcackArgs = args }

	sub := NewSubscription("subName").With(Job).Calls(callback)
	state.Subscribe(sub)

	state.SetNested(Path{Job, Title}, "Developer")
	state.SetNested(Path{Job, Title}, "Manager")

	expCallArgs := map[KeyString]interface{}{Job: map[KeyString]interface{}{Title: "Manager"}}
	assert.Equal(t, expCallArgs, callbcackArgs)
}

func Test_Get_Simple(t *testing.T) {
	state := DefaultState()
	state.Set(Job, "Developer")
	value, found := state.Get(Job)
	assert.True(t, found)
	assert.Equal(t, "Developer", value)
}

func Test_Get_Missing(t *testing.T) {
	state := DefaultState()
	value, found := state.Get(MissingKey)
	assert.False(t, found)
	assert.Nil(t, value)
}

func Test_GetNested_NoKeys(t *testing.T) {
	state := DefaultState()
	value, found := state.GetNested()
	assert.False(t, found)
	assert.Nil(t, value)
}

func Test_GetNested_LenghtOne(t *testing.T) {
	state := DefaultState()
	value, found := state.GetNested(Name)
	assert.True(t, found)
	assert.Equal(t, "John", value)
}

func Test_GetNested_LenghtTwo(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Developer")
	value, found := state.GetNested(Job, Title)
	assert.True(t, found)
	assert.Equal(t, "Developer", value)
}

func Test_GetNested_LenghtTwoMissing(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job}, "Developer")
	value, found := state.GetNested(Job, Title)
	assert.False(t, found)
	assert.Nil(t, value)
}

func Test_SetNested_Zero(t *testing.T) {
	state := NewState()
	state.SetNested(Path{}, "Mike")
	assert.Equal(t, CallbackArgs{}, state.AsMap())
}

func Test_SetNested_One(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Name}, "Mike")
	value, found := state.Get(Name)
	assert.True(t, found)
	assert.Equal(t, "Mike", value)
}

func Test_SetNested_Two(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Developer")
	value, found := state.GetNested(Job, Title)
	assert.True(t, found)
	assert.Equal(t, "Developer", value)
}

func Test_SetNested_OverrideExisting(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job}, "Foo")
	state.SetNested(Path{Job, Title}, "Manager")
	value, found := state.GetNested(Job, Title)
	assert.True(t, found)
	assert.Equal(t, "Manager", value)
	topValue, topFound := state.Get(Job)
	assert.True(t, topFound)
	assert.Equal(t, map[KeyString]interface{}{Title: "Manager"}, topValue)
}

func Test_SetNested_KeepExistingNestedKeys(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Developer")
	state.SetNested(Path{Job, Compensation, Ammount}, 1000)
	value, found := state.GetNested(Job, Compensation, Ammount)
	assert.True(t, found)
	assert.Equal(t, 1000, value)
	titleValue, titleFound := state.GetNested(Job, Title)
	assert.True(t, titleFound)
	assert.Equal(t, "Developer", titleValue)
	jobValue, jobFound := state.Get(Job)
	assert.True(t, jobFound)
	assert.Equal(
		t,
		map[KeyString]interface{}{
			Title:        "Developer",
			Compensation: map[KeyString]interface{}{Ammount: 1000},
		},
		jobValue,
	)
}

func Test_ClearNested_Empty(t *testing.T) {
	state := DefaultState()
	state.ClearNested(Path{})
	assert.Equal(t, DefaultState().AsMap(), state.AsMap())
}

func Test_ClearNested_OneLong(t *testing.T) {
	state := DefaultState()
	state.ClearNested(Path{Name})
	value, found := state.Get(Name)
	assert.False(t, found)
	assert.Nil(t, value)
}

func Test_ClearNested_TwoLong(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Developer")
	state.ClearNested(Path{Job, Title})
	value, found := state.GetNested(Job, Title)
	assert.False(t, found)
	assert.Nil(t, value)
}

func Test_ClearNested_FiresSubscriptions(t *testing.T) {
	state := DefaultState()
	state.SetNested(Path{Job, Title}, "Developer")
	callCount := 0
	callback := func(args CallbackArgs) { callCount++ }

	sub := NewSubscription("subName").WithNested(Job, Title).Calls(callback)
	state.Subscribe(sub)

	state.ClearNested(Path{Job, Title})

	assert.Equal(t, 1, callCount)
}

func Test_GetArg(t *testing.T) {
	state := DefaultState()
	callCount := 0
	callback := func(args CallbackArgs) {
		callCount++
		age, ageFound := GetArg[int](args, Age)
		assert.True(t, ageFound)
		assert.Equal(t, 30, age)
		name, nameFound := GetArg[string](args, Name)
		assert.True(t, nameFound)
		assert.Equal(t, "Mike", name)
		fake, fakeFound := GetArg[interface{}](args, KeyString("fake"))
		assert.False(t, fakeFound)
		assert.Nil(t, fake)
		wrongType, wrongTypeFound := GetArg[int](args, Name)
		assert.False(t, wrongTypeFound)
		assert.Equal(t, 0, wrongType)
	}
	sub := NewSubscription("subName").With(Name).With(Age).Calls(callback)
	state.Subscribe(sub)
	state.Set(Name, "Mike")
}
