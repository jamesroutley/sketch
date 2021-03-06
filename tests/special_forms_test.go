package sketchtest

import "testing"

func TestSpecialForm_Fn(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "fn defines a new function",
			input:    "(fn (x) (x))",
			expected: "#<function>",
		},
		{
			name:     "fn defines a new function, which is callable",
			input:    "((fn (x) (+ 1 x)) 1)",
			expected: "2",
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_FnVariadicArgs(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "& sumbol bundles any provided args into a list",
			input:    "((fn (& y) y) 1)",
			expected: "(1)",
		},
		{
			name:     "& sumbol bundles any provided args into a list (2 args)",
			input:    "((fn (& y) y) 1 2)",
			expected: "(1 2)",
		},
		{
			name:     "& sumbol bundles any provided args into a list (0 args)",
			input:    "((fn (& y) y))",
			expected: "()",
		},
		{
			name:     "you can mix named and variadic args",
			input:    "((fn (a b & y) (list a b y)) 1 2 3 4)",
			expected: "(1 2 (3 4))",
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_Def(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "def adds an item to the environment",
			input:    "(do (def a 1) a)",
			expected: "1",
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_Quote(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "quote stops evaluation of a form",
			input:    "(quote (1 1 1))",
			expected: "(1 1 1)",
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_QuasiquoteExpand(t *testing.T) {
	cases := []*TestCase{
		{
			name: "quasiquoteexpand expands quasiquote",
			input: `
(do
	(def a "world")
	(quasiquote (hello (unquote a))))`,
			expected: `(hello "world")`,
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_Defmacro(t *testing.T) {
	cases := []*TestCase{
		{
			name: "defmacro defines a new macro",
			input: `
(defmacro nil!  (fn
	(name)
	(quasiquote (def (unquote name) nil))))`,
			expected: "#<function>",
		},
		{
			name: "defmacro defines a new macro, which is callable",
			input: `
(do
	(defmacro nil!  (fn
		(name)
		(quasiquote (def (unquote name) nil))))
	(def a 1)
	(nil! a)
	a
)`,
			expected: "nil",
		},
	}
	runTests(t, cases)
}

func TestSpecialForm_Macroexpand(t *testing.T) {
	cases := []*TestCase{
		{
			name: "macroexpand expands and prints a macro, without evaluating it",
			input: `
(do
	(defmacro nil! (fn (name) (quasiquote (def (unquote name) nil))))
	(macroexpand (nil! a))
)
`,
			expected: "(def a nil)",
		},
	}
	runTests(t, cases)
}
