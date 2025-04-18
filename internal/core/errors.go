package core

import (
	"errors"
)

var (
	ErrDuplicateExternalDependencyName = errors.New("duplicate external dependency name")
	ErrInvalidVariablePointer          = errors.New("must be a valid pointer")
	ErrVariableCanNotBeBound           = errors.New("variable cannot be bound")
)
