package sketchtest

import (
	"testing"
)

func TestIf(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "If true, return first arg",
			input:    "(if true 1 2)",
			expected: "1",
		},
		{
			name:     "If false, return second arg",
			input:    "(if false 1 2)",
			expected: "2",
		},
		{
			name:     "If false, and no second arg, return nil",
			input:    "(if false 1)",
			expected: "nil",
		},
	}
	runTests(t, cases)
}

func TestLet(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "let evaluates second arg in the newly created environment",
			input:    "(let ((a 1)) a)",
			expected: "1",
		},
		{
			name:     "let evaluates the even arguments in the parameter list",
			input:    "(let ((a (+ 1 1))) a)",
			expected: "2",
		},
		{
			name:     "later arguments in the parameter list can refer to earlier ones",
			input:    "(let ((a 1) (b (+ 1 a))) b)",
			expected: "2",
		},
	}
	runTests(t, cases)
}

func TestDef(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "def defines a new value",
			input:    "(do (def a 1) a)",
			expected: "1",
		},
	}
	runTests(t, cases)
}

func TestFn(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "fn defines a function closure",
			input:    "(fn (a) a)",
			expected: "#<function>",
		},
		{
			name:     "fn defines a function closure, which can be called",
			input:    "((fn (a) a) 100)",
			expected: "100",
		},
	}
	runTests(t, cases)
}

func TestRecursion(t *testing.T) {
	// This test is slow to run, because it recurses so deep. Because we only
	// use it to test that tail call optimisation works, we only run it
	// manually.
	t.Skip()
	cases := []*TestCase{
		{
			name: "deep recursion - this will overflow if `if expression` TCO not implemented",
			input: `
(do
	(def count-to (fn (num) (if (= num 0) nil (count-to (- num 1)))))
	(count-to 5000000)
)`,
			expected: "nil",
		},
	}
	runTests(t, cases)
}

func TestReadString(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "read string with no escaped chars",
			input:    `"hello world"`,
			expected: `"hello world"`,
		},
		{
			name:     "escaped double quote",
			input:    `"hello \" world"`,
			expected: `"hello \" world"`,
		},
	}
	runTests(t, cases)
}

func TestQuote(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "quote",
			input:    "(quote (1 1))",
			expected: "(1 1)",
		},
	}
	runTests(t, cases)
}
