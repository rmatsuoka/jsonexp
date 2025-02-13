package jsonexp

import (
	"cmp"
	"maps"
	"slices"
)

type valueExp interface {
	matchValue(Value) bool
}

type objectExp map[string]valueExp

var _ valueExp = objectExp{}

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

func (e objectExp) matchValue(value Value) bool {
	obj, ok := value.(Object)
	if !ok {
		return false
	}
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
		if !expv.matchValue(v) {
			return false
		}
		delete(restKey, k)
	}

	return len(restKey) == 0
}

type arrayExp []valueExp

var _ valueExp = arrayExp{}

func (e arrayExp) matchValue(value Value) bool {
	arr, ok := value.(Array)
	if !ok {
		return false
	}
	if len(e) != len(arr) {
		return false
	}
	for i := range e {
		if !e[i].matchValue(arr[i]) {
			return false
		}
	}
	return true
}

type numberExp float64

func (e numberExp) matchValue(v Value) bool {
	return float64(e) == v
}

type booleanExp bool

func (e booleanExp) matchValue(v Value) bool {
	return bool(e) == v
}

type nullExp struct{}

func (nullExp) matchValue(v Value) bool {
	return v == nil
}
