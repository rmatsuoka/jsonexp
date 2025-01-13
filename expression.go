package jsonexp

import (
	"cmp"
	"encoding/json"
	"maps"
	"slices"

	"github.com/rmatsuoka/jsonexp/internal/diff"
)

type (
	// expValue is expObject | expArray | *textExp | expNumber | expBoolean | nil
	expValue any

	expObject map[string]expValue
	expArray  = []expValue
	// expString does not exist. string will be replaced by textExp
	expNumber  = float64
	expBoolean = bool
)

func (e expObject) get(key string) (expValue, bool) {
	v, ok := e[key]
	if !ok {
		v, ok = e["..."]
	}
	return v, ok
}

func (e expObject) equalLen(l int) bool {
	if _, ok := e["..."]; ok {
		return true
	}
	return len(e) == l
}

func (e expObject) sortedKeys() []string {
	return slices.SortedFunc(maps.Keys(e), func(s1, s2 string) int {
		if s1 == s2 {
			return 0
		}
		if s1 == "..." {
			return -1
		}
		if s2 == "..." {
			return 1
		}
		return cmp.Compare(s1, s2)
	})
}

func (e expObject) Match(obj object) bool {
	if !e.equalLen(len(obj)) {
		return false
	}
	restKey := collectKey(maps.Keys(e), true)
	delete(restKey, "...")

	for k, v := range obj {
		expv, ok := e.get(k)
		if !ok {
			return false
		}
		if !equalValue(expv, v) {
			return false
		}
		delete(restKey, k)
	}

	return len(restKey) == 0
}

type Expression struct {
	value expValue
}

func Parse(b []byte) (*Expression, error) {
	var raw any
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return nil, err
	}
	value, err := toExpValue(raw)
	if err != nil {
		return nil, err
	}
	return &Expression{
		value: value,
	}, err
}

func toExpValue(raw any) (expValue, error) {
	var err error
	switch raw := raw.(type) {
	case map[string]any:
		expo := make(expObject, len(raw))
		for k, v := range raw {
			expo[k], err = toExpValue(v)
			if err != nil {
				return nil, err
			}
		}
		return expo, nil
	case []any:
		exp := make(expArray, len(raw))
		for i, v := range raw {
			exp[i], err = toExpValue(v)
			if err != nil {
				return nil, err
			}
		}
		return exp, nil
	case string:
		return parseTextExp(raw)
	case float64:
		return raw, nil
	case boolean:
		return raw, nil
	case nil:
		return raw, nil
	default:
		panic("unexpected type")
	}
}

func listDiff(exp expValue, value value, at path) (diffs []path) {
	switch exp := exp.(type) {
	case expObject:
		obj, ok := value.(object)
		if !ok {
			diffs = append(diffs, at)
			return diffs
		}
		diffs = append(diffs, listDiffObject(exp, obj, at)...)
	case expArray:
		arr, ok := value.(array)
		if !ok {
			diffs = append(diffs, at)
			return diffs
		}
		diffs = append(diffs, listDiffArray(exp, arr, at)...)
	case *textExp:
		if !exp.Match(value) {
			diffs = append(diffs, at)
		}
	case expNumber:
		if exp != value {
			diffs = append(diffs, at)
		}
	case expBoolean:
		if exp != value {
			diffs = append(diffs, at)
		}
	case nil:
		if value != nil {
			diffs = append(diffs, at)
		}
	default:
		panic("unreachable")
	}
	return diffs
}

func listDiffObject(exp expObject, obj object, at path) (diffs []path) {
	restKeys := collectKey(maps.Keys(exp), true)
	for k := range obj {
		at := at.CloneAppend(objectKey(k))
		expv, ok := exp.get(k)
		if !ok {
			diffs = append(diffs, at)
			continue
		}
		diffs = append(diffs, listDiff(expv, obj[k], at)...)
		delete(restKeys, k)
	}

	for k := range restKeys {
		if k == "..." {
			continue
		}
		diffs = append(diffs, at.CloneAppend(objectKey(k)))
	}
	return diffs
}

func listDiffArray(exp expArray, arr array, at path) (diffs []path) {
	ds := diff.Slice(len(exp), len(arr), func(ix, iy int) bool {
		return equalValue(exp[ix], arr[iy])
	})
	for _, d := range ds {
		if d.Op == diff.OpSubStitution {
			diffs = append(diffs, listDiff(exp[d.Xi], arr[d.Yi], at.CloneAppend(arrayIndex(d.Xi)))...)
			continue
		}
		diffs = append(diffs, at.CloneAppend(arrayIndex(d.Xi)))
	}
	return diffs
}
