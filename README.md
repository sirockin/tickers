## Tickers (Work in Progress)

A demonstration of the testing/synctest package.

### About synctest

Introduced in Go 1.25 the [testing/synctest](https://pkg.go.dev/testing/synctest) package allows test code to use the standard time package in a deterministic way, without waiting for real time to pass and without flakiness while background goroutines do their work.

The clock is mocked by wrapping the test function with `synctest.Test`. Calling `synctest.Wait` ensures that goroutines started inside the test complete their work.

### About this repository

We implement and test a new ticker type - `tickers.Exponential` which implements an exponential backoff ticker.


