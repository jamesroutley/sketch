package core

import (
	"strings"

	"github.com/jamesroutley/sketch/sketch/types"
	"github.com/jamesroutley/sketch/sketch/validation"
)

func hashMapSet(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("hashmap-set", 3, args); err != nil {
		return nil, err
	}

	hashmap, err := validation.HashMapArg("hashmap-set", args[0], 0)
	if err != nil {
		return nil, err
	}

	if err := types.ValidHashMapKey(args[1]); err != nil {
		return nil, err
	}

	newHashMap := hashmap.Set(args[1], args[2])
	return newHashMap, nil
}

func hashMapGet(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgsRange("hashmap-get", 2, 3, args); err != nil {
		return nil, err
	}

	hashmap, err := validation.HashMapArg("hashmap-get", args[0], 0)
	if err != nil {
		return nil, err
	}

	defaultProvided := len(args) == 3

	value, err := hashmap.Get(args[1])
	if err != nil {
		if strings.Contains(err.Error(), "map doesn't contain key") && defaultProvided {
			return args[2], nil
		}
		return nil, err
	}

	return value, nil
}

func hashMapKeys(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("hashmap-keys", 1, args); err != nil {
		return nil, err
	}

	hashmap, err := validation.HashMapArg("hashmap-keys", args[0], 0)
	if err != nil {
		return nil, err
	}

	return &types.SketchList{
		Items: hashmap.Keys(),
	}, nil
}

func hashMapValues(args ...types.SketchType) (types.SketchType, error) {
	if err := validation.NArgs("hashmap-values", 1, args); err != nil {
		return nil, err
	}

	hashmap, err := validation.HashMapArg("hashmap-values", args[0], 0)
	if err != nil {
		return nil, err
	}

	return &types.SketchList{
		Items: hashmap.Values(),
	}, nil
}
