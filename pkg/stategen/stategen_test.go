package stategen_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/transfig"
	. "github.com/vitorqb/transfig/pkg/stategen"
)

type TestStruct struct{}

func renderToString(t *testing.T, f *jen.File) string {
	var buff bytes.Buffer
	err := f.Render(&buff)
	assert.NoError(t, err)
	return buff.String()
}

func Test_StateStruct_ZeroPath(t *testing.T) {
	f := jen.NewFile("foo")
	path := Path{}
	f.Add(StateStruct(path))
	result := renderToString(t, f)
	assert.Contains(t, result, "import transfig \"github.com/vitorqb/transfig\"")
	assert.Contains(t, result, "type NewState struct {\n\t*transfig.State\n}\n")
}

func Test_StateStruct_TwoPath(t *testing.T) {
	f := jen.NewFile("foo")
	path := Path{"Foo", "Bar"}
	f.Add(StateStruct(path))
	result := renderToString(t, f)
	assert.Contains(t, result, "import transfig \"github.com/vitorqb/transfig\"")
	assert.Contains(t, result, "type FooBarNewState struct {\n\t*transfig.State\n}\n")
}

func Test_ConstructorFunc(t *testing.T) {
	f := jen.NewFile("foo")
	path := Path{}
	f.Add(ConstructorFunc(path))
	result := renderToString(t, f)
	assert.Contains(t, result, "func New(s *transfig.State) *NewState {\n\treturn &NewState{s}\n}\n")
}

func Test_SubStateGetter_PathLenOne(t *testing.T) {
	f := jen.NewFile("foo")
	path := Path{"Foo"}
	f.Add(SubStateGetter(path))
	result := renderToString(t, f)
	assert.Contains(t, result, "func (s *NewState) Foo() *FooNewState {\n\treturn &FooNewState{s.State}\n}")
}

func Test_TypeFor_Base(t *testing.T) {
	f := jen.NewFile("foo")
	nodeType := reflect.TypeFor[string]()
	f.Add(jen.Var().Id("v").Add(TypeFor(nodeType)))
	result := renderToString(t, f)
	assert.Contains(t, result, "var v string")
}

func Test_TypeFor_Slice(t *testing.T) {
	f := jen.NewFile("foo")
	nodeType := reflect.SliceOf(reflect.TypeFor[string]())
	f.Add(jen.Var().Id("v").Add(TypeFor(nodeType)))
	result := renderToString(t, f)
	assert.Contains(t, result, "var v []string")
}

func Test_TypeFor_SliceOfCustomStruct(t *testing.T) {
	f := jen.NewFile("foo")
	nodeType := reflect.SliceOf(reflect.TypeFor[TestStruct]())
	f.Add(jen.Var().Id("v").Add(TypeFor(nodeType)))
	result := renderToString(t, f)
	assert.Contains(t, result, "import stategentest \"github.com/vitorqb/transfig/pkg/stategen_test\"")
	assert.Contains(t, result, "var v []stategentest.TestStruct")
}
