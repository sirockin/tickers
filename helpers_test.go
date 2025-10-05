package tickers_test

// receivedValue checks if a value has been received on the channel
// it returns the value and true if a value is received, if not, the zero value
func receivedValue[T any](ch <-chan T) (T, bool) {
	select {
	case val := <-ch:
		return val, true
	default:
		var zero T
		return zero, false
	}
}
