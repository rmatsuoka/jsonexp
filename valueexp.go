package jsonexp

import (
	"cmp"
	"maps"
	"slices"
)

type valueExp interface {
	valueExp()
}

type objectExp map[string]valueExp

func (objectExp) valueExp() {}

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

type arrayExp []valueExp

func (arrayExp) valueExp() {}

func (e arrayExp) match(arr Array) bool {
	if len(e) != len(arr) {
		return false
	}
	for i := range e {
		if !equalValue(e[i], arr[i]) {
			return false
		}
	}
	return true
}

type numberExp float64

func (numberExp) valueExp() {}

func (e numberExp) match(v any) bool {
	return float64(e) == v
}

type booleanExp bool

func (e booleanExp) match(v any) bool {
	return bool(e) == v
}

func (booleanExp) valueExp() {}
