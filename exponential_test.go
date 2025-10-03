package tickers_test

import (
	"testing"
	"testing/synctest"
	"time"
	"github.com/stretchr/testify/assert"

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
					synctest.Wait()	// wait for background activity to complete
					assert.False(t, channelHasValue(ticker.C), "Ticker channel should not have value before interval %d", i)
					time.Sleep(1 * time.Millisecond)
					synctest.Wait()	// wait for background activity to complete
					assert.True(t, channelHasValue(ticker.C), "Ticker channel should have value after interval %d", i)
				}
			})
		})
	}
}

func TestExponentialStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ticker := tickers.NewExponential(1*time.Second, 2)
		time.Sleep(999 * time.Millisecond)
		synctest.Wait() // wait for background activity to complete
		assert.False(t, channelHasValue(ticker.C), "Ticker channel should not have value immediately after creation")
		ticker.Stop()
		time.Sleep(100 * time.Hour)	// I can provide a very long wait here because time is mocked
		synctest.Wait() // wait for background activity to complete
		assert.False(t, channelHasValue(ticker.C), "Ticker channel should not have value after Stop()")
	})
}
