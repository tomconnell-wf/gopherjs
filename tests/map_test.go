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

func Test_MapRangeMutation(t *testing.T) {
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

func Test_MapRangeNil(t *testing.T) {
	// Tests that loops on nil maps do not error.
	i := 0
	var m map[string]int
	for k, v := range m {
		_, _ = k, v
		i++
	}

	if i > 0 {
		t.Error(`Loops happened on a nil map`)
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

func Test_MapDelete(t *testing.T) {
	var nilMap map[string]string
	m := map[string]string{`key`: `value`}

	delete(nilMap, `key`) // noop
	delete(m, `key`)
	if m[`key`] == `value` {
		t.Error(`entry should have been deleted`)
	}
	delete(m, `key`) // noop
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
