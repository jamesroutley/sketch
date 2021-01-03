package core

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/types"
)

func ValidateNArgs(fnName string, n int, args []types.SketchType) error {
	if numArgs := len(args); numArgs != n {
		return fmt.Errorf("the function %s expects %d arguments, but got %d", fnName, n, numArgs)
	}
	return nil
}

func ValidateNIntArgs(fnName string, n int, args []types.SketchType) ([]*types.SketchInt, error) {
	if err := ValidateNArgs(fnName, n, args); err != nil {
		return nil, err
	}
	numbers := make([]*types.SketchInt, len(args))
	for i, arg := range args {
		if err := ValidateArgType(fnName, arg, "int", i); err != nil {
			return nil, err
		}
		numbers[i] = arg.(*types.SketchInt)
	}
	return numbers, nil
}

func ValidateListArg(
	fnName string, arg types.SketchType, position int,
) (*types.SketchList, error) {
	if err := ValidateArgType(fnName, arg, "list", position); err != nil {
		return nil, err
	}
	return arg.(*types.SketchList), nil
}

func ValidateIntArg(
	fnName string, arg types.SketchType, position int,
) (*types.SketchInt, error) {
	if err := ValidateArgType(fnName, arg, "int", position); err != nil {
		return nil, err
	}
	return arg.(*types.SketchInt), nil
}

func ValidateFunctionArg(
	fnName string, arg types.SketchType, position int,
) (*types.SketchFunction, error) {
	if err := ValidateArgType(fnName, arg, "function", position); err != nil {
		return nil, err
	}
	return arg.(*types.SketchFunction), nil
}

func ValidateArgType(
	fnName string, arg types.SketchType, expectedType string, position int,
) error {
	if arg.Type() != expectedType {
		oneIndexedPosition := position + 1
		return fmt.Errorf(
			"the function %s expects the %s argument `%s` to be type %s, got type %s",
			fnName, toOrdinal(oneIndexedPosition), arg, expectedType, arg.Type())
	}
	return nil
}

// toOrdinal takes an integer and returns its ordinal form - e.g. 1 -> 1st
func toOrdinal(n int) string {
	lastTwoDigits := n % 100
	switch lastTwoDigits {
	case 11, 12, 13:
		return fmt.Sprintf("%dth", n)
	}
	lastDigit := n % 10
	switch lastDigit {
	case 1:
		return fmt.Sprintf("%dst", n)
	case 2:
		return fmt.Sprintf("%dnd", n)
	case 3:
		return fmt.Sprintf("%drd", n)
	}
	return fmt.Sprintf("%dth", n)
}

// 1st 2nd 3rd 4th 5th 5th 7th
