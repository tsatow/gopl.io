package main

import (
	"bytes"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Bool:
		buf.WriteByte('t')
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
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		// array、sliceでゼロ値をめぐる挙動が異なるのはかなり不親切。
		// arrayを常に出力するか、arrayのゼロ値の定義に合わせるかで悩む。
		// ここはsliceもarrayのゼロ値の定義に従う。
		nonZerosCount := 0
		for i := 0; i < v.Len(); i++ {
			if !isZeroValue(v.Index(i)) {
				nonZerosCount++
			}
		}
		if nonZerosCount > 0 {
			buf.WriteByte('(')
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					buf.WriteByte(' ')
				}
				if err := encode(buf, v.Index(i)); err != nil {
					return err
				}
			}
			buf.WriteByte(')')
		}
	case reflect.Struct:
		structBuf := bytes.NewBuffer(make([]byte, 0))
		numNonZeroField := 0
		for i := 0; i < v.NumField(); i++ {
			if !isZeroValue(v.Field(i)) {
				numNonZeroField++
				if i > 0 {
					structBuf.WriteByte(' ')
				}
				fmt.Fprintf(structBuf, "(%s ", v.Type().Field(i).Name)
				if err := encode(structBuf, v.Field(i)); err != nil {
					return err
				}
				structBuf.WriteByte(')')
			}
		}
		if numNonZeroField > 0 {
			buf.WriteByte('(')
			buf.Write(structBuf.Bytes())
			buf.WriteByte(')')
		}
	case reflect.Map:
		if v.Len() != 0 {
			buf.WriteByte('(')
			for i, key := range v.MapKeys() {
				if i > 0 {
					buf.WriteByte(' ')
				}
				buf.WriteByte('(')
				if err := encode(buf, key); err != nil {
					return err
				}
				buf.WriteByte(' ')
				if err := encode(buf, v.MapIndex(key)); err != nil {
					return err
				}
				buf.WriteByte(')')
			}
			buf.WriteByte(')')
		}
	case reflect.Interface:
		buf.WriteString(fmt.Sprintf(`("%s" (`, v.Type()))
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}
