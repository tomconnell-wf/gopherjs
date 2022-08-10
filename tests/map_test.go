package tests

import (
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
