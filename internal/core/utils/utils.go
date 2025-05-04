package utils

import (
	"time"
)

// NewNeverTriggeringTimer returns a timer that is set to a very long duration
// (e.g., 100 years) and effectively will never fire on its own.
// This can be used as a placeholder or to block until explicitly reset or stopped.
func NewNeverTriggeringTimer() (t *time.Timer) {
	return time.NewTimer(100 * 365 * 24 * time.Hour)
}

// CallSpend measures and returns the time duration that the given function takes to execute.
// It captures the start time, runs the provided function, and returns the elapsed time.
//
// Example:
//
//	duration := CallSpend(func() {
//	    doSomething()
//	})
//	fmt.Println("Execution took:", duration)
func CallSpend(fn func()) time.Duration {
	st := time.Now()
	fn()
	return time.Since(st)
}
