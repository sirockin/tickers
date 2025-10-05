package tickers_test

import (
	"fmt"
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
					got := <-ticker.C
					elapsed := time.Since(start)
					expected := time.Now()
					if got != expected {
						t.Fatalf("received value should be same as sent time: at interval %d expected channel value to be %v but got %v", i, expected, got)
					}
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
		initialDuration      time.Duration
		factor               float64
		jitter               time.Duration
		expectedMinIntervals []time.Duration
	}{
		"factor 2": {
			initialDuration: 1 * time.Second,
			factor:          2,
			jitter:          500 * time.Millisecond,
			expectedMinIntervals: []time.Duration{
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
			expectedMinIntervals: []time.Duration{
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
				start := time.Now()
				defer ticker.Stop()
				for i, minInterval := range tc.expectedMinIntervals {
					got := <-ticker.C
					elapsed := time.Since(start)
					expected := time.Now()
					if got != expected {
						t.Fatalf("received value should be same as sent time: at interval %d expected channel value to be %v but got %v", i, expected, got)
					}
					maxInterval := minInterval + tc.jitter
					if elapsed < minInterval || elapsed > maxInterval {
						t.Fatalf("at interval %d expected interval between %v and %v, got %v", i, minInterval, maxInterval, elapsed)
					}
					start = time.Now()
				}
			})
		})
	}
}

func TestExponentialStop(t *testing.T) {
	cases := []int{0, 2, 40}
	for _, numIntervals := range cases {
		t.Run(fmt.Sprintf("Stop after %d intervals", numIntervals), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ticker := tickers.NewExponential(1*time.Second, 2)

				for i := 0; i < numIntervals; i++ {
					fmt.Printf("Getting next interval: %d\n", i)
					<-ticker.C
				}
				ticker.Stop()

				synctest.Wait() // Wait for goroutines to exit
				if _, ok := receivedValue(ticker.C); ok {
					t.Fatalf("Ticker channel should not have value after Stop()")
				}
			})
		})
	}
}
