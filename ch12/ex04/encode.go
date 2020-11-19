package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), ""); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value, offset string) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Bool:
		if v.Bool() {
			buf.WriteByte('t')
		} else {
			buf.WriteString("nil")
		}
	case reflect.Float32, reflect.Float64:
		// 合ってるかな？
		fmt.Fprintf(buf, "%v", v.Float())
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		fmt.Fprintf(buf, "#C(%v %v)", real(c), imag(c))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem(), offset)
	case reflect.Array, reflect.Slice:
		newOffset := ""
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				newOffset = offset + " "
				fmt.Fprintf(buf, "\n%s", newOffset)
			}
			if err := encode(buf, v.Index(i), newOffset); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		fieldNameOffset := ""
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			fieldName := v.Type().Field(i).Name
			fieldValueOffset := offset+strings.Repeat(" ", len(fieldName)+3)
			if i > 0 {
				buf.WriteByte('\n')
				fieldNameOffset = offset + " "
			}
			fmt.Fprintf(buf, "%s(%s ", fieldNameOffset, fieldName)
			if err := encode(buf, v.Field(i), fieldValueOffset); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Map:
		keyOffset := ""
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				keyOffset = offset + " "
				fmt.Fprintf(buf, "\n%s", keyOffset)
			}
			buf.WriteByte('(')
			if err := encode(buf, key, keyOffset); err != nil {
				return err
			}
			buf.WriteByte(' ')
			// TODO 構造体とか複数行にまたがる場合は縦配置にしたい
			if err := encode(buf, v.MapIndex(key), keyOffset+" "); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Interface:
		typeName := fmt.Sprintf(`("%s" (`, v.Type())
		buf.WriteString(typeName)
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			v.Type()
			if err := encode(buf, v.Field(i), offset+strings.Repeat(" ", len(typeName)+2)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
