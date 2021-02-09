// Package errors implements a custom error type for use in Sketch. This error
// type contains things like the call stack leading up to an error. This error
// type is only returned by Eval. It's meant to represent an error that happens
// during the evaluation of Sketch code, not an error that happens during the
// non-evaluation of Sketch code.
// This error is maybe closer to an Exception type like in Python?
package errors

// type StackFrame struct {
// 	// Name of the function which pushed this f
// 	FunctionName string
// }

type Error struct {
	Err   error
	Stack []string
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func Wrap(err error, stack []string) error {
	xerr, ok := err.(*Error)
	if ok {
		// already wrapped - append to the stack
		xerr.Stack = append(xerr.Stack, stack...)
	}
	return &Error{
		Err:   err,
		Stack: stack,
	}
}
