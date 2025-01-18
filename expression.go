package jsonexp

import (
	"cmp"
	"encoding/json"
	"maps"
	"slices"

	"github.com/rmatsuoka/jsonexp/internal/diff"
)

type (
	// valueExp is expObject | expArray | *textExp | expNumber | expBoolean | nil
	valueExp any

	objectExp map[string]valueExp
	arrayExp  = []valueExp
	// expString does not exist. string will be replaced by textExp
	numberExp  = float64
	booleanExp = bool
)

func (e objectExp) get(key string) (valueExp, bool) {
	v, ok := e[key]
	if !ok {
		v, ok = e["..."]
	}
	return v, ok
}

func (e objectExp) equalLen(l int) bool {
	if _, ok := e["..."]; ok {
		return true
	}
	return len(e) == l
}

func (e objectExp) sortedKeys() []string {
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

func (e objectExp) match(obj Object) bool {
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
	value valueExp
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

func toExpValue(raw any) (valueExp, error) {
	var err error
	switch raw := raw.(type) {
	case map[string]any:
		expo := make(objectExp, len(raw))
		for k, v := range raw {
			expo[k], err = toExpValue(v)
			if err != nil {
				return nil, err
			}
		}
		return expo, nil
	case []any:
		exp := make(arrayExp, len(raw))
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
	case bool:
		return raw, nil
	case nil:
		return raw, nil
	default:
		panic("unexpected type")
	}
}

func listDiff(exp valueExp, value Value, at Path) (diffs []Diff) {
	switch exp := exp.(type) {
	case objectExp:
		obj, ok := value.(Object)
		if !ok {
			diffs = append(diffs, Diff{
				At:   at,
				Type: OpSubStitution,
			})
			return diffs
		}
		diffs = append(diffs, listDiffObject(exp, obj, at)...)
	case arrayExp:
		arr, ok := value.(Array)
		if !ok {
			diffs = append(diffs, Diff{
				At:   at,
				Type: OpSubStitution,
			})
			return diffs
		}
		diffs = append(diffs, listDiffArray(exp, arr, at)...)
	case *textExp:
		if !exp.Match(value) {
			diffs = append(diffs, Diff{
				At:   at,
				Type: OpSubStitution,
			})
		}
	case numberExp:
		if exp != value {
			diffs = append(diffs, Diff{
				At:   at,
				Type: OpSubStitution,
			})
		}
	case booleanExp:
		if exp != value {
			diffs = append(diffs, Diff{
				At:   at,
				Type: OpSubStitution,
			})
		}
	case nil:
		if value != nil {
			diffs = append(diffs, Diff{
				At:   at,
				Type: OpSubStitution,
			})
		}
	default:
		panic("unreachable")
	}
	return diffs
}

func listDiffObject(exp objectExp, obj Object, at Path) (diffs []Diff) {
	restKeys := collectKey(maps.Keys(exp), true)
	for k := range obj {
		at := at.CloneAppend(ObjectKey(k))
		expv, ok := exp.get(k)
		if !ok {
			diffs = append(diffs, Diff{
				At:   at,
				Type: OpInsertion,
			})
			continue
		}
		diffs = append(diffs, listDiff(expv, obj[k], at)...)
		delete(restKeys, k)
	}

	for k := range restKeys {
		if k == "..." {
			continue
		}
		diffs = append(diffs, Diff{
			At:   at.CloneAppend(ObjectKey(k)),
			Type: OpDeletion,
		})
	}
	return diffs
}

func listDiffArray(exp arrayExp, arr Array, at Path) (diffs []Diff) {
	ds := diff.Slice(len(exp), len(arr), func(ix, iy int) bool {
		return equalValue(exp[ix], arr[iy])
	})
	for _, d := range ds {
		if d.Op == diff.OpSubStitution {
			diffs = append(diffs, listDiff(exp[d.Xi], arr[d.Yi], at.CloneAppend(ArrayIndex(d.Xi)))...)
			continue
		}
		diffs = append(diffs, Diff{
			At:   at.CloneAppend(ArrayIndex(d.Yi)),
			Type: fromDiffOp(d.Op),
		})
	}
	return diffs
}
