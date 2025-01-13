package jsonexp

import "slices"

type Diff struct {
	At   path
	Type Operation
}

type Operation int

const (
	OpInsertion Operation = iota + 1
	OpDeletion
	OpSubStitution
)

func (o Operation) String() string {
	switch o {
	case OpInsertion:
		return "insertion"
	case OpDeletion:
		return "deletion"
	case OpSubStitution:
		return "substitution"
	default:
		panic("unreachable")
	}
}

func SortDiff(ds []Diff) {
	slices.SortFunc(ds, func(x, y Diff) int {
		return x.At.Compare(y.At)
	})
}

func SearchDiff(ds []Diff, at path) (int, bool) {
	return slices.BinarySearchFunc(ds, at, func(x Diff, p path) int {
		return x.At.Compare(p)
	})
}
