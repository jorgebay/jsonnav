package mappath

import (
	"github.com/samber/lo"
)

// PathValue represents the result of a path search expression over unmarshalled json.
// For example: `value.Get("a.b")` will access the value in the object at the path "a.b".
//
// Get() and `Set()` methods partially supports GJSON syntax: https://github.com/tidwall/gjson/blob/master/SYNTAX.md
//
// Not supported expressions:
// - Nested conditions like "friends.#(nets.#(=="fb"))#"
// - Tilde comparison
//
// Background: We used to parse json over and over again using gjson library as part of config version compatibility.
// This had a significant impact on CPU and memory resource usage.
// This structure tries to bridge the gap between dealing with raw maps and avoiding to use an external library and
// parsing over and over.
type PathValue interface {
	// Exists returns true if value exists and it's not null.
	Exists() bool

	// IsEmpty returns true if the value is
	//   - a null JSON value
	//   - an object containing no items
	//   - an array of zero length
	//   - an empty string
	IsEmpty() bool

	// IsNull returns true when the value is representation of the JSON null value.
	IsNull() bool

	// IsArray returns true if the value is a JSON array.
	IsArray() bool

	// IsObject returns true if the value is a JSON object/map.
	IsObject() bool

	// IsString returns true if the value is a string scalar.
	IsString() bool

	// IsFloat returns true if the value is a float/number scalar.
	IsFloat() bool

	// IsBool returns true if the value is a bool scalar.
	IsBool() bool

	// Bool returns a boolean representation.
	Bool() bool

	// Float returns a float64 representation.
	Float() float64

	// Int returns an integer representation.
	Int() int64

	// String returns a string representation of the value.
	// If the internal value is a bool or float64 scalar, it will be converted to string.
	// If the internal value is a JSON object or array, it will return empty string.
	String() string

	// Value returns one of these types:
	//
	//	bool, for JSON booleans
	//	float64, for JSON numbers
	//	Number, for JSON numbers
	//	string, for JSON string literals
	//	nil, for JSON null
	//	map[string]any, for JSON objects
	//	[]any, for JSON arrays
	Value() any

	// Get searches for the specified path.
	Get(path string) PathValue

	// Set sets the value in the provided path and returns the modified instance.
	// If the path does not exist, it will be created.
	Set(path string, rawValue any) PathValue

	// Delete deletes the value in the provided path and returns the modified instance.
	Delete(path string) PathValue

	// Array returns back an array of values.
	// If the result represents a null value or is non-existent, then an empty
	// array will be returned.
	// If the result is not a JSON array, the return value will be an
	// array containing one result.
	Array() PathValueSlice

	// Map returns back a map of values. The result should be a JSON object.
	// If the result is not a JSON object, the return value will be an empty map.
	Map() map[string]PathValue
}

// An internal type to mark the set operation as a delete.
type deleteType int

const deleteValue deleteType = -1

// EmptyOrEqualToDefault checks whether the given value is either empty or equal to the provided default value.
// It supports various types including bool, string, float, objects, and arrays. It returns `true` if `value` is
// empty or equal to `defaultValue` and`false` otherwise.
//
// - Returns true if:
//   - `value` does not exist, is empty, or equals `defaultValue`.
//   - For objects, it recursively checks all keys in the object.
//   - For arrays, it recursively checks all elements.
func EmptyOrEqualToDefault(value, defaultValue PathValue) bool {
	if !value.Exists() || !defaultValue.Exists() || value.IsEmpty() {
		return true
	}
	if value.IsBool() || value.IsString() || value.IsFloat() {
		if value.Value() == defaultValue.Value() {
			return true
		}
	}
	if value.IsObject() && defaultValue.IsObject() {
		defaultValueMap := defaultValue.Map()
		for key, defaultItem := range defaultValueMap {
			if !lo.HasKey(value.Map(), key) {
				return false
			}
			item := value.Get(key)
			if !EmptyOrEqualToDefault(item, defaultItem) {
				return false
			}
		}
		return true
	}
	if value.IsArray() && defaultValue.IsArray() {
		valueArray := value.Array()
		defaultValueArray := defaultValue.Array()
		for i, defaultItem := range defaultValueArray {
			item := valueArray[i]
			if !EmptyOrEqualToDefault(item, defaultItem) {
				return false
			}
		}
		return true
	}
	return false
}
