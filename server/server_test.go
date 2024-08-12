// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT
package server

import (
	"os"
	"testing"
)

func TestApplicationServer(t *testing.T) {
	t.Run("Server should shutdown after being interrupted", func(t *testing.T) {
		// arrange
		unit := NewApplicationServer(":5001", "")
		quit := make(chan os.Signal)
		done := make(chan bool)
		// action shutdown
		go unit.GracefullShutdown(quit, done)
		quit <- os.Interrupt
		// verify
		<-done
		// nothing to verify
	})
}
