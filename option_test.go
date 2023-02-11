package goldson

import (
	"testing"
)

func Test_matchPath(t *testing.T) {
	if !MatchPath(`level1.level\.2.level3`, []string{"level1", "level.2", "level3"}) {
		t.Error("the result should be true")
	}
	if !MatchPath("level1.level2.level3", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be true")
	}
	if !MatchPath("level1.le*2.level3", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be true")
	}
	if !MatchPath("level1.le?el2.level?", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be true")
	}
	if !MatchPath("*.level2.level3", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be true")
	}
	if !MatchPath("level1.*.level3", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be true")
	}
	if !MatchPath("level1.level2.*", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be true")
	}

	if MatchPath("level1.level2.level3", []string{"level1", "level2"}) {
		t.Error("the result should be false")
	}
	if MatchPath("level1.level2", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be false")
	}
	if MatchPath("level1.lev*1.level3", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be false")
	}
	if MatchPath("level1.lev?2.level3", []string{"level1", "level2", "level3"}) {
		t.Error("the result should be false")
	}
}

func Test_Ignore(t *testing.T) {
	opt := Ignore("level1.lev?2.level3")
	if opt(t, []string{"level1", "level2", "level3"}, nil, nil) {
		t.Error("the result should be false")
	}

	if !opt(t, []string{"level1", "leve2", "level3"}, nil, nil) {
		t.Error("the result should be true")
	}
}

func Test_Sort(t *testing.T) {
	opt := Sort("level1.lev?l2.level3", func(a, b any) bool {
		aVal := a.(int)
		bVal := b.(int)

		return aVal < bVal
	})

	golden := []any{1, 2, 3}
	actual := []any{1, 3, 2}

	if !opt(t, []string{"level1", "level2", "level3"}, golden, actual) {
		t.Error("the result should be true")
	}

	if !opt(&testing.T{}, []string{"level1", "level2", "level3"}, []any{1, 2, 3}, []any{1, 3}) {
		t.Error("the result should be true")
	}
}
