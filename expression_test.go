package jsonexp

import (
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

	t.Log(e.DiffValue(val).Text())
}
