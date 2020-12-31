package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type MalType interface {
	String() string
	PrettyPrint(indent int) string
}

type EnvType interface {
	Set(string, MalType)
	Get(string) (MalType, error)
	Find(string) (EnvType, error)
}

type MalList struct {
	Items []MalType
}

func (l *MalList) String() string {
	itemStrings := make([]string, len(l.Items))
	for i, item := range l.Items {
		itemStrings[i] = item.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(itemStrings, " "))
}

func (l *MalList) PrettyPrint(indent int) string {
	items := l.Items
	if len(items) == 0 {
		return l.String()
	}

	// Shorter line length. Try getting the whole form. If it's less than 80
	// chars (including indent), print it all on one line
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
		// if i != len(args)-1 {
		// 	fmt.Fprintf(&b, "\n")
		// }
	}
	fmt.Fprintf(&b, ")")
	return b.String()

	// operator, ok := items[0].(*MalSymbol)
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

type MalInt struct {
	Value int
}

func (i *MalInt) String() string {
	return strconv.Itoa(i.Value)
}

func (i *MalInt) PrettyPrint(indent int) string {
	return strconv.Itoa(i.Value)
}

type MalSymbol struct {
	Value string
}

func (s *MalSymbol) String() string {
	return s.Value
}

func (s *MalSymbol) PrettyPrint(indent int) string {
	return s.String()
}

type MalFunction struct {
	Func              func(args ...MalType) (MalType, error)
	TailCallOptimised bool
	AST               MalType
	Params            []*MalSymbol
	Env               EnvType
	IsMacro           bool
}

func (f *MalFunction) String() string {
	return "#<function>"
}

func (f *MalFunction) PrettyPrint(indent int) string {
	return f.String()
}

type MalBoolean struct {
	Value bool
}

func (b *MalBoolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b *MalBoolean) PrettyPrint(indent int) string {
	return b.String()
}

type MalNil struct{}

func (n *MalNil) String() string {
	return "nil"
}

func (n *MalNil) PrettyPrint(indent int) string {
	return n.String()
}

type MalString struct {
	Value string
}

func (s *MalString) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

func (s *MalString) PrettyPrint(indent int) string {
	return s.String()
}

func getIndent(indent int) string {
	return strings.Repeat("  ", indent)
}
