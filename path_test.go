package jsonexp

import (
	"encoding/json"
	"reflect"
	"testing"
)

func mustUnmarshal(js string) any {
	var value any
	err := json.Unmarshal([]byte(js), &value)
	if err != nil {
		panic(err)
	}
	return value
}

func Test_path(t *testing.T) {

	tests := []struct {
		p    Path
		want string
	}{
		{NewPath(), ""},
		{NewPath(3), "[3]"},
		{NewPath("a"), ".a"},
		{NewPath(1, "x", 5), "[1].x[5]"},
		{NewPath("p", "q", "r"), ".p.q.r"},
	}

	for _, test := range tests {
		if test.p.String() != test.want {
			t.Errorf("p.String() returns %s, want %s", test.p.String(), test.want)
		}
	}
}

func Test_query(t *testing.T) {
	t.Run("nil path", func(t *testing.T) {
		tests := []any{1, "z", nil, map[string]any{"x": 1}, []any{}}
		p := (Path)(nil)
		for _, value := range tests {
			got, err := p.query(value)
			if err != nil {
				t.Errorf("query(%+v, nil) returns unexpected non-nil error: %v", value, err)
			}
			if !reflect.DeepEqual(value, got) {
				t.Errorf("query(%v, nil)  returns %+v, want %+v", value, got, value)
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		value := mustUnmarshal(`
{
  "first_name": "John",
  "last_name": "Smith",
  "is_alive": true,
  "age": 27,
  "address": {
    "street_address": "21 2nd Street",
    "city": "New York",
    "state": "NY",
    "postal_code": "10021-3100"
  },
  "phone_numbers": [
    {
      "type": "home",
      "number": "212 555-1234"
    },
    {
      "type": "office",
      "number": "646 555-4567"
    }
  ],
  "children": [
    "Catherine",
    "Thomas",
    "Trevor"
  ],
  "spouse": null
}
		`)

		tests := []struct {
			path Path
			want any
			err  error
		}{
			{NewPath("age"), 27.0, nil},
			{NewPath("spouse"), nil, nil},
			{NewPath("address", "city"), "New York", nil},
			{NewPath("phone_numbers", 1, "type"), "office", nil},
		}

		for _, test := range tests {
			got, err := test.path.query(value)
			if err != test.err {
				t.Errorf("query(value, %v) returns unexpceted error: %v", test.path, err)
				continue
			}
			if err != nil {
				continue
			}
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("query(value, %v) returns %v, want %v", test.path, got, test.want)
			}
		}
	})
}
