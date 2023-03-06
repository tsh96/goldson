package goldson

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Walker will skip that iteration if Option return true
type Option func(t testing.TB, path []string, golden, actual any, options ...Option) (skip bool)

// A PathPattern is intended to be easily expressed as a series of
// components separated by a `.` character.
//
// A key may contain the special wildcard characters `*` and `?`. The `*`
// will match on any zero+ characters, and `?` matches on any one character.
//
// `/` is escape character
func Ignore(pathPattern string) Option {
	return func(t testing.TB, path []string, golden, actual any, options ...Option) (skip bool) {
		t.Helper()
		return MatchPath(pathPattern, path)
	}
}

// A PathPattern is intended to be easily expressed as a series of
// components separated by a `.` character.
//
// A key may contain the special wildcard characters `*` and `?`. The `*`
// will match on any zero+ characters, and `?` matches on any one character.
//
// `/` is escape character
func Sort(pathPattern string, less func(a, b any) bool) Option {
	return func(t testing.TB, path []string, golden, actual any, options ...Option) (skip bool) {
		t.Helper()
		if !MatchPath(pathPattern, path) {
			return false
		}

		goldenArr, ok := golden.([]any)
		if !ok {
			t.Error("["+strings.Join(path, ".")+"]", "the type of golden json is not array")
			return true
		}
		actualArr, ok := actual.([]any)
		if !ok {
			t.Error("["+strings.Join(path, ".")+"]", "the type of actual json is not array")
			return true
		}

		if len(goldenArr) != len(actualArr) {
			t.Error("["+strings.Join(path, ".")+"]", "length of array does not match, golden:", len(goldenArr), "actual:", len(actualArr))
			return true
		}

		newGoldenArr := append([]any{}, goldenArr...)
		newActualArr := append([]any{}, actualArr...)

		sort.Slice(newGoldenArr, func(i, j int) bool {
			return less(newGoldenArr[i], newGoldenArr[j])
		})

		sort.Slice(newActualArr, func(i, j int) bool {
			return less(newActualArr[i], newActualArr[j])
		})

		Walk(t, path, newGoldenArr, newActualArr, options...)

		return true
	}
}

func MatchPath(pattern string, path []string) bool {
	parsedPattern := []string{}
	lastBreak := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '\\' {
			i++
		} else if pattern[i] == '.' {
			parsedPattern = append(parsedPattern, pattern[lastBreak:i])
			lastBreak = i + 1
		}
	}
	if len(pattern) > lastBreak {
		parsedPattern = append(parsedPattern, pattern[lastBreak:])
	}

	if len(path) != len(parsedPattern) {
		return false
	}

	for i, pattern := range parsedPattern {
		pattern = regexp.MustCompile(`([^\\])\*`).ReplaceAllString(pattern, `$1.*`)
		pattern = regexp.MustCompile(`([^\\])\?`).ReplaceAllString(pattern, `$1.`)
		match := regexp.MustCompile("^" + pattern + "$").MatchString(path[i])

		if !match {
			return false
		}
	}

	return true
}
