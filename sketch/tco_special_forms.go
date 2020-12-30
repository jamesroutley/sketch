package sketch

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/types"
)

// TODO: none of the evaluators use the `operator` param at the moment. It's
// possible we can remove it, but it might come in useful when we eventually
// print where in the source code an error comes from.
type tcoSpecialFormEvaluator func(
	operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, newEnv *environment.Env, err error)

// Some special forms end in an evaluation. We could implement this by
// recusively calling `Eval` (it's recusive because evalTCOSpecialForm is
// called by Eval), but that can lead to stack overflow issues. Instead, we
// tail call optimise by returning a new AST to evaluate, and a new environment
// to invaluate it in. Eval loops back to the beginning of the function and
// re-runs itself using these new params.
func evalTCOSpecialForm(
	ast types.MalType, env *environment.Env,
) (evaluated bool, newAST types.MalType, newEnv *environment.Env, err error) {
	tok, ok := ast.(*types.MalList)
	if !ok {
		return false, nil, nil, nil
	}
	items := tok.Items
	if len(items) == 0 {
		return false, nil, nil, nil
	}

	operator, ok := items[0].(*types.MalSymbol)
	if !ok {
		return false, nil, nil, nil
	}

	args := items[1:]
	var evaluator tcoSpecialFormEvaluator

	switch operator.Value {
	case "let":
		evaluator = evalLet
	case "if":
		evaluator = evalIf
	case "do":
		evaluator = evalDo
	case "quasiquote":
		evaluator = evalQuasiquote

	default:
		return false, nil, nil, nil
	}

	newAST, newEnv, err = evaluator(operator, args, env)
	return true, newAST, newEnv, err
}

// evalLet evaluates the `let` special form
// Creates a new environment with certain variables set, then evaluates a
// statement in that environment.
// e.g:
//
// > (let (a 1 b (+ a 1)) b)
// 2 ; a == b, b == a+1 == 2
func evalLet(
	operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, newEnv *environment.Env, err error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf("let takes 2 args")
	}
	bindingList, ok := args[0].(*types.MalList)
	if !ok {
		return nil, nil, fmt.Errorf("let: first arg isn't a list")
	}
	if len(bindingList.Items)%2 != 0 {
		return nil, nil, fmt.Errorf("let: first arg doesn't have an even number of items")
	}

	childEnv := env.ChildEnv()
	for i := 0; i < len(bindingList.Items); i += 2 {
		key, ok := bindingList.Items[i].(*types.MalSymbol)
		if !ok {
			return nil, nil, fmt.Errorf("let: binding list: arg %d isn't a symbol", i)
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
}

// evalDo evaluates the `do` special form
// Evaluates the elements in the arg list and returns the final result.
// For TCO, we eval all but the last argument here, then return the last
// argument to be evaluated in the main Eval loop.
func evalDo(
	operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, newEnv *environment.Env, err error) {
	for _, arg := range args[:len(args)-1] {
		if _, err := Eval(arg, env); err != nil {
			return nil, nil, err
		}
	}
	return args[len(args)-1], env, nil
}

// evalIf evaluates the `if` special form.
// Evaluate first param. If not `nil` or `false`, return the second param
// to be evaluated. If it is, return the third param to be evaluated, or
// `nil` if none is supplied. If none is supplied, the `nil` value is
// evalulated, but just evaluates to `nil`.
func evalIf(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, newEnv *environment.Env, err error) {
	if numArgs := len(args); numArgs != 2 && numArgs != 3 {
		return nil, nil, fmt.Errorf("if statements must have two or three arguments, got %d", numArgs)
	}
	condition, err := Eval(args[0], env)
	if err != nil {
		return nil, nil, err
	}
	if isTruthy(condition) {
		return args[1], env, nil
	}

	if len(args) == 3 {
		return args[2], env, nil
	}

	return &types.MalNil{}, env, nil
}

func evalQuasiquote(operator *types.MalSymbol, args []types.MalType, env *environment.Env,
) (newAST types.MalType, newEnv *environment.Env, err error) {
	ast, err := quasiquote(args[0])
	if err != nil {
		return nil, nil, err
	}
	return ast, env, nil
}
