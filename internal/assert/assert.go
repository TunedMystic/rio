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

// Error checks if the error `got` is equivalent to the error `want`.
// If `want` is nil, `got` must also be nil.
// If `want` is not nil, `got` must also not be nil, and their `Error()` messages must be identical.
func Error(t *testing.T, got, want error) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Errorf("want nil error, got %#v", got)
		}
		return
	}

	if got == nil {
		t.Errorf("want error %#v, got nil", want)
		return
	}

	if got.Error() != want.Error() {
		t.Errorf("error messages not equal:\n  want: %q\n  got:  %q", want.Error(), got.Error())
	}
}

// True checks if the given condition is true.
func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("want true, got false")
	}
}

// False checks if the given condition is false.
func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("want false, got true")
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
