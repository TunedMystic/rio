package rio

import (
	"net/http"
	"testing"
)

func BenchmarkEvent(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := NewServer()
		s.HandleFunc("/", handleIndex)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
