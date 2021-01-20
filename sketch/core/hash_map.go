package core

import (
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
