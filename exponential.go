package tickers

import (
	"math/rand/v2"
	"time"
)

type Exponential struct {
	C        chan time.Time
	done     chan struct{}
	interval time.Duration
	factor   float64
	jitter   time.Duration
}

type ExponentialOption func(*Exponential)

func WithJitter(jitter time.Duration) ExponentialOption {
	return func(e *Exponential) {
		e.jitter = jitter
	}
}

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
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				e.C <- t
				ticker.Stop()
				e.interval = time.Duration(float64(e.interval) * e.factor)
				nextInterval = e.interval
				if e.jitter > 0 {
					nextInterval += time.Duration(rand.Int64N(int64(e.jitter)))
				}
				ticker = time.NewTimer(nextInterval)
			case <-e.done:
				return
			}
		}
	}()

	return e
}

func (e *Exponential) Stop() {
	close(e.done)
}
