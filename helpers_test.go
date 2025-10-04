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

// receivedValue checks if a value has been received on the channel 
// it returns the value and true if a value is received, if not, the zero value
func receivedValue(ch <-chan time.Time) (time.Time, bool) {
	select {
	case val := <-ch:
		return val, true
	default:
		var zero time.Time
		return zero, false
	}
}
