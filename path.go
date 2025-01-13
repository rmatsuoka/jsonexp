package jsonexp

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

type Path []Node

func NewPath(nodes ...any) Path {
	p := Path{}
	for _, n := range nodes {
		switch n := n.(type) {
		case string:
			p = append(p, ObjectKey(n))
		case int:
			p = append(p, ArrayIndex(n))
		default:
			panic("not string or int")
		}
	}
	return p
}

func (p Path) Clone() Path {
	n := make(Path, len(p))
	copy(n, p)
	return n
}

func (p Path) CloneAppend(node Node) Path {
	return append(p.Clone(), node)
}

func (p Path) String() string {
	b := strings.Builder{}
	for _, n := range p {
		switch n := n.(type) {
		case ObjectKey:
			b.WriteString("." + string(n))
		case ArrayIndex:
			fmt.Fprintf(&b, "[%d]", int(n))
		}
	}
	return b.String()
}

func (p Path) Equal(q Path) bool {
	return slices.EqualFunc(p, q, func(x, y Node) bool {
		return x == y
	})
}

func (p Path) Compare(q Path) int {
	return slices.CompareFunc(p, q, func(x, y Node) int {
		// define
		// object > number
		switch x := x.(type) {
		case ObjectKey:
			y, ok := y.(ObjectKey)
			if !ok {
				return -1
			}
			return cmp.Compare(x, y)
		case ArrayIndex:
			y, ok := y.(ArrayIndex)
			if !ok {
				return 1
			}
			return cmp.Compare(x, y)
		default:
			panic("unreachable")
		}
	})
}

func (p Path) IsAncestorOf(child Path) bool {
	if len(p) >= len(child) {
		return false
	}
	for i := range p {
		if p[i] != child[i] {
			return false
		}
	}
	return true
}

type Node interface {
	Node()
}

type ObjectKey string

func (ObjectKey) Node() {}

type ArrayIndex int

func (ArrayIndex) Node() {}

func (p Path) Query(value Value) (Value, error) {
	for _, n := range p {
		switch n := n.(type) {
		case ObjectKey:
			obj, ok := value.(Object)
			if !ok {
				return nil, fmt.Errorf("is not object")
			}
			value, ok = obj[string(n)]
			if !ok {
				return nil, fmt.Errorf("not found")
			}
		case ArrayIndex:
			arr, ok := value.(Array)
			if !ok {
				return nil, fmt.Errorf("is not array")
			}
			if n < 0 && len(arr) <= int(n) {
				return nil, fmt.Errorf("index is out of range")
			}
			value = arr[n]
		default:
			return nil, fmt.Errorf("wrong path")
		}
	}
	return value, nil
}
