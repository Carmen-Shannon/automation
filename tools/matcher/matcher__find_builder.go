package matcher

import "time"

type findBuilderOption struct {
	MSEThreshold float64
	Timeout      time.Duration
}

type FindBuilderOption func(*findBuilderOption)

func MseThresholdOpt(threshold float64) FindBuilderOption {
	return func(opts *findBuilderOption) {
		opts.MSEThreshold = threshold
	}
}

func TimeoutOpt(timeout time.Duration) FindBuilderOption {
	return func(opts *findBuilderOption) {
		opts.Timeout = timeout
	}
}
