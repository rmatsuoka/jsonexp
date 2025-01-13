package jsonexp

import (
	"slices"

	"github.com/rmatsuoka/jsonexp/internal/diff"
)

type Diff struct {
	At   Path
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

func fromDiffOp(o diff.Operation) Operation {
	switch o {
	case diff.OpInsertion:
		return OpInsertion
	case diff.OpDeletion:
		return OpDeletion
	case diff.OpSubStitution:
		return OpSubStitution
	default:
		panic("unreachable")
	}
}

func SortDiff(ds []Diff) {
	slices.SortFunc(ds, func(x, y Diff) int {
		return x.At.Compare(y.At)
	})
}

func SearchDiff(ds []Diff, at Path) (int, bool) {
	return slices.BinarySearchFunc(ds, at, func(x Diff, p Path) int {
		return x.At.Compare(p)
	})
}
