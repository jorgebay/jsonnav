package jsonnav

import (
	"strconv"
	"strings"
)

// Map represents a JSON object.
type Map struct {
	m map[string]any
}

// Exists returns true if the value is defined.
func (m *Map) Exists() bool {
	return true
}

// IsEmpty returns true if the map does not have any items.
func (m *Map) IsEmpty() bool {
	return len(m.m) == 0
}

// IsNull returns false for maps.
func (*Map) IsNull() bool {
	return false
}

// IsArray returns false for maps.
func (*Map) IsArray() bool {
	return false
}

// IsObject returns true for maps.
func (*Map) IsObject() bool {
	return true
}

// IsString returns false for maps.
func (*Map) IsString() bool {
	return false
}

// IsFloat returns false for maps.
func (*Map) IsFloat() bool {
	return false
}

// IsBool returns false for maps.
func (*Map) IsBool() bool {
	return false
}

// IsInt returns false for maps.
func (m *Map) Bool() bool {
	return false
}

// Float returns 0 for maps.
func (m *Map) Float() float64 {
	return 0
}

// Int returns 0 for maps.
func (m *Map) Int() int64 {
	return 0
}

// String returns an empty string for maps.
func (m *Map) String() string {
	return ""
}

// Value returns the underlying map.
func (m *Map) Value() any {
	return m.m
}

// Get searches json for the specified path.
func (m *Map) Get(path string) Value {
	if len(path) == 0 {
		panic("invalid zero length")
	}
	key, remainingPath, _ := strings.Cut(path, ".")
	equalityIndex := strings.Index(key, "=")
	if equalityIndex != -1 {
		// Only string conditions are supported
		condition := key[equalityIndex+1:]
		key = key[:equalityIndex]

		if !areEqualFromCondition(m.m[key], condition) {
			return undefinedScalar
		}

		// Condition for map matched, return itself
		if remainingPath == "" {
			return m
		}
		return m.Get(remainingPath)
	}

	var value Value
	if rawValue, ok := m.m[key]; ok {
		value = mustToPathValue(rawValue)
	} else {
		value = undefinedScalar
	}

	if remainingPath == "" {
		return value
	}

	return value.Get(remainingPath)
}

// areEqualFromCondition determines whether 2 values are equal according to gjson syntax.
func areEqualFromCondition(value any, expected string) bool {
	if value == expected {
		return true
	}

	switch value.(type) {
	case bool:
		if typedExpected, err := strconv.ParseBool(expected); err == nil && value == typedExpected {
			return true
		}
	case float64:
		if typedExpected, err := strconv.ParseFloat(expected, 64); err == nil && value == typedExpected {
			return true
		}
	}

	return false
}

// Set updates the value at the specified path.
func (m *Map) Set(path string, rawValue any) Value {
	key, remainingPath, _ := strings.Cut(path, ".")
	if remainingPath == "" {
		m.setLeaf(key, rawValue)
		return m
	}
	if _, ok := m.m[key]; !ok {
		// Insert a branch
		m.m[key] = createRawChild(remainingPath)
	}

	m.m[key] = mustToPathValue(m.m[key]).Set(remainingPath, rawValue).Value()
	return m
}

// Delete removes the value at the specified path.
func (m *Map) Delete(path string) Value {
	return m.Set(path, deleteValue)
}

func createRawChild(remainingPath string) any {
	nextKey, _, _ := strings.Cut(remainingPath, ".")
	childIsSlice := false
	if _, err := strconv.Atoi(nextKey); err == nil {
		childIsSlice = true
	}
	if childIsSlice {
		return []any{}
	}
	return make(map[string]any)
}

func (m *Map) setLeaf(key string, rawValue any) {
	if rawValue == deleteValue {
		delete(m.m, key)
		return
	}
	m.m[key] = toJSONValue(rawValue)
}

func toJSONValue(rawValue any) any {
	if intValue, ok := rawValue.(int); ok {
		// We use ints and float64 in an indistinctive way across our codebase
		// Only float64 is valid json numbers
		return float64(intValue)
	}
	return rawValue
}

// Array returns an empty slice for maps.
func (m *Map) Array() Slice {
	return Slice{}
}

// Map returns the underlying map.
func (m *Map) Map() map[string]Value {
	newMap := make(map[string]Value, len(m.m))
	for k, v := range m.m {
		newMap[k] = mustToPathValue(v)
	}

	return newMap
}
