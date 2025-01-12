package diff

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
