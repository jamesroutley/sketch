// Package core implements Sketch's builtin functions and variables.
package core

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/types"
)

var EnvironmentItems = map[string]types.SketchType{}

// TODO: move to separate file and pull with go-bindata
var SketchCode = `
(def not (fn (a) (if a false true)))

(def load-file (fn (f) (eval (read-string (+ "(do " (slurp f) "\nnil)")))))
`

func register(symbol string, f func(...types.SketchType) (types.SketchType, error)) {
	EnvironmentItems[symbol] = &types.SketchFunction{Func: f}
}

func init() {
	register("+", add)
	register("-", subtract)
	register("*", multiply)
	register("/", divide)
	register("prn", prn)
	register("list", list)
	register("list?", isList)
	register("empty?", isEmpty)
	register("count", count)
	register("nth", nth)
	register("=", equals)
	register("<", lt)
	register("<=", lte)
	register(">", gt)
	register(">=", gte)
	register("read-string", readString)
	register("slurp", slurp)
	register("cons", cons)
	register("concat", concat)
	register("map", sketchMap)
	register("filter", filter)
	register("first", first)
	register("rest", rest)
	register("and", and)
	register("or", or)
}

func ArgsToSketchInt(args []types.SketchType) ([]*types.SketchInt, error) {
	numbers := make([]*types.SketchInt, len(args))
	for i, arg := range args {
		number, ok := arg.(*types.SketchInt)
		if !ok {
			return nil, fmt.Errorf("could not cast type to int")
		}
		numbers[i] = number
	}
	return numbers, nil
}
