package printer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jamesroutley/sketch/sketch/types"
)

func PrStr(t types.SketchType) string {
	return t.String()
}

// When reading a program from a file, we implicitly wrap it in a (do ...)
// expression. When printing that (unevaluated) program AST, we don't want to
// include the (do) expression. This prints the ast, removing a top level `do`
// if there is one
func PrettyPrintTopLevelDo(ast types.SketchType) string {
	list, ok := ast.(*types.SketchList)
	if !ok {
		return PrettyPrint(ast)
	}

	if len(list.Items) == 0 {
		return PrettyPrint(ast)
	}

	symbol, ok := list.Items[0].(*types.SketchSymbol)
	if !ok {
		return PrettyPrint(ast)
	}

	if symbol.Value != "do" {
		return PrettyPrint(ast)
	}

	// Okay, it's a `do` expression. Print each of the expressions in it,
	// separating each with a newline.
	var b bytes.Buffer

	items := list.Items[1:]

	for i, expr := range items {
		fmt.Fprintln(&b, PrettyPrint(expr))

		// Comments above top level statements should 'stick' to them
		if expr.Type() == "comment" {
			continue
		}

		// Imports should be on consecutive lines. We check if the next item is
		// an import rather than this one so the final import has a newline
		// below it
		if i+1 < len(items) && isImportExpression(items[i+1]) {
			continue
		}

		fmt.Fprint(&b, "\n")
	}

	return strings.TrimSpace(b.String())
}

func isImportExpression(expr types.SketchType) bool {
	list, ok := expr.(*types.SketchList)
	if !ok {
		return false
	}
	items := list.Items
	if len(items) == 0 {
		return false
	}

	operator, ok := items[0].(*types.SketchSymbol)
	if !ok {
		return false
	}

	return operator.Value == "import"
}
