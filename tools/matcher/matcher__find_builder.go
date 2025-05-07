package matcher

import "time"

type findBuilderOption struct {
	Threshold float64
	Timeout   time.Duration
}

// FindBuilderOption is the builder option function for matcher package and it's associated uses.
type FindBuilderOption func(*findBuilderOption)

// ThresholdOpt sets the threshold for the MSE matching algorithm.
// This can be configured so that matches require less certainty or more to return a result.
// Depending on the size of the template and the scan, this can be as low as 10.0 or as high as 5000.0.
//
// Parameters:
//   - threshold: The threshold value for the MSE matching algorithm. This is a float64 value that determines how strict the matching should be.
//     A lower value means a stricter match, while a higher value means a more lenient match.
func ThresholdOpt(threshold float64) FindBuilderOption {
	return func(opts *findBuilderOption) {
		opts.Threshold = threshold
	}
}

// TimeoutOpt sets the timeout for the matching operation.
// This is in any duration format in time.Duration.
// This allows the matching flow to run until this threshold is reached, at which point it will stop the worker pool and return an error.
//
// Parameters:
//   - timeout: The timeout duration for the matching operation. This is a time.Duration value that determines how long the matching should wait before timing out.
func TimeoutOpt(timeout time.Duration) FindBuilderOption {
	return func(opts *findBuilderOption) {
		opts.Timeout = timeout
	}
}
