package sketchtest

import "testing"

func TestStringToList(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "string without spaces",
			input:    `(string-to-list "abc")`,
			expected: `("a" "b" "c")`,
		},
		{
			name:     "empty string",
			input:    `(string-to-list "")`,
			expected: `()`,
		},
		{
			name:     "non-compount emoji",
			input:    `(string-to-list "🎛")`,
			expected: `("🎛")`,
		},
	}
	runTests(t, cases)
}

func TestFilter(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "filter even",
			input:    `(filter (fn (x) (= (modulo x 2) 0)) (list 1 2 3 4 5 6))`,
			expected: `(2 4 6)`,
		},
	}
	runTests(t, cases)
}

func TestMap(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "map add1",
			input:    "(map (fn (x) (+ x 1)) (list 1 2 3 4 5))",
			expected: "(2 3 4 5 6)",
		},
	}
	runTests(t, cases)
}

func TestCond(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "empty cond",
			input:    "(macroexpand (cond))",
			expected: "nil",
		},
		{
			name:     "cond with two cases",
			input:    "(macroexpand (cond (false 1) (true 2)))",
			expected: "(if false 1 (cond (true 2)))",
		},
		{
			name:     "cond with two cases",
			input:    "(cond (false 1) (true 2))",
			expected: "2",
		},
	}
	runTests(t, cases)
}

func TestCore(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "min",
			input:    "(min (list 0 1 2 3 4))",
			expected: "0",
		},
		{
			name:     "min negative",
			input:    "(min (list -1 4))",
			expected: "-1",
		},
		{
			name:     "min one item in list",
			input:    "(min (list 4))",
			expected: "4",
		},

		{
			name:     "max",
			input:    "(max (list -1 4))",
			expected: "4",
		},
		{
			name:     "max one item in list",
			input:    "(max (list 4))",
			expected: "4",
		},

		{
			name:     "apply",
			input:    "(apply + (list 1 2 3 4))",
			expected: "10",
		},

		{
			name:     "flatten",
			input:    "(flatten (list 1 2 (list 3 (list 4 5))))",
			expected: "(1 2 3 4 5)",
		},

		{
			name:     "hashset",
			input:    "(hashset 1 2 3)",
			expected: "{1 true 2 true 3 true}",
		},
		{
			name:     "empty hashset",
			input:    "(hashset)",
			expected: "{}",
		},

		{
			name:     "hashset-get: found",
			input:    "(hashset-get (hashset 1) 1)",
			expected: "true",
		},
		{
			name:     "hashset-get: not found",
			input:    "(hashset-get (hashset 1) 2)",
			expected: "false",
		},

		{
			name:     "range: one arg",
			input:    "(range 5)",
			expected: "(0 1 2 3 4)",
		},
		{
			name:     "range: two args",
			input:    "(range -2 3)",
			expected: "(-2 -1 0 1 2)",
		},
	}
	runTests(t, cases)
}
