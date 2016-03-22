package zero

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// Copyright © 2016 Alexandr Razumov
// Copyright © 2016 Zenhotels
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Package zero provides a zero relation for arbitrary values.

import (
	"reflect"
	"unsafe"
)

func isZero(v reflect.Value, seen map[comparison]bool) bool {
	// cycle check
	if v.CanAddr() {
		ptr := unsafe.Pointer(v.UnsafeAddr())
		c := comparison{ptr, v.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool() == false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0

	case reflect.Ptr, reflect.Interface:
		return isZero(v.Elem(), seen)

	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i), seen) {
				return false
			}
		}
		return true

	case reflect.Slice, reflect.String, reflect.Map:
		return v.Len() == 0

	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if !isZero(v.Field(i), seen) {
				return false
			}
		}
		return true
	// reflect.Chan, reflect.UnsafePointer, reflect.Func
	default:
		return v.IsNil()
	}
}

// IsZero reports whether v is zero struct
func IsZero(v interface{}) bool {
	seen := make(map[comparison]bool)
	return isZero(reflect.ValueOf(v), seen)
}

type comparison struct {
	v unsafe.Pointer
	t reflect.Type
}
