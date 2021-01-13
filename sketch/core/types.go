package core

import (
	"fmt"
	"strconv"

	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
)

const integerDocstring = `
TODO
`

func integer(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("int", 1, args); err != nil {
		return nil, err
	}
	switch arg := args[0].(type) {
	case *types.SketchInt:
		return arg, nil
	case *types.SketchString:
		i, err := strconv.Atoi(arg.Value)
		if err != nil {
			return nil, err
		}
		return &types.SketchInt{
			Value: i,
		}, nil
	default:
		return nil, fmt.Errorf("int: unable to convert type %s to an int", arg.Type())
	}
}
