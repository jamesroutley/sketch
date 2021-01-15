package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/jamesroutley/sketch/sketch/printer"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
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
		Value: len(list.Items),
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

	if n.Value >= len(list.Items) {
		return nil, fmt.Errorf(
			"nth: index out of range - %d, with length %d", n.Value, len(list.Items),
		)
	}

	return list.Items[n.Value], nil
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
	function, err := validation.FunctionArg("map", args[0], 0)
	if err != nil {
		return nil, err
	}
	list, err := validation.ListArg("map", args[1], 1)
	if err != nil {
		return nil, err
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
	function, err := validation.FunctionArg("filter", args[0], 0)
	if err != nil {
		return nil, err
	}
	list, err := validation.ListArg("filter", args[1], 1)
	if err != nil {
		return nil, err
	}

	// Short circuit
	if len(list.Items) == 0 {
		return list, nil
	}

	g := new(errgroup.Group)
	filteredItems := make([]types.SketchType, len(list.Items))
	for i, item := range list.Items {
		i := i
		item := item
		g.Go(func() error {
			passed, err := function.Func(item)
			if err != nil {
				return err
			}
			// Only add item to the filtered array if it's passed the filter
			// function. We add it back to it's original position, and we'll
			// filter out nil values in the array later. We do this to preserve
			// the ordering of the list
			if IsTruthy(passed) {
				filteredItems[i] = item
			}

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	items := make([]types.SketchType, 0, len(list.Items))
	for _, item := range filteredItems {
		if item == nil {
			continue
		}
		items = append(items, item)
	}

	return &types.SketchList{
		Items: items,
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
		return &types.SketchList{}, nil
	case *types.SketchList:
		if len(arg.Items) == 0 {
			return &types.SketchList{}, nil
		}
		return &types.SketchList{Items: arg.Items[1:]}, nil
	case *types.SketchString:
		runes := []rune(arg.Value)
		if len(runes) <= 1 {
			return &types.SketchList{}, nil
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
		Items: chars,
	}, nil
}

func length(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("length", 1, args); err != nil {
		return nil, err
	}

	itemLength := 0
	switch arg := args[0].(type) {
	case *types.SketchList:
		itemLength = len(arg.Items)
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
