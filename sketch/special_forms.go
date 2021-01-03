package sketch

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/types"
)

type specialFormEvaluator func(
	operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error)

func evalSpecialForm(
	ast types.SketchType, env *environment.Env,
) (evaluated bool, newAST types.SketchType, err error) {
	tok, ok := ast.(*types.SketchList)
	if !ok {
		return false, nil, nil
	}
	items := tok.Items
	if len(items) == 0 {
		return false, nil, nil
	}

	operator, ok := items[0].(*types.SketchSymbol)
	if !ok {
		return false, nil, nil
	}

	args := items[1:]
	var evaluator specialFormEvaluator

	switch operator.Value {
	case "fn":
		evaluator = evalFn
	case "def":
		evaluator = evalDef
	case "quote":
		evaluator = evalQuote
	case "quasiquoteexpand":
		evaluator = evalQuasiquoteExpand
	case "defmacro":
		evaluator = evalDefmacro
	case "macroexpand":
		evaluator = evalMacroexpand

	default:
		return false, nil, nil
	}

	newAST, err = evaluator(operator, args, env)
	return true, newAST, err
}

// Create a new function.
//
// e.g:
// > (def add1 (fn (a) (+ a 1)))
// #<function>
// > (add1 2)
// 3
func evalFn(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("fn statements must have two arguments, got %d", len(args))
	}

	// arguments is the first argument supplied to the fn function (e.g.
	// `(a)` in the example above)
	arguments, ok := args[0].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("fn statements must have a list as the first arg")
	}
	// Cast it from a list of SketchType to a list of SketchSymbol
	binds := make([]*types.SketchSymbol, len(arguments.Items))
	for i, a := range arguments.Items {
		bind, ok := a.(*types.SketchSymbol)
		if !ok {
			// TODO: improve this - say which argument isn't a symbol
			return nil, fmt.Errorf("fn statements must have a list of symbols as the first arg")
		}
		binds[i] = bind
	}

	// TODO: recomment this
	return &types.SketchFunction{
		TailCallOptimised: true,
		AST:               args[1],
		Params:            binds,
		Env:               env,
		// This Go function is what's run when the Lisp function is
		// run. When the Lisp function is run, we create a new environment,
		// which binds the Lisp function's arguments to the parameters
		// defined when the function was defined.
		Func: func(exprs ...types.SketchType) (types.SketchType, error) {
			childEnv, err := environment.NewChildEnv(
				env, binds, exprs,
			)
			if err != nil {
				return nil, err
			}
			return Eval(args[1], childEnv)
		},
	}, nil
}

// Assigns a value to a symbol in the current environment
// e.g:
//
// > (def a 10)
// 10
// > a
// 10
func evalDef(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("def takes 2 args")
	}
	key, ok := args[0].(*types.SketchSymbol)
	if !ok {
		return nil, fmt.Errorf("def: first arg isn't a symbol")
	}
	value, err := Eval(args[1], env)
	if err != nil {
		return nil, err
	}
	env.Set(key.Value, value)
	return value, nil
}

func evalQuote(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	return args[0], nil
}

// evalQuasiquoteExpand evaluates the `quasiquoteexpand` macro.
// This macro is used to test the internal implementation of quasiquote
func evalQuasiquoteExpand(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	return quasiquote(args[0])
}

// Creates a new macro
func evalDefmacro(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("defmacro takes 2 args")
	}
	key, ok := args[0].(*types.SketchSymbol)
	if !ok {
		return nil, fmt.Errorf("defmacro: first arg isn't a symbol")
	}
	value, err := Eval(args[1], env)
	if err != nil {
		return nil, err
	}
	function, ok := value.(*types.SketchFunction)
	if !ok {
		return nil, fmt.Errorf("defmacro: second arg isn't a function definition")
	}
	function.IsMacro = true
	env.Set(key.Value, function)
	return function, nil
}

func evalMacroexpand(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	return macroExpand(args[0], env)
}
