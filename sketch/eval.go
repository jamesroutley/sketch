package sketch

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/types"
)

// Eval evaulates a piece of parsed code.
// The way code is evaluated depends on its structure.
//
// 1. Special forms - these are language-level features, which behave
// differently to normal functions. These include `def` and `let`. For
// example, certain elements in the argument list might be evaluated
// differently (or not at all)
// 2. Symbols: evaluated to their corresponding value in the environment `env`
// 3. Lists: by default, they're treated as function calls - each item is
// evaluated, and the first item (the function itself) is called with the rest
// of the items as arguments.
func Eval(ast types.MalType, env *environment.Env) (types.MalType, error) {
	// This whlie loop enables tail call optimisation (TCO), where we mutate
	// `ast` and `env` and jump back to the top, rather than recursively
	// calling `Eval`. This stops a stack frame from being pushed, and lets us
	// recurse to depths that would otherwise cause a stack overflow.
	for {
		// First - check if ast is a list. If it isn't we can evaluate it as an
		// atom and return.
		// N.B: a lot of mutation goes on in this function. We use these scoping
		// {} parens to try and minimise cross contamination between sections
		{
			list, ok := ast.(*types.MalList)
			if !ok {
				return evalAST(ast, env)
			}

			// If the list is empty - return it. Empty lists eval to themselves
			//
			// > ()
			// ()
			//
			// (these cryptic symbols indicate that what you'd see on the REPL when
			// evaluating an empty list)
			if len(list.Items) == 0 {
				return ast, nil
			}
		}

		// Ok, AST is a list. Lists can contain function calls, special forms
		// and macros. Here we handle those cases.

		// First, macros. A macro modifies Lisp source code, so we need to
		// expand them before we continue evaluating.
		{
			expandedAST, err := macroExpand(ast, env)
			if err != nil {
				return nil, err
			}

			// Check if the ast is still a list after the macro expansion. If it
			// isn't, we just return evalAST, like we did for non-lists above.
			// If it is, continue.
			switch expandedAST.(type) {
			case *types.MalList:
				ast = expandedAST
				// continue
			default:
				return evalAST(expandedAST, env)
			}
		}

		// Process tail call optimised special forms. Some special forms call
		// Eval. For example, an `if` expression evaluates one of the two
		// expressions provided, depending on whether the condition is true or
		// false. Instead of recusively calling Eval, they return a new `ast`
		// and `env`, and we loop back to the top of this function.
		{
			evaluated, newAST, newEnv, err := evalTCOSpecialForm(ast, env)
			if err != nil {
				return nil, err
			}
			if evaluated {
				// TCO
				ast = newAST
				env = newEnv
				continue
			}
		}

		// Process non tail call optimised special forms
		{
			evaluated, newAST, err := evalSpecialForm(ast, env)
			if err != nil {
				return nil, err
			}
			if evaluated {
				return newAST, err
			}
		}

		// Okay - our list has had any macros expanded, and isn't a (tail call
		// optimised) special form. Evaluate it according to Lisp rules -
		// evaluate all elements in the list, then call the first as a
		// function, with the rest as arguments.
		// Evaluating the list converts the symbol at the head of the list to
		// a function, and evaluating the rest of the arguments 'collapses'
		// them down - e.g. if one arg is `(+ 1 1)`, it'll eval down to 2.
		// Evaluating the arguments first makes this Lisp 'eager' (i.e. not
		// lazy). Not evaluating them up front would give us a lazy language,
		// which has interesting and different properties.
		{
			newAST, err := evalAST(ast, env)
			if err != nil {
				return nil, err
			}

			list, ok := newAST.(*types.MalList)
			if !ok {
				return nil, fmt.Errorf("list did not evaluate to a list")
			}

			function, ok := list.Items[0].(*types.MalFunction)
			if !ok {
				return nil, fmt.Errorf("first item in list isn't a function")
			}

			if !function.TailCallOptimised {
				return function.Func(list.Items[1:]...)
			}

			// Function is tail call optimised.
			// Construct the correct environment it should be run in
			childEnv := environment.NewChildEnv(
				function.Env.(*environment.Env), function.Params, list.Items[1:],
			)
			// TCO
			ast = function.AST
			env = childEnv
			continue
		}
	}
}

