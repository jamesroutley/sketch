package sketch

import (
	"errors"
	"testing"
)

func TestSpecialForm_Let(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "let creates a new environment, with variables bound",
			input:    "(let* (c 2) 2)",
			expected: "2",
		},
		{
			name: "bound variables aren't defined outside let",
			input: `
(do
	(let* (c 2) 2)
	c)`,
			expectedError: errors.New("`c` is undefined"),
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_If(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "if true",
			input:    "(if true 1 2)",
			expected: "1",
		},
		{
			name:     "if true",
			input:    "(if false 1 2)",
			expected: "2",
		},
		// In these two tests, 'undefined' is an unbound symbol. If that branch
		// were evaluated, the interpreter would throw an error
		// "`undefined` is undefined"
		{
			name:     "false branch isn't evaluated if true",
			input:    "(if true 1 undefined)",
			expected: "1",
		},
		{
			name:     "true branch isn't evaluated if false",
			input:    "(if false undefined 1)",
			expected: "1",
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_Do(t *testing.T) {
	cases := []*TestCase{
		// In this test, we implictly assert that each arg is evaluated because
		// the a will be undefined if the `(def ...)` line isn't evaluated
		{
			name: "do evaluates all args, and returns the last",
			input: `
(do
	(def! a 1)
	a)`,
			expected: "1",
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_Quasiquote(t *testing.T) {
	cases := []*TestCase{
		{
			name: "quasiquote ",
			input: `
(do
	(def! a "world")
	(quasiquote ("hello" (unquote a))))`,
			expected: `("hello" "world")`,
		},
	}
	runTests(t, cases)
}
