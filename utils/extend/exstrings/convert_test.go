package exstrings

import (
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestUnsafeToBytes(t *testing.T) {
	s := "hello word"
	b := UnsafeToBytes("hello word")

	bptr := (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data
	sptr := (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	if bptr != sptr {
		t.Fatalf("bptr=%x sptr=%x", bptr, sptr)
	}

	if string(b) != s {
		t.Fatalf("string(b)=%s s=%s", string(b), s)
	}

	s = strings.Repeat("A", 3)
	b = UnsafeToBytes(s)
	b[0] = 'A'
	b[1] = 'B'
	b[2] = 'C'

	bptr = (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data
	sptr = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	if bptr != sptr {
		t.Fatalf("bptr=%x sptr=%x", bptr, sptr)
	}

	if string(b) != s {
		t.Fatalf("string(b)=%s s=%s", string(b), s)
	}
}
