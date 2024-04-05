package transfig

import (
	"fmt"
	"reflect"

	jen "github.com/dave/jennifer/jen"
)

// GenNode represents a node in the state tree
type GenNode map[string]interface{}

// TypeWrapper represents a wrapper representing a type
type TypeWrapper interface {
	String() string
}
type StringType string
func (s StringType) String() string { return string(s) }
type ReflectType struct { t reflect.Type }
func (r ReflectType) String() string { return r.t.String() }
func (r ReflectType) Type() reflect.Type { return r.t }

func NewTypeWrapper(s interface{}) TypeWrapper {
	if v, ok := s.(string); ok {
		return StringType(v)
	}
	if v, ok := s.(reflect.Type); ok {
		return ReflectType{v}
	}
	return nil
}

// StateGen generates a new state tree
func StateGen(
	rootNode GenNode,
	packagePath string,
	filepath string,
) error {
	f := jen.NewFilePath(packagePath)
	gen(Path{}, rootNode, f)
	return f.Save(filepath)
}

func gen(path Path, rootNode GenNode, f *jen.File) {
	structName := ""
	for _, k := range path {
		structName = structName + string(k)
	}
	structName = structName + "NewState"
	f.Type().Id(structName).Struct(
		jen.Id("wrappedState").Op("*").Qual("github.com/vitorqb/transfig", "State"),
	)
	for key, node := range rootNode {
		if nodeAsNode, ok := node.(GenNode); ok {
			nextStructName := key + structName
			f.Func().Params(jen.Id("s").Id(structName)).Id(key).Params().Op("*").Id(nextStructName).Block(
				jen.Return(jen.Op("&").Id(nextStructName)).Values(jen.Id("s").Dot("wrappedState")),
			)
			gen(append(path, KeyString(key)), nodeAsNode, f)
		}
		nodeTypeWrapper := NewTypeWrapper(node)
		if nodeTypeWrapper == nil {
			continue
		}
		jenPath := []jen.Code{}
		for _, k := range path {
			jenPath = append(jenPath, jen.Lit(string(k)))
		}
		jenPath = append(jenPath, jen.Lit(string(key)))
		if reflectType, ok := nodeTypeWrapper.(ReflectType); ok {
			if pkgPath := reflectType.t.PkgPath(); pkgPath != "" {
				fmt.Println(pkgPath)
				f.ImportName(pkgPath, "")
			}
		}

		// Getter
		f.Func().Params(jen.Id("s").Id(structName)).Id(key).Params().Id(nodeTypeWrapper.String()).Block(
			jen.List(jen.Id("v"), jen.Id("_")).Op(":=").Id("s").Dot("wrappedState").Dot("GetNested").Call(jenPath...),
			jen.Return(jen.Id("v").Assert(jen.Id(nodeTypeWrapper.String()))),
		)

		// Setter
		f.Func().Params(jen.Id("s").Id(structName)).Id("Set" + key).Params(jen.Id("v").Id(nodeTypeWrapper.String())).Block(
			jen.Id("s").Dot("wrappedState").Dot("SetNested").Call(jen.Qual("github.com/vitorqb/transfig", "Path").Values(jenPath...), jen.Id("v")),
		)
	}
}
