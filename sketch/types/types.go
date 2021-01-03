package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type SketchType interface {
	String() string
	// PrettyPrint returns a 'pretty' version of the type. For non-lists, this
	// is just the type's value. For lists, see SketchList.PrettyPrint docstring.
	// This function powers sketchfmt
	// TODO: Because we currently don't currently read comments, they aren't
	// preserved during pretty printing.
	PrettyPrint(indent int) string
	// Returns a human readable name for the type
	Type() string
}

type EnvType interface {
	Set(string, SketchType)
	Get(string) (SketchType, error)
	Find(string) (EnvType, error)
}

type SketchList struct {
	Items []SketchType
}

func (l *SketchList) String() string {
	itemStrings := make([]string, len(l.Items))
	for i, item := range l.Items {
		itemStrings[i] = item.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(itemStrings, " "))
}

func (l *SketchList) Type() string {
	return "list"
}

// PrettyPrint returns a 'pretty' version of the list. The rules are:
// Check if the whole list + the indentation can fit in an 80 char line. If so,
// return the list on a single line.
// Else, return the list with each item printed on a separate line, in the form
// (a
//   b
//   c)
//
// TODO: we should consider having stricter formatting for certain special
// forms. For example, it might be nice to always print `case` statements on
// different lines.
func (l *SketchList) PrettyPrint(indent int) string {
	items := l.Items
	if len(items) == 0 {
		return l.String()
	}

	// If the whole list fits on an  80 char line (including indent), print it
	// all on one line
	trial := l.String()
	if len(trial)+(indent*2) < 80 {
		return trial
	}

	var b bytes.Buffer
	fmt.Fprintf(&b, "(%s", items[0])

	args := items[1:]
	for _, arg := range args {
		fmt.Fprintf(&b, "\n")
		fmt.Fprintf(&b, "%s%s", getIndent(indent+1), arg.PrettyPrint(indent+1))
	}
	fmt.Fprintf(&b, ")")
	return b.String()

	// operator, ok := items[0].(*SketchSymbol)
	// if !ok {
	// 	return l.String()
	// }

	// switch operator.Value {
	// case "do":
	// 	var b bytes.Buffer
	// 	fmt.Fprintf(&b, "(do\n")
	// 	for _, arg := range args {
	// 		fmt.Fprintf(&b, "%s%s\n", getIndent(indent+1), arg.PrettyPrint(indent+1))
	// 	}
	// 	fmt.Fprintf(&b, "%s)", getIndent(indent))
	// 	return b.String()

	// case "if":
	// 	var b bytes.Buffer
	// 	fmt.Fprintf(&b, "(if %s\n", args[0])
	// 	for _, arg := range args[1:] {
	// 		fmt.Fprintf(&b, "%s%s\n", getIndent(indent+1), arg.PrettyPrint(indent+1))
	// 	}
	// 	fmt.Fprintf(&b, "%s)", getIndent(indent))
	// 	return b.String()
	// }
	// return l.String()
}

type SketchInt struct {
	Value int
}

func (i *SketchInt) String() string {
	return strconv.Itoa(i.Value)
}

func (i *SketchInt) Type() string {
	return "int"
}

func (i *SketchInt) PrettyPrint(indent int) string {
	return strconv.Itoa(i.Value)
}

type SketchSymbol struct {
	Value string
}

func (s *SketchSymbol) String() string {
	return s.Value
}

func (s *SketchSymbol) Type() string {
	return "symbol"
}

func (s *SketchSymbol) PrettyPrint(indent int) string {
	return s.String()
}

type SketchFunction struct {
	Func              func(args ...SketchType) (SketchType, error)
	TailCallOptimised bool
	AST               SketchType
	Params            []*SketchSymbol
	Env               EnvType
	IsMacro           bool
}

func (f *SketchFunction) String() string {
	return "#<function>"
}

func (f *SketchFunction) Type() string {
	return "function"
}

func (f *SketchFunction) PrettyPrint(indent int) string {
	return f.String()
}

type SketchBoolean struct {
	Value bool
}

func (b *SketchBoolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b *SketchBoolean) Type() string {
	return "boolean"
}

func (b *SketchBoolean) PrettyPrint(indent int) string {
	return b.String()
}

type SketchNil struct{}

func (n *SketchNil) String() string {
	return "nil"
}

func (n *SketchNil) Type() string {
	return "nil"
}

func (n *SketchNil) PrettyPrint(indent int) string {
	return n.String()
}

type SketchString struct {
	Value string
}

func (s *SketchString) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

func (s *SketchString) Type() string {
	return "string"
}

func (s *SketchString) PrettyPrint(indent int) string {
	return s.String()
}

func getIndent(indent int) string {
	return strings.Repeat("  ", indent)
}
