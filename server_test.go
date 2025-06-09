package rio

import (
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func TestServer(t *testing.T) {
	t.Run("NewServer", func(t *testing.T) {
		server := NewServer()
		assert.Equal(t, len(server.middleware), 3)
	})
}
