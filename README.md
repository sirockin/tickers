## Tickers

[![CI](https://github.com/sirockin/tickers/actions/workflows/ci.yml/badge.svg)](https://github.com/sirockin/tickers/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sirockin/tickers)](https://goreportcard.com/report/github.com/sirockin/tickers)
[![GoDoc](https://pkg.go.dev/badge/github.com/sirockin/tickers.svg)](https://pkg.go.dev/github.com/sirockin/tickers)
[![codecov](https://codecov.io/gh/sirockin/tickers/branch/main/graph/badge.svg)](https://codecov.io/gh/sirockin/tickers)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sirockin/tickers)](https://github.com/sirockin/tickers)
[![License](https://img.shields.io/github/license/sirockin/tickers)](https://github.com/sirockin/tickers/blob/main/LICENSE)

An implementation of an exponential Ticker - `tickers.Exponential` with a similar interface to `time.Ticker` but providing an exponentially increasing delay with optional Jitter.

This may be useful in itself but its main purpose is to demonstrate asynchronous testing using the [testing/synctest](https://pkg.go.dev/testing/synctest) package.

### About synctest

Introduced in Go 1.25 the [testing/synctest](https://pkg.go.dev/testing/synctest) package allows test code to use the standard time package in a deterministic way, without waiting for real time to pass and without flakiness while background goroutines do their work.

The clock is mocked by wrapping the test function with `synctest.Test`. Calling `synctest.Wait` ensures that goroutines started inside the test complete their work.



