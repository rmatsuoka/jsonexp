package jsonexp

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	e, err := Parse([]byte(`
{
  "hello": {
    "a": "((number))",
	"num": 1,
	"arr": [
	  {
	    "obj": {
		  "boolean": false,
	      "a": 1,
		  "b": 2,
		  "c": 3
	    }
      },
	  {"x": "y"},
	  {"y": 1},
	  {"z": "w"}
	],
	"nil": null,
	"aobj": {
	  "inner": "inner",
	  "num": 3
	},
	"x": "y",
	"...": "..."
  }
}	
	`))
	if err != nil {
		t.Fatal(err)
	}

	val := mustUnmarshal(`
{
  "hello": {
    "a": 3,
    "num": 5,
	"arr": 	[
	  {
	    "obj": {
	      "a": 1,
		  "b": 2,
		  "c": 3
	    }
      },
	  {"x": "y"},
	  {"y": 1},
	  {"z": "w"}
	],
	"nil": null,
	"aobj": {
	  "inner": "inner",
	  "num": 3
    },
	"x": "y",
	"z": "..."
  }
}`)

	diffs := diffValue(e.value, val, Path{})
	t.Logf("%+v", diffs)

	SortDiff(diffs)
	b := strings.Builder{}
	dt := diffTexter{
		w:  &b,
		ds: diffs,
		di: 0,
	}
	dt.diffTextValue(Path{}, e.value, val, "")
	t.Log(b.String())
}
