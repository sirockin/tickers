package tickers

import "time"

type ExponentialTicker struct{}

func NewExponentialTicker(initialDuration time.Duration, factor float64) *ExponentialTicker {
	return &ExponentialTicker{}
}
