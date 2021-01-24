// Package evaluator implements Sketch's evaluator.
package evaluator

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/core"
	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/jamesroutley/sketch/sketch/types"
)

// RootEnvironment initialises a root environment loaded with all the built in
// functions and variables defined in the core package. This environment is
// used as the context in which Sketch code is evaluated.
func RootEnvironment() (*environment.Env, error) {
	env := environment.NewEnv()
	for key, value := range core.EnvironmentItems {
		env.Set(key, value)
	}

	if core.SketchCode == "" {
		return env, nil
	}
	ast, err := reader.Read(fmt.Sprintf("(do %s)", core.SketchCode))
	if err != nil {
		return nil, err
	}
	_, err = Eval(ast, env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func Evaluate(ast types.SketchType) (types.SketchType, error) {
	env, err := RootEnvironment()
	if err != nil {
		return nil, err
	}
	return Eval(ast, env)
}

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
func Eval(ast types.SketchType, env *environment.Env) (types.SketchType, error) {
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
			list, ok := ast.(*types.SketchList)
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
			case *types.SketchList:
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

			list, ok := newAST.(*types.SketchList)
			if !ok {
				return nil, fmt.Errorf("list did not evaluate to a list")
			}

			function, ok := list.Items[0].(*types.SketchFunction)
			if !ok {
				return nil, fmt.Errorf("Error evaluating list %s. I expected the first item in the list to be a function, but it's a %s.", list, list.Items[0].Type())
			}

			if !function.TailCallOptimised {
				return function.Func(list.Items[1:]...)
			}

			// Function is tail call optimised.
			// Construct the correct environment it should be run in
			childEnv, err := environment.NewFunctionEnv(
				function.Env.(*environment.Env), function.Params, list.Items[1:],
			)
			if err != nil {
				return nil, err
			}
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
func evalAST(ast types.SketchType, env *environment.Env) (types.SketchType, error) {
	switch tok := ast.(type) {
	case *types.SketchSymbol:
		value, err := env.Get(tok.Value)
		if err != nil {
			return nil, err
		}
		return value, nil
	case *types.SketchList:
		items := make([]types.SketchType, len(tok.Items))
		for i, item := range tok.Items {
			evaluated, err := Eval(item, env)
			if err != nil {
				return nil, err
			}
			items[i] = evaluated
		}
		return &types.SketchList{
			Items: items,
		}, nil
	case *types.SketchHashMap:
		keys := tok.Keys()
		items := make([]types.SketchType, 0, len(keys)*2)
		for _, key := range keys {
			value, err := tok.Get(key)
			if err != nil {
				// This is bad
				panic(err)
			}
			evaluated, err := Eval(value, env)
			if err != nil {
				return nil, err
			}
			items = append(items, key, evaluated)
		}
		return types.NewSketchHashMap(items)
	}
	return ast, nil
}
