package atomics

import (
	"testing"
)

func TestAtomicBool(t *testing.T) {
	ab := NewAtomicBool()
	if ab.Get() {
		t.Errorf("Default Value should be false")
	}
	ab.Set(true)
	if !ab.Get() {
		t.Errorf("Expected True")
	}

	ab.Set(false)
	for i := 0; i < 1000; i++ {
		go func() {
			// Flaky because it's not atomic.
			ab.Set(!ab.Get())
		}()
	}

	if ab.Get() {
		t.Errorf("Expected True")
	}
}
