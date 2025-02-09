package jsonnav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalMap(t *testing.T) {
	t.Run("should return a valid map", func(t *testing.T) {
		result, err := UnmarshalMap(`{"string": "value", "float64": 123, "bool": true}`)
		require.NoError(t, err)
		require.Equal(t, result.m, map[string]any{
			"string":  "value",
			"float64": 123.0,
			"bool":    true,
		})
		require.True(t, result.Exists())
		require.False(t, result.IsArray())
		require.False(t, result.Bool())
	})

	t.Run("should fail when it's not a valid map", func(t *testing.T) {
		_, err := UnmarshalMap(`true`)
		require.ErrorContains(t, err, "cannot unmarshal bool into Go value of type")
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("should return a valid value from slice", func(t *testing.T) {
		result, err := Unmarshal(`["a", "b"]`)
		require.NoError(t, err)
		require.Equal(t, result.Value(), []any{"a", "b"})
		require.True(t, result.Exists())
		require.True(t, result.IsArray())
		require.False(t, result.Bool() || result.IsObject() || result.IsString() || result.IsFloat())
	})

	t.Run("should return a valid value from map", func(t *testing.T) {
		result, err := Unmarshal(`{"string": "value", "float64": 123, "bool": true}`)
		require.NoError(t, err)
		require.Equal(t, result.Value(), map[string]any{
			"string":  "value",
			"float64": 123.0,
			"bool":    true,
		})
		require.True(t, result.Exists())
		require.False(t, result.IsArray())
		require.False(t, result.Bool() || result.IsString() || result.IsFloat())
	})

	t.Run("should return a valid value from scalar", func(t *testing.T) {
		result, err := Unmarshal(`"value"`)
		require.NoError(t, err)
		require.Equal(t, result.Value(), "value")
		require.True(t, result.Exists())
		require.True(t, result.IsString())
		require.False(t, result.IsArray() || result.IsObject() || result.Bool() || result.IsFloat())
	})

	t.Run("should return a valid value from null", func(t *testing.T) {
		result, err := Unmarshal(`null`)
		require.NoError(t, err)
		require.True(t, result.IsNull())
		require.True(t, result.Exists())
		require.IsType(t, &scalar{}, result)
	})
}

func TestFrom(t *testing.T) {
	t.Run("should return a valid map", func(t *testing.T) {
		v := map[string]any{
			"string":  "value",
			"float64": 123.0,
			"bool":    true,
		}
		result := From(v)
		require.True(t, result.IsObject())
		require.True(t, result.Exists())
		require.False(t, result.IsArray())
		require.False(t, result.Bool())
		require.Equal(t, result.Value(), v)
		require.Equal(t, result.Get("string"), &scalar{"value"})
		require.Equal(t, result.Get("float64"), &scalar{123.0})
		require.Equal(t, result.Get("bool"), &scalar{true})
	})

	t.Run("should return a valid slice", func(t *testing.T) {
		v := []any{"a", "b"}
		result := From(v)
		require.True(t, result.IsArray())
		require.True(t, result.Exists())
		require.False(t, result.IsObject())
		require.False(t, result.Bool())
		require.Equal(t, result.Value(), v)
		require.Equal(t, result.Array(), Slice{From("a"), From("b")})
	})

	t.Run("should return a valid scalar", func(t *testing.T) {
		for _, value := range []any{true, false, 123.0, "value", ""} {
			switch v := value.(type) {
			case bool:
				require.Equal(t, From(v), &scalar{v})
			case float64:
				require.Equal(t, From(v), &scalar{v})
			case string:
				require.Equal(t, From(v), &scalar{v})
			}
		}
	})
}
