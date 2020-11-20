package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			var keyStr string
			switch key.Kind() {
			case reflect.Array:
				buf := bytes.NewBufferString(fmt.Sprintf("%s[%d]{ ", key.Type(), key.Len()))
				for i := 0; i < key.Len(); i++ {
					// 構造体や配列のネストは考えたいけど時間がない
					buf.WriteString(fmt.Sprintf("%s, ", formatAtom(key.Index(i))))
				}
				buf.WriteString("}")
				keyStr = buf.String()
			case reflect.Struct:
				buf := bytes.NewBufferString(fmt.Sprintf("%s{ ", key.Type()))
				for i := 0; i < key.NumField(); i++ {
					// 構造体や配列のネストは考えたいけど時間がない
					buf.WriteString(fmt.Sprintf("%s: %s, ", key.Type().Field(i).Name, formatAtom(key.Field(i))))
				}
				buf.WriteString("}")
				keyStr = buf.String()
			default:
				keyStr = formatAtom(key)
			}
			display(fmt.Sprintf("%s[%s]", path, keyStr), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default:
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + " value"
	}
}