package mappath

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
