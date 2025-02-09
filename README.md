# JSONNAV

jsonnav is a [Go package](https://pkg.go.dev/github.com/jorgebay/jsonnav#section-documentation) for accessing,
navigating and manipulating values from an untyped json document.

## Features

- Retrieving values from deeply nested json documents safely.
- No error checking needed when properties does not exist.
- Iterate over arrays and objects.
- Supports GJSON syntax for navigating the json document.
- Set or delete values in place.

It supports retrieving deeply nested values safely, as well as iterating over arrays and objects.

## Installing

```shell
go get github.com/jorgebay/jsonnav
```

[![Build](https://github.com/jorgebay/jsonnav/actions/workflows/test.yml/badge.svg)](https://github.com/jorgebay/jsonnav/actions/workflows/test.yml)

## Usage

```go
v, err := jsonnav.Unmarshal(`{"name":{"first":"Jimi","last":"Hendrix"},"age":27}`)
v.Get("name").Get("first").String() // "Jimi"
v.Get("name.last").String() // "Hendrix"
v.Get("age").Float() // 27.0
v.Get("path").Get("does").Get("not").Get("exist").Exists() // false
v.Get("path.does.not.exist").Exists() // false
```

It uses [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md) for navigating the json document.

The library uses Golang built-in json marshallers. In case you want to use a custom marshaller, you can use
`jsonnav.From[T]()` or `jsonnav.FromAny()` by providing the actual value.

```go
v, err := jsonnav.From(map[string]any{"name": "John", "age": 30})
v.IsObject() // true
v.Get("name").String() // "John"
```

### Accessing values that may not exist

It's safe to access values that may not exist. The library will return a scalar with `nil` underlying value.

```go
v, err := jsonnav.Unmarshal(`{
    "name":{"first":"Jimi","last":"Hendrix"},
    "instruments":[{"name":"guitar"}]
}`)
v.Get("birth").Get("date").Exists() // false
v.Get("instruments").Array().At(0).Get("name").String() // "guitar"
v.Get("instruments").Array().At(1).Get("name").String() // ""
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

## License

jsonnav is distributed under [MIT License](https://opensource.org/license/MIT).
