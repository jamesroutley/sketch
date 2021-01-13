package evaluator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesroutley/sketch/sketch/reader"
	"github.com/jamesroutley/sketch/sketch/stdlib/file"
	"github.com/jamesroutley/sketch/sketch/stdlib/str"
	"github.com/jamesroutley/sketch/sketch/types"
)

type registeredModule struct {
	EnvironmentItems map[string]types.SketchType
	SketchCode       string
}

var registeredModules = map[string]*registeredModule{}

func registerModule(name string, items map[string]types.SketchType, code string) {
	registeredModules[name] = &registeredModule{
		EnvironmentItems: items,
		SketchCode:       code,
	}
}

func init() {
	registerModule("string", str.EnvironmentItems, str.SketchCode)
	registerModule("file", file.EnvironmentItems, file.SketchCode)
}

func loadStdlibModule(name string) (*types.SketchModule, error) {
	rawModule, ok := registeredModules[name]
	if !ok {
		return nil, fmt.Errorf("could not find stdlib module %s", name)
	}

	env, err := RootEnvironment()
	if err != nil {
		return nil, err
	}

	// Add the items defined in Go to the environment
	exportedItems := make([]string, 0, len(rawModule.EnvironmentItems))
	for key, value := range rawModule.EnvironmentItems {
		env.Set(key, value)
		exportedItems = append(exportedItems, key)
	}

	if strings.TrimSpace(rawModule.SketchCode) == "" {
		return &types.SketchModule{
			Environment: env,
			SourceFile:  name,
			Exported:    exportedItems,
			DefaultName: name,
			// Maybe we shoulnd't set this here
			Name: name,
		}, nil
	}
	// Pull the exported module from any Sketch code
	ast, err := reader.ReadStr(fmt.Sprintf("(do %s)", rawModule.SketchCode))
	if err != nil {
		return nil, err
	}
	evaluated, err := Eval(ast, env)
	if err != nil {
		if err.Error() == "read comment" {
			return nil, fmt.Errorf("to be importable, %s must end in an `export-as` statement", name)
		}
		return nil, err
	}

	module, ok := evaluated.(*types.SketchModule)
	if !ok {
		return nil, fmt.Errorf("to be importable, %s must end in an `export-as` statement", name)
	}

	module.Exported = append(module.Exported, exportedItems...)
	return module, nil
}

func importModule(path string) (*types.SketchModule, error) {
	if _, ok := registeredModules[path]; ok {
		return loadStdlibModule(path)
	}

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return nil, fmt.Errorf("import: $GOPATH not set")
	}

	fullPath := filepath.Join(goPath, "src", path)

	moduleEnv, err := RootEnvironment()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(fullPath)
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
			return nil, fmt.Errorf("to be importable, %s must end in an `export-as` statement", path)
		}
		return nil, err
	}

	module, ok := evaluated.(*types.SketchModule)
	if !ok {
		return nil, fmt.Errorf("to be importable, %s must end in an `export-as` statement", path)
	}

	return module, nil
}
