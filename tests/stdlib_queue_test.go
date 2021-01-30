package sketchtest

import "testing"

func TestQueue(t *testing.T) {
	cases := []*TestCase{
		{
			name:     "New queue",
			input:    "(queue.new 1 2 3)",
			expected: "((1 2 3) ())",
		},
		{
			name:     "Queue put",
			input:    "(queue.put (queue.new) 1)",
			expected: "(() (1))",
		},
		{
			name:     "Queue head, with items there",
			input:    "(queue.head (queue.new 2))",
			expected: "2",
		},
		{
			name:     "Queue head, force rebalance",
			input:    "(queue.head (queue.put (queue.new) 1))",
			expected: "1",
		},
		{
			name:     "Queue tail, with items there",
			input:    "(queue.tail (queue.new 1 2))",
			expected: "((2) ())",
		},
		{
			name:     "Queue tail, force rebalance",
			input:    "(queue.tail (queue.put (queue.put (queue.new) 1) 2))",
			expected: "((2) ())",
		},
	}

	runTestsWithImports(t, cases, "queue")
}
