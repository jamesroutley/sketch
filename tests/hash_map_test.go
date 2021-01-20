package sketchtest

import "testing"

func TestHashMap(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "Hash map literal",
			input:    "{1 2}",
			expected: "{1 2}",
		},
		{
			name:     "Hash map set",
			input:    "(do (def h {}) (hashmap-set h 1 2))",
			expected: "{1 2}",
		},
		{
			name:     "Hash map set twice",
			input:    "(do (def h {}) (def h (hashmap-set h 1 2)) (hashmap-set h 3 4))",
			expected: "{1 2 3 4}",
		},
	}
	runTests(t, cases)
}
