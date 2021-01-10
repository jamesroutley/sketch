package sketchtest

import "testing"

func TestStrings(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "strings.join joins strings",
			input:    `(strings.join (list "a" "b") "-")`,
			expected: `"a-b"`,
		},
		{
			name:     "strings.join returns the element if one element in list",
			input:    `(strings.join (list "hello") "-")`,
			expected: `"hello"`,
		},
		{
			name:     "strings.join returns empty string for empty list",
			input:    `(strings.join () "-")`,
			expected: `""`,
		},
	}
	runTestsWithImports(t, cases, "strings")
}
