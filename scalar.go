package jsonnav

import (
	"strconv"
)

// scalar represents either a boolean, float64 or string type. Inner value can be nil.
type scalar struct {
	v any
}

// Represents an undefined JSON value.
var undefinedScalar Value = &scalar{v: nil}

func (s *scalar) Exists() bool {
	return s != undefinedScalar
}

func (s *scalar) IsArray() bool {
	return false
}

func (*scalar) IsObject() bool {
	return false
}

func (s *scalar) IsString() bool {
	if s.v == nil {
		return false
	}
	if _, ok := s.v.(string); ok {
		return true
	}
	return false
}

func (s *scalar) IsFloat() bool {
	if s.v == nil {
		return false
	}
	if _, ok := s.v.(float64); ok {
		return true
	}
	return false
}

func (s *scalar) IsBool() bool {
	if s.v == nil {
		return false
	}
	if _, ok := s.v.(bool); ok {
		return true
	}
	return false
}

func (s *scalar) IsEmpty() bool {
	return s.v == nil || s.v == ""
}

func (s *scalar) IsNull() bool {
	return s.v == nil
}

func (s *scalar) Bool() bool {
	return s.v == true
}

func (s *scalar) Float() float64 {
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

func (s *scalar) Int() int64 {
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

func (s *scalar) String() string {
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

func (s *scalar) Value() any {
	return s.v
}

func (s *scalar) Get(path string) Value {
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

func (s *scalar) Set(_ string, _ any) Value {
	// Scalar values can't be set by path: noop
	return s
}

func (s *scalar) Delete(_ string) Value {
	// Scalar values can't be deleted by path: noop
	return s
}

func (s *scalar) Array() Slice {
	if s.v == nil {
		return []Value{}
	}
	return []Value{s}
}

func (s *scalar) Map() map[string]Value {
	return map[string]Value{}
}
