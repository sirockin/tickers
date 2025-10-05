package tickers_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/sirockin/tickers"
)

func TestExponentialIntervals(t *testing.T) {
	cases := map[string]struct {
		initialDuration   time.Duration
		factor            float64
		expectedIntervals []time.Duration
	}{
		"factor 2": {
			initialDuration: 1 * time.Second,
			factor:          2,
			expectedIntervals: []time.Duration{
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				8 * time.Second,
			},
		},
		"factor 3": {
			initialDuration: 500 * time.Millisecond,
			factor:          3,
			expectedIntervals: []time.Duration{
				500 * time.Millisecond,
				1500 * time.Millisecond,
				4500 * time.Millisecond,
				13500 * time.Millisecond,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Run inside synctest to use mocked time
			synctest.Test(t, func(t *testing.T) {
				ticker := tickers.NewExponential(tc.initialDuration, tc.factor)
				defer ticker.Stop()
				for i, expectedInterval := range tc.expectedIntervals {
					time.Sleep(expectedInterval - 1*time.Millisecond)
					synctest.Wait() // Make sure any goroutines are unlocked
					if _, ok := receivedValue(ticker.C); ok {
						t.Fatalf("Ticker channel should not have value before interval %d", i)
					}
					time.Sleep(1 * time.Millisecond)
					// beforeWait := time.Now()
					synctest.Wait() // Make sure any goroutines are unlocked
					// if time.Since(beforeWait) == 0 {
					// 	t.Fatalf("some time should have elapsed")
					// }
					got, ok := receivedValue(ticker.C)
					if !ok {
						t.Fatalf("Ticker channel should have value after interval has elapsed")
					}
					expected := time.Now()
					if got != expected {
						t.Fatalf("Expected: %v but got %v", expected, got)
					}
				}
			})
		})
	}
}

func TestExponentialStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ticker := tickers.NewExponential(1*time.Second, 2)
		time.Sleep(999 * time.Millisecond)
		synctest.Wait()
		ticker.Stop()
		time.Sleep(100 * time.Hour) // We can provide a very long wait here because time is mocked
		synctest.Wait()
		if _, ok := receivedValue(ticker.C); ok {
			t.Fatalf("Ticker channel should not have value after Stop()")
		}
	})
}
