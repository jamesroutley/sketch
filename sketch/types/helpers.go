package types

import "fmt"

func ValidHashMapKey(arg SketchType) error {
	switch arg.(type) {
	case *SketchInt, *SketchString, *SketchSymbol, *SketchList, *SketchBoolean:
		return nil
	}
	return fmt.Errorf("hash map argument %s has type %s - can't use this as a hash map key", arg.String(), arg.Type())
}
