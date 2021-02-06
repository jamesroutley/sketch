package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
)

func prn(args ...types.SketchType) (types.SketchType, error) {
	ss := make([]string, len(args))
	for i, arg := range args {
		ss[i] = printer.PrStr(arg)
	}
	fmt.Println(strings.Join(ss, " "))
	return &types.SketchNil{}, nil
}

func list(args ...types.SketchType) (types.SketchType, error) {
	return &types.SketchList{
		List: types.NewList(args),
	}, nil
}

func isList(args ...types.SketchType) (types.SketchType, error) {
	_, ok := args[0].(*types.SketchList)
	return &types.SketchBoolean{
		Value: ok,
	}, nil
}

func isEmpty(args ...types.SketchType) (types.SketchType, error) {
	list, ok := args[0].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("first argument to empty? isn't a list")
	}

	return &types.SketchBoolean{
		Value: list.List.Empty(),
	}, nil
}

func count(args ...types.SketchType) (types.SketchType, error) {
	if args[0].Type() == "int" {
		return &types.SketchInt{
			Value: 0,
		}, nil
	}
	list, err := validation.ListArg("count", args[0], 0)
	if err != nil {
		return nil, err
	}

	return &types.SketchInt{
		Value: list.List.Length(),
	}, nil
}

func nth(args ...types.SketchType) (types.SketchType, error) {
	list, err := validation.ListArg("nth", args[0], 0)
	if err != nil {
		return nil, err
	}
	n, err := validation.IntArg("nth", args[1], 1)
	if err != nil {
		return nil, err
	}

	items := list.List.ToSlice()

	if n.Value >= len(items) {
		return nil, fmt.Errorf(
			"nth: index out of range - %d, with length %d, %s", n.Value, len(items), items,
		)
	}

	return items[n.Value], nil
}

func equals(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("=", 2, args); err != nil {
		return nil, err
	}

	return &types.SketchBoolean{
		Value: equalsInternal(args[0], args[1]),
	}, nil
}

func equalsInternal(aa types.SketchType, bb types.SketchType) bool {
	if reflect.TypeOf(aa) != reflect.TypeOf(bb) {
		return false
	}

	switch a := aa.(type) {
	case *types.SketchList:
		b := bb.(*types.SketchList)
		aSlice := a.List.ToSlice()
		bSlice := b.List.ToSlice()
		if len(aSlice) != len(bSlice) {
			return false
		}

		for i := range aSlice {
			if !equalsInternal(aSlice[i], bSlice[i]) {
				return false
			}
		}

	case *types.SketchInt:
		b := bb.(*types.SketchInt)
		return a.Value == b.Value

	case *types.SketchBoolean:
		b := bb.(*types.SketchBoolean)
		return a.Value == b.Value

	case *types.SketchSymbol:
		b := bb.(*types.SketchSymbol)
		return a.Value == b.Value

	case *types.SketchString:
		b := bb.(*types.SketchString)
		return a.Value == b.Value

	case *types.SketchNil:
		// Nils don't have values, so they're always equal
		return true

	default:
		log.Fatalf("equals unimplemented for type %T", a)
	}

	return true
}

func readString(args ...types.SketchType) (types.SketchType, error) {
	arg, ok := args[0].(*types.SketchString)
	if !ok {
		return nil, fmt.Errorf("read-string takes a string")
	}

	return reader.Read(arg.Value)
}

func slurp(args ...types.SketchType) (types.SketchType, error) {
	filename, ok := args[0].(*types.SketchString)
	if !ok {
		return nil, fmt.Errorf("slurp takes a string")
	}

	data, err := ioutil.ReadFile(filename.Value)
	if err != nil {
		return nil, err
	}

	return &types.SketchString{
		Value: string(data),
	}, nil
}

// cons prepends arg1 onto the list at arg2
// >(cons 1 (quote (2 3)))
// (1 2 3)
func cons(args ...types.SketchType) (types.SketchType, error) {
	list, ok := args[1].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("cons takes a list as its second argument")
	}
	return &types.SketchList{
		List: list.List.Conj(args[0]),
	}, nil
}

