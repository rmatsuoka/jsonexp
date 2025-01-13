package jsonexp

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
)

func DiffText(w io.Writer, ds []Diff, exp expValue, val value) {

}

func diffTextValue(w io.Writer, ds []Diff, exp expValue, val value, prefix string, at path) {
	// switch exp := exp.(type) {
	// case expObject:
	// }
}

func diffTextObject(w io.Writer, ds []Diff, exp expObject, obj object, prefix string, at path) {
	keys := slices.Sorted(maps.Keys(obj))

	fmt.Fprintf(w, "{")
	di := 0
	for _, k := range keys {
		atKey := at.CloneAppend(objectKey(k))
		for ; ds[di].At.Compare(atKey) < 0; di++ {
			key := ds[di].At[len(atKey):][0].(objectKey)
			fmt.Fprintf(w, "- %s %s: %v\n", prefix, key, exp[string(key)])
		}

		wrote := false
		if atKey.HasPrefix(ds[di].At) {
			if ds[di].At.Equal(atKey) {
				fmt.Fprintf(w, "+ %s  %s: %J\n", prefix, k, jsonFormatter{obj[k]})
			} else {
				diffTextValue(w, ds[di:], exp[k], obj[k], prefix+"  ", atKey)
			}
			di++
			wrote = true
			for ; atKey.HasPrefix(ds[di].At); di++ {
				// nothing
			}
		}
		if !wrote {
			fmt.Fprintf(w, "  %s  %s: %J\n", prefix, k, jsonFormatter{obj[k]})
		}
	}

	for ; at.HasPrefix(ds[di].At); di++ {
		key := ds[di].At[len(at):][0].(objectKey)
		fmt.Fprintf(w, "- %s %s: %J\n", prefix, key, jsonFormatter{exp})
	}
}

type jsonFormatter struct {
	value any
}

func (f jsonFormatter) Format(s fmt.State, verb rune) {
	if verb != 'J' {
		panic("verb must be %J")
	}
	buf, err := json.Marshal(f.value)
	if err != nil {
		panic(err)
	}
	s.Write(buf)
}
