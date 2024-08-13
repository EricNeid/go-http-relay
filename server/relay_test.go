// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EricNeid/go-webserver/internal/verify"
)

func TestRelay_get200(t *testing.T) {
	// arrange
	relayCallReceived := false
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		relayCallReceived = true
		w.WriteHeader(200)
	}))
	defer destination.Close()
	unit := NewApplicationServer(":5001", "", destination.URL)
	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()
	defer rec.Result().Body.Close()
	// action
	unit.relay(rec, req)
	// verify
	verify.Equals(t, 200, rec.Result().StatusCode)
	verify.Assert(t, relayCallReceived, "relay does not received call")
}

func TestRelay_get204(t *testing.T) {
	// arrange
	relayCallReceived := false
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		relayCallReceived = true
		w.WriteHeader(204)
	}))
	defer destination.Close()
	unit := NewApplicationServer(":5001", "", destination.URL)
	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()
	// action
	unit.relay(rec, req)
	// verify
	verify.Equals(t, 204, rec.Result().StatusCode)
	verify.Assert(t, relayCallReceived, "relay does not received call")
}
