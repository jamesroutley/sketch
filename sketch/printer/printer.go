package printer

import "github.com/jamesroutley/sketch/sketch/types"

func PrStr(t types.MalType) string {
	return t.String()
}

func PrettyPrint(ast types.MalType) string {
	return ast.PrettyPrint(0)
}
