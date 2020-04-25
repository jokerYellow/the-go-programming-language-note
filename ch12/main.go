package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	a := 10
	fmt.Println(formatAtom(a))
}

func any(i interface{}) string {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		b := bytes.Buffer{}
		b.WriteByte('[')
		items := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			items[i] = decode(v.Index(i))
		}
		b.WriteString(strings.Join(items, " "))
		b.WriteByte(']')
		return b.String()
	default:
		return formatAtom(i)
	}
}

func decode(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return "1"
		}
		return "0"
	case reflect.Complex64, reflect.Complex128:
		bitSize := 64
		if value.Kind() == reflect.Complex64 {
			bitSize = 32
		}
		b := bytes.Buffer{}
		c := value.Complex()
		b.WriteByte('(')
		b.WriteString(strconv.FormatFloat(real(c), 'g', -1, bitSize))
		if imag(c) >= 0 {
			b.WriteByte('+')
		}
		b.WriteString(strconv.FormatFloat(imag(c), 'g', -1, bitSize))
		b.WriteByte('i')
		b.WriteByte(')')
		return b.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Float32, reflect.Float64:
		bitSize := 64
		if value.Kind() == reflect.Float32 {
			bitSize = 32
		}
		return strconv.FormatFloat(value.Float(), 'g', -1, bitSize)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.String:
		return value.String()
	case reflect.Interface:
		return decode(value.Elem())
	case reflect.Invalid:
		return "invalid"
	}
	return ""
}

func formatAtom(i interface{}) string {
	value := reflect.ValueOf(i)
	return decode(value)
}
