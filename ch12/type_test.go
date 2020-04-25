package main

import (
	"testing"
)

type item struct {
	value  interface{}
	expect string
}

func TestAny(t *testing.T) {
	cases := []item{
		{[]interface{}{float32(12.2), 2, 3}, "[12.2 2 3]"},
		{[]int{1, 2, 3}, "[1 2 3]"},
		{[]float32{1.1, 2.2, 311, 222.2}, "[1.1 2.2 311 222.2]"},
		{float32(11111111111111111112), "11111111111111111112"},
	}
	for _, c := range cases {
		o := any(c.value)
		expect := c.expect
		if o != expect {
			t.Errorf("error: any(%v)=%s", c.value, o)
		} else {
			t.Logf("any(%v)=%s", c.value, o)
		}
	}
}

func TestFormatAtom(t *testing.T) {
	cases := []item{
		{false, "0"},
		{true, "1"},
		{int(10), "10"},
		{int8(127), "127"},
		{int16(110), "110"},
		{int32(110), "110"},
		{int64(110), "110"},
		{int64(-110), "-110"},
		{uint8(127), "127"},
		{uint16(110), "110"},
		{uint32(110), "110"},
		{uint64(110), "110"},
		{uint64(18446744073709551615), "18446744073709551615"},
		{uint64(18446744073709551614), "18446744073709551614"},
		{float32(11.2331), "11.2331"},
		{float64(11.233), "11.233"},
		{"@", "@"},
		{complex(float32(12), float32(12)), "(12+12i)"},
		{complex(float32(12.12), float32(12.12)), "(12.12+12.12i)"},
		{complex(12, 12), "(12+12i)"},
		{complex(12, -12), "(12-12i)"},
		{complex(12, 0), "(12+0i)"},
		{complex(12, -0), "(12+0i)"},
	}
	for _, c := range cases {
		o := formatAtom(c.value)
		expect := c.expect
		if o != expect {
			t.Errorf("error: formatAtom(%v)=%s", c.value, o)
		} else {
			t.Logf("formatAtom(%v)=%s", c.value, o)
		}
	}
}
