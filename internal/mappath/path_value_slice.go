package mappath

import (
	"fmt"
	"strconv"
	"strings"
)

// PathValueSlice represents an array of values.
type PathValueSlice []PathValue

var _ PathValue = PathValueSlice{}

func (s PathValueSlice) Exists() bool {
	return true
}

func (s PathValueSlice) IsEmpty() bool {
	return len(s) == 0
}

func (s PathValueSlice) IsNull() bool {
	return false
}

func (s PathValueSlice) IsArray() bool {
	return true
}

func (s PathValueSlice) IsObject() bool {
	return false
}

func (s PathValueSlice) IsString() bool {
	return false
}

func (s PathValueSlice) IsFloat() bool {
	return false
}

func (s PathValueSlice) IsBool() bool {
	return false
}

func (s PathValueSlice) Bool() bool {
	return false
}

func (s PathValueSlice) Float() float64 {
	return 0
}

func (s PathValueSlice) Int() int64 {
	return 0
}

func (s PathValueSlice) String() string {
	return ""
}

func (s PathValueSlice) Value() any {
	//	[]any, for JSON arrays
	values := make([]any, 0, len(s))
	for _, v := range s {
		values = append(values, v.Value())
	}

	return values
}

func (s PathValueSlice) Get(path string) PathValue {
	if path == "#" {
		return &scalarValue{v: len(s)}
	}
	if strings.HasPrefix(path, "#(") {
		// Apply the condition
		key, remainingPath, _ := strings.Cut(path, ".")
		childPath := key[2:]
		shouldReturnList := false
		if strings.HasSuffix(childPath, ")#") {
			childPath = childPath[:len(childPath)-2] // remove trailing ")#" chars
			shouldReturnList = true
		} else {
			childPath = childPath[:len(childPath)-1] // remove the parenthesis
		}

		slice := s.applyChildConditionPath(childPath)
		var result PathValue = slice

		if !shouldReturnList {
			// Return the first result
			if len(slice) == 0 {
				result = undefinedScalar
			} else {
				result = slice[0]
			}
		}

		if remainingPath != "" {
			result = result.Get(remainingPath)
		}

		return result
	}
	if strings.HasPrefix(path, "#.") {
		// Apply the selection to the slice
		childPath := path[2:]
		newSlice := make(PathValueSlice, 0, len(s))
		for _, value := range s {
			v := value.Get(childPath)
			if v.Exists() {
				newSlice = append(newSlice, v)
			}
		}
		return newSlice
	}

	// path is not supported on a slice
	return undefinedScalar
}

func (s PathValueSlice) applyChildConditionPath(childPath string) PathValueSlice {
	newSlice := make(PathValueSlice, 0, len(s))
	for _, value := range s {
		conditionValue := value.Get(childPath)
		if conditionValue.Exists() {
			// Return the complete child value if the condition matches
			newSlice = append(newSlice, value)
		}
	}
	return newSlice
}

func (s PathValueSlice) Set(path string, rawValue any) PathValue {
	key, remainingPath, _ := strings.Cut(path, ".")

	// Apply the rawValue to all slice elements
	if key == "#" {
		result := s
		for index, element := range result {
			result[index] = element.Set(remainingPath, rawValue)
		}
		return result
	}

	index, err := strconv.Atoi(key)
	if err != nil {
		// expected an index for a PathValueSlice: noop
		return s
	}
	result := s
	if remainingPath == "" {
		if rawValue == deleteValue {
			// Remove by index
			if index < len(result) {
				result = append(result[:index], result[index+1:]...)
			}
			return result
		}
		result = growSliceIfNeeded(result, index)
		// Edit in place
		result[index] = toPathValue(toJSONValue(rawValue))
		return result
	}

	// Set in subpath
	result = growSliceIfNeeded(result, index)
	child := result[index]
	if child.IsNull() {
		child = toPathValue(createRawChild(remainingPath))
	}
	result[index] = child.Set(remainingPath, rawValue)

	return result
}

func (s PathValueSlice) Delete(path string) PathValue {
	return s.Set(path, deleteValue)
}

func growSliceIfNeeded(slice PathValueSlice, index int) PathValueSlice {
	initialLen := len(slice)
	if initialLen <= index {
		for i := 0; i < index+1-initialLen; i++ {
			slice = append(slice, &scalarValue{v: nil})
		}
	}

	return slice
}

func (s PathValueSlice) Array() PathValueSlice {
	return s
}

func (s PathValueSlice) Map() map[string]PathValue {
	return map[string]PathValue{}
}

// toPathValue returns a path value from a given json type (float64, string, bool, nil, map and array).
func toPathValue(jsonValue any) PathValue {
	if jsonValue == nil {
		return &scalarValue{v: nil}
	}

	switch v := jsonValue.(type) {
	case float64:
		return &scalarValue{v: v}
	case string:
		return &scalarValue{v: v}
	case bool:
		return &scalarValue{v: v}
	case map[string]any:
		return &Map{m: v}
	case []any:
		newSlice := make(PathValueSlice, 0, len(v))
		for _, item := range v {
			newSlice = append(newSlice, toPathValue(item))
		}
		return newSlice
	default:
		panic(fmt.Sprintf("type %T not supported, is it coming from json?", jsonValue)) // likely a developer issue
	}
}
