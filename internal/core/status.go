package core

import (
	"sync/atomic"
)

// IsRunning returns true if the unit is running.
func (u *Unit) IsRunning() bool {
	return atomic.LoadInt32(&u.running) == 1
}
