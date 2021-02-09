package types

type StackFrame struct {
	Next         *StackFrame
	FunctionName string
}

type CallStack struct {
	Head *StackFrame
}

func NewCallStack() *CallStack {
	return &CallStack{}
}

func (c *CallStack) Push(frame *StackFrame) *CallStack {
	frame.Next = c.Head
	return &CallStack{
		Head: frame,
	}
}
