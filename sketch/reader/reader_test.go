package reader

import "testing"

func TestReadAtom(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "comment",
			input:    "; abc",
			expected: sComment("abc"),
		},
	}

	runTests(t, cases)
}

func TestReaderMacro(t *testing.T) {
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

func TestReadString(t *testing.T) {
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
	}

	runTests(t, cases)
}
