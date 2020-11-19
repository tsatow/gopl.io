package equal

import (
	"reflect"
	"unsafe"
)

func isCirculation(x reflect.Value, seen map[unsafe.Pointer]bool) bool {
	if !x.IsValid() {
		return false
	}
	if x.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		if seen[xptr] {
			return true
		}
		seen[xptr] = true
	}

	switch x.Kind() {
	case reflect.Ptr, reflect.Interface:
		return isCirculation(x.Elem(), seen)

	case reflect.Array, reflect.Slice:
		for i := 0; i < x.Len(); i++ {
			if isCirculation(x.Index(i), seen) {
				return true
			}
		}
		return false

	case reflect.Map:
		for _, k := range x.MapKeys() {
			if isCirculation(x.MapIndex(k), seen) {
				return true
			}
		}
		return false

	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			if isCirculation(x.Field(i), seen) {
				return true
			}
		}
		return false

	default:
		return false
	}
}

func IsCirculation(x interface{}) bool {
	seen := make(map[unsafe.Pointer]bool)
	return isCirculation(reflect.ValueOf(x), seen)
}
