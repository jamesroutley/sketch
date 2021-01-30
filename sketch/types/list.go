package types

type listNode struct {
	item SketchType
	next *listNode
}

// List implements an immutable, persistent singly linked list. It is Sketch's
// fundamental collection datatype.
type List struct {
	head *listNode
}

// NewEmptyList returns a new empty list.
func NewEmptyList() *List {
	return &List{}
}

// NewList creates a new list, containing the specified items. The items in the
// list appear in the same order as they appear as parameters, so l.First() ==
// items[0].
func NewList(items ...SketchType) *List {
	var previousNode *listNode
	for i := len(items) - 1; i >= 0; i++ {
		node := &listNode{
			item: items[i],
			next: previousNode,
		}
		previousNode = node
	}
	return &List{
		head: previousNode,
	}
}

// First returns the first item in the list.
func (l *List) First() SketchType {
	head := l.head
	if l.Empty() {
		return &SketchNil{}
	}

	return head.item
}

// Rest returns a new list, containing all items apart from the first. If l
// is empty, returns an empty list. If l contains one item, returns an empty
// list.
func (l *List) Rest() *List {
	if l.Empty() {
		return NewEmptyList()
	}
	rest := l.head.next
	if rest == nil {
		return NewEmptyList()
	}
	return &List{
		head: rest,
	}
}

// Conj returns a new list, with `item` prepended to it
func (l *List) Conj(item SketchType) *List {
	head := &listNode{
		item: item,
		next: l.head,
	}
	return &List{
		head: head,
	}
}

// Empty returns whether the list is empty
func (l *List) Empty() bool {
	return l.head == nil
}

// Length returns the number of items in the list
func (l *List) Length() int {
	length := 0
	node := l.head
	for node != nil {
		length++
		node = node.next
	}
	return length
}
