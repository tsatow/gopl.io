package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type Encoder struct {
	w   io.Writer
	err error
}

func (encoder *Encoder) Encode(v interface{}) error {
	return encode(encoder.w, reflect.ValueOf(v))
}

func NewEncoder(w io.Writer) Encoder {
	return Encoder{w: w}
}

func encode(w io.Writer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		io.WriteString(w, "nil")
	case reflect.Bool:
		io.WriteString(w, "t")
	case reflect.Float32, reflect.Float64:
		// 合ってるかな？
		fmt.Fprintf(w, "%v", v.Float())
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		fmt.Fprintf(w, "#C(%v %v)", real(c), imag(c))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(w, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(w, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(w, "%q", v.String())
	case reflect.Ptr:
		return encode(w, v.Elem())
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
			io.WriteString(w, "(")
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					io.WriteString(w, " ")
				}
				if err := encode(w, v.Index(i)); err != nil {
					return err
				}
			}
			io.WriteString(w, ")")
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
				v.Field(i).IsNil()
				fmt.Fprintf(structBuf, "(%s ", v.Type().Field(i).Name)
				if err := encode(structBuf, v.Field(i)); err != nil {
					return err
				}
				structBuf.WriteByte(')')
			}
		}
		if numNonZeroField > 0 {
			io.WriteString(w, "(")
			w.Write(structBuf.Bytes())
			io.WriteString(w, ")")
		}
	case reflect.Map:
		if v.Len() != 0 {
			io.WriteString(w, "(")
			for i, key := range v.MapKeys() {
				if i > 0 {
					io.WriteString(w, " ")
				}
				io.WriteString(w, "(")
				if err := encode(w, key); err != nil {
					return err
				}
				io.WriteString(w, " ")
				if err := encode(w, v.MapIndex(key)); err != nil {
					return err
				}
				io.WriteString(w, ")")
			}
			io.WriteString(w, ")")
		}
	case reflect.Interface:
		fmt.Fprintf(w, `("%s" (`, v.Type())
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				io.WriteString(w, " ")
			}
			if err := encode(w, v.Field(i)); err != nil {
				return err
			}
		}
		io.WriteString(w, ")")
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
