package jsonexp

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
)

func DiffText(w io.Writer, ds []DiffLine, exp valueExp, val Value) {

}

type diffTexter struct {
	w  io.Writer
	ds []DiffLine
	di int
}

func (t *diffTexter) cur() DiffLine {
	return t.ds[t.di]
}

func (t *diffTexter) next() {
	t.di++
}
func (t *diffTexter) rest() bool {
	return t.di < len(t.ds)
}

// diffTextValue must be called only if `at.isParent(t.cur().At)`
func (t *diffTexter) diffTextValue(at Path, exp valueExp, val Value, prefix string) {
	switch exp := exp.(type) {
	case objectExp:
		t.diffTextObject(at, exp, val.(Object), prefix)
	case arrayExp:
		t.diffTextArray(at, exp, val.(Array), prefix)
	default:
		panic("unreachable")
	}
}

// diffTextValue must be called if `at.isParent(t.cur().At)`
func (t *diffTexter) diffTextObject(at Path, exp objectExp, obj Object, prefix string) {
	keys := slices.Sorted(maps.Keys(obj))

	fmt.Fprintf(t.w, "{\n")
	for _, k := range keys {
		keyAt := at.CloneAppend(ObjectKey(k))
		for ; t.rest() && t.cur().At.Compare(keyAt) < 0; t.next() {
			key := string(t.cur().At[len(at)].(ObjectKey))
			fmt.Fprintf(t.w, "cx - %s  %s: %J\n", prefix, key, jsonFormatter{exp[key]})
		}
		// log.Printf("keyAt = %s, t.ds[t.di].At = %s, keyAt.Equal(..) = %t, keyAt.isParent(..) = %t", keyAt, t.ds[t.di].At, keyAt.Equal(t.ds[t.di].At), keyAt.isParentOf(t.ds[t.di].At))
		if t.rest() && keyAt.Equal(t.cur().At) {
			switch t.ds[t.di].Type {
			case OpSubStitution:
				e, _ := exp.get(k)
				fmt.Fprintf(t.w, "eq - %s  %s: %J\n", prefix, k, jsonFormatter{e})
				fmt.Fprintf(t.w, "eq + %s  %s: %J\n", prefix, k, jsonFormatter{obj[k]})
			case OpInsertion:
				fmt.Fprintf(t.w, "ei + %s  %s: %J\n", prefix, k, jsonFormatter{obj[k]})
			default:
				panic("unreacable")
			}
			t.next()
		} else if t.rest() && keyAt.IsAncestorOf(t.cur().At) {
			fmt.Fprintf(t.w, "pa   %s  %s:", prefix, k)
			t.diffTextValue(keyAt, exp[k], obj[k], prefix+"  ")
			t.rest()
		} else {
			fmt.Fprintf(t.w, "no   %s  %s: %J\n", prefix, k, jsonFormatter{obj[k]})
		}
	}

	for ; t.rest() && at.IsAncestorOf(t.cur().At); t.next() {
		key := string(t.cur().At[len(at)].(ObjectKey))
		switch t.cur().Type {
		case OpDeletion:
			fmt.Fprintf(t.w, "rd - %s  %s: %J\n", prefix, key, jsonFormatter{exp[key]})
		case OpInsertion:
			fmt.Fprintf(t.w, "ri + %s  %s: %J\n", prefix, key, jsonFormatter{exp[key]})
		}
	}
	fmt.Fprintf(t.w, "     %s}\n", prefix)
}

func (t *diffTexter) diffTextArray(at Path, exp arrayExp, arr Array, prefix string) {
	fmt.Fprintf(t.w, "[\n")
	for i := range arr {
		iAt := at.CloneAppend(ArrayIndex(i))
		for ; t.rest() && t.cur().At.Compare(iAt) < 0; t.next() {
			index := int(t.cur().At[len(at)].(ArrayIndex))
			fmt.Fprintf(t.w, "cx - %s  %J\n", prefix, jsonFormatter{exp[index]})
		}
		if t.rest() && iAt.Equal(t.cur().At) {
			switch t.ds[t.di].Type {
			case OpSubStitution:
				fmt.Fprintf(t.w, "eq - %s  %J\n", prefix, jsonFormatter{exp[i]})
				fmt.Fprintf(t.w, "eq + %s  %J\n", prefix, jsonFormatter{arr[i]})
			case OpInsertion:
				fmt.Fprintf(t.w, "ei + %s  %J\n", prefix, jsonFormatter{arr[i]})
			default:
				panic("unreacable")
			}
			t.next()
		} else if t.rest() && iAt.IsAncestorOf(t.cur().At) {
			fmt.Fprintf(t.w, "pa   %s  ", prefix)
			t.diffTextValue(iAt, exp[i], arr[i], prefix+"  ")
			t.rest()
		} else {
			fmt.Fprintf(t.w, "no   %s  %J\n", prefix, jsonFormatter{arr[i]})
		}
	}
	for ; t.rest() && at.IsAncestorOf(t.cur().At); t.next() {
		index := int(t.cur().At[len(at)].(ArrayIndex))
		switch t.cur().Type {
		case OpDeletion:
			fmt.Fprintf(t.w, "rd - %s  %J\n", prefix, jsonFormatter{exp[index]})
		case OpInsertion:
			fmt.Fprintf(t.w, "ri + %s  %J\n", prefix, jsonFormatter{exp[index]})
		}
	}
	fmt.Fprintf(t.w, "     %s]\n", prefix)
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

var _ fmt.Formatter = jsonFormatter{}
