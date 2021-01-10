// Package core implements Sketch's builtin functions and variables.
package core

import (
	"github.com/jamesroutley/sketch/sketch/types"
)

var EnvironmentItems = map[string]types.SketchType{}

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
