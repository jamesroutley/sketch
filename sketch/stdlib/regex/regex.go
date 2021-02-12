package regex

import (
	"regexp"

	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
)

var EnvironmentItems = map[string]types.SketchType{}

func register(symbol string, f func(...types.SketchType) (types.SketchType, error)) {
	EnvironmentItems[symbol] = &types.SketchFunction{
		Func:      f,
		BoundName: symbol,
	}
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

	var matches []types.SketchType
	for _, f := range found {
		var matchItems []types.SketchType
		for _, m := range f {
			matchItems = append(matchItems, &types.SketchString{
				Value: m,
			})
		}
		// matches.Items = append(matches.Items, match)
		matches = append(matches, &types.SketchList{
			List: types.NewList(matchItems),
		})
	}
	return &types.SketchList{
		List: types.NewList(matches),
	}, nil
}
