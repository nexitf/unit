package core

import (
	"sync/atomic"
)

// IsRunning returns true if the unit is running.
func (u *Unit) IsRunning() bool {
	return atomic.LoadInt32(&u.running) == 1
}

// Errored returns true if the unit has errored.
func (u *Unit) Errored() bool {
	return atomic.LoadInt32(&u.errored) == 1
}
