package tickers_test

import "time"

func valueHasArrived(ch <-chan time.Time) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func receivedValue(ch <-chan time.Time) (time.Time, bool) {
	select {
	case val := <-ch:
		return val, true
	default:
		var zero time.Time
		return zero, false
	}
}
