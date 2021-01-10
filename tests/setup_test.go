package sketchtest

import (
	"fmt"
	"testing"

	"github.com/jamesroutley/sketch/sketch"
	"github.com/jamesroutley/sketch/sketch/evaluator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// All these tests take the same form
type TestCase struct {
	name          string
	input         string
	expected      string
	expectedError error
}

func runTests(t *testing.T, cases []*TestCase) {
	runTestsWithImports(t, cases)
}

func runTestsWithImports(t *testing.T, cases []*TestCase, imports ...string) {
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env, err := evaluator.RootEnvironment()
			require.NoError(t, err)

			for _, module := range imports {
				_, err := sketch.Rep(fmt.Sprintf(`(import "%s")`, module), env)
				require.NoError(t, err)
			}

			actual, err := sketch.Rep(tc.input, env)
			if tc.expectedError != nil {
				// TODO: assert on error message
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
