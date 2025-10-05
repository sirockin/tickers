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
				start := time.Now()
				defer ticker.Stop()
				for i, expectedInterval := range tc.expectedIntervals {
					<-ticker.C
					elapsed := time.Since(start)
					if elapsed != expectedInterval {
						t.Fatalf("at interval %d expected interval of %v, got %v", i, expectedInterval, elapsed)
					}
					start = time.Now()
				}
			})
		})
	}
}

func TestExponentialWithJitter(t *testing.T) {
	cases := map[string]struct {
		initialDuration   time.Duration
		factor            float64
		jitter            time.Duration
		expectedIntervals []time.Duration
	}{
		"factor 2": {
			initialDuration: 1 * time.Second,
			factor:          2,
			jitter:          500 * time.Millisecond,
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
			jitter:          250 * time.Millisecond,
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
				ticker := tickers.NewExponential(tc.initialDuration, tc.factor, tickers.WithJitter(tc.jitter))
				defer ticker.Stop()
				for i, expectedInterval := range tc.expectedIntervals {
					time.Sleep(expectedInterval - 1*time.Millisecond)
					earliest := time.Now()
					synctest.Wait() // Make sure any goroutines are unlocked
					if _, ok := receivedValue(ticker.C); ok {
						t.Fatalf("Ticker channel should not have value before interval %d", i)
					}
					val := <-ticker.C
					now := time.Now()
					if val != now {
						t.Fatalf("expected %v but got %v", now, val)
					}
					if time.Since(earliest) > tc.jitter {
						t.Fatalf("Received too late")
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
