package testutil

import (
	"sync"
	"testing"
	"time"
)

func TestErrorAfterSuccess(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
	}()
	if err := ErrorAfter(1*time.Second, &wg); err != nil {
		t.Errorf("ErrorAfter should have succeeded, but failed with error: %v", err)
	}
}

func TestErrorAfterError(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}()
	if err := ErrorAfter(10*time.Millisecond, &wg); err == nil {
		t.Error("ErrorAfter should have errored")
	}
}

func TestErrorAfterFuncSuccess(t *testing.T) {
	err := ErrorAfterFunc(100*time.Millisecond, func() {
		time.Sleep(1 * time.Millisecond)
	})
	if err != nil {
		t.Errorf("ErrorAfterFunc should have succeeded, but failed with error: %v", err)
	}
}

func TestErrorAfterFuncError(t *testing.T) {
	err := ErrorAfterFunc(1*time.Millisecond, func() {
		time.Sleep(100 * time.Millisecond)
	})
	if err == nil {
		t.Error("ErrorAfterFunc should have errored")
	}
}

// Wrapper for testing.T with Fatal and Fatalf mocked
type mockFatal struct {
	testing.T
	// DidFatal records whether a call to Fatal or Fatalf occurred
	DidFatal bool
	// mut should be locked when accessing DidFatal
	mut sync.Mutex
}

func (m *mockFatal) Fatal(...interface{}) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.DidFatal = true
}

func (m *mockFatal) Fatalf(string, ...interface{}) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.DidFatal = true
}

func TestFatalAfterFuncSuccess(t *testing.T) {
	var m mockFatal

	// Should not fatal
	FatalAfterFunc(&m, 100*time.Millisecond, func() {
		time.Sleep(1 * time.Millisecond)
	})

	m.mut.Lock()
	defer m.mut.Unlock()
	if m.DidFatal {
		t.Error("FatalAfter should have succeeded")
	}
}

func TestFatalAfterFuncFatal(t *testing.T) {
	var m mockFatal

	// Should fatal
	FatalAfterFunc(&m, 1*time.Millisecond, func() {
		time.Sleep(100 * time.Millisecond)
	})

	m.mut.Lock()
	defer m.mut.Unlock()
	if !m.DidFatal {
		t.Error("FatalAfter should have fataled")
	}
}

func ExampleFatalAfterFunc_success() {
	var t testing.T
	// This test will pass
	FatalAfterFunc(&t, 10*time.Millisecond, func() {
		time.Sleep(1 * time.Millisecond)
	})
}

func ExampleFatalAfterFunc_fatal() {
	var t testing.T
	// This test will fatal
	FatalAfterFunc(&t, 10*time.Millisecond, func() {
		time.Sleep(1 * time.Millisecond)
	})
}

func TestFatalAfterSuccess(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
	}()

	var m mockFatal
	// Should not fatal
	FatalAfter(&m, 1*time.Second, &wg)
	m.mut.Lock()
	defer m.mut.Unlock()
	if m.DidFatal {
		t.Error("FatalAfter should have succeeded")
	}
}

func TestFatalAfterFatal(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}()

	var m mockFatal
	FatalAfter(&m, 10*time.Millisecond, &wg)

	// Don't allow test to exit before fatal occurs
	time.Sleep(20 * time.Millisecond)

	m.mut.Lock()
	defer m.mut.Unlock()
	if !m.DidFatal {
		t.Error("FatalAfter should have fataled")
	}
}
