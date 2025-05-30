package core

import (
	"github.com/nexitf/unit/internal/errors"
)

var (
	ErrDuplicateExternalName  = errors.New("duplicate external dependency name")
	ErrInvalidVariablePointer = errors.New("must be a valid pointer")
	ErrPluginNotFound         = errors.New("plugin not found")
	ErrAlreadyRunning         = errors.New("unit already running")
	ErrProcessTerminated      = errors.New("process terminated")
)
