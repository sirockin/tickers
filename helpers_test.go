package tickers_test

import "time"

func channelHasValue(ch <-chan time.Time) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
