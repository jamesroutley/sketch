package errors

import (
	"github.com/jamesroutley/sketch/sketch/environment"
)

type Error struct {
	Err error
	Env *environment.Env
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func Wrap(err error, env *environment.Env) error {
	_, ok := err.(*Error)
	if ok {
		// already wrapped
		return err
	}
	return &Error{
		Err: err,
		Env: env,
	}
}
