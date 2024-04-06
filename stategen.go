package transfig

import (
	"reflect"

	jen "github.com/dave/jennifer/jen"
)

// GenNode represents a node in the state tree
type GenNode map[string]interface{}

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
		jen.Op("*").Qual("github.com/vitorqb/transfig", "State"),
	)
	if len(path) == 0 {
		f.Func().Id("New").Params(jen.Id("s").Op("*").Qual("github.com/vitorqb/transfig", "State")).Op("*").Id(structName).Block(
			jen.Return(jen.Op("&").Id(structName)).Values(jen.Id("s")),
		)
	}
	for key, node := range rootNode {
		if nodeAsNode, ok := node.(GenNode); ok {
			nextStructName := ""
			for _, k := range append(path, KeyString(key)) {
				nextStructName = nextStructName + string(k)
			}
			nextStructName = nextStructName + "NewState"
			f.Func().Params(jen.Id("s").Op("*").Id(structName)).Id(key).Params().Op("*").Id(nextStructName).Block(
				jen.Return(jen.Op("&").Id(nextStructName)).Values(jen.Id("s").Dot("State")),
			)
			gen(append(path, KeyString(key)), nodeAsNode, f)
		}
		jenPath := []jen.Code{}
		for _, k := range path {
			jenPath = append(jenPath, jen.Lit(string(k)))
		}
		jenPath = append(jenPath, jen.Lit(string(key)))
		if nodeAsType, ok := node.(reflect.Type); ok {
			isSlice := nodeAsType.Kind() == reflect.Slice
			pkgPath := nodeAsType.PkgPath()
			if isSlice {
				pkgPath = nodeAsType.Elem().PkgPath()
			}
			typeName := nodeAsType.Name()
			if isSlice {
				typeName = nodeAsType.Elem().Name()
			}

			// Getter
			getterReturnType := jen.Empty()
			if isSlice {
				getterReturnType.Add(jen.Index())
			}
			if pkgPath == "" {
				getterReturnType.Add(jen.Id(nodeAsType.String()))
			} else {
				getterReturnType.Add(jen.Qual(pkgPath, typeName))
			}
			getterReturnType = jen.Params(getterReturnType, jen.Bool())
			getter := f.Func().Params(jen.Id("s").Op("*").Id(structName)).Id(key).Params().Add(getterReturnType)
			getterAssertParam := jen.Empty()
			if isSlice {
				getterAssertParam.Add(jen.Index())
			}
			if pkgPath == "" {
				getterAssertParam.Add(jen.Id(nodeAsType.String()))
			} else {
				getterAssertParam.Add(jen.Qual(pkgPath, typeName))
			}
			getter = getter.Block(
				jen.List(jen.Id("v"), jen.Id("f")).Op(":=").Id("s").Dot("GetNested").Call(jenPath...),
				jen.If(jen.Op("!").Id("f")).Block(
					jen.Var().Id("v").Add(getterAssertParam),
					jen.Return(jen.Id("v"), jen.Id("f")),
				),
				jen.Return(jen.List(jen.Id("v").Assert(getterAssertParam)), jen.Id("f")),
			)

			// Setter
			setter := f.Func().Params(jen.Id("s").Op("*").Id(structName)).Id("Set" + key).Params(jen.Id("v").Add(getterAssertParam))
			setter = setter.Block(
				jen.Id("s").Dot("SetNested").Call(jen.Qual("github.com/vitorqb/transfig", "Path").Values(jenPath...), jen.Id("v")),
			)
		}
	}
}
