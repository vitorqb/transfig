package main

import (
	"reflect"
	"time"

	. "github.com/vitorqb/transfig/pkg/stategen"
)

func main() {
	_ = StateGen(GenNode{
		"CurrentPhase":   reflect.TypeFor[string](),
		"PreviousPahses": reflect.SliceOf(reflect.TypeFor[string]()),
		"Transaction": GenNode{
			"Date":        reflect.TypeFor[time.Time](),
			"Description": reflect.TypeFor[string](),
			"Tags":        reflect.SliceOf(reflect.TypeFor[string]()),
			"Postings":    reflect.SliceOf(reflect.TypeFor[string]()),
		},
	}, "foo", "tmp/foo.go")
}
