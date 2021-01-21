package sketchtest

import "testing"

func TestFold(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "Fold left",
			input:    "(fold-left + 0 (list 1 2 3 4))",
			expected: "10",
		},
		{
			name:     "Reduce",
			input:    "(reduce + (list 1 2 3 4))",
			expected: "10",
		},
	}

	runTests(t, cases)
}
