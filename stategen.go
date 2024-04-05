package transfig

import (
	jen "github.com/dave/jennifer/jen"
	"reflect"
)

// GenNode represents a node in the state tree
type GenNode map[string]interface{}

// StateGen generates a new state tree
func StateGen(rootNode GenNode, filepath string) error {
	f := jen.NewFilePath(filepath)
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
		if nodeAsType, ok := node.(reflect.Type); ok {
			jenPath := []jen.Code{}
			for _, k := range path {
				jenPath = append(jenPath, jen.Lit(string(k)))
			}
			jenPath = append(jenPath, jen.Lit(string(key)))

			// Getter
			f.Func().Params(jen.Id("s").Id(structName)).Id(key).Params().Id(nodeAsType.String()).Block(
				jen.List(jen.Id("v"), jen.Id("_")).Op(":=").Id("s").Dot("wrappedState").Dot("GetNested").Call(jenPath...),
				jen.Return(jen.Id("v").Assert(jen.Id(nodeAsType.String()))),
			)

			// Setter
			f.Func().Params(jen.Id("s").Id(structName)).Id("Set" + key).Params(jen.Id("v").Id(nodeAsType.String())).Block(
				jen.Id("s").Dot("wrappedState").Dot("SetNested").Call(jen.Qual("github.com/vitorqb/transfig", "Path").Values(jenPath...), jen.Id("v")),
			)
		}
	}
}
