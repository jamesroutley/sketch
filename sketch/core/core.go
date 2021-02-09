// Package core implements Sketch's builtin functions and variables.
package core

import (
	"github.com/jamesroutley/sketch/sketch/types"
)

var EnvironmentItems = map[string]types.SketchType{}

func register(symbol string, f func(...types.SketchType) (types.SketchType, error)) {
	EnvironmentItems[symbol] = &types.SketchFunction{
		Func:      f,
		BoundName: symbol,
	}
}

func init() {
	register("prn", prn)
	register("list", list)
	register("list?", isList)
	register("empty?", isEmpty)
	register("count", count)
	register("nth", nth)
	register("read-string", readString)
	register("slurp", slurp)
	register("cons", cons)
	register("concat", concat)
	register("first", first)
	register("rest", rest)
	register("and", and)
	register("or", or)
	register("string-to-list", stringToList)
	register("length", length)

	register("int", integer)

	register("+", add)
	register("-", subtract)
	register("*", multiply)
	register("/", divide)
	register("=", equals)
	register("<", lt)
	register("<=", lte)
	register(">", gt)
	register(">=", gte)
	register("modulo", modulo)

	register("apply", apply)

	register("hashmap", hashMap)
	register("hashmap-set", hashMapSet)
	register("hashmap-get", hashMapGet)
	register("hashmap-keys", hashMapKeys)
	register("hashmap-values", hashMapValues)

	register("map", sketchMap)
	register("filter", filter)
	register("fold-left", foldLeft)
	register("flatten", flatten)
	register("range", sketchRange)
}
