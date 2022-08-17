package tests

import (
	"github.com/gopherjs/gopherjs/js"
	"strings"
	"testing"
)

// These tests exercise the api of maps and built-in functions that operate on maps
func Test_MapLiteral(t *testing.T) {
	myMap := map[string]int{`test`: 0, `key`: 1, `charm`: 2}

	assertMapApi(t, myMap)
}

func Test_MapLiteralAssign(t *testing.T) {
	myMap := map[string]int{}
	myMap[`test`] = 0
	myMap[`key`] = 1
	myMap[`charm`] = 2

	assertMapApi(t, myMap)
}

func Test_MapMake(t *testing.T) {
	myMap := make(map[string]int)
	myMap[`test`] = 0
	myMap[`key`] = 1
	myMap[`charm`] = 2

	assertMapApi(t, myMap)
}

func Test_MapMakeSizeHint(t *testing.T) {
	myMap := make(map[string]int, 3)
	myMap[`test`] = 0
	myMap[`key`] = 1
	myMap[`charm`] = 2

	assertMapApi(t, myMap)
}

func Test_MapNew(t *testing.T) {
	myMap := new(map[string]int)
	if *myMap != nil {
		t.Error(`map should be nil when made with new()`)
	}
}

func Test_MapType(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error(`assignment on nil map should panic`)
		} else {
			str := err.(error).Error()
			if !strings.Contains(str, `assignment to entry in nil map`) {
				t.Error(`assigning to a nil map should panic`)
			}
		}
	}()

	var myMap map[string]int
	if myMap != nil {
		t.Error(`map should be nil when declared with var`)
	}

	myMap[`key`] = 666
}

func Test_MapLenPrecedence(t *testing.T) {
	// This test verifies that the expression len(m) compiles to is evaluated
	// correctly in the context of the enclosing expression.
	m := make(map[string]string)

	if len(m) != 0 {
		t.Error(`inline len should have been 0`)
	}

	i := len(m)
	if i != 0 {
		t.Error(`assigned len should have been 0`)
	}
}

func Test_MapRange(t *testing.T) {
	// This test verifies that all of a map is iterated, even if the map is modified during iteration.

	myMap := map[string]int{`one`: 1, `two`: 2, `three`: 3}

	var seenKeys []string

	for k, _ := range myMap {
		seenKeys = append(seenKeys, k)
		if k == `two` {
			delete(myMap, k)
		}
	}

	if len(seenKeys) != 3 {
		t.Error(`iteration seenKeys len was not 3`)
	}
}

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

func assertMapApi(t *testing.T, myMap map[string]int) {
	if len(myMap) != 3 {
		t.Error(`initial len of map should be 3`)
	}

	var keys []string
	var values []int

	for k, v := range myMap {
		keys = append(keys, k)
		values = append(values, v)
	}

	if len(keys) != 3 || !containsString(keys, `test`) || !containsString(keys, `key`) || !containsString(keys, `charm`) {
		t.Error(`range did not contain the correct keys`)
	}

	if len(values) != 3 || !containsInt(values, 0) || !containsInt(values, 1) || !containsInt(values, 2) {
		t.Error(`range did not contain the correct values`)
	}

	if myMap[`test`] != 0 {
		t.Error(`value should be 0`)
	}
	if myMap[`key`] != 1 {
		t.Error(`value should be 1`)
	}
	if myMap[`missing`] != 0 {
		t.Error(`absent key value should be 0`)
	}

	charm, found := myMap[`charm`]
	if charm != 2 {
		t.Error(`value should be 2`)
	}
	if !found {
		t.Error(`key should be found`)
	}

	missing2, found := myMap[`missing`]
	if missing2 != 0 {
		t.Error(`absent key value should be 0`)
	}
	if found {
		t.Error(`absent key should not be found`)
	}

	delete(myMap, `missing`)
	if len(myMap) != 3 {
		t.Error(`len after noop delete should still be 3`)
	}

	delete(myMap, `charm`)
	if len(myMap) != 2 {
		t.Error(`len after delete should still be 2`)
	}

	myMap[`add`] = 3
	if len(myMap) != 3 {
		t.Error(`len after assign by key should be 3`)
	}
	if myMap[`add`] != 3 {
		t.Error(`value should be 3`)
	}

	myMap[`add`] = 10
	if len(myMap) != 3 {
		t.Error(`len after update by key should be 3`)
	}
	if myMap[`add`] != 10 {
		t.Error(`value should be 10`)
	}

	myMap2 := myMap
	if len(myMap2) != len(myMap) {
		t.Error(`copy should be equivalent`)
	}
}

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// These benchmarks test various Map operations, and include a slice benchmark for reference.
const size = 10000

func makeMap(size int) map[int]string {
	myMap := make(map[int]string, size)
	for i := 0; i < size; i++ {
		myMap[i] = `data`
	}

	return myMap
}

func makeSlice(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i
	}

	return slice
}

func BenchmarkSliceLen(b *testing.B) {
	slice := makeSlice(size)

	for i := 0; i < b.N; i++ {
		if len(slice) > 0 {
		}
	}
}

func BenchmarkMapLen(b *testing.B) {
	myMap := makeMap(size)

	for i := 0; i < b.N; i++ {
		if len(myMap) > 0 {
		}
	}
}

func BenchmarkMapNilCheck(b *testing.B) {
	myMap := makeMap(size)

	for i := 0; i < b.N; i++ {
		if myMap != nil {
		}
	}
}

func BenchmarkMapNilElementCheck(b *testing.B) {
	myMap := makeMap(size)

	for i := 0; i < b.N; i++ {
		if myMap[0] != `` {
		}
	}
}

func BenchmarkSliceRange(b *testing.B) {
	slice := makeSlice(size)

	for i := 0; i < b.N; i++ {
		for range slice {
		}
	}
}

func BenchmarkMapRange(b *testing.B) {
	myMap := makeMap(size)

	for i := 0; i < b.N; i++ {
		for range myMap {
		}
	}
}
