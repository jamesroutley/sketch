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
