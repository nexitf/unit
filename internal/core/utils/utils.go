package utils

import (
	"time"
)

// NewStoppedTimer creates a timer that will never fire.
func NewStoppedTimer() (t *time.Timer) {
	defer func() {
		t.Stop()
	}()
	return time.NewTimer(100 * 365 * 24 * time.Hour)
}
