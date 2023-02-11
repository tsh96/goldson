# Goldson

Goldson is a simple and extensible golden test utility for json data.

## Install

```bash
go get github.com/tsh96/goldson
```

## Example

```go
func TestSimple(t *testing.T) {
	golden := []byte(`{"name": {"first": "Tom", "last": "Anderson"} }`)
	actual := []byte(`{"name": {"first": "Tom", "last": "Anderson"} }`)

	TestFromBytes(t, golden, actual) //ok
}

func TestError(t *testing.T) {
	golden := []byte(`{"name": {"first": "Tom", "last": "Anderson"} }`)
	actual := []byte(`{"name": {"first": "To", "last": "Anders"} }`)

	TestFromBytes(t, golden, actual) // error
	// [name.first] golden: "Tom" actual: "To"
	// [name.last] golden: "Anderson" actual: "Anders"
}

func TestWithOptions(t *testing.T) {
	golden := []byte(`{"colors": ["green", "red", "blue"] }`)
	actual := []byte(`{"colors": ["red", "green", "blue"] }`)

	TestFromBytes(t, golden, actual,
    Sort("colors", func(a, b any) bool { return a.(string) < b.(string) }),
  ) //ok

	TestFromBytes(t, golden, actual) //error
	// [colors.0] golden: "green" actual: "red"
	// [colors.1] golden: "red" actual: "green"
}

func TestFile(t *testing.T) {
	actual := []byte(`{"colors": ["red", "green", "blue"] }`)

	TestFromFile(t, "golden.json", actual,
    Sort("colors", func(a, b any) bool { return a.(string) < b.(string) }),
  ) //ok
}
```

## Extensible

You are allowed to implement the type `Option` to create your own option.

```go
type Option func(t testing.TB, path []string, golden, actual any, options ...Option) (skip bool)
```
