package rio

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// Server
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

func TestServer(t *testing.T) {
	t.Run("NewServer", func(t *testing.T) {
		server := NewServer()
		assert(t, len(server.middleware), 3)
	})
}

// ------------------------------------------------------------------
// ------------------------------------------------------------------
//
//
//
// Helpers
//
//
//
// ------------------------------------------------------------------
// ------------------------------------------------------------------

// assert checks that value a is equal to expected.
func assert(t *testing.T, a interface{}, expected interface{}) {
	if a == expected {
		return
	}

	_, filename, line, _ := runtime.Caller(1)
	msg := "%s:%d expected %v (type %v), got %v (type %v)\n"
	fmt.Printf(msg, filepath.Base(filename), line, expected, reflect.TypeOf(expected), a, reflect.TypeOf(a))
	t.FailNow()
}
