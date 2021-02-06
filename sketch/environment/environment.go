// Package environment defines Sketch's environments, which are objects which
// store the mapping of symbols to values bound in a particular scope. For
// example, when the Sketch interpreter starts, we create a root environment,
// into which we load all the core functions. When the user defines a function,
// we create a new environment in which the arguments the function was called
// with are bound to the function's parameter names.
package environment

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/types"
)

type Env struct {
	Outer *Env
	Data  map[string]types.SketchType
}

func NewEnv() *Env {
	return &Env{
		Outer: nil,
		Data:  map[string]types.SketchType{},
	}
}

// NewFunctionEnv creates a new environment with `parent` as its outer
// environment. It also takes a a list of arguments, which should be bound to
// the symbols in `parameters` one by one.
func NewFunctionEnv(parent *Env, parameters []*types.SketchSymbol, arguments []types.SketchType) (*Env, error) {
	env := &Env{
		Outer: parent,
		Data:  map[string]types.SketchType{},
	}

	if err := validateBindList(parameters, arguments); err != nil {
		return nil, err
	}

	for i, symbol := range parameters {
		// Variadic arguments. Bind the remaining arguments to the symbol after
		// the &.
		if symbol.Value == "&" {
			// Validate that only one parameter is specified after the &
			collectorSymbols := parameters[i+1:]
			switch len(collectorSymbols) {
			case 1:
				// continue
			case 0:
				return nil, fmt.Errorf("variadic arguments: no collector specified")
			default:
				return nil, fmt.Errorf("variadic arguments: you can only specify one collector argument")
			}

			collectorSymbol := collectorSymbols[0]
			env.Set(collectorSymbol.Value, &types.SketchList{
				List: types.NewList(arguments[i:]),
			})
			return env, nil
		}
		env.Set(symbol.Value, arguments[i])
	}
	return env, nil
}

func validateBindList(parameters []*types.SketchSymbol, arguments []types.SketchType) error {
	variadicArguments := false
	numRequiredArgs := 0 // Only valid if variadicArguments == true
	for i, param := range parameters {
		if param.Value != "&" {
			continue
		}
		variadicArguments = true
		if collectors := parameters[i+1:]; len(collectors) != 1 {
			return fmt.Errorf("there can only be one collector symbol after the & in a function definition, got %s", collectors)
		}
		numRequiredArgs = i
	}
	if !variadicArguments {
		if len(parameters) != len(arguments) {
			return fmt.Errorf("can't create env - num parameters (%d) != num arguments (%d)", len(parameters), len(arguments))
		}
		return nil
	}

	// variadicArguments == true from here on in

	if numRequiredArgs > len(arguments) {
		return fmt.Errorf("can't create env - num required parameters (%d) > num arguments (%d)", numRequiredArgs, len(arguments))
	}

	return nil
}

func (e *Env) Set(key string, value types.SketchType) {
	e.Data[key] = value
}

// TODO: would be nice to switch this to return `comma ok`, rather than an error
func (e *Env) Find(key string) (types.EnvType, error) {
	if e == nil {
		return nil, fmt.Errorf("`%s` is undefined", key)
	}
	if _, ok := e.Data[key]; ok {
		return e, nil
	}
	return e.Outer.Find(key)
}

func (e *Env) Get(key string) (types.SketchType, error) {
	env, err := e.Find(key)
	if err != nil {
		return nil, err
	}
	return env.(*Env).Data[key], nil
}

func (e *Env) ChildEnv() *Env {
	return &Env{
		Outer: e,
		Data:  map[string]types.SketchType{},
	}
}
