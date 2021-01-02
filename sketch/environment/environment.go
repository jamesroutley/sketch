package environment

import (
	"fmt"

	"github.com/jamesroutley/sketch/sketch/types"
)

type Env struct {
	Outer *Env
	Data  map[string]types.MalType
}

func NewEnv() *Env {
	return &Env{
		Outer: nil,
		Data:  map[string]types.MalType{},
	}
}

// NewChildEnv creates a new environment with `parent` as its outer
// environment. It also takes a a list of arguments, which should be bound to
// the symbols in `parameters` one by one.
func NewChildEnv(parent *Env, parameters []*types.MalSymbol, arguments []types.MalType) (*Env, error) {
	env := &Env{
		Outer: parent,
		Data:  map[string]types.MalType{},
	}

	variadicArguments := false
	for _, param := range parameters {
		if param.Value == "&" {
			variadicArguments = true
			break
		}
	}
	if !variadicArguments && (len(parameters) != len(arguments)) {
		return nil, fmt.Errorf("can't create env - num parameters (%d) != num arguments (%d)", len(parameters), len(arguments))
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
			env.Set(collectorSymbol.Value, &types.MalList{
				Items: arguments[i:],
			})
			return env, nil
		}
		env.Set(symbol.Value, arguments[i])
	}
	return env, nil
}

func (e *Env) Set(key string, value types.MalType) {
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

func (e *Env) Get(key string) (types.MalType, error) {
	env, err := e.Find(key)
	if err != nil {
		return nil, err
	}
	return env.(*Env).Data[key], nil
}

func (e *Env) ChildEnv() *Env {
	return &Env{
		Outer: e,
		Data:  map[string]types.MalType{},
	}
}
