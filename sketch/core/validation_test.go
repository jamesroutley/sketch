package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToOrdinal(t *testing.T) {
	cases := []struct {
		n        int
		expected string
	}{
		{1, "1st"},
		{2, "2nd"},
		{3, "3rd"},
		{4, "4th"},
		{5, "5th"},
		{6, "6th"},
		{7, "7th"},
		{8, "8th"},
		{9, "9th"},
		{10, "10th"},
		{11, "11th"},
		{12, "12th"},
		{13, "13th"},
		{14, "14th"},
		{15, "15th"},
		{16, "16th"},
		{17, "17th"},
		{18, "18th"},
		{19, "19th"},
		{20, "20th"},
		{21, "21st"},
		{22, "22nd"},
		{23, "23rd"},
		{24, "24th"},
		{110, "110th"},
		{111, "111th"},
		{131, "131st"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprint(tc.n), func(t *testing.T) {
			t.Parallel()
			actual := toOrdinal(tc.n)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
