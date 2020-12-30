package sketch

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/types"
)

func isTCOSpecialForm(ast types.MalType) (operator *types.MalSymbol, args []types.MalType, ok bool) {
	tok, ok := ast.(*types.MalList)
	if !ok {
		return nil, nil, false
	}
	items := tok.Items
	if len(items) == 0 {
		return nil, nil, false
	}

	operator, ok = items[0].(*types.MalSymbol)
	if !ok {
		return nil, nil, false
	}

	switch operator.Value {
	case "let*", "if", "do", "quasiquote":
		return operator, items[1:], true
	}

	return nil, nil, false
}

func isSpecialForm(ast types.MalType) (operator *types.MalSymbol, args []types.MalType, ok bool) {
	tok, ok := ast.(*types.MalList)
	if !ok {
		return nil, nil, false
	}
	items := tok.Items
	if len(items) == 0 {
		return nil, nil, false
	}

	operator, ok = items[0].(*types.MalSymbol)
	if !ok {
		return nil, nil, false
	}

	switch operator.Value {
	case "fn*", "def!", "quote", "quasiquoteexpand", "defmacro!", "macroexpand":
		return operator, items[1:], true
	}

	return nil, nil, false
}

// Some special forms end in an evaluation. We could implement this by
// recusively calling `Eval` (it's recusive because evalTCOSpecialForm is
// called by Eval), but that can lead to stack overflow issues. Instead, we
// tail call optimise by returning a new AST to evaluate, and a new environment
// to invaluate it in. Eval loops back to the beginning of the function and
// re-runs itself using these new params.
func evalTCOSpecialForm(
	operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, newEnv *environment.Env, err error) {
	switch operator.Value {

	// Creates a new environment with certain variables set, then evaluates a
	// statement in that environment.
	// e.g:
	//
	// > (let* (a 1 b (+ a 1)) b)
	// 2 ; a == b, b == a+1 == 2
	case "let*":
		if len(args) != 2 {
			return nil, nil, fmt.Errorf("let* takes 2 args")
		}
		bindingList, ok := args[0].(*types.MalList)
		if !ok {
			return nil, nil, fmt.Errorf("let*: first arg isn't a list")
		}
		if len(bindingList.Items)%2 != 0 {
			return nil, nil, fmt.Errorf("let*: first arg doesn't have an even number of items")
		}

		childEnv := env.ChildEnv()
		for i := 0; i < len(bindingList.Items); i += 2 {
			key, ok := bindingList.Items[i].(*types.MalSymbol)
			if !ok {
				return nil, nil, fmt.Errorf("let*: binding list: arg %d isn't a symbol", i)
			}
			value, err := Eval(bindingList.Items[i+1], childEnv)
			if err != nil {
				return nil, nil, err
			}
			childEnv.Set(key.Value, value)
		}

		// Finally, return the last arg as the new AST to be evaluated, and the
		// newly constructed env as the environment
		return args[1], childEnv, nil

	// Evaluates the elements in the arg list and returns the final result.
	// For TCO, we eval all but the last argument here, then return the last
	// argument to be evaluated in the main Eval loop.
	case "do":
		for _, arg := range args[:len(args)-1] {
			var err error
			_, err = Eval(arg, env)
			if err != nil {
				return nil, nil, err
			}
		}
		return args[len(args)-1], env, nil

	// Evaluate first param. If not `nil` or `false`, return the second param
	// to be evaluated. If it is, return the third param to be evaluated, or
	// `nil` if none is supplied. If none is supplied, the `nil` value is
	// evalulated, but just evaluates to `nil`.
	case "if":
		if numArgs := len(args); numArgs != 2 && numArgs != 3 {
			return nil, nil, fmt.Errorf("if statements must have two or three arguments, got %d", numArgs)
		}
		condition, err := Eval(args[0], env)
		if err != nil {
			return nil, nil, err
		}
		if IsTruthy(condition) {
			return args[1], env, nil
		}

		if len(args) == 3 {
			return args[2], env, nil
		}

		return &types.MalNil{}, env, nil

	case "quasiquote":
		ast, err := quasiquote(args[0])
		if err != nil {
			return nil, nil, err
		}
		return ast, env, nil

	default:
		return nil, nil, fmt.Errorf("unexpected tail call optimised special form: %s", operator.Value)
	}
}

func evalSpecialForm(
	operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (types.MalType, error) {
	switch operator.Value {

	// Assigns a value to a symbol in the current environment
	// e.g:
	//
	// > (def a 10)
	// 10
	// > a
	// 10
	case "def!":
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

	// Create a new function.
	//
	// e.g:
	// > (def! add1 (fn* (a) (+ a 1)))
	// #<function>
	// > (add1 2)
	// 3
	case "fn*":
		if len(args) != 2 {
			return nil, fmt.Errorf("fn* statements must have two arguments, got %d", len(args))
		}

		// arguments is the first argument supplied to the fn* function (e.g.
		// `(a)` in the example above)
		arguments, ok := args[0].(*types.MalList)
		if !ok {
			return nil, fmt.Errorf("fn* statements must have a list as the first arg")
		}
		// Cast it from a list of MalType to a list of MalSymbol
		binds := make([]*types.MalSymbol, len(arguments.Items))
		for i, a := range arguments.Items {
			bind, ok := a.(*types.MalSymbol)
			if !ok {
				// TODO: improve this - say which argument isn't a symbol
				return nil, fmt.Errorf("fn* statements must have a list of symbols as the first arg")
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

	case "quote":
		return args[0], nil

	case "quasiquoteexpand":
		return quasiquote(args[0])

	// Creates a new macro
	case "defmacro!":
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

	// Macroexpand expands a macro and returns the expanded form. Useful for
	// debugging macros
	case "macroexpand":
		return macroExpand(args[0], env)

	// XXX: if you add a case here, you also need to add it to `isSpecialForm`

	default:
		return nil, fmt.Errorf("unexpected special form: %s", operator.Value)
	}
}
