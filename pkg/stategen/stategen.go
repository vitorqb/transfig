package stategen

import (
	"reflect"

	jen "github.com/dave/jennifer/jen"
	. "github.com/vitorqb/transfig"
)

const TransfigImportPath = "github.com/vitorqb/transfig"

// GenNode represents a node in the state tree
type GenNode map[string]interface{}

// StateGen generates code for a state that wraps an `State` object into a
// struct with getters and setters for each nested object in the state tree.
func StateGen(rootNode GenNode, packagePath string, filepath string) error {
	f := jen.NewFilePath(packagePath)
	gen(Path{}, rootNode, f)
	return f.Save(filepath)
}

// gen is a recursive code generator function used by StateGen. `path` is the
// current path in the state tree, `rootNode` is the current node in the state
// tree, and `f` is the file being generated.
func gen(path Path, rootNode GenNode, f *jen.File) {
	f.Add(stateStruct(path))
	if len(path) == 0 {
		f.Add(constructorFunc(path))
	}
	for key, node := range rootNode {
		if nodeAsNode, ok := node.(GenNode); ok {
			newPath := append(path, KeyString(key))
			f.Add(subStateGetter(newPath))
			gen(newPath, nodeAsNode, f)
			return
		}
		jenPath := []jen.Code{}
		for _, k := range path {
			jenPath = append(jenPath, jen.Lit(string(k)))
		}
		jenPath = append(jenPath, jen.Lit(string(key)))
		if nodeAsType, ok := node.(reflect.Type); ok {
			varType := typeFor(nodeAsType)
			// Getter
			f.Func().Params(jen.Id("s").Op("*").Id(stateStructName(path))).Id(key).Params().Params(varType, jen.Bool()).Block(
				jen.List(jen.Id("v"), jen.Id("f")).Op(":=").Id("s").Dot("GetNested").Call(jenPath...),
				jen.If(jen.Op("!").Id("f")).Block(
					jen.Var().Id("zero").Add(varType),
					jen.Return(jen.Id("zero"), jen.Id("f")),
				),
				jen.Return(jen.List(jen.Id("v").Assert(varType)), jen.Id("f")),
			)

			// Setter
			f.Func().Params(jen.Id("s").Op("*").Id(stateStructName(path))).Id("Set" + key).Params(jen.Id("v").Add(varType)).Block(
				jen.Id("s").Dot("SetNested").Call(jen.Qual(TransfigImportPath, "Path").Values(jenPath...), jen.Id("v")),
			)
		}
	}
}

func stateStructName(path Path) string {
	structName := ""
	for _, k := range path {
		structName = structName + string(k)
	}
	return structName + "NewState"
}

func stateStruct(path Path) *jen.Statement {
	return jen.Type().Id(stateStructName(path)).Struct(
		jen.Op("*").Qual(TransfigImportPath, "State"),
	)
}

func constructorFunc(path Path) *jen.Statement {
	return jen.Func().Id("New").Params(jen.Id("s").Op("*").Qual(TransfigImportPath, "State")).Op("*").Id(stateStructName(path)).Block(
		jen.Return(jen.Op("&").Id(stateStructName(path))).Values(jen.Id("s")),
	)
}

func subStateGetter(newPath Path) *jen.Statement {
	key := string(newPath[len(newPath)-1])
	rootPath := newPath[:len(newPath)-1]
	subStructName := stateStructName(newPath)
	return jen.Func().Params(jen.Id("s").Op("*").Id(stateStructName(rootPath))).Id(key).Params().Op("*").Id(subStructName).Block(
		jen.Return(jen.Op("&").Id(subStructName)).Values(jen.Id("s").Dot("State")),
	)
}

func typeFor(node reflect.Type) (o *jen.Statement) {
	o = jen.Empty()
	if node.Kind() == reflect.Slice {
		o.Add(jen.Index())
		node = node.Elem()
	}
	if pkgPath := node.PkgPath(); pkgPath != "" {
		o.Add(jen.Qual(pkgPath, node.Name()))
	} else {
		o.Add(jen.Id(node.String()))
	}
	return
}
