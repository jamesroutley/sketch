package sketchtest

import "testing"

func TestStringJoin(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "string.join joins strings",
			input:    `(string.join (list "a" "b") "-")`,
			expected: `"a-b"`,
		},
		{
			name:     "string.join returns the element if one element in list",
			input:    `(string.join (list "hello") "-")`,
			expected: `"hello"`,
		},
		{
			name:     "string.join returns empty string for empty list",
			input:    `(string.join () "-")`,
			expected: `""`,
		},
	}
	runTestsWithImports(t, cases, "string")
}

func TestStringSplit(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "string.join joins strings",
			input:    `(string.split "a-b" "-")`,
			expected: `("a" "b")`,
		},
		{
			name:     "string.join joins strings",
			input:    `(string.split "a-b" "")`,
			expected: `("a" "-" "b")`,
		},
	}
	runTestsWithImports(t, cases, "string")
}
