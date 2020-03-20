package starlarkrepl

import (
	"reflect"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestAutoComplete(t *testing.T) {
	mod := &starlarkstruct.Module{
		Name: "hello",
		Members: starlark.StringDict{
			"world": starlark.String("world"),
			"dict": func() *starlark.Dict {
				d := starlark.NewDict(2)
				d.SetKey(starlark.String("key"), starlark.String("value"))
				return d
			}(),
		},
	}

	for _, tt := range []struct {
		name    string
		globals starlark.StringDict
		line    string
		want    []string
	}{{
		name: "simple",
		globals: map[string]starlark.Value{
			"abc": starlark.String("hello"),
		},
		line: "a",
		want: []string{"abc", "all", "any"},
	}, {
		name: "nest",
		globals: map[string]starlark.Value{
			"hello": mod,
		},
		line: "hello.wo",
		want: []string{"hello.world"},
		//}, {
		//	name: "dict",
		//	globals: map[string]starlark.Value{
		//		"hello": mod,
		//	},
		//	line: "hello.dict[\"",
		//	want: []string{"hello.dict[\"key\"]"},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			cmplter := newCompleter(tt.globals)
			got := cmplter(tt.line)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("%v != %v", tt.want, got)
			}
		})
	}
}
