package ginche

import (
	"testing"
)

func TestString(t *testing.T) {
	str := "test"
	pointer := String(str)
	if *pointer != str {
		t.Errorf("String() = %v, want %v", pointer, str)
	}
}

func TestBool(t *testing.T) {
	b := true
	pointer := Bool(b)
	if *pointer != b {
		t.Errorf("Bool() = %v, want %v", pointer, b)
	}
}

func TestInt(t *testing.T) {
	i := 1
	pointer := Int(i)
	if *pointer != i {
		t.Errorf("Int() = %v, want %v", pointer, i)
	}
}

func TestInt64(t *testing.T) {
	i := int64(1)
	pointer := Int64(i)
	if *pointer != i {
		t.Errorf("Int64() = %v, want %v", pointer, i)
	}
}

func TestFloat64(t *testing.T) {
	f := float64(1)
	pointer := Float64(f)
	if *pointer != f {
		t.Errorf("Float64() = %v, want %v", pointer, f)
	}
}

func TestDuration(t *testing.T) {
	d := Duration(1)
	pointer := Duration(*d)
	if *pointer != *d {
		t.Errorf("Duration() = %v, want %v", pointer, d)
	}
}
