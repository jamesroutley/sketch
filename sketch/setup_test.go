package sketch

import (
	"testing"

	"github.com/jamesroutley/sketch/sketch/core"
	"github.com/jamesroutley/sketch/sketch/environment"
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
			env := environment.NewEnv()
			for _, item := range core.Namespace {
				env.Set(item.Symbol.Value, item.Func)
			}
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
