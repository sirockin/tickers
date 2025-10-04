# Testing Asynchronous Operations in Go using `testing/synctest` 

## Introduction

Imagine we wanted to test go's  `time.NewTimer` function. Our first test might look something like this:

```go
func TestTimerFiresAtSpecifiedTime(t *testing.T) {
    interval := 2 * time.Second
    ticker := time.NewTimer(interval)

    time.Sleep(interval - 1*time.Millisecond)
    if valueHasArrived(ticker.C) {
        t.Fatalf("Timer channel should not have value before interval has elapsed")
    }
    time.Sleep(1 * time.Millisecond)
    if !valueHasArrived(ticker.C) {
        t.Fatalf("Ticker channel should have value once interval has elapsed")
    }
}

// valueHasArrived immediately returns true if a value can be read from the channel, false if not
func valueHasArrived(ch <-chan time.Time) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
```

We run it several times and find that it fails, sometimes at L12 sometimes L16. Thinking about it we realize that there are race conditions which mean that the timer may arrive before and after sleep. We decide to solve this by checking a bit earlier than the interval expiry time and a bit after:

```go
func TestTimerFiresAtSpecifiedTime(t *testing.T) {
	interval := 2 * time.Second
	beforeTolerance := 10 * time.Millisecond
	afterTolerance := 10 * time.Millisecond
	ticker := time.NewTimer(interval)

	time.Sleep(interval - beforeTolerance)
	if valueHasArrived(ticker.C) {
		t.Fatalf("Ticker channel should not have value before the interval has expired")
	}
	time.Sleep(beforeTolerance+afterTolerance)
	if !valueHasArrived(ticker.C) {
		t.Fatalf("Ticker channel should have value after interval has elapsed")
	}
}
```

We run it and it works. Great! Just to make sure, we run it 10 times:
```sh
$ go test -timeout 30s -count 10 -run ^TestTimerFiresAtSpecifiedTime$ .
ok      github.com/somebody/tickers     20.443s
```
Again that worked but it's taken over 20s since each test takes at least 2s. We then push it to CI and it fails. Of course - the environment's less performant, race conditions will be exacerbated. 

So we have a flakey slow test suite, which will get slower as we add more tests and cases.

These are common problems when testing asynchronous operations in any language. One of the ways of dealing with this has been to allow the injection of a mock clock into our implementation code and there are some `go` libraries {which ones} which provide this. Unfortunately they make the implementation code more complicated and have other issues such as {expand on this}

## Introducing `testing/synctest`

Fortunately `go` v1.25 an new `testing/synctest` library that tackles these issues. Let's rewrite our test:

```go
package tickers_test

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestTimerFiresAtSpecifiedTime(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		interval := 1 * time.Second
		ticker := time.NewTimer(interval)

		time.Sleep(interval - 1*time.Millisecond)   // Advance to 1s before
		if valueHasArrived(ticker.C) {
			t.Fatalf("Ticker channel should not have value before the interval has expired")
		}
		time.Sleep(1 * time.Millisecond)  // Advance to the exact interval
		if !valueHasArrived(ticker.C) {
			t.Fatalf("Ticker channel should have value after interval has elapsed")
		}
	})
}
```

So what's changed? [TODO]




