package tickers

import (
    "time"
)

type Exponential struct {
    C        chan time.Time
    done     chan struct{}
    interval time.Duration
    factor   float64
}

func NewExponential(initialDuration time.Duration, factor float64) *Exponential {
    e := &Exponential{
        C:        make(chan time.Time),
        done:     make(chan struct{}),
        interval: initialDuration,
        factor:   factor,
    }

    go func() {
        ticker := time.NewTicker(e.interval)
        defer ticker.Stop()

        for {
            select {
            case t := <-ticker.C:
                e.C <- t
                ticker.Stop()
                e.interval = time.Duration(float64(e.interval) * e.factor)
                ticker = time.NewTicker(e.interval)
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