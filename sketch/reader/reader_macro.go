package reader

import (
	"strings"

	"github.com/jamesroutley/sketch/sketch/types"
)

func stripComments2(ast types.SketchType) types.SketchType {
	list, ok := ast.(*types.SketchList)
	if !ok {
		return ast
	}
	var newItems []types.SketchType
	for _, item := range list.Items {
		switch item := item.(type) {
		case *types.SketchComment:
			// skip
		case *types.SketchList:
			newItem := stripComments2(item)
			newItems = append(newItems, newItem)
		default:
			newItems = append(newItems, item)
		}
	}

	list.Items = newItems
	return list
}

func expandModuleLookup(ast types.SketchType) types.SketchType {
	switch ast := ast.(type) {
	case *types.SketchSymbol:
		symbol := ast.Value
		// Expand module lookup symbols into the `module-lookup` function.
		// E.g: strings.join -> (module-lookup strings join)
		if strings.Contains(symbol, ".") {
			parts := strings.SplitN(symbol, ".", 2)
			return &types.SketchList{
				Items: []types.SketchType{
					&types.SketchSymbol{Value: "module-lookup"},
					&types.SketchSymbol{Value: parts[0]},
					&types.SketchSymbol{Value: parts[1]},
				},
			}
		}
	case *types.SketchList:
		for i, item := range ast.Items {
			ast.Items[i] = expandModuleLookup(item)
		}
		return ast
	}
	return ast
}
