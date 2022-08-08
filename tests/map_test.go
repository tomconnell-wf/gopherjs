package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// These tests exercise the api of maps and built-in functions that operate on maps
func Test_MapLiteral(t *testing.T) {
	myMap := map[string]int{`test`: 0, `key`: 1, `charm`: 2}

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
	assert.Nil(t, *myMap)
	assert.IsType(t, map[string]int{}, *myMap)
}

func Test_MapType(t *testing.T) {
	var myMap map[string]int
	assert.Nil(t, myMap)
	assert.IsType(t, map[string]int{}, myMap)

	// assert.PanicsWithError has more than just the message `assignment to entry in nil map` in gopherjs runtime.
	// There's the stack, also.  There doesn't seem to be an assert api to check a message substring.
	assert.Panics(t, func() { myMap[`key`] = 666 })
}

func assertMapApi(t *testing.T, myMap map[string]int) {
	assert.IsType(t, map[string]int{}, myMap, `type assertion`)
	assert.Equal(t, 3, len(myMap), `initial len`)

	assert.Equal(t, 0, myMap[`test`], `access by key 1`)
	assert.Equal(t, 1, myMap[`key`], `access by key 2`)
	assert.Equal(t, 0, myMap[`missing`], `access by key 3`)

	charm, found := myMap[`charm`]
	assert.Equal(t, 2, charm, `tuple access by key, found`)
	assert.True(t, found, `tuple access by key, found`)
	missing2, found := myMap[`missing`]
	assert.Equal(t, 0, missing2, `tuple access by missing key, found`)
	assert.False(t, found, `tuple access by missing key, found`)

	delete(myMap, `missing`)
	assert.Equal(t, 3, len(myMap), `noop delete`)

	delete(myMap, `charm`)
	assert.Equal(t, 2, len(myMap), `delete`)

	myMap[`add`] = 3
	assert.Equal(t, 3, len(myMap), `assign by key len`)
	assert.Equal(t, 3, myMap[`add`], `assign by key`)

	myMap[`add`] = 10
	assert.Equal(t, 3, len(myMap), `update by key len`)
	assert.Equal(t, 10, myMap[`add`], `update by key`)

	copy := myMap
	assert.Equal(t, copy, myMap, `reference equality`)
}
