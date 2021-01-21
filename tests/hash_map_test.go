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

		{
			name:     "Hash map get",
			input:    "(hashmap-get {1 2} 1)",
			expected: "2",
		},
		{
			name:     "Hash map get with default",
			input:    `(hashmap-get {"a" true} "b" false)`,
			expected: "false",
		},

		{
			name:     "Hash map keys",
			input:    `(hashmap-keys {1 2 3 4})`,
			expected: "(1 3)",
		},
		{
			name:     "Hash map values",
			input:    `(hashmap-values {1 2 3 4})`,
			expected: "(2 4)",
		},
	}
	runTests(t, cases)
}
