// Package tickers provides specialized ticker implementations with configurable timing behaviors.
package tickers

import (
	"math/rand/v2"
	"time"
)

// Exponential is a ticker which provides a channel sending time values
// with an exponentially increasing interval.
type Exponential struct {
	// C is the channel on which the ticks are delivered.
	C        chan time.Time
	done     chan struct{}
	interval time.Duration
	factor   float64
	jitter   time.Duration
}

// ExponentialOption is a function which can be provided to the NewExponential function
// to configure the exponential ticker.
type ExponentialOption func(*Exponential)

// WithJitter returns an ExponentialOption that adds random jitter to each interval.
// The jitter is a random duration between 0 and the provided jitter value that is added to each interval.
func WithJitter(jitter time.Duration) ExponentialOption {
	return func(e *Exponential) {
		e.jitter = jitter
	}
}

// NewExponential creates and starts a new exponential ticker with the given initial duration and exponential factor.
// The ticker will send the current time on its channel at exponentially increasing intervals.
// The first tick occurs after initialDuration, the second after initialDuration * factor,
// the third after initialDuration * factor^2, and so on.
// Optional ExponentialOption functions can be provided to configure the ticker (e.g., WithJitter).
func NewExponential(initialDuration time.Duration, factor float64, opts ...ExponentialOption) *Exponential {
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
		nextInterval := e.interval
		if e.jitter > 0 {
			nextInterval += time.Duration(rand.Int64N(int64(e.jitter)))
		}
		ticker := time.NewTimer(nextInterval)

		for {
			select {
			case t := <-ticker.C:
				select {
				case e.C <- t:
					e.interval = time.Duration(float64(e.interval) * e.factor)
					nextInterval = e.interval
					if e.jitter > 0 {
						nextInterval += time.Duration(rand.Int64N(int64(e.jitter)))
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
func (e *Exponential) Stop() {
	close(e.done)
}
