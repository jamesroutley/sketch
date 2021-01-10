package evaluator

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
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
	case "import":
		evaluator = evalImport
	case "export-as":
		evaluator = evalExportAs
	case "module-lookup":
		evaluator = evalModuleLookup

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
	if numArgs := len(args); numArgs != 2 && numArgs != 3 {
		return nil, fmt.Errorf("fn statements must have two or three arguments, got %d", numArgs)
	}

	// Functions can optionally have a docstring set as its first argument
	var docstring string
	docstringSet := len(args) == 3
	if docstringSet {
		str, ok := args[0].(*types.SketchString)
		if !ok {
			return nil, fmt.Errorf("if a fn expression has three arguments, the 1st should be a docstrig with type string. Got %s", args[0].Type())
		}
		docstring = str.Value
		// We've processed the first argument pop it off the list for the rest
		// of the processing
		args = args[1:]
	}

	// arguments is the first argument supplied to the fn function (e.g.
	// `(a)` in the example above)
	arguments, ok := args[0].(*types.SketchList)
	if !ok {
		err := fmt.Errorf("fn statements must have a list as the 1st arg, got %s", args[0].Type())
		if docstringSet {
			err = fmt.Errorf("fn statements must have a list as the 2nd arg, got %s", args[0].Type())
		}
		return nil, err
	}

	// Cast it from a list of SketchType to a list of SketchSymbol
	binds := make([]*types.SketchSymbol, len(arguments.Items))
	for i, a := range arguments.Items {
		bind, ok := a.(*types.SketchSymbol)
		if !ok {
			return nil, fmt.Errorf("fn statements must have a list of symbols as the first arg, the parameter list. Parameter %d (`%s`) has type %s", i, a.String(), a.Type())
		}
		binds[i] = bind
	}

	return &types.SketchFunction{
		// All functions are by default tail call optimised. This means,
		// instead of calling this object's Func() method (which recursively
		// calls Eval), in the Eval loop, we create a new env using `Params`,
		// and jump to the top of the Eval loop, setting `env` to be this new
		// environment, and `ast` to be this function's AST value.
		TailCallOptimised: true,
		AST:               args[1],
		Params:            binds,
		Env:               env,
		Docs:              docstring,
		// This is the non-tail call optimised function. We don't call this
		// during normal execution, but it's sometimes useful to be able to
		// execute a function from the function type itself. We do this in the
		// stdlib function `map`, where we want to execute the passed in
		// function, but don't have access to the Eval loop to tail call
		// optimise it.
		Func: func(exprs ...types.SketchType) (types.SketchType, error) {
			childEnv, err := environment.NewFunctionEnv(
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
	if err := validation.NArgs("def", 2, args); err != nil {
		return nil, err
	}
	key, err := validation.SymbolArg("def", args[0], 0)
	if err != nil {
		return nil, err
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
	if err := validation.NArgs("quote", 1, args); err != nil {
		return nil, err
	}
	return args[0], nil
}

// evalQuasiquoteExpand evaluates the `quasiquoteexpand` macro.
// This macro is used to test the internal implementation of quasiquote
func evalQuasiquoteExpand(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if err := validation.NArgs("quasiquoteexpand", 1, args); err != nil {
		return nil, err
	}
	return quasiquote(args[0])
}

// Creates a new macro
func evalDefmacro(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if err := validation.NArgs("defmacro", 2, args); err != nil {
		return nil, err
	}
	key, err := validation.SymbolArg("defmacro", args[0], 0)
	if err != nil {
		return nil, err
	}
	value, err := Eval(args[1], env)
	if err != nil {
		return nil, err
	}
	function, err := validation.FunctionArg("defmacro", value, 1)
	if err != nil {
		return nil, err
	}
	function.IsMacro = true
	env.Set(key.Value, function)
	return function, nil
}

func evalMacroexpand(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if err := validation.NArgs("macroexpand", 1, args); err != nil {
		return nil, err
	}
	return macroExpand(args[0], env)
}

// evalImport imports a module.
func evalImport(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if err := validation.NArgs("import", 1, args); err != nil {
		return nil, err
	}
	relativePath, err := validation.StringArg("import", args[0], 0)
	if err != nil {
		return nil, err
	}

	module, err := importModule(relativePath.Value)
	if err != nil {
		return nil, err
	}

	env.Set(module.DefaultName, module)

	return module, nil
}

func evalExportAs(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if err := validation.NArgs("export-as", 2, args); err != nil {
		return nil, err
	}
	defaultName, err := validation.SymbolArg("export-as", args[0], 0)
	if err != nil {
		return nil, err
	}

	exports, err := validation.ListArg("export-as", args[1], 1)
	if err != nil {
		return nil, err
	}

	var exported []string
	for i, item := range exports.Items {
		export, ok := item.(*types.SketchSymbol)
		if !ok {
			return nil, fmt.Errorf(
				"the function export-as expects the second argument to be a list of symbols to export, but the %s item is type %s",
				validation.ToOrdinal(i), item.Type())
		}

		// Check exported symbol is in the environment
		_, err := env.Get(export.Value)
		if err != nil {
			return nil, fmt.Errorf("export-as: cannot export: %w", err)
		}
		exported = append(exported, export.Value)
	}

	module := &types.SketchModule{
		Environment: env,
		SourceFile:  "TODO",
		Exported:    exported,
		DefaultName: defaultName.Value,
		// Maybe we shoulnd't set this here
		Name: defaultName.Value,
	}

	// - Creates a new `SketchModule` object using the current environment. It
	//   validates that all exported symbols are present in that environment
	// - Returns the module
	return module, nil
}

func evalModuleLookup(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	if err := validation.NArgs("module-lookup", 2, args); err != nil {
		return nil, err
	}
	moduleName, err := validation.SymbolArg("module-lookup", args[0], 0)
	if err != nil {
		return nil, err
	}

	valueName, err := validation.SymbolArg("module-lookup", args[1], 1)
	if err != nil {
		return nil, err
	}

	module, err := env.Get(moduleName.Value)
	if err != nil {
		return nil, err
	}

	m, ok := module.(*types.SketchModule)
	if !ok {
		return nil, fmt.Errorf("module-lookup: %s isn't a module, got %s", moduleName.Value, module.Type())
	}

	return m.Environment.Get(valueName.Value)
}
