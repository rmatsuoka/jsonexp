package jsonexp

import (
	"slices"

	"github.com/rmatsuoka/jsonexp/internal/diff"
)

type DiffLine struct {
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

func SortDiffLines(ds []DiffLine) {
	slices.SortFunc(ds, func(x, y DiffLine) int {
		return x.At.Compare(y.At)
	})
}