// concat takes a number of lists and concatenates them together
// > (concat (list 1 2) (list 3 4))
// (1 2 3 4)
func concat(args ...types.SketchType) (types.SketchType, error) {
	var allItems []types.SketchType

	for _, arg := range args {
		list, ok := arg.(*types.SketchList)
		if !ok {
			return nil, fmt.Errorf("concat takes lists as arguments")
		}
		allItems = append(allItems, list.List.ToSlice()...)
	}

	return &types.SketchList{
		List: types.NewList(allItems),
	}, nil
}

func first(args ...types.SketchType) (types.SketchType, error) {
	switch arg := args[0].(type) {
	case *types.SketchNil:
		return arg, nil
	case *types.SketchList:
		return arg.List.First(), nil
	case *types.SketchString:
		runes := []rune(arg.Value)
		if len(runes) == 0 {
			return &types.SketchNil{}, nil
		}
		return &types.SketchString{
			Value: string(runes[0]),
		}, nil
	default:
		return nil, fmt.Errorf("first arg to first must be a list, got %s %s", arg.Type(), arg.String())
	}
}

func rest(args ...types.SketchType) (types.SketchType, error) {
	switch arg := args[0].(type) {
	case *types.SketchNil:
		return &types.SketchList{List: types.NewEmptyList()}, nil
	case *types.SketchList:
		return &types.SketchList{
			List: arg.List.Rest(),
		}, nil
	case *types.SketchString:
		runes := []rune(arg.Value)
		if len(runes) <= 1 {
			return &types.SketchList{List: types.NewEmptyList()}, nil
		}
		return &types.SketchString{
			Value: string(runes[1:]),
		}, nil
	default:
		return nil, fmt.Errorf("first arg to rest must be a list, got %s", arg.Type())
	}
}

func and(args ...types.SketchType) (types.SketchType, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("and takes at least one arg, got %d", len(args))
	}
	result := true
	for _, arg := range args {
		if !IsTruthy(arg) {
			result = false
			break
		}
	}
	return &types.SketchBoolean{
		Value: result,
	}, nil
}

func or(args ...types.SketchType) (types.SketchType, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("or takes at least one arg, got %d", len(args))
	}
	result := false
	for _, arg := range args {
		if IsTruthy(arg) {
			result = true
			break
		}
	}
	return &types.SketchBoolean{
		Value: result,
	}, nil
}

func stringToList(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("string-to-list", 1, args); err != nil {
		return nil, err
	}
	str, err := validation.StringArg("string-to-list", args[0], 0)
	if err != nil {
		return nil, err
	}

	var chars []types.SketchType
	for _, r := range str.Value {
		chars = append(chars, &types.SketchString{
			Value: string(r),
		})
	}

	return &types.SketchList{
		List: types.NewList(chars),
	}, nil
}

func length(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("length", 1, args); err != nil {
		return nil, err
	}

	itemLength := 0
	switch arg := args[0].(type) {
	case *types.SketchList:
		itemLength = arg.List.Length()
	case *types.SketchString:
		runes := []rune(arg.Value)
		itemLength = len(runes)
	default:
		return nil, fmt.Errorf("length called with type %s, only supports list and string", arg.Type())
	}

	return &types.SketchInt{
		Value: itemLength,
	}, nil
}

// apply
// (apply + (list 1 2 3)) is equivalent to (+ 1 2 3)
func apply(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("apply", 2, args); err != nil {
		return nil, err
	}
	function, err := validation.FunctionArg("apply", args[0], 0)
	if err != nil {
		return nil, err
	}
	list, err := validation.ListArg("apply", args[1], 1)
	if err != nil {
		return nil, err
	}

	return function.Func(list.List.ToSlice()...)
}

// func list(args ...types.SketchType) (types.SketchType, error) {
// }

// IsTruthy returns a type's truthiness. Currently: it's falsy if the type is
// `nil` or the boolean 'false'. All other values are truthy.
func IsTruthy(t types.SketchType) bool {
	switch token := t.(type) {
	case *types.SketchNil:
		return false
	case *types.SketchBoolean:
		return token.Value
	}
	return true
}
