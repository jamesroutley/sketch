package sketch

import (
	"testing"

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
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env, err := evaluator.RootEnvironment()
			require.NoError(t, err)
			actual, err := Rep(tc.input, env)
			if tc.expectedError != nil {
				// TODO: assert on error message
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
	t.Parallel()

}
