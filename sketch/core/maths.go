package core

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
)

func add(args ...types.SketchType) (types.SketchType, error) {
	switch a := args[0].(type) {
	case *types.SketchInt:
		sum := a.Value
		for _, arg := range args[1:] {
			b, ok := arg.(*types.SketchInt)
			if !ok {
				return nil, fmt.Errorf("addition between different types")
			}
			sum += b.Value
		}
		return &types.SketchInt{
			Value: sum,
		}, nil
	case *types.SketchString:
		sum := a.Value
		for _, arg := range args[1:] {
			b, ok := arg.(*types.SketchString)
			if !ok {
				return nil, fmt.Errorf("addition between different types")
			}
			sum += b.Value
		}
		return &types.SketchString{
			Value: sum,
		}, nil
	}
	return nil, fmt.Errorf("unsupported first arg to +")
}

func subtract(args ...types.SketchType) (types.SketchType, error) {
	numbers, err := validation.NIntArgs("-", 2, args)
	if err != nil {
		return nil, err
	}
	return &types.SketchInt{
		Value: numbers[0].Value - numbers[1].Value,
	}, nil
}

func multiply(args ...types.SketchType) (types.SketchType, error) {
	numbers, err := validation.NIntArgs("*", 2, args)
	if err != nil {
		return nil, err
	}
	return &types.SketchInt{
		Value: numbers[0].Value * numbers[1].Value,
	}, nil
}

func divide(args ...types.SketchType) (types.SketchType, error) {
	numbers, err := validation.NIntArgs("/", 2, args)
	if err != nil {
		return nil, err
	}
	return &types.SketchInt{
		Value: numbers[0].Value / numbers[1].Value,
	}, nil
}

func lt(args ...types.SketchType) (types.SketchType, error) {
	numbers, err := validation.NIntArgs("<", 2, args)
	if err != nil {
		return nil, err
	}
	return &types.SketchBoolean{
		Value: numbers[0].Value < numbers[1].Value,
	}, nil
}

func lte(args ...types.SketchType) (types.SketchType, error) {
	numbers, err := validation.NIntArgs("<=", 2, args)
	if err != nil {
		return nil, err
	}
	return &types.SketchBoolean{
		Value: numbers[0].Value <= numbers[1].Value,
	}, nil
}

func gt(args ...types.SketchType) (types.SketchType, error) {
	numbers, err := validation.NIntArgs(">", 2, args)
	if err != nil {
		return nil, err
	}
	return &types.SketchBoolean{
		Value: numbers[0].Value > numbers[1].Value,
	}, nil
}

func gte(args ...types.SketchType) (types.SketchType, error) {
	numbers, err := validation.NIntArgs(">=", 2, args)
	if err != nil {
		return nil, err
	}
	return &types.SketchBoolean{
		Value: numbers[0].Value >= numbers[1].Value,
	}, nil
}
