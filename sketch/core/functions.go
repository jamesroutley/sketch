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

func prn(args ...types.MalType) (types.MalType, error) {
	fmt.Println(printer.PrStr(args[0]))
	return &types.MalNil{}, nil
}

func list(args ...types.MalType) (types.MalType, error) {
	return &types.MalList{
		Items: args,
	}, nil
}

func isList(args ...types.MalType) (types.MalType, error) {
	_, ok := args[0].(*types.MalList)
	return &types.MalBoolean{
		Value: ok,
	}, nil
}

func isEmpty(args ...types.MalType) (types.MalType, error) {
	list, ok := args[0].(*types.MalList)
	if !ok {
		return nil, fmt.Errorf("first argument to empty? isn't a list")
	}

	return &types.MalBoolean{
		Value: len(list.Items) == 0,
	}, nil
}

func count(args ...types.MalType) (types.MalType, error) {
	if _, ok := args[0].(*types.MalNil); ok {
		return &types.MalInt{
			Value: 0,
		}, nil
	}
	list, ok := args[0].(*types.MalList)
	if !ok {
		return nil, fmt.Errorf("first argument to count isn't a list")
	}

	return &types.MalInt{
		Value: len(list.Items),
	}, nil
}

func equals(args ...types.MalType) (types.MalType, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("equals requires 2 args - got %d", len(args))
	}

	return &types.MalBoolean{
		Value: equalsInternal(args[0], args[1]),
	}, nil
}

func equalsInternal(aa types.MalType, bb types.MalType) bool {
	if reflect.TypeOf(aa) != reflect.TypeOf(bb) {
		return false
	}

	switch a := aa.(type) {
	case *types.MalList:
		b := bb.(*types.MalList)
		if len(a.Items) != len(b.Items) {
			return false
		}

		for i := range a.Items {
			if !equalsInternal(a.Items[i], b.Items[i]) {
				return false
			}
		}

	case *types.MalInt:
		b := bb.(*types.MalInt)
		return a.Value == b.Value

	case *types.MalBoolean:
		b := bb.(*types.MalBoolean)
		return a.Value == b.Value

	case *types.MalSymbol:
		b := bb.(*types.MalSymbol)
		return a.Value == b.Value

	case *types.MalString:
		b := bb.(*types.MalString)
		return a.Value == b.Value

	case *types.MalNil:
		// Nils don't have values, so they're always equal
		return true

	default:
		log.Fatalf("equals unimplemented for type %T", a)
	}

	return true
}

func lt(args ...types.MalType) (types.MalType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToMalInt(args)
	if err != nil {
		return nil, err
	}

	return &types.MalBoolean{
		Value: numbers[0].Value < numbers[1].Value,
	}, nil
}

func lte(args ...types.MalType) (types.MalType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToMalInt(args)
	if err != nil {
		return nil, err
	}

	return &types.MalBoolean{
		Value: numbers[0].Value <= numbers[1].Value,
	}, nil
}

func gt(args ...types.MalType) (types.MalType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToMalInt(args)
	if err != nil {
		return nil, err
	}

	return &types.MalBoolean{
		Value: numbers[0].Value > numbers[1].Value,
	}, nil
}

func gte(args ...types.MalType) (types.MalType, error) {
	if err := ValidateNArgs(2, args); err != nil {
		return nil, err
	}
	numbers, err := ArgsToMalInt(args)
	if err != nil {
		return nil, err
	}

	return &types.MalBoolean{
		Value: numbers[0].Value >= numbers[1].Value,
	}, nil
}

func readString(args ...types.MalType) (types.MalType, error) {
	arg, ok := args[0].(*types.MalString)
	if !ok {
		return nil, fmt.Errorf("read-string takes a string")
	}

	return reader.ReadStr(arg.Value)
}

