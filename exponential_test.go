package tickers_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStandardTicker(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ticker := time.NewTicker(2 * time.Second)

		time.Sleep(1900 * time.Millisecond)
		assert.False(t, channelHasValue(ticker.C), "Ticker channel should not have value immediately after creation")
		time.Sleep(200 * time.Millisecond)
		assert.True(t, channelHasValue(ticker.C), "Ticker channel should have value after 2.1s")
	})
}

func channelHasValue(ch <-chan time.Time) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
