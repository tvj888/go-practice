package main

import (
	"go.uber.org/goleak"
	"testing"
	"time"
)

func TestLeak(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Эта горутина "утечёт"
	go func() {
		time.Sleep(time.Hour)
	}()
}
