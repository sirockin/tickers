// Package tickers provides specialized ticker implementations with configurable timing behaviors.
package tickers

import (
	"math"
	"math/rand/v2"
	"sync"
	"time"
)

// Exponential is a ticker which provides a channel sending time values
// with an exponentially increasing interval.
// If the receiver is not ready when a tick is due, the tick is dropped
// and the interval continues to increase for the next tick.
type Exponential struct {
	// C is the channel on which the ticks are delivered.
	// Do not send to or close this channel; use Stop() to stop the ticker.
	C        chan time.Time
	done     chan struct{}
	stopOnce sync.Once
	interval time.Duration
	factor   float64
	jitter   time.Duration
}

// ExponentialOption is a function which can be provided to the NewExponential function
// to configure the exponential ticker.
type ExponentialOption func(*Exponential)

// WithJitter returns an ExponentialOption that adds random jitter to each interval.
// The jitter is a random duration between 0 and the provided jitter value that is added to each interval.
// Negative jitter values are coerced to zero.
func WithJitter(jitter time.Duration) ExponentialOption {
	return func(e *Exponential) {
		if jitter <= 0 {
			return
		}
		e.jitter = jitter
	}
}

// NewExponential creates and starts a new exponential ticker with the given initial duration and exponential factor.
// The ticker will send the current time on its channel at exponentially increasing intervals.
// The first tick occurs after initialDuration, the second after initialDuration * factor,
// the third after initialDuration * factor^2, and so on.
// Optional ExponentialOption functions can be provided to configure the ticker (e.g., WithJitter).
// Panics if initialDuration is not positive, if factor is not greater than 1.0, or if factor is NaN or infinite.
func NewExponential(initialDuration time.Duration, factor float64, opts ...ExponentialOption) *Exponential {
	if initialDuration <= 0 {
		panic("initialDuration must be positive")
	}
	if factor <= 1.0 {
		panic("factor must be greater than 1.0")
	}
	if math.IsNaN(factor) || math.IsInf(factor, 0) {
		panic("factor must be a finite number")
	}

	e := &Exponential{
		C:        make(chan time.Time),
		done:     make(chan struct{}),
		interval: initialDuration,
		factor:   factor,
	}

	for _, opt := range opts {
		opt(e)
	}

	go func() {
		currentInterval := e.interval
		nextInterval := currentInterval
		if e.jitter > 0 {
			nextInterval += time.Duration(rand.Int64N(int64(e.jitter))) // #nosec G404 -- using math/rand for timing jitter is acceptable
		}
		ticker := time.NewTimer(nextInterval)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				select {
				case e.C <- t:
					// Calculate next interval with overflow protection
					newInterval := float64(currentInterval) * e.factor
					if newInterval > float64(math.MaxInt64) {
						newInterval = float64(math.MaxInt64)
					}
					currentInterval = time.Duration(newInterval)

					nextInterval = currentInterval
					if e.jitter > 0 {
						nextInterval += time.Duration(rand.Int64N(int64(e.jitter))) // #nosec G404 -- using math/rand for timing jitter is acceptable
					}
					ticker.Reset(nextInterval)
				case <-e.done:
					return
				}
			case <-e.done:
				return
			}
		}
	}()

	return e
}

// Stop stops the exponential ticker and releases associated resources.
// After calling Stop, no more ticks will be sent on the ticker's channel.
// Stop is safe to call multiple times.
func (e *Exponential) Stop() {
	e.stopOnce.Do(func() {
		close(e.done)
	})
}
