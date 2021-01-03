package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/jamesroutley/sketch/sketch/types"
	"golang.org/x/sync/errgroup"
)

func prn(args ...types.SketchType) (types.SketchType, error) {
	fmt.Println(printer.PrStr(args[0]))
	return &types.SketchNil{}, nil
}

func list(args ...types.SketchType) (types.SketchType, error) {
	return &types.SketchList{
		Items: args,
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
		Value: len(list.Items) == 0,
	}, nil
}

func count(args ...types.SketchType) (types.SketchType, error) {
	if _, ok := args[0].(*types.SketchNil); ok {
		return &types.SketchInt{
			Value: 0,
		}, nil
	}
	list, ok := args[0].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("first argument to count isn't a list")
	}

	return &types.SketchInt{
		Value: len(list.Items),
	}, nil
}

func nth(args ...types.SketchType) (types.SketchType, error) {
	list, ok := args[0].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("first argument to nth isn't a list")
	}
	index, ok := args[1].(*types.SketchInt)
	if !ok {
		return nil, fmt.Errorf("second argument to nth isn't an integer")
	}

	return list.Items[index.Value], nil
}

func equals(args ...types.SketchType) (types.SketchType, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("equals requires 2 args - got %d", len(args))
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
		if len(a.Items) != len(b.Items) {
			return false
		}

		for i := range a.Items {
			if !equalsInternal(a.Items[i], b.Items[i]) {
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

func lt(args ...types.SketchType) (types.SketchType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToSketchInt(args)
	if err != nil {
		return nil, err
	}

	return &types.SketchBoolean{
		Value: numbers[0].Value < numbers[1].Value,
	}, nil
}

func lte(args ...types.SketchType) (types.SketchType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToSketchInt(args)
	if err != nil {
		return nil, err
	}

	return &types.SketchBoolean{
		Value: numbers[0].Value <= numbers[1].Value,
	}, nil
}

func gt(args ...types.SketchType) (types.SketchType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToSketchInt(args)
	if err != nil {
		return nil, err
	}

	return &types.SketchBoolean{
		Value: numbers[0].Value > numbers[1].Value,
	}, nil
}

func gte(args ...types.SketchType) (types.SketchType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToSketchInt(args)
	if err != nil {
		return nil, err
	}

	return &types.SketchBoolean{
		Value: numbers[0].Value >= numbers[1].Value,
	}, nil
}

func readString(args ...types.SketchType) (types.SketchType, error) {
	arg, ok := args[0].(*types.SketchString)
	if !ok {
		return nil, fmt.Errorf("read-string takes a string")
	}

	return reader.ReadStr(arg.Value)
}

func slurp(args ...types.SketchType) (types.SketchType, error) {
	filename, ok := args[0].(*types.SketchString)
	if !ok {
		return nil, fmt.Errorf("read-string takes a string")
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
	items := append([]types.SketchType{args[0]}, list.Items...)
	return &types.SketchList{
		Items: items,
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
		allItems = append(allItems, list.Items...)
	}

	return &types.SketchList{
		Items: allItems,
	}, nil
}

// sketchMap implements map - i.e. run func for all items in a list
func sketchMap(args ...types.SketchType) (types.SketchType, error) {
	function, ok := args[0].(*types.SketchFunction)
	if !ok {
		return nil, fmt.Errorf("first arg to map must be a function")
	}
	list, ok := args[1].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("second arg to map must be a list")
	}

	// Short circuit
	if len(list.Items) == 0 {
		return list, nil
	}

	g := new(errgroup.Group)
	mappedItems := make([]types.SketchType, len(list.Items))
	for i, item := range list.Items {
		i := i
		item := item
		g.Go(func() error {
			mappedItem, err := function.Func(item)
			if err != nil {
				return err
			}
			mappedItems[i] = mappedItem
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &types.SketchList{
		Items: mappedItems,
	}, nil
}

func filter(args ...types.SketchType) (types.SketchType, error) {
	function, ok := args[0].(*types.SketchFunction)
	if !ok {
		return nil, fmt.Errorf("first arg to map must be a function")
	}
	list, ok := args[1].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("second arg to map must be a list")
	}

	// Short circuit
	if len(list.Items) == 0 {
		return list, nil
	}

	g := new(errgroup.Group)
	filteredItems := make([]types.SketchType, 0, len(list.Items))
	for _, item := range list.Items {
		item := item
		g.Go(func() error {
			passed, err := function.Func(item)
			if err != nil {
				return err
			}
			if isTruthy(passed) {
				filteredItems = append(filteredItems, item)
			}

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &types.SketchList{
		Items: filteredItems,
	}, nil
}

func first(args ...types.SketchType) (types.SketchType, error) {
	switch arg := args[0].(type) {
	case *types.SketchNil:
		return arg, nil
	case *types.SketchList:
		if len(arg.Items) == 0 {
			return &types.SketchNil{}, nil
		}
		return arg.Items[0], nil
	default:
		return nil, fmt.Errorf("first arg to first must be a list, got %s %s", arg.Type(), arg.String())
	}
}

func rest(args ...types.SketchType) (types.SketchType, error) {
	switch arg := args[0].(type) {
	case *types.SketchNil:
		return &types.SketchList{}, nil
	case *types.SketchList:
		if len(arg.Items) == 0 {
			return &types.SketchList{}, nil
		}
		return &types.SketchList{Items: arg.Items[1:]}, nil
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
		if !isTruthy(arg) {
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
		if isTruthy(arg) {
			result = true
			break
		}
	}
	return &types.SketchBoolean{
		Value: result,
	}, nil
}

// func list(args ...types.SketchType) (types.SketchType, error) {
// }

// isTruthy returns a type's truthiness. Currently: it's falsy if the type is
// `nil` or the boolean 'false'. All other values are truthy.
func isTruthy(t types.SketchType) bool {
	switch token := t.(type) {
	case *types.SketchNil:
		return false
	case *types.SketchBoolean:
		return token.Value
	}
	return true
}
