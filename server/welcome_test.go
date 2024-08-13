// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT
package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EricNeid/go-webserver/internal/verify"
)

func TestWelcome(t *testing.T) {
	// arrange
	request := httptest.NewRequest("GET", "/", http.NoBody)
	responseRecorder := httptest.NewRecorder()
	// action
	welcome(responseRecorder, request)
	// verify
	verify.Assert(t, responseRecorder.Code == 200, fmt.Sprintf("Status code is %d\n", responseRecorder.Code))
	verify.Assert(t, responseRecorder.Body.String() == "Hello, World!", fmt.Sprintf("Body is %s\n", responseRecorder.Body.String()))
}
