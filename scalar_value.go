package jsonnav

import (
	"strconv"
)

// scalarValue represents either a boolean, float64 or string type. Inner value can be nil.
type scalarValue struct {
	v any
}

// Represents an undefined JSON value.
var undefinedScalar PathValue = &scalarValue{v: nil}

func (s *scalarValue) Exists() bool {
	return s != undefinedScalar
}

func (s *scalarValue) IsArray() bool {
	return false
}

func (*scalarValue) IsObject() bool {
	return false
}

func (s *scalarValue) IsString() bool {
	if s.v == nil {
		return false
	}
	if _, ok := s.v.(string); ok {
		return true
	}
	return false
}

func (s *scalarValue) IsFloat() bool {
	if s.v == nil {
		return false
	}
	if _, ok := s.v.(float64); ok {
		return true
	}
	return false
}

func (s *scalarValue) IsBool() bool {
	if s.v == nil {
		return false
	}
	if _, ok := s.v.(bool); ok {
		return true
	}
	return false
}

func (s *scalarValue) IsEmpty() bool {
	return s.v == nil || s.v == ""
}

func (s *scalarValue) IsNull() bool {
	return s.v == nil
}

func (s *scalarValue) Bool() bool {
	return s.v == true
}

func (s *scalarValue) Float() float64 {
	if s.v == nil {
		return 0
	}
	switch v := s.v.(type) {
	case float64:
		return v
	case string:
		n, _ := strconv.ParseFloat(v, 64)
		return n
	default:
		return 0
	}
}

func (s *scalarValue) Int() int64 {
	if s.v == nil {
		return 0
	}
	switch v := s.v.(type) {
	case float64:
		return int64(v)
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	default:
		return 0
	}
}

func (s *scalarValue) String() string {
	if s.v == nil {
		return ""
	}
	switch v := s.v.(type) {
	case float64:
		return strconv.FormatFloat(v, 'E', -1, 64)
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

func (s *scalarValue) Value() any {
	return s.v
}

func (s *scalarValue) Get(path string) PathValue {
	if path[0] == '=' {
		// Equality check
		if path[1] == '"' { // string literal
			stringValue := path[2 : len(path)-1] // remove quotes
			if s.v == stringValue {
				return s
			}
		}
	}

	// No nested values
	return undefinedScalar
}

func (s *scalarValue) Set(_ string, _ any) PathValue {
	// Scalar values can't be set by path: noop
	return s
}

func (s *scalarValue) Delete(_ string) PathValue {
	// Scalar values can't be deleted by path: noop
	return s
}

func (s *scalarValue) Array() PathValueSlice {
	if s.v == nil {
		return []PathValue{}
	}
	return []PathValue{s}
}

func (s *scalarValue) Map() map[string]PathValue {
	return map[string]PathValue{}
}