func slurp(args ...types.MalType) (types.MalType, error) {
	filename, ok := args[0].(*types.MalString)
	if !ok {
		return nil, fmt.Errorf("read-string takes a string")
	}

	data, err := ioutil.ReadFile(filename.Value)
	if err != nil {
		return nil, err
	}

	return &types.MalString{
		Value: string(data),
	}, nil
}

// cons prepends arg1 onto the list at arg2
// >(cons 1 (quote (2 3)))
// (1 2 3)
func cons(args ...types.MalType) (types.MalType, error) {
	list, ok := args[1].(*types.MalList)
	if !ok {
		return nil, fmt.Errorf("cons takes a list as its second argument")
	}
	items := append([]types.MalType{args[0]}, list.Items...)
	return &types.MalList{
		Items: items,
	}, nil
}

// concat takes a number of lists and concatenates them together
// > (concat (list 1 2) (list 3 4))
// (1 2 3 4)
func concat(args ...types.MalType) (types.MalType, error) {
	var allItems []types.MalType

	for _, arg := range args {
		list, ok := arg.(*types.MalList)
		if !ok {
			return nil, fmt.Errorf("concat takes lists as arguments")
		}
		allItems = append(allItems, list.Items...)
	}

	return &types.MalList{
		Items: allItems,
	}, nil
}

// sketchMap implements map - i.e. run func for all items in a list
func sketchMap(args ...types.MalType) (types.MalType, error) {
	function, ok := args[0].(*types.MalFunction)
	if !ok {
		return nil, fmt.Errorf("first arg to map must be a function")
	}
	list, ok := args[1].(*types.MalList)
	if !ok {
		return nil, fmt.Errorf("second arg to map must be a list")
	}

	// Short circuit
	if len(list.Items) == 0 {
		return list, nil
	}

	g := new(errgroup.Group)
	mappedItems := make([]types.MalType, len(list.Items))
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

	return &types.MalList{
		Items: mappedItems,
	}, nil
}

func filter(args ...types.MalType) (types.MalType, error) {
	function, ok := args[0].(*types.MalFunction)
	if !ok {
		return nil, fmt.Errorf("first arg to map must be a function")
	}
	list, ok := args[1].(*types.MalList)
	if !ok {
		return nil, fmt.Errorf("second arg to map must be a list")
	}

	// Short circuit
	if len(list.Items) == 0 {
		return list, nil
	}

	g := new(errgroup.Group)
	filteredItems := make([]types.MalType, 0, len(list.Items))
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

	return &types.MalList{
		Items: filteredItems,
	}, nil
}

func first(args ...types.MalType) (types.MalType, error) {
	switch arg := args[0].(type) {
	case *types.MalNil:
		return arg, nil
	case *types.MalList:
		if len(arg.Items) == 0 {
			return &types.MalNil{}, nil
		}
		return arg.Items[0], nil
	default:
		return nil, fmt.Errorf("first arg to first must be a list")
	}
}

func rest(args ...types.MalType) (types.MalType, error) {
	switch arg := args[0].(type) {
	case *types.MalNil:
		return &types.MalList{}, nil
	case *types.MalList:
		if len(arg.Items) == 0 {
			return &types.MalList{}, nil
		}
		return &types.MalList{Items: arg.Items[1:]}, nil
	default:
		return nil, fmt.Errorf("first arg to first must be a list")
	}
}

func and(args ...types.MalType) (types.MalType, error) {
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
	return &types.MalBoolean{
		Value: result,
	}, nil
}

func or(args ...types.MalType) (types.MalType, error) {
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
	return &types.MalBoolean{
		Value: result,
	}, nil
}

// func list(args ...types.MalType) (types.MalType, error) {
// }

// isTruthy returns a type's truthiness. Currently: it's falsy if the type is
// `nil` or the boolean 'false'. All other values are truthy.
func isTruthy(t types.MalType) bool {
	switch token := t.(type) {
	case *types.MalNil:
		return false
	case *types.MalBoolean:
		return token.Value
	}
	return true
}
