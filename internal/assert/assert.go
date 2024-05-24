package assert

import (
	"reflect"
	"testing"
)

// Equal checks that value a is equal to expected.
func Equal(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %#v, got %#v", want, got)
	}
}

// Panic checks if the given function f panics.
func Panic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Error("the code did not panic")
		}
	}()
	f() // function that panics
}
