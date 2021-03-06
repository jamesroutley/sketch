package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead_ModuleLookup(t *testing.T) {
	cases := []*TestCase{
		{
			name:  "module-lookup reader macro",
			input: `(strings.join (list "hello" "world"), " ")`,
			expected: sList(
				sList(
					sSym("module-lookup"), sSym("strings"), sSym("join"),
				),
				sList(sSym("list"), sStr("hello"), sStr("world")),
				sStr(" "),
			),
		},
	}

	runTests(t, cases)
}

func TestRead_Strings(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "string",
			input:    `"hello"`,
			expected: sStr(`hello`),
		},
		{
			name:     "string with newline",
			input:    "\"hello\nworld\"",
			expected: sStr("hello\nworld"),
		},
		{
			name:     "newline",
			input:    `"\n"`,
			expected: sStr("\n"),
		},
	}

	runTests(t, cases)
}

func TestRead_Comments(t *testing.T) {
	cases := []*TestCase{
		{
			name: "comment",
			input: `(hello
; comment
)`,
			expected: sList(sSym("hello")),
		},
	}

	runTests(t, cases)
}

func TestRead_HashMap(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "hash map with two items",
			input:    `{a b}`,
			expected: sHashMap(sSym("a"), sSym("b")),
		},
		{
			name:     "lists can be keys",
			input:    `{(a b) c}`,
			expected: sHashMap(sList(sSym("a"), sSym("b")), sSym("c")),
		},
	}

	runTests(t, cases)
}

func TestRead_List(t *testing.T) {
	cases := []*TestCase{
		{
			name:  "list",
			input: "(fold-left + 0 (list 1 2 3 4))",
			expected: sList(
				sSym("fold-left"),
				sSym("+"),
				sInt(0),
				sList(sSym("list"), sInt(1), sInt(2), sInt(3), sInt(4)),
			),
		},
	}
	runTests(t, cases)
}

func TestReadWithoutReaderMacros(t *testing.T) {
	cases := []*TestCase{
		{
			name: "comment",
			input: `(hello
; comment
)`,
			expected: sList(sSym("hello"), sComment("comment")),
		},
		{
			name:     "module lookup",
			input:    `(string.join abc)`,
			expected: sList(sSym("string.join"), sSym("abc")),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := ReadWithoutReaderMacros(tc.input)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