// evalAST implements the evaluation rules for normal expressions. Any special
// cases are handed above us, in the Eval function. This function is an
// implementation detail of Eval, and shoulnd't be called apart from by it.
func evalAST(ast types.MalType, env *environment.Env) (types.MalType, error) {
	switch tok := ast.(type) {
	case *types.MalSymbol:
		value, err := env.Get(tok.Value)
		if err != nil {
			return nil, err
		}
		return value, nil
	case *types.MalList:
		items := make([]types.MalType, len(tok.Items))
		for i, item := range tok.Items {
			evaluated, err := Eval(item, env)
			if err != nil {
				return nil, err
			}
			items[i] = evaluated
		}
		return &types.MalList{
			Items: items,
		}, nil
	}
	return ast, nil
}

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

func quasiquote(ast types.MalType) (types.MalType, error) {
	list, ok := ast.(*types.MalList)
	if !ok {
		// `ast` isn't a list, which means it can't be an unquoted form. Return
		// its quoted form. Here, we quote it regardless of its type.
		// Quoting forms such as ints and strings is redundant - quoting
		// prevents the evaluator from evaluating a form, but those forms
		// evaluate to themselves, so there's no difference caused by quoting
		// them. However, there's also no harm harm in doing so.
		// This return statements returns the AST version of (quote <ast>)
		return &types.MalList{
			Items: []types.MalType{
				&types.MalSymbol{Value: "quote"},
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
	if symbol, ok := items[0].(*types.MalSymbol); ok && symbol.Value == "unquote" {
		return list.Items[1], nil
	}

	// Okay - ast is a list, than hasn't been unquoted
	quasiquoted := &types.MalList{}

	for i := len(items) - 1; i >= 0; i-- {
		element := items[i]

		// TODO: implement `splice-unquote` functionality
		if args, ok := isSpliceUnquoteForm(element); ok {
			quasiquoted = &types.MalList{
				Items: []types.MalType{
					&types.MalSymbol{Value: "concat"},
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

		quasiquoted = &types.MalList{
			Items: []types.MalType{
				&types.MalSymbol{Value: "cons"},
				quasiqutoedElement,
				quasiquoted,
			},
		}

	}
	return quasiquoted, nil
}

func isSpliceUnquoteForm(ast types.MalType) (spliceUnquoteArgs []types.MalType, ok bool) {
	list, ok := ast.(*types.MalList)
	if !ok {
		return nil, false
	}
	items := list.Items
	if len(items) == 0 {
		return nil, false
	}
	symbol, ok := items[0].(*types.MalSymbol)
	if !ok {
		return nil, false
	}
	if symbol.Value == "splice-unquote" {
		return items[1:], true
	}
	return nil, false
}

func isMacroCall(ast types.MalType, env *environment.Env) bool {
	list, ok := ast.(*types.MalList)
	if !ok {
		return false
	}
	items := list.Items
	if len(items) == 0 {
		return false
	}
	symbol, ok := items[0].(*types.MalSymbol)
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
	function, ok := value.(*types.MalFunction)
	if !ok {
		return false
	}
	return function.IsMacro
}

func macroExpand(ast types.MalType, env *environment.Env) (types.MalType, error) {
	for isMacroCall(ast, env) {
		// TODO: isMacroCall could return the macro function, which would save
		// the casting below
		//
		// Don't check the ok value here because we've validated that ast is a
		// list in isMacroCall. If it's not, something very strange has
		// happened
		list := ast.(*types.MalList)
		// Again, we've already checked this - skip ok checking
		macroName := list.Items[0].(*types.MalSymbol)

		macroNameValue, err := env.Get(macroName.Value)
		if err != nil {
			// Shouldn't happen - we've already validated this above
			return nil, err
		}
		macroFunc := macroNameValue.(*types.MalFunction)

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
