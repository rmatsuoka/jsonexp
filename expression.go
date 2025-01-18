package jsonexp

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/rmatsuoka/jsonexp/internal/diff"
)

type Expression struct {
	exp valueExp
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
		exp: value,
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
		return numberExp(raw), nil
	case bool:
		return booleanExp(raw), nil
	case nil:
		return nullExp{}, nil
	default:
		panic("unexpected type")
	}
}

func (e *Expression) Diff(jsontext string) (*Diff, error) {
	var v Value
	err := json.Unmarshal([]byte(jsontext), &v)
	if err != nil {
		return nil, fmt.Errorf("jsonexp: %w", err)
	}

	return &Diff{
		exp:   e.exp,
		value: v,
		diffs: diffValue(e.exp, v, Path{}),
	}, nil
}

func (e *Expression) DiffValue(v Value) *Diff {
	return &Diff{
		exp:   e.exp,
		value: v,
		diffs: diffValue(e.exp, v, Path{}),
	}
}

func diffValue(exp valueExp, value Value, parent Path) (diffs []DiffLine) {
	switch exp := exp.(type) {
	case objectExp:
		obj, ok := value.(Object)
		if !ok {
			diffs = append(diffs, DiffLine{
				At:   parent,
				Type: OpSubStitution,
			})
			return diffs
		}
		diffs = append(diffs, diffObject(exp, obj, parent)...)
	case arrayExp:
		arr, ok := value.(Array)
		if !ok {
			diffs = append(diffs, DiffLine{
				At:   parent,
				Type: OpSubStitution,
			})
			return diffs
		}
		diffs = append(diffs, diffArray(exp, arr, parent)...)
	default:
		if !exp.matchValue(value) {
			diffs = append(diffs, DiffLine{
				At:   parent,
				Type: OpSubStitution,
			})
		}
	}
	return diffs
}

func diffObject(exp objectExp, obj Object, parent Path) (diffs []DiffLine) {
	restKeys := collectKey(maps.Keys(exp), true)
	for k := range obj {
		at := parent.CloneAppend(ObjectKey(k))
		expv, ok := exp.get(k)
		if !ok {
			diffs = append(diffs, DiffLine{
				At:   at,
				Type: OpInsertion,
			})
			continue
		}
		diffs = append(diffs, diffValue(expv, obj[k], at)...)
		delete(restKeys, k)
	}

	for k := range restKeys {
		if k == "..." {
			continue
		}
		diffs = append(diffs, DiffLine{
			At:   parent.CloneAppend(ObjectKey(k)),
			Type: OpDeletion,
		})
	}
	return diffs
}

func diffArray(exp arrayExp, arr Array, parent Path) (diffs []DiffLine) {
	ds := diff.Slice(len(exp), len(arr), func(ix, iy int) bool {
		return exp[ix].matchValue(arr[iy])
	})
	for _, d := range ds {
		if d.Op == diff.OpSubStitution {
			diffs = append(diffs, diffValue(exp[d.Xi], arr[d.Yi], parent.CloneAppend(ArrayIndex(d.Xi)))...)
			continue
		}
		diffs = append(diffs, DiffLine{
			At:   parent.CloneAppend(ArrayIndex(d.Yi)),
			Type: fromDiffOp(d.Op),
		})
	}
	return diffs
}
