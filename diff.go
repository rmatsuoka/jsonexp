package jsonexp

import (
	"io"
	"strings"
)

type Diff struct {
	exp   valueExp
	value Value
	diffs []DiffLine
}

func (d *Diff) Lines() []DiffLine { return d.diffs }

func (d *Diff) HasDiff() bool { return len(d.diffs) != 0 }

func (d *Diff) Text() string {
	b := strings.Builder{}
	d.WriteText(&b) // strings.Builder does not return non-nil error.
	return b.String()
}

func (d *Diff) WriteText(w io.Writer) error {
	return diffText(w, d.diffs, d.exp, d.value)
}
