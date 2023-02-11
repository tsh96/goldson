package goldson

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestFromFile auto create the file and write the actual data into it if not exists.
// if file exists, it call function `TestFromBytes`. See also `TestFromBytes`
func TestFromFile(t testing.TB, filepath string, actual []byte, options ...Option) {
	goldenBytes, err := os.ReadFile(filepath)
	if os.IsNotExist(err) {
		os.WriteFile(filepath, actual, 0666)
		return
	} else if err != nil {
		t.Fatal(err)
	}

	TestFromBytes(t, goldenBytes, actual, options...)
}

// TestFromBytes compare the diff between golden data and actual data.
// options can change the behavior of the comparison, for example `Ignore`, `Sort`
func TestFromBytes(t testing.TB, golden, actual []byte, options ...Option) {
	var goldenJson any
	if err := json.Unmarshal(golden, &goldenJson); err != nil {
		t.Fatal(err)
	}

	var actualJson any
	if err := json.Unmarshal(actual, &actualJson); err != nil {
		t.Fatal(err)
	}

	Walk(t, []string{}, goldenJson, actualJson, options...)
}

func Walk(t testing.TB, path []string, golden, actual any, options ...Option) {
	pathStr := strings.Join(path, ".")

	for _, option := range options {
		if option(t, path, golden, actual) {
			return
		}
	}

	switch golden := golden.(type) {
	case bool, int, float64, string, nil:
		if golden != actual {
			t.Errorf("[%v] golden: %v actual: %v", pathStr, toJson(t, golden), toJson(t, actual))
		}

	case map[string]any:
		if actualVal, ok := actual.(map[string]any); ok {
			keys := map[string]struct{}{}

			for key := range golden {
				keys[key] = struct{}{}
			}

			for key := range actualVal {
				keys[key] = struct{}{}
			}

			for key := range keys {
				Walk(t, append(path, key), golden[key], actualVal[key], options...)
			}
		} else {
			t.Errorf("[%v] golden: %v actual: %v", pathStr, toJson(t, golden), toJson(t, actual))
		}

	case []any:
		if actualVal, ok := actual.([]any); ok {
			if len(actualVal) != len(golden) {
				t.Errorf("[%v] length of array does not match,  golden: %v actual: %v", pathStr, toJson(t, golden), toJson(t, actual))
			} else {
				for key, val := range golden {
					Walk(t, append(path, fmt.Sprint(key)), val, actualVal[key], options...)
				}
			}
		} else {
			t.Errorf("[%v] golden: %v actual: %v", pathStr, toJson(t, golden), toJson(t, actual))
		}

	default:
		t.Fatal("undefined type")
	}
}

func toJson(t testing.TB, value any) string {
	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal("fail to convert to json:", value)
	}
	return string(result)
}
