package mappath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testJSON = `{
	"string": "value",
	"emptyString": "",
	"float64": 123,
	"bool": true,
	"bool2": false,
	"nil": null,
	"object": {"a": 1.1},
	"array": ["a", "b"],
	"nestedObject": {"name": "John", "attrs": {"age": 31}},
	"nestedArray": [{"name": "Alice", "attrs": {"age": 33}}, {"name": "Bob"}]
}`

func TestMap_Unmarshal(t *testing.T) {
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
		require.ErrorContains(t, err, "error found in #1 byte")
	})
}

func TestMap_Get(t *testing.T) {
	value, err := UnmarshalMap(testJSON)
	require.NoError(t, err)

	t.Run("should return a simple scalar value", func(t *testing.T) {
		require.Equal(t, "value", value.Get("string").String())
		require.Equal(t, 123.0, value.Get("float64").Float())
		require.Equal(t, int64(123), value.Get("float64").Int())
		require.Equal(t, true, value.Get("bool").Bool())
		require.Equal(t, false, value.Get("bool2").Bool())
		require.Equal(t, &scalarValue{v: nil}, value.Get("nil"))
		require.True(t, value.Get("nil").Exists())
		require.False(t, value.Get("NOT_EXISTS").Exists())
		require.Equal(t, &Map{m: map[string]any{"a": 1.1}}, value.Get("object"))
		require.True(t, value.Get("array").IsArray())
		require.Equal(t, PathValueSlice{&scalarValue{v: "a"}, &scalarValue{v: "b"}}, value.Get("array").Array())
	})
	t.Run("should return nested values", func(t *testing.T) {
		require.False(t, value.Get("string.z").Exists())
		require.Equal(t, "John", value.Get("nestedObject.name").String())
		require.Equal(t, int64(31), value.Get("nestedObject.attrs.age").Int())
		require.True(t, value.Get("nestedArray.#.name").IsArray())
		require.Equal(t,
			PathValueSlice{&scalarValue{v: "Alice"}, &scalarValue{v: "Bob"}},
			value.Get("nestedArray.#.name"))
	})
	t.Run("should support conditions", func(t *testing.T) {
		require.Equal(t, "value", value.Get("string").Get(`="value"`).String())
		require.False(t, value.Get("string").Get(`="zzz"`).Exists())
		require.Equal(t, value, value.Get("string=value"))
		require.True(t, value.Get("nestedArray.#(name=Bob)").IsObject())
		require.Equal(t, "Bob", value.Get("nestedArray.#(name=Bob)").Get("name").String())
		require.True(t, value.Get("nestedArray.#(name=Bob)#").IsArray())
		require.Equal(t, 1, len(value.Get("nestedArray.#(name=Bob)#").Array()))
		require.Equal(t, "Bob", value.Get("nestedArray.#(name=Bob)#").Array()[0].Get("name").String())
		require.Equal(t, 1, len(value.Get("nestedArray.#(attrs)#").Array()))
		require.Equal(t, "Alice", value.Get("nestedArray.#(attrs)#").Array()[0].Get("name").String())
		require.Equal(
			t,
			PathValueSlice{&scalarValue{v: 33.0}},
			value.Get("nestedArray.#(name=Alice)#").Get("#.attrs.age"))
		require.Equal(t, &scalarValue{v: 33.0}, value.Get("nestedArray.#(name=Alice).attrs.age"))
		require.Equal(t, &scalarValue{v: "a"}, value.Get(`array.#(="a")`))
	})
}

func TestMap_Set(t *testing.T) {
	t.Run("should set leaf values", func(t *testing.T) {
		value := MustUnmarshalMap(testJSON)
		value.Set("new", "new value")
		value.Set("string", "value2")
		require.Equal(t, "new value", value.Get("new").String())
		require.Equal(t, "value2", value.Get("string").String())
		require.Equal(t, 123.0, value.Get("float64").Float())
		require.Equal(t, int64(123), value.Get("float64").Int())
		require.Equal(t, true, value.Get("bool").Bool())
	})

	t.Run("should set branch values", func(t *testing.T) {
		value := MustUnmarshalMap(testJSON)
		require.Equal(t, 1.1, value.Get("object.a").Float())
		value.Set("object.a", 2.0)
		value.Set("object.b", 3.1)
		value.Set("object.c.prop1.nested", 4.0)
		require.Equal(t, 2.0, value.Get("object.a").Float())
		require.Equal(t, 3.1, value.Get("object.b").Float())
		require.True(t, value.Get("object.c").Exists())
		require.True(t, value.Get("object.c.prop1").Exists())
		require.Equal(t, 4.0, value.Get("object.c.prop1.nested").Float())
		require.Equal(t, "value", value.Get("string").String())
	})

	t.Run("should set using indices", func(t *testing.T) {
		value := MustUnmarshalMap(testJSON)
		value.Set("array.1", "c")
		require.Equal(t, []any{"a", "c"}, value.Get("array").Value())

		value.Set("nestedArray.1.lastName", "Smith")
		require.Equal(
			t,
			map[string]any{"name": "Bob", "lastName": "Smith"},
			value.Get("nestedArray.#(name=Bob)").Value())
	})

	t.Run("should set all elements of an array when using #", func(t *testing.T) {
		value := MustUnmarshalMap(testJSON)

		value.Set("nestedArray.#.name", "Chuck Norris")
		require.Equal(
			t,
			[]any{"Chuck Norris", "Chuck Norris"},
			value.Get("nestedArray.#.name").Value())

		// nested path
		newAge := float64(42)
		value.Set("nestedArray.#.attrs.age", newAge)
		require.Equal(
			t,
			[]any{newAge, newAge},
			value.Get("nestedArray.#.attrs.age").Value())

		// delete value
		value.Set("nestedArray.#.attrs", deleteValue)
		require.Equal(
			t,
			[]any{map[string]any{"name": "Chuck Norris"}, map[string]any{"name": "Chuck Norris"}},
			value.Get("nestedArray").Value())
	})
}

func TestMap_Delete(t *testing.T) {
	t.Run("should delete values", func(t *testing.T) {
		value := MustUnmarshalMap(testJSON)
		value.Delete("new")
		value.Delete("string")
		value.Delete("array.0")
		value.Delete("nestedArray.0.attrs")
		require.False(t, value.Get("new").Exists())
		require.False(t, value.Get("string").Exists())
		require.Equal(t, true, value.Get("bool").Bool())
		require.Equal(t, []any{"b"}, value.Get("array").Value())
		require.Equal(t, map[string]any{"name": "Alice"}, value.Get("nestedArray.#(name=Alice)").Value())

		// Delete all nested elements in an array
		value = MustUnmarshalMap(testJSON)
		value.Delete("nestedArray.#.name")
		require.Equal(t,
			[]any{
				map[string]any{"attrs": map[string]any{"age": float64(33)}},
				map[string]any{},
			},
			value.Get("nestedArray").Value())
	})
}
