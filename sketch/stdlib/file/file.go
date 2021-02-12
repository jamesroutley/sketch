package file

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
)

var EnvironmentItems = map[string]types.SketchType{}

func register(symbol string, f func(...types.SketchType) (types.SketchType, error)) {
	EnvironmentItems[symbol] = &types.SketchFunction{
		Func:      f,
		BoundName: symbol,
	}
}

func init() {
	register("read-all", readAll)
	register("read-lines", readLines)
}

func readAll(args ...types.SketchType) (types.SketchType, error) {
	filename, err := validation.StringArg("read-all", args[0], 0)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename.Value)
	if err != nil {
		return nil, err
	}

	return &types.SketchString{
		Value: string(data),
	}, nil
}

func readLines(args ...types.SketchType) (types.SketchType, error) {
	filename, err := validation.StringArg("read-lines", args[0], 0)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename.Value)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var items []types.SketchType
	for scanner.Scan() {
		line := scanner.Text()
		items = append(items, &types.SketchString{
			Value: line,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &types.SketchList{
		List: types.NewList(items),
	}, nil
}
