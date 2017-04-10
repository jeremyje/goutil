package periodicfunction

import (
	"testing"
	"time"
)

var singleCallTests = []struct {
	interval     time.Duration
	sample       time.Duration
	expectedLow  int
	expectedHigh int
}{
	{10 * time.Millisecond, 500 * time.Millisecond, 49, 51},
}

func Test500Ms(t *testing.T) {
	for _, tt := range singleCallTests {
		counter := 0
		stableValue := 0

		pf := &PeriodicFunction{
			Function: func() { counter = counter + 1 },
			Interval: tt.interval,
		}
		pf.Start()
		time.Sleep(tt.sample)
		pf.Stop()
		stableValue = counter
		if !(counter >= tt.expectedLow && counter <= tt.expectedHigh) {
			t.Errorf("Expected {%d..%d} not %d", tt.expectedLow, tt.expectedHigh, counter)
		}

		time.Sleep(tt.sample)
		if counter != stableValue {
			t.Errorf("%d extra calls", counter-stableValue)
		}
	}
}
