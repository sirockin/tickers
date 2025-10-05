## Tickers

An implementation of an exponential Ticker - `tickers.Exponential` with a similar interface to `time.Ticker` but providing an exponentially increasing delay with optional Jitter.

This may be useful in itself but its main purpose is to demonstrate asynchronous testing using the [testing/synctest](https://pkg.go.dev/testing/synctest) package.

### About synctest

Introduced in Go 1.25 the [testing/synctest](https://pkg.go.dev/testing/synctest) package allows test code to use the standard time package in a deterministic way, without waiting for real time to pass and without flakiness while background goroutines do their work.

The clock is mocked by wrapping the test function with `synctest.Test`. Calling `synctest.Wait` ensures that goroutines started inside the test complete their work.



