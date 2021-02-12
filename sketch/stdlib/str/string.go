package str

import (
	"strings"

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
	register("split", split)
	register("fields", fields)
}

func split(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("spit", 2, args); err != nil {
		return nil, err
	}

	s, err := validation.StringArg("split", args[0], 0)
	if err != nil {
		return nil, err
	}

	separator, err := validation.StringArg("separator", args[1], 1)
	if err != nil {
		return nil, err
	}

	split := strings.Split(s.Value, separator.Value)

	items := make([]types.SketchType, len(split))
	for i, item := range split {
		items[i] = &types.SketchString{
			Value: item,
		}
	}

	return &types.SketchList{
		List: types.NewList(items),
	}, nil
}

func fields(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("fields", 1, args); err != nil {
		return nil, err
	}

	s, err := validation.StringArg("fields", args[0], 0)
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(s.Value)
	items := make([]types.SketchType, len(fields))
	for i, field := range fields {
		items[i] = &types.SketchString{Value: field}
	}

	return &types.SketchList{
		List: types.NewList(items),
	}, nil
}
