package evaluator

import (
	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/types"
)

func quasiquote(ast types.SketchType) (types.SketchType, error) {
	list, ok := ast.(*types.SketchList)
	if !ok {
		// `ast` isn't a list, which means it can't be an unquoted form. Return
		// its quoted form. Here, we quote it regardless of its type.
		// Quoting forms such as ints and strings is redundant - quoting
		// prevents the evaluator from evaluating a form, but those forms
		// evaluate to themselves, so there's no difference caused by quoting
		// them. However, there's also no harm harm in doing so.
		// This return statements returns the AST version of (quote <ast>)
		return &types.SketchList{
			Items: []types.SketchType{
				&types.SketchSymbol{Value: "quote"},
				ast,
			},
		}, nil
	}

	// Okay - ast is a list
	items := list.Items

	// If the list has no items, return it unmodified
	if len(items) == 0 {
		return ast, nil
	}

	// If the first item in the list is the function `unquote`, return the
	// first argument without quoting it.
	if symbol, ok := items[0].(*types.SketchSymbol); ok && symbol.Value == "unquote" {
		return list.Items[1], nil
	}

	// Okay - ast is a list, than hasn't been unquoted
	quasiquoted := &types.SketchList{}

	for i := len(items) - 1; i >= 0; i-- {
		element := items[i]

		if args, ok := isSpliceUnquoteForm(element); ok {
			quasiquoted = &types.SketchList{
				Items: []types.SketchType{
					&types.SketchSymbol{Value: "concat"},
					args[0],
					quasiquoted,
				},
			}
			continue
		}

		quasiqutoedElement, err := quasiquote(element)
		if err != nil {
			return nil, err
		}

		quasiquoted = &types.SketchList{
			Items: []types.SketchType{
				&types.SketchSymbol{Value: "cons"},
				quasiqutoedElement,
				quasiquoted,
			},
		}

	}
	return quasiquoted, nil
}

func isSpliceUnquoteForm(ast types.SketchType) (spliceUnquoteArgs []types.SketchType, ok bool) {
	list, ok := ast.(*types.SketchList)
	if !ok {
		return nil, false
	}
	items := list.Items
	if len(items) == 0 {
		return nil, false
	}
	symbol, ok := items[0].(*types.SketchSymbol)
	if !ok {
		return nil, false
	}
	if symbol.Value == "splice-unquote" {
		return items[1:], true
	}
	return nil, false
}

func isMacroCall(ast types.SketchType, env *environment.Env) bool {
	list, ok := ast.(*types.SketchList)
	if !ok {
		return false
	}
	items := list.Items
	if len(items) == 0 {
		return false
	}
	symbol, ok := items[0].(*types.SketchSymbol)
	if !ok {
		return false
	}
	value, err := env.Get(symbol.Value)
	if err != nil {
		// This looks dangerous, but is okay - the only error this function
		// returns is a not found when the symbol isn't defined in any
		// environment
		return false
	}
	function, ok := value.(*types.SketchFunction)
	if !ok {
		return false
	}
	return function.IsMacro
}

func macroExpand(ast types.SketchType, env *environment.Env) (types.SketchType, error) {
	for isMacroCall(ast, env) {
		// TODO: isMacroCall could return the macro function, which would save
		// the casting below
		//
		// Don't check the ok value here because we've validated that ast is a
		// list in isMacroCall. If it's not, something very strange has
		// happened
		list := ast.(*types.SketchList)
		// Again, we've already checked this - skip ok checking
		macroName := list.Items[0].(*types.SketchSymbol)

		macroNameValue, err := env.Get(macroName.Value)
		if err != nil {
			// Shouldn't happen - we've already validated this above
			return nil, err
		}
		macroFunc := macroNameValue.(*types.SketchFunction)

		newAst, err := macroFunc.Func(list.Items[1:]...)
		if err != nil {
			return nil, err
		}

		// Set the evaulated macro to `ast` and loop back - this lets us
		// iteratively expand nested macros
		ast = newAst
	}
	return ast, nil
}
