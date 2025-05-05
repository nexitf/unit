package stats

import (
	"time"

	"github.com/nexitf/unit/internal/core/utils"
)

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
	return utils.CallSpend(fn)
}
