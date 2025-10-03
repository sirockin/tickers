## Tickers (Work in Progress)

A demonstration of the testing/synctest package.

### About synctest

Introduced in Go 1.25 the [testing/synctest](https://pkg.go.dev/testing/synctest) package allows test code to use the standard time package in a deterministic way, without waiting for real time to pass and without flakiness while background goroutines do their work.

The clock is mocked by wrapping the test function with `synctest.Test`. Calling `synctest.Wait` ensures that goroutines started inside the test complete their work.

### About this repository

We demonstrate a simple comparison test for the standard `time.Ticker` using both real time and `synctest`. The test using real time is slow and flaky, while the test using `synctest` is fast and reliable.

We also test implement a new ticker type - `tickers.Exponential` which implements an exponential backoff ticker. This is tested using `synctest` only.


