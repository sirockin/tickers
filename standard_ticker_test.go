package tickers_test

import (
	"testing"
	"testing/synctest"
	"time"
)

// Test passes but uses real time so is slow and flaky
func TestStandardTickerWithRealClock(t *testing.T) {
	ticker := time.NewTicker(2 * time.Second)

	time.Sleep(1900 * time.Millisecond)
	if channelHasValue(ticker.C) {
		t.Fatalf("Ticker channel should not have value immediately after creation")
	}
	time.Sleep(200 * time.Millisecond)
	if !channelHasValue(ticker.C) {
		t.Fatalf("Ticker channel should have value after 2.1s")
	}
}

// Using synctest to control time makes the test fast and reliable
func TestStandardTickerWithSyncTest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ticker := time.NewTicker(2 * time.Second)

		time.Sleep(1999 * time.Millisecond)
		synctest.Wait() // wait for background activity to complete
		if channelHasValue(ticker.C) {
			t.Fatalf("Ticker channel should not have value immediately after creation")
		}
		time.Sleep(1 * time.Millisecond)
		synctest.Wait() // wait for background activity to complete
		if !channelHasValue(ticker.C) {
			t.Fatalf("Ticker channel should have value after 2.1s")
		}
	})
}

