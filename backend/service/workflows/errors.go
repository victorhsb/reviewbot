package workflows

import "errors"

type UserError struct {
	Message string
	Err     error
}

func (e *UserError) Error() string {
	return e.Err.Error()
}

var (
	// ErrExit is used to signal that the user wants to exit the workflow.
	ErrExit = errors.New("exit flow")
)
