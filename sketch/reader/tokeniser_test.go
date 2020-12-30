package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize_Comments(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "whole line commented",
			input:    "; (+ 1 1)",
			expected: []string{"; (+ 1 1)"},
		},
		{
			name: "one line commented, next not",
			input: `; (+ 1 1)
1`,
			// TODO: as this expected case demonstrates - it's a bit strange
			// that the tokenizer doesn't strip whitespace from the tokens
			expected: []string{"; (+ 1 1)", "\n1"},
		},
		{
			name:     "comment in middle of line",
			input:    `1 ; 2`,
			expected: []string{"1", " ; 2"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Tokenize(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestTokenize(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "list with single character items",
			input:    "(+ 1 1)",
			expected: []string{"(", "+", " 1", " 1", ")"},
		},
		{
			name:     "list with multiple character items",
			input:    "(add one two)",
			expected: []string{"(", "add", " one", " two", ")"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Tokenize(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
