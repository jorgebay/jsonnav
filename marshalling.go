package jsonnav

import (
	"encoding/json"
	"fmt"
)

// UnmarshalMap parses the json and returns the map.
func UnmarshalMap(v string) (*Map, error) {
	result := Map{}
	if err := json.Unmarshal([]byte(v), &result.m); err != nil {
		return nil, err
	}
	return &result, nil
}

// Unmarshal parses the json and returns the Value.
func Unmarshal(jsonString string) (Value, error) {
	var value any
	if err := json.Unmarshal([]byte(jsonString), &value); err != nil {
		return nil, err
	}

	return toPathValue(value)
}

// MarshalMap returns the json string for the provided map.
func MarshalMap(m *Map) (string, error) {
	blob, err := json.Marshal(m.m)
	if err != nil {
		return "", err
	}

	return string(blob), nil
}

// Marshal returns the json string for the provided Value.
func Marshal(value Value) (string, error) {
	blob, err := json.Marshal(value.Value())
	if err != nil {
		return "", err
	}

	return string(blob), nil
}

// MustUnmarshalMap is a non-fallible version of UnmarshalMap() used for static variables and tests.
func MustUnmarshalMap(v string) *Map {
	return must(UnmarshalMap(v))
}

// FromJSONMap creates a new Map from a map[string]any.
// It expects that the internal map is already a valid json map (composed only by arrays, maps and scalars).
//
// Note that the provided map is not copied, so any modification in the original map will reflect in the Map.
func FromJSONMap(m map[string]any) *Map {
	return &Map{m: m}
}

// JSONValue is the type constraint for a json value.
type JSONValue interface {
	float64 | string | bool | map[string]any | []any
}

// From creates a new Value from a JSONValue.
//
// It panics if the value is not float64, string, bool, map or array.
// In the case of maps and slices it expects the child values to be composed only by valid json values
// (float64, string, bool, map and slice). Note that the provided value is not copied, so any modification in the
// original map/slice will reflect in the Value.
func From[T JSONValue](value T) Value {
	return mustToPathValue(value)
}

// FromAny creates a new Value from an any value.
//
// In the case of maps and slices it expects the child values to be composed only by valid json values
// (float64, string, bool, map and slice). Note that the provided value is not copied, so any modification in the
// original map/slice will reflect in the Value.
func FromAny(value any) Value {
	return mustToPathValue(value)
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// MustUnmarshalScalar parses the json and returns the scalar value.
func MustUnmarshalScalar(v string) Value {
	var result any
	if err := json.Unmarshal([]byte(v), &result); err != nil {
		panic(err)
	}

	if result == nil {
		return &scalar{v: nil}
	}

	switch result.(type) {
	case string:
		return &scalar{v: result}
	case bool:
		return &scalar{v: result}
	case float64:
		return &scalar{v: result}
	}

	panic("invalid scalar")
}

// mustToPathValue returns a path value from a given json type (float64, string, bool, nil, map and array).
func mustToPathValue(jsonValue any) Value {
	return must(toPathValue(jsonValue))
}

func toPathValue(jsonValue any) (Value, error) {
	if jsonValue == nil {
		return &scalar{v: nil}, nil
	}

	switch v := jsonValue.(type) {
	case float64:
		return &scalar{v: v}, nil
	case string:
		return &scalar{v: v}, nil
	case bool:
		return &scalar{v: v}, nil
	case map[string]any:
		return &Map{m: v}, nil
	case []any:
		newSlice := make(Slice, 0, len(v))
		for _, item := range v {
			newSlice = append(newSlice, mustToPathValue(item))
		}
		return newSlice, nil
	default:
		// A developer issue
		return nil, fmt.Errorf("type %T not supported, only values from json decoding are supported", jsonValue)
	}
}
