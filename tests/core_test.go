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
