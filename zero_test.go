package zero

import (
	"bytes"
	"fmt"
	"testing"
)

type TestDetailParam struct {
	ID int
}

type TestDetailSubStructure struct {
	ID     int
	Params []TestDetailParam
}

type TestDetail struct {
	ID   int
	Data TestDetailSubStructure
}

type Detail interface{}

func TestZero(t *testing.T) {
	one, zeroInt := 1, 0
	ch1 := make(chan int)
	var zeroChan chan int

	type mystring string

	var iface1, ifaceZero interface{} = &one, &zeroInt

	var (
		zeroDetail1 Detail = &struct{}{}
		zeroDetail2 Detail = &TestDetail{}
		zeroDetail3 Detail = struct{}{}
		zeroDetail4 Detail = &TestDetail{}
		zeroDetail5 Detail = &TestDetail{Data: TestDetailSubStructure{Params: nil}}
		zeroDetail6 Detail = &TestDetail{Data: TestDetailSubStructure{Params: make([]TestDetailParam, 0, 10)}}

		nonZeroDetail1 Detail = &TestDetail{Data: TestDetailSubStructure{Params: []TestDetailParam{TestDetailParam{55}}}}
	)

	for i, test := range []struct {
		v    interface{}
		want bool
	}{
		// basic types
		{0, true},
		{1, false},
		{1.0, false},
		{"foo", false},
		{"", true},
		{mystring(""), true},     // different types
		{mystring("foo"), false}, // different types
		// slices
		{[]string{"foo"}, false},
		{[]string(nil), true},
		{[]string{}, true},
		// maps
		{map[string][]int{"foo": {1, 2, 3}}, false},
		{map[string][]int{"foo": {1, 2, 3}}, false},
		{map[string][]int{}, true},
		{map[string][]int(nil), true},
		// pointers
		{&one, false},
		{&zeroInt, true},
		{new(bytes.Buffer), true},
		// functions
		{(func())(nil), true},
		{func() {}, false},
		// arrays
		{[...]int{1, 2, 3}, false},
		// channels
		{ch1, false},
		{zeroChan, true},
		// interfaces
		{&iface1, false},
		{&ifaceZero, true},
		// special case for structures
		{zeroDetail1, true},
		{zeroDetail2, true},
		{zeroDetail3, true},
		{zeroDetail4, true},
		{zeroDetail5, true},
		{zeroDetail6, true},
		{nonZeroDetail1, false},
	} {
		if IsZero(test.v) != test.want {
			t.Errorf("Zero(%v)[%d] = %t", test.v, i, !test.want)
		}
	}
}

func Example_equal() {
	fmt.Println(IsZero([]int{1, 2, 3}))      // "false"
	fmt.Println(IsZero([]string{"bar"}))     // "false"
	fmt.Println(IsZero([]string(nil)))       // "true"
	fmt.Println(IsZero([]string{}))          // "true"
	fmt.Println(IsZero(map[string]int{}))    // "true"
	fmt.Println(IsZero(map[string]int(nil))) // "true"

	// Output:
	// false
	// false
	// true
	// true
	// true
	// true
}
