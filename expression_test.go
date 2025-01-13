package jsonexp

import "testing"

func TestParse(t *testing.T) {
	e, err := Parse([]byte(`
{
  "hello": {
	"num": 1,
	"arr": [
	  {"x": "y"},
	  {"y": 1},
	  {"z": "w"}
	],
	"nil": 3
  }
}	
	`))
	if err != nil {
		t.Fatal(err)
	}

	val := mustUnmarshal(`
{
  "hello": {
    "num": 3,
	"arr": [
	  {"x": "y"},
	  {"y": "x"},
	  {"z": "w"},
	  {"x": "x"}
	],
	"nil": null
  }
}`)

	diffs := listDiff(e.value, val, path{})
	t.Logf("%+v", diffs)
}
