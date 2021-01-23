package core

import (
	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
	"golang.org/x/sync/errgroup"
)

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

func foldLeft(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("fold-left", 3, args); err != nil {
		return nil, err
	}

	function, err := validation.FunctionArg("fold-left", args[0], 0)
	if err != nil {
		return nil, err
	}

	list, err := validation.ListArg("fold-left", args[2], 2)
	if err != nil {
		return nil, err
	}

	collector := args[1]
	for _, item := range list.Items {
		result, err := function.Func(collector, item)
		if err != nil {
			return nil, err
		}
		collector = result
	}

	return collector, nil
}

func flatten(args ...types.SketchType) (types.SketchType, error) {
	list, err := validation.ListArg("flatten", args[0], 0)
	if err != nil {
		return nil, err
	}

	flattened := flattenRecur(list.Items)

	return &types.SketchList{
		Items: flattened,
	}, nil
}

func flattenRecur(items []types.SketchType) []types.SketchType {
	var flattened []types.SketchType
	for _, item := range items {
		l, ok := item.(*types.SketchList)
		if ok {
			flattened = append(flattened, flattenRecur(l.Items)...)
			continue
		}

		flattened = append(flattened, item)
	}
	return flattened
}

// (range 5) => (0 1 2 3 4)
// (range 1 5) => (1 2 3 4)
func sketchRange(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgsRange("range", 1, 2, args); err != nil {
		return nil, err
	}

	var lower, upper int

	switch len(args) {
	case 1:
		lower = 0
		rawUpper, err := validation.IntArg("range", args[0], 0)
		if err != nil {
			return nil, err
		}
		upper = rawUpper.Value
	case 2:
		rawLower, err := validation.IntArg("range", args[0], 0)
		if err != nil {
			return nil, err
		}
		rawUpper, err := validation.IntArg("range", args[1], 1)
		if err != nil {
			return nil, err
		}
		lower = rawLower.Value
		upper = rawUpper.Value
	}

	var items []types.SketchType
	for i := lower; i < upper; i++ {
		items = append(items, &types.SketchInt{Value: i})
	}

	return &types.SketchList{
		Items: items,
	}, nil
}
