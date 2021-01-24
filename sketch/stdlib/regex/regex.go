package regex

import (
	"regexp"

	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
)

var EnvironmentItems = map[string]types.SketchType{}

func register(symbol string, f func(...types.SketchType) (types.SketchType, error)) {
	EnvironmentItems[symbol] = &types.SketchFunction{Func: f}
}

func init() {
	register("find", find)
}

func find(args ...types.SketchType) (types.SketchType, error) {
	pattern, err := validation.StringArg("find", args[0], 0)
	if err != nil {
		return nil, err
	}
	str, err := validation.StringArg("find", args[1], 1)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, err
	}
	found := re.FindAllStringSubmatch(str.Value, -1)

	matches := &types.SketchList{
		Items: []types.SketchType{},
	}
	for _, f := range found {
		match := &types.SketchList{
			Items: []types.SketchType{},
		}
		for _, m := range f {
			match.Items = append(match.Items, &types.SketchString{
				Value: m,
			})
		}
		matches.Items = append(matches.Items, match)
	}
	return matches, nil
}
