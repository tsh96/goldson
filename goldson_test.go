package goldson

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type TestMock struct {
	testing.TB
	Errors []error
}

func (t *TestMock) Error(args ...any) {
	t.Errors = append(t.Errors, errors.New(fmt.Sprint(args...)))
}

func (t *TestMock) Errorf(format string, args ...any) {
	t.Errors = append(t.Errors, fmt.Errorf(format, args...))
}

func Test_Walk(t *testing.T) {
	t.Run("the order of array does not matched and have no option", func(t *testing.T) {
		golden := map[string]any{}
		output := map[string]any{}

		if err := json.Unmarshal([]byte(goldenJson), &golden); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(actualJson), &output); err != nil {
			t.Fatal(err)
		}

		tm := &TestMock{}
		Walk(tm, []string{}, golden, output)
		if len(tm.Errors) == 0 {
			t.Error("error should occurred")
		}
	})

	t.Run("the order of array does not matched and a sort option is passed", func(t *testing.T) {
		golden := map[string]any{}
		output := map[string]any{}

		if err := json.Unmarshal([]byte(goldenJson), &golden); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(actualJson), &output); err != nil {
			t.Fatal(err)
		}

		tm := &TestMock{}
		Walk(tm, []string{}, golden, output,
			Sort("friends.*.nets", func(a, b any) bool {
				return strings.Compare(a.(string), b.(string)) < 0
			}),
		)

		if len(tm.Errors) != 0 {
			t.Error("error should not occurred")
		}
	})
}

const goldenJson = `
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`

const actualJson = `
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["tw", "fb"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`

func TestSimple(t *testing.T) {
	golden := []byte(`{"name": {"first": "Tom", "last": "Anderson"} }`)
	actual := []byte(`{"name": {"first": "Tom", "last": "Anderson"} }`)

	TestFromBytes(t, golden, actual) //ok
}

func TestWithOptions(t *testing.T) {
	golden := []byte(`{"colors": ["green", "red", "blue"] }`)
	actual := []byte(`{"colors": ["red", "green", "blue"] }`)

	TestFromBytes(t, golden, actual, Sort("colors", func(a, b any) bool {
		return a.(string) < b.(string)
	})) //ok
}

func TestFile(t *testing.T) {
	actual := []byte(`{"colors": ["red", "green", "blue"] }`)

	TestFromFile(t, "golden.json", actual, Sort("colors", func(a, b any) bool {
		return a.(string) < b.(string)
	})) // ok, auto create golden.json if not exists
}
