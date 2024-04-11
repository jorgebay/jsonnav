package mappath

import (
	"strconv"
	"strings"

	"github.com/kong/koko/internal/json"
)

type Map struct {
	m map[string]any
}

func (m *Map) Exists() bool {
	return true
}

func (m *Map) IsEmpty() bool {
	return len(m.m) == 0
}

func (*Map) IsNull() bool {
	return false
}

func (*Map) IsArray() bool {
	return false
}

func (*Map) IsObject() bool {
	return true
}

func (*Map) IsString() bool {
	return false
}

func (*Map) IsFloat() bool {
	return false
}

func (*Map) IsBool() bool {
	return false
}

func (m *Map) Bool() bool {
	return false
}

func (m *Map) Float() float64 {
	return 0
}

func (m *Map) Int() int64 {
	return 0
}

func (m *Map) String() string {
	return ""
}

func (m *Map) Value() any {
	return m.m
}

// Get searches json for the specified path.
func (m *Map) Get(path string) PathValue {
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

	var value PathValue
	if rawValue, ok := m.m[key]; ok {
		value = toPathValue(rawValue)
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

func (m *Map) Set(path string, rawValue any) PathValue {
	key, remainingPath, _ := strings.Cut(path, ".")
	if remainingPath == "" {
		m.setLeaf(key, rawValue)
		return m
	}
	if _, ok := m.m[key]; !ok {
		// Insert a branch
		m.m[key] = createRawChild(remainingPath)
	}

	m.m[key] = toPathValue(m.m[key]).Set(remainingPath, rawValue).Value()
	return m
}

func (m *Map) Delete(path string) PathValue {
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

func (m *Map) Array() PathValueSlice {
	return PathValueSlice{}
}

func (m *Map) Map() map[string]PathValue {
	newMap := make(map[string]PathValue, len(m.m))
	for k, v := range m.m {
		newMap[k] = toPathValue(v)
	}

	return newMap
}

// UnmarshalMap parses the json and returns the map.
func UnmarshalMap(v string) (*Map, error) {
	result := Map{}
	if err := json.Unmarshal([]byte(v), &result.m); err != nil {
		return nil, err
	}
	return &result, nil
}

// MarshalMap returns the json string for the provided map.
func MarshalMap(m *Map) (string, error) {
	blob, err := json.Marshal(m.m)
	if err != nil {
		return "", err
	}

	return string(blob), nil
}

// MustUnmarshalMap is a non-fallible version of UnmarshalMap() used for static variables and tests.
func MustUnmarshalMap(v string) *Map {
	result, err := UnmarshalMap(v)
	if err != nil {
		panic(err)
	}
	return result
}

// MustUnmarshalScalar parses the json and returns the scalar value used for tests.
func MustUnmarshalScalar(v string) PathValue {
	var result any
	if err := json.Unmarshal([]byte(v), &result); err != nil {
		panic(err)
	}

	if result == nil {
		return &scalarValue{v: nil}
	}

	switch result.(type) {
	case string:
		return &scalarValue{v: result}
	case bool:
		return &scalarValue{v: result}
	case float64:
		return &scalarValue{v: result}
	}

	panic("invalid scalar")
}
