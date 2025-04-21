package core

import (
	"errors"
)

var (
	ErrDuplicateExternalName  = errors.New("duplicate external dependency name")
	ErrInvalidVariablePointer = errors.New("must be a valid pointer")
	ErrVariableCanNotBeBound  = errors.New("variable cannot be bound")
	ErrProcessExited          = errors.New("process exited")
	ErrPluginNotFound         = errors.New("plugin not found")
	ErrAlreadyRunning         = errors.New("unit already running")
)
