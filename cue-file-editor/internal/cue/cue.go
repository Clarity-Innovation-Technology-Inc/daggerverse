package cue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"

	"github.com/kubevela/workflow/pkg/cue/model/sets"
)

func Decode(filename string) (cue.Value, error) {
	currentValue := cue.Value{}

	ctx := cuecontext.New()
	entrypoints := []string{filename}

	bis := load.Instances(entrypoints, nil)
	for _, bi := range bis {
		// check for errors on the instance
		// these are typically parsing errors
		if bi.Err != nil {
			return cue.Value{}, fmt.Errorf("error during load: %v", bi.Err)
		}

		value := ctx.BuildInstance(bi).Value()
		if value.Err() != nil {
			return cue.Value{}, fmt.Errorf("error during build: %v", value.Err())
		}

		// Validate the value
		err := value.Validate()
		if err != nil {
			return cue.Value{}, fmt.Errorf("error during validation: %v", err)
		}

		currentValue = value
	}
	return currentValue, nil
}

func UpdateValues(oldVal cue.Value, keyPath, newVal string) (cue.Value, error) {
	return replace(oldVal, keyPath, newVal)
}

// Replace replaces the value at the given path with the given value.
func replace(v cue.Value, path string, value interface{}) (cue.Value, error) {
	p := cue.ParsePath(path)

	var emptyValue string
	switch value.(type) {
	case string:
		emptyValue = `string`
	case int:
		emptyValue = `int`
	case float64:
		emptyValue = `number`
	default:
		emptyValue = "[...]"
	}

	emptyBase := v.Context().CompileString(fmt.Sprintf(`{ %s: %s }`, strings.ReplaceAll(path, ".", ":"), emptyValue))
	n := emptyBase.FillPath(p, value)

	ret, err := sets.StrategyUnify(v, n, sets.UnifyByJSONMergePatch{})
	if err != nil {
		return cue.Value{}, err
	}

	return ret, nil
}
