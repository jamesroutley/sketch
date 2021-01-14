// Package types defines the objects used to represent datatypes in Sketch
package types

import (
	"fmt"
	"strconv"
	"strings"
)

type SketchType interface {
	String() string
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

type SketchInt struct {
	Value int
}

func (i *SketchInt) String() string {
	return strconv.Itoa(i.Value)
}

func (i *SketchInt) Type() string {
	return "int"
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

type SketchFunction struct {
	Func              func(args ...SketchType) (SketchType, error)
	TailCallOptimised bool
	AST               SketchType
	Params            []*SketchSymbol
	Env               EnvType
	IsMacro           bool
	Docs              string
}

func (f *SketchFunction) String() string {
	return "#<function>"
}

func (f *SketchFunction) Type() string {
	return "function"
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

type SketchNil struct{}

func (n *SketchNil) String() string {
	return "nil"
}

func (n *SketchNil) Type() string {
	return "nil"
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

// SketchComment represents a comment in source code.
type SketchComment struct {
	Value string
}

func (c *SketchComment) String() string {
	// TODO: I think it will be necessary to return a newline after this
	return fmt.Sprintf("; %s", c.Value)
}

func (c *SketchComment) Type() string {
	return "comment"
}

type SketchModule struct {
	Environment EnvType
	SourceFile  string // filepath relative to $GOPATH/src
	Exported    []string
	// The name of the module, as specified in the `export-as` statement
	DefaultName string
	// The name used to refer to this module in the scope it's being used.
	// This will be different from DefaultName if (import-as) is used.
	Name string
}

func (m *SketchModule) String() string {
	return "#<module>"
}

func (m *SketchModule) Type() string {
	return "module"
}
