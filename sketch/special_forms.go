package sketch

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jamesroutley/sketch/sketch/environment"
	"github.com/jamesroutley/sketch/sketch/reader"
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
	if len(args) == 3 {
		str, ok := args[0].(*types.SketchString)
		if !ok {
			return nil, fmt.Errorf("if a fn expression has three arguments, the first is expected to be a docstrig, with type string. Got %s", args[0].Type())
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
		return nil, fmt.Errorf("fn statements must have a list as the first arg, got %s", args[0].Type())
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

// evalImport imports a module.
func evalImport(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	relativePath, ok := args[0].(*types.SketchString)
	if !ok {
		return nil, fmt.Errorf("import: first arg isn't a string")
	}

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return nil, fmt.Errorf("import: $GOPATH not set")
	}

	path := filepath.Join(goPath, "src", relativePath.Value)

	moduleEnv, err := rootEnvironment()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ast, err := reader.ReadStr(fmt.Sprintf(`(do %s)`, data))
	if err != nil {
		return nil, err
	}
	evaluated, err := Eval(ast, moduleEnv)
	if err != nil {
		if err.Error() == "read comment" {
			return nil, fmt.Errorf("to be importable, %s must end in an `export-as` statement", relativePath)
		}
		return nil, err
	}

	module, ok := evaluated.(*types.SketchModule)
	if !ok {
		fmt.Println(evaluated.Type())
		return nil, fmt.Errorf("to be importable, %s must end in an `export-as` statement", relativePath)
	}

	env.Set(module.DefaultName, module)

	// - Locates the corresponding source code file
	// - Evaluates it in the context of a new environment
	// - Checks that a module is returned. This will happen if the last line in the
	//   file is an `export-as` call. If not, error.
	// - Binds the module to the importing environment, using `module.ExportedName`,
	//   or the name specified in an `import-as` expression
	// - Return the imported module. Or nil?
	return module, nil
}

func evalExportAs(operator *types.SketchSymbol, args []types.SketchType, env *environment.Env,
) (newAST types.SketchType, err error) {
	defaultName, ok := args[0].(*types.SketchSymbol)
	if !ok {
		return nil, fmt.Errorf("export-as: first arg isn't a symbol")
	}

	exports, ok := args[1].(*types.SketchList)
	if !ok {
		return nil, fmt.Errorf("export-as: second arg isn't a list")
	}

	var exported []string
	for i, export := range exports.Items {
		export, ok := export.(*types.SketchSymbol)
		if !ok {
			return nil, fmt.Errorf("export-as: export %d isn't a symbol", i)
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
	moduleName, ok := args[0].(*types.SketchSymbol)
	if !ok {
		return nil, fmt.Errorf("export-as: first arg isn't a symbol")
	}
	key, ok := args[1].(*types.SketchSymbol)
	if !ok {
		return nil, fmt.Errorf("export-as: second arg isn't a symbol")
	}

	module, err := env.Get(moduleName.Value)
	if err != nil {
		return nil, err
	}

	m, ok := module.(*types.SketchModule)
	if !ok {
		return nil, fmt.Errorf("module-lookup: %s isn't a module, got %s", moduleName.Value, module.Type())
	}

	return m.Environment.Get(key.Value)
}
