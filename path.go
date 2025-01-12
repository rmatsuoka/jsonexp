package jsonexp

import (
	"fmt"
	"strconv"
	"strings"
)

type path []key

func newPath(keys ...any) []key {
	p := path{}
	for _, k := range keys {
		switch k := k.(type) {
		case string:
			p = append(p, objectKey(k))
		case int:
			p = append(p, arrayIndex(k))
		default:
			panic("not string or int")
		}
	}
	return p
}

func (p path) Clone() path {
	n := make(path, len(p))
	copy(n, p)
	return n
}

func (p path) CloneAppend(key key) path {
	return append(p.Clone(), key)
}

func (p path) String() string {
	b := strings.Builder{}
	for _, k := range p {
		b.WriteString(k.String())
	}
	return b.String()
}

type key interface {
	key()
	String() string
}

type objectKey string

func (objectKey) key() {}
func (k objectKey) String() string {
	return "." + string(k)
}

type arrayIndex int

func (arrayIndex) key() {}

func (k arrayIndex) String() string {
	return "[" + strconv.Itoa(int(k)) + "]"
}

func query(value any, p path) (any, error) {
	for _, k := range p {
		switch k := k.(type) {
		case objectKey:
			obj, ok := value.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("is not object")
			}
			value, ok = obj[string(k)]
			if !ok {
				return nil, fmt.Errorf("not found")
			}
		case arrayIndex:
			arr, ok := value.([]any)
			if !ok {
				return nil, fmt.Errorf("is not array")
			}
			if k < 0 && len(arr) <= int(k) {
				return nil, fmt.Errorf("index is out of range")
			}
			value = arr[k]
		default:
			return nil, fmt.Errorf("wrong path")
		}
	}
	return value, nil
}
