package jsonexp

import "testing"

func Test_textExp(t *testing.T) {
	tests := []struct {
		exp     string
		match   []any
		unmatch []any
	}{
		{"", []any{""}, []any{nil, 0., []any{}, map[string]any{}}},
		{"x", []any{"x"}, []any{nil, "z", 0.}},
		{"==", []any{""}, []any{nil, "=="}},
		{"==x", []any{"x"}, []any{nil, "==x"}},
		{"==  3", []any{"  3"}, []any{nil, 3., "3"}},
		{"==((number))", []any{"((number))"}, []any{5., 8.}},
		{"((number))", []any{0., 0.1, -1., 4.}, []any{nil, "((number))", "number"}},
		{"((array))", []any{[]any{}, []any{3., 1.}}, []any{nil, "array", map[string]any{"x": "y"}}},
		{"((any))", []any{nil, false, "x", map[string]any{"x": 2}, []any{}}, []any{}},
	}

	for _, test := range tests {
		e, err := parseTextExp(test.exp)
		if err != nil {
			t.Errorf("parseTextExp(%s) returns unexpcted error: %v", test.exp, err)
		}
		for _, match := range test.match {
			if !e.Match(match) {
				t.Errorf("%s does not match %s, but should", test.exp, match)
			}
		}
		for _, unmatch := range test.unmatch {
			if e.Match(unmatch) {
				t.Errorf("%s does match %s, but should not", test.exp, unmatch)
			}
		}
	}
}
