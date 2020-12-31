package printer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jamesroutley/sketch/sketch/types"
)

func PrStr(t types.MalType) string {
	return t.String()
}

func PrettyPrint(ast types.MalType) string {
	return ast.PrettyPrint(0)
}

// When reading a program from a file, we implicitly wrap it in a (do ...)
// expression. When printing that (unevaluated) program AST, we don't want to
// include the (do) expression. This prints the ast, removing a top level `do`
// if there is one
func PrettyPrintTopLevelDo(ast types.MalType) string {
	list, ok := ast.(*types.MalList)
	if !ok {
		return PrettyPrint(ast)
	}

	if len(list.Items) == 0 {
		return PrettyPrint(ast)
	}

	symbol, ok := list.Items[0].(*types.MalSymbol)
	if !ok {
		return PrettyPrint(ast)
	}

	if symbol.Value != "do" {
		return PrettyPrint(ast)
	}

	// Okay, it's a `do` expression. Print each of the expressions in it,
	// separating each with a newline.
	var b bytes.Buffer

	for _, expr := range list.Items[1:] {
		fmt.Fprintln(&b, PrettyPrint(expr))
		fmt.Fprint(&b, "\n")
	}

	return strings.TrimSpace(b.String())
}
