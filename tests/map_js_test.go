//go:build js && !wasm
// +build js,!wasm

package tests

import (
	"github.com/gopherjs/gopherjs/js"
	"testing"
)

func Test_MapWrapper(t *testing.T) {
	// This tests that various map types, and a map as a function argument and return,
	// wrap and unwrap correctly.
	type Dummy struct {
		Msg string
	}

	type StructWithMap struct {
		StringMap map[string]string
		IntMap    map[int]int
		DummyMap  map[string]*Dummy
		MapFunc   func(map[string]string) map[string]string
	}

	dummyMap := map[string]*Dummy{`key`: {Msg: `value`}}
	swm := &StructWithMap{
		StringMap: map[string]string{`key`: `value`},
		IntMap:    map[int]int{1: 2},
		DummyMap:  dummyMap,
		MapFunc: func(m map[string]string) map[string]string {
			return m
		},
	}
	swmWrapper := js.MakeFullWrapper(swm)
	swmUnwrapped := swmWrapper.Interface().(*StructWithMap)
	mapFuncArg := map[string]string{`key2`: `value2`}

	if swmWrapper.Get(`StringMap`).Get(`key`).String() != swm.StringMap[`key`] {
		t.Error(`StringMap did not match`)
	}
	if swmWrapper.Get(`IntMap`).Get(`1`).Int() != swm.IntMap[1] {
		t.Error(`IntMap did not match`)

	}
	if swmWrapper.Get(`DummyMap`).Get(`key`).Get(`Msg`).String() != swm.DummyMap[`key`].Msg {
		t.Error(`DummyMap did not match`)

	}
	if swmWrapper.Call(`MapFunc`, mapFuncArg).Get(`key2`).String() != mapFuncArg[`key2`] {
		t.Error(`MapFunc did not match`)

	}

	if swmUnwrapped.StringMap[`key`] != swm.StringMap[`key`] {
		t.Error(`Unwrapped StringMap did not match`)
	}
	if swmUnwrapped.IntMap[1] != swm.IntMap[1] {
		t.Error(`Unwrapped IntMap did not match`)
	}
	if swmUnwrapped.DummyMap[`key`].Msg != swm.DummyMap[`key`].Msg {
		t.Error(`Unwrapped DummyMap did not match`)
	}
	if swmUnwrapped.MapFunc(mapFuncArg)[`key2`] != swm.MapFunc(mapFuncArg)[`key2`] {
		t.Error(`Unwrapped MapFunc did not match`)
	}
}

func Test_MapStructObjectWrapper(t *testing.T) {
	// This tests that maps work as expected when wrapping a Struct with js.Object field containing a map.
	// js.Object fields' content should be passed to JS, so this is basically doubly-wrapping a map in two structs.

	stringMap := map[string]string{`key`: `value`}

	// You cannot wrap a map directly, so put it in a stuct.
	type StructWithMap struct {
		Map map[string]string
	}

	swm := &StructWithMap{Map: stringMap}
	swmWrapped := js.MakeFullWrapper(swm)

	// Now that map is wrapped in a struct, wrap the js.object in *another* struct.
	type StructWithObject struct {
		Wrappedswm *js.Object // This Object contains StructWithMap.
	}

	swo := &StructWithObject{Wrappedswm: swmWrapped}
	swoWrapper := js.MakeFullWrapper(swo)
	swmUnwrapped := swoWrapper.Interface().(*StructWithObject)

	// Using "Get(`Map`)" shows that the first wrapping was unchanged.
	if swoWrapper.Get(`Wrappedswm`).Get(`Map`).Get(`key`).String() != stringMap[`key`] {
		t.Error(`Wrapped Wrappedswm value did not match`)
	}

	if swmUnwrapped.Wrappedswm.Get(`Map`).Get(`key`).String() != stringMap[`key`] {
		t.Error(`Unwrapped Wrappedswm value did not match`)
	}
}

func Test_MapEmbeddedObject(t *testing.T) {
	o := js.Global.Get("JSON").Call("parse", `{"props": {"one": 1, "two": 2}}`)

	type data struct {
		*js.Object
		Props map[string]int `js:"props"`
	}

	d := data{Object: o}
	if d.Props["one"] != 1 {
		t.Error(`key "one" value should be 1`)
	}
	if d.Props["two"] != 2 {
		t.Error(`key "two" value should be 2`)
	}

}
