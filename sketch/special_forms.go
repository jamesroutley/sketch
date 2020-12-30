package sketch

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/types"
)

type specialFormEvaluator func(
	operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, err error)

func evalSpecialForm(
	ast types.MalType, env *environment.Env,
) (evaluated bool, newAST types.MalType, err error) {
	tok, ok := ast.(*types.MalList)
	if !ok {
		return false, nil, nil
	}
	items := tok.Items
	if len(items) == 0 {
		return false, nil, nil
	}

	operator, ok := items[0].(*types.MalSymbol)
	if !ok {
		return false, nil, nil
	}

	args := items[1:]
	var evaluator specialFormEvaluator

	switch operator.Value {
	case "fn":
		evaluator = evalFn
	case "def!":
		evaluator = evalDef
	case "quote":
		evaluator = evalQuote
	case "quasiquoteexpand":
		evaluator = evalQuasiquoteExpand
	case "defmacro!":
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
// > (def! add1 (fn (a) (+ a 1)))
// #<function>
// > (add1 2)
// 3
func evalFn(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("fn statements must have two arguments, got %d", len(args))
	}

	// arguments is the first argument supplied to the fn function (e.g.
	// `(a)` in the example above)
	arguments, ok := args[0].(*types.MalList)
	if !ok {
		return nil, fmt.Errorf("fn statements must have a list as the first arg")
	}
	// Cast it from a list of MalType to a list of MalSymbol
	binds := make([]*types.MalSymbol, len(arguments.Items))
	for i, a := range arguments.Items {
		bind, ok := a.(*types.MalSymbol)
		if !ok {
			// TODO: improve this - say which argument isn't a symbol
			return nil, fmt.Errorf("fn statements must have a list of symbols as the first arg")
		}
		binds[i] = bind
	}

	// TODO: recomment this
	return &types.MalFunction{
		TailCallOptimised: true,
		AST:               args[1],
		Params:            binds,
		Env:               env,
		// This Go function is what's run when the Lisp function is
		// run. When the Lisp function is run, we create a new environment,
		// which binds the Lisp function's arguments to the parameters
		// defined when the function was defined.
		Func: func(exprs ...types.MalType) (types.MalType, error) {
			childEnv := environment.NewChildEnv(
				env, binds, exprs,
			)
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
func evalDef(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("def! takes 2 args")
	}
	key, ok := args[0].(*types.MalSymbol)
	if !ok {
		return nil, fmt.Errorf("def!: first arg isn't a symbol")
	}
	value, err := Eval(args[1], env)
	if err != nil {
		return nil, err
	}
	env.Set(key.Value, value)
	return value, nil
}

func evalQuote(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, err error) {
	return args[0], nil
}

// evalQuasiquoteExpand evaluates the `quasiquoteexpand` macro.
// This macro is used to test the internal implementation of quasiquote
func evalQuasiquoteExpand(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, err error) {
	return quasiquote(args[0])
}

// Creates a new macro
func evalDefmacro(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("defmacro! takes 2 args")
	}
	key, ok := args[0].(*types.MalSymbol)
	if !ok {
		return nil, fmt.Errorf("defmacro!: first arg isn't a symbol")
	}
	value, err := Eval(args[1], env)
	if err != nil {
		return nil, err
	}
	function, ok := value.(*types.MalFunction)
	if !ok {
		return nil, fmt.Errorf("defmacro!: second arg isn't a function definition")
	}
	function.IsMacro = true
	env.Set(key.Value, function)
	return function, nil
}

func evalMacroexpand(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, err error) {
	return macroExpand(args[0], env)
}
