package zero

import (
	"bytes"
	"fmt"
	"reflect"
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
	ID     int
	Detail Detail
	Data   TestDetailSubStructure
}

type Detail interface{}

func TestZero(t *testing.T) {
	one, zeroInt := 1, 0
	ch1 := make(chan int)
	var zeroChan chan int

	type myString string

	var interface1, interfaceZero interface{} = &one, &zeroInt

	var (
		zeroDetail1 Detail = &struct{}{}
		zeroDetail2 Detail = &TestDetail{}
		zeroDetail3 Detail = struct{}{}
		zeroDetail4 Detail = &TestDetail{}
		zeroDetail5 Detail = &TestDetail{Data: TestDetailSubStructure{Params: nil}}
		zeroDetail6 Detail = &TestDetail{Data: TestDetailSubStructure{
			Params: make([]TestDetailParam, 0, 10)},
		}

		nonZeroDetail1 Detail = &TestDetail{Data: TestDetailSubStructure{
			Params: []TestDetailParam{TestDetailParam{55}}},
		}
		nonZeroDetail2 Detail = &TestDetail{Data: TestDetailSubStructure{ID: 1234}}
		nonZeroDetail3 Detail = &TestDetail{ID: 1234}
		nonZeroDetail4 Detail = &TestDetail{Detail: nonZeroDetail3}
	)

	for i, test := range []struct {
		v    interface{}
		want bool
	}{
		// basic types
		{0, true},
		{complex(0, 0), true},
		{1, false},
		{1.0, false},
		{true, false},
		{0.0, true},
		{"foo", false},
		{"", true},
		{myString(""), true},     // different types
		{myString("foo"), false}, // different types
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
		{&interface1, false},
		{&interfaceZero, true},
		// special case for structures
		{zeroDetail1, true},
		{zeroDetail2, true},
		{zeroDetail3, true},
		{zeroDetail4, true},
		{zeroDetail5, true},
		{zeroDetail6, true},
		{nonZeroDetail1, false},
		{nonZeroDetail2, false},
		{nonZeroDetail3, false},
		{nonZeroDetail4, false},
	} {
		if IsZero(test.v) != test.want {
			t.Errorf("Zero(%v)[%d] = %t", test.v, i, !test.want)
		}
	}
}

func BenchmarkDetail(b *testing.B) {
	var nonZeroDetail1 Detail = &TestDetail{Data: TestDetailSubStructure{
		Params: []TestDetailParam{TestDetailParam{55}}},
	}
	for i := 0; i < b.N; i++ {
		IsZero(nonZeroDetail1)
	}
}

func BenchmarkDetailSimple(b *testing.B) {
	var nonZeroDetail1 Detail = &TestDetailParam{}
	for i := 0; i < b.N; i++ {
		IsZero(nonZeroDetail1)
	}
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func IsEmptyDetail(detail Detail) bool {
	v := reflect.ValueOf(detail).Elem()
	for i := 0; i < v.NumField(); i++ {
		if !isEmptyValue(v.Field(i)) {
			return false
		}
	}
	return true
}

func BenchmarkIsEmpty(b *testing.B) {
	var nonZeroDetail1 Detail = &TestDetail{Data: TestDetailSubStructure{Params: []TestDetailParam{TestDetailParam{55}}}}
	for i := 0; i < b.N; i++ {
		IsEmptyDetail(nonZeroDetail1)
	}
}

func BenchmarkIsEmptySimple(b *testing.B) {
	var nonZeroDetail1 Detail = &TestDetailParam{}
	for i := 0; i < b.N; i++ {
		IsEmptyDetail(nonZeroDetail1)
	}
}

type Structure struct {
	ID int
}

func ExampleStructure() {
	zeroStructure := Structure{}
	zeroStructurePointer := &zeroStructure
	nonZero := Structure{ID: 1}
	nonZeroPointer := &nonZero
	fmt.Println(IsZero(zeroStructure))        // true
	fmt.Println(IsZero(zeroStructurePointer)) // true
	fmt.Println(IsZero(nonZero))              // false
	fmt.Println(IsZero(nonZeroPointer))       // false
	// Output:
	// true
	// true
	// false
	// false
}
