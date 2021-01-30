package printer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jamesroutley/sketch/sketch/types"
)

func PrettyPrint(ast types.SketchType) string {
	return prettyPrint(ast, 0)
}

func prettyPrint(ast types.SketchType, indent int) string {
	switch ast := ast.(type) {
	case *types.SketchList:
		return prettyPrintList(ast, indent)
	default:
		return ast.String()
	}
}

// prettyPrintList returns a 'pretty' version of the list. The rules are:
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
func prettyPrintList(list *types.SketchList, indent int) string {
	items := list.List.ToSlice()
	if len(items) == 0 {
		return list.String()
	}

	containsComment := false
	for _, item := range items {
		if item.Type() == "comment" {
			containsComment = true
			break
		}
	}

	// If the whole list fits on an  80 char line (including indent), print it
	// all on one line. This is only safe to do if the list doesn't contain a
	// comment, because they need a newline after them to not accidentally
	// comment out too much
	if !containsComment {
		trial := list.String()
		if len(trial)+(indent*2) < 80 {
			return trial
		}
	}

	var b bytes.Buffer
	fmt.Fprintf(&b, "(%s", items[0])

	args := items[1:]
	for i, arg := range args {
		if arg.Type() == "comment" {
			// Inline comments 'stick' to the previous form. Print it on the
			// same line as that
			fmt.Fprintf(&b, " %s", prettyPrint(arg, indent+1))
			// If the comment is after the last list item, we need to insert
			// a newline, so the right paren printed after this for loop isn't
			// commented out. We also print some indentation to get the right
			// paren to be on the same line as the left paren.
			// This styling isn't ideal, but is at least syntactically correct.
			if i == len(args)-1 {
				fmt.Fprintf(&b, "\n%s", getIndent(indent))
			}
			continue
		}
		fmt.Fprintf(&b, "\n")
		fmt.Fprintf(&b, "%s%s", getIndent(indent+1), prettyPrint(arg, indent+1))
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

// TODO: switch to tabs
func getIndent(indent int) string {
	return strings.Repeat("  ", indent)
}
