package jsonexp

import (
	"encoding/json"
	"maps"
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
		for i, v := range raw {
			raw[i], err = toExpValue(v)
			if err != nil {
				return nil, err
			}
		}
		return raw, nil
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
	panic("x")
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

func listDiffArray(exp expObject) {
	panic("x")
}
