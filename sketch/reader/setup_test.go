package reader

import (
	"testing"

	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name     string
	input    string
	expected types.SketchType
}

// The reader outputs an AST. Here, we implement some functions which make it
// terser to hand write ASTs

func sList(items ...types.SketchType) *types.SketchList {
	return &types.SketchList{Items: items}
}

func sSym(val string) *types.SketchSymbol {
	return &types.SketchSymbol{Value: val}
}

func sStr(val string) *types.SketchString {
	return &types.SketchString{Value: val}
}

func sComment(val string) *types.SketchComment {
	return &types.SketchComment{Value: val}
}

func sHashMap(vals ...types.SketchType) *types.SketchHashMap {
	m, err := types.NewSketchHashMap(vals)
	if err != nil {
		// Not ideal to panic, but I like not having to pass a *testing.T into
		// this constructor
		panic(err)
	}

	return m
}

func runTests(t *testing.T, cases []*TestCase) {
	t.Helper()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := Read(tc.input)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
