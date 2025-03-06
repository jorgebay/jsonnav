package jsonnav

import (
	"strconv"
	"strings"
)

// Slice represents an array of values.
type Slice []Value

var _ Value = Slice{}

// Exists returns true for slices.
func (s Slice) Exists() bool {
	return true
}

// IsEmpty returns true if the slice does not have any items.
func (s Slice) IsEmpty() bool {
	return len(s) == 0
}

// IsNull returns false for slices.
func (s Slice) IsNull() bool {
	return false
}

// IsArray returns true for slices.
func (s Slice) IsArray() bool {
	return true
}

// IsObject returns false for slices.
func (s Slice) IsObject() bool {
	return false
}

// IsString returns false for slices.
func (s Slice) IsString() bool {
	return false
}

// IsFloat returns false for slices.
func (s Slice) IsFloat() bool {
	return false
}

// IsBool returns false for slices.
func (s Slice) IsBool() bool {
	return false
}

// Bool returns false for slices.
func (s Slice) Bool() bool {
	return false
}

// Float returns 0 for slices.
func (s Slice) Float() float64 {
	return 0
}

// Int returns 0 for slices.
func (s Slice) Int() int64 {
	return 0
}

// String returns an empty string for slices.
func (s Slice) String() string {
	return ""
}

// At returns the value at the specified index.
// If the index is out of range, it returns an scalar with an internal nil value.
func (s Slice) At(index int) Value {
	if index < 0 {
		panic("index out of range")
	}

	if index >= len(s) {
		return undefinedScalar
	}
	return s[index]
}

// Value returns the inner values as a slice.
func (s Slice) Value() any {
	//	[]any, for JSON arrays
	values := make([]any, 0, len(s))
	for _, v := range s {
		values = append(values, v.Value())
	}

	return values
}

// Get searches for the specified path within the slice.
func (s Slice) Get(path string) Value {
	if path == "#" {
		return &scalar{v: len(s)}
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
		var result Value = slice

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
		newSlice := make(Slice, 0, len(s))
		for _, value := range s {
			v := value.Get(childPath)
			if v.Exists() {
				newSlice = append(newSlice, v)
			}
		}
		return newSlice
	}

	indexString, remainingPath, _ := strings.Cut(path, ".")
	index, err := strconv.Atoi(indexString)
	if err == nil && index >= 0 && index < len(s) {
		result := s[index]
		if remainingPath != "" {
			return result.Get(remainingPath)
		}
		return result
	}

	// path is not supported on a slice
	return undefinedScalar
}

func (s Slice) applyChildConditionPath(childPath string) Slice {
	newSlice := make(Slice, 0, len(s))
	for _, value := range s {
		conditionValue := value.Get(childPath)
		if conditionValue.Exists() {
			// Return the complete child value if the condition matches
			newSlice = append(newSlice, value)
		}
	}
	return newSlice
}

// Set sets the value at the specified path.
func (s Slice) Set(path string, rawValue any) Value {
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
		result[index] = mustToPathValue(toJSONValue(rawValue))
		return result
	}

	// Set in subpath
	result = growSliceIfNeeded(result, index)
	child := result[index]
	if child.IsNull() {
		child = mustToPathValue(createRawChild(remainingPath))
	}
	result[index] = child.Set(remainingPath, rawValue)

	return result
}

// Delete removes the value at the specified path.
func (s Slice) Delete(path string) Value {
	return s.Set(path, deleteValue)
}

func growSliceIfNeeded(slice Slice, index int) Slice {
	initialLen := len(slice)
	if initialLen <= index {
		for i := 0; i < index+1-initialLen; i++ {
			slice = append(slice, &scalar{v: nil})
		}
	}

	return slice
}

// Array returns the same instance.
func (s Slice) Array() Slice {
	return s
}

// Map returns an empty map for slices.
func (s Slice) Map() map[string]Value {
	return map[string]Value{}
}
