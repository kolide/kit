// Package testutil provides utilities for use in tests.
package testutil

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// ErrorAfter will return an error after the provided timeout has passed if
// the provided WaitGroup has not unblocked. All Add() calls to the WaitGroup
// must be performed before calling this function.
func ErrorAfter(timeout time.Duration, wg *sync.WaitGroup) error {
	done := make(chan struct{})
	go func() {
		// Turn the blocking wg.Wait() into a selectable channel close
		wg.Wait()
		close(done)
	}()
	select {
	case <-time.After(timeout):
		return fmt.Errorf("test exceeded timeout: %v", timeout)
	case <-done:
		return nil
	}
}

// ErrorAfterFunc will return an error after the provided timeout has passed if
// the provided function has not returned.
func ErrorAfterFunc(timeout time.Duration, f func()) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
	return ErrorAfter(timeout, &wg)
}

// FatalAfterFunc will fatal the test after the provided timeout has passed if
// the provided function has not returned.
func FatalAfterFunc(t testing.TB, timeout time.Duration, f func()) {
	if err := ErrorAfterFunc(timeout, f); err != nil {
		t.Fatal(err)
	}
}

// FatalAfter will fatal the test after the provided timeout has passed if the
// provided WaitGroup has not unblocked. All Add() calls to the WaitGroup must
// be performed before calling this function.
func FatalAfter(t testing.TB, timeout time.Duration, wg *sync.WaitGroup) {
	if err := ErrorAfter(timeout, wg); err != nil {
		t.Fatal(err)
	}
}
