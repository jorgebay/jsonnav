# JSONNAV

jsonnav is a [Go package](https://pkg.go.dev/github.com/jorgebay/jsonnav#section-documentation) for accessing,
navigating and manipulating values from an untyped json document.

[![Build](https://github.com/jorgebay/jsonnav/actions/workflows/test.yml/badge.svg)](https://github.com/jorgebay/jsonnav/actions/workflows/test.yml)

## Features

- Retrieve values from deeply nested json documents safely.
- Built-in type check functions and conversions.
- Iterate over arrays and objects.
- Supports [GJSON][gjson] syntax for navigating the json document.
- Set or delete values in place.

## Installing

```shell
go get github.com/jorgebay/jsonnav
```

## Usage

```go
v, err := jsonnav.Unmarshal(`{"name":{"first":"Jimi","last":"Hendrix"},"age":27}`)
v.Get("name").Get("first").String() // "Jimi"
v.Get("name.last").String() // "Hendrix"
v.Get("age").Float() // 27.0
v.Get("path").Get("does").Get("not").Get("exist").Exists() // false
v.Get("path.does.not.exist").Exists() // false
```

It uses [GJSON syntax][gjson] for navigating the json document.

### Accessing values that may not exist

It's safe to access values that may not exist. The library will return a scalar `Value` representation
with `nil` underlying value.

```go
v, err := jsonnav.Unmarshal(`{
    "name": {"first": "Jimi", "last": "Hendrix"},
    "instruments": [{"name": "guitar"}]
}`)
v.Get("birth").Get("date").Exists() // false
v.Get("birth.date").Exists() // false
v.Get("instruments").Array().At(0).Get("name").String() // "guitar"
v.Get("instruments.0.name").String() // "guitar"
v.Get("instruments").Array().At(1).Get("name").String() // ""
v.Get("instruments.1.name").String() // ""
```

### Setting and deleting values

You can set or delete values in place using `Set()` and `Delete()` methods.

```go
v, err := jsonnav.Unmarshal(`{"name":{"first":"Jimi","last":"Hendrix"},"age":27}`)
v.Get("name").Set("middle", "Marshall")
v.Set("birth.date", "1942-11-27")
v.Delete("age")

v.Get("birth").Get("date").String()  // "1942-11-27"
v.Get("name").Get("middle").String() // "Marshall"
```

### Type checks and conversions

The library provides built-in functions for type checks and conversions that are safely free of errors and panics.

#### Type check functions

```go
v, err := jsonnav.Unmarshal(`{"name":"Jimi","age":27}`)
v.Get("name").IsString()     // true
v.Get("age").IsFloat()       // true
v.Get("age").IsBool()        // false
v.Get("age").IsObject()      // false
v.Get("age").IsArray()       // false
v.Get("not_found").IsNull()  // true
v.Get("not_found").IsEmpty() // true
```

#### Typed getters

```go
v, err := jsonnav.Unmarshal(`{
    "name": "Jimi",
    "age": 27,
    "instruments": ["guitar"],
    "hall_of_fame": true
}`)
v.Get("name").String()       // "Jimi"
v.Get("age").Float()         // 27.0
v.Get("hall_of_fame").Bool() // true
v.Get("instruments").Array() // a slice of 1 Value with underlying value "guitar"
```

When the value doesn't match the expected type or it does not exist, it will default to a zero value of the
expected type and do conversions for scalars.

- `String()` returns the string representation of float and bool values, otherwise an empty string.
- `Float()` returns the float representation of string values, for other types it returns 0.0.
- `Bool()` returns the bool representation of string values, for other types it returns false.
- `Array()` returns an empty slice for non-array values.

### Iterating over arrays

You can iterate over arrays using the `Array()` method.

```go
v, err := jsonnav.Unmarshal(`{
    "instruments": [
        {"name": "guitar"},
        {"name": "bass"}
    ]
}`)

for _, instrument := range v.Get("instruments").Array() {
    fmt.Println(instrument.Get("name").String())
}
```

You can also collect internal properties of an array using `#` gjson wildcard.

```go
for _, instrumentName := range v.Get("instruments.#.name").Array() {
    fmt.Println(instrumentName.String())
}
```

### Parsing

The library uses Golang built-in json marshallers. In case you want to use a custom marshaller, you can use
`jsonnav.From[T]()` or `jsonnav.FromAny()` by providing the actual value.

```go
v, err := jsonnav.From(map[string]any{"name": "John", "age": 30})
v.Get("name").String() // "John"
```

## License

jsonnav is distributed under [MIT License](https://opensource.org/license/MIT).

[gjson]: https://github.com/tidwall/gjson/blob/master/SYNTAX.md