package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRelay(t *testing.T) {
	// arrange
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer destination.Close()
	unit := NewApplicationServer(":5001", "", destination.URL)
	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()
	// action
	unit.relay(rec, req)

	// verify
}
