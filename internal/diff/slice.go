package diff

import (
	"cmp"
	"fmt"
	"io"
	"slices"
)

type SliceDiff struct {
	Xi, Yi int
	Op     Operation
}

type SliceDiffs []SliceDiff

func (d SliceDiffs) Text(w io.Writer, nx, ny int, xstr, ystr func(int) string) {
	slices.SortFunc(d, func(p, q SliceDiff) int {
		return cmp.Or(cmp.Compare(p.Xi, q.Xi), cmp.Compare(p.Yi, q.Yi))
	})

	i := 0
	for ; d[i].Xi == -1; i++ {
		fmt.Fprintf(w, "+ %s\n", ystr(d[i].Yi))
	}

	for xi := range nx {
		first := true
		for ; i < len(d) && d[i].Xi == xi; i++ {
			switch d[i].Op {
			case OpDeletion:
				fmt.Fprintf(w, "- %s\n", xstr(xi))
			case OpInsertion:
				if first {
					fmt.Fprintf(w, "  %s\n", xstr(xi))
				}
				fmt.Fprintf(w, "+ %s\n", ystr(d[i].Yi))
			case OpSubStitution:
				fmt.Fprintf(w, "- %s\n+ %s\n", xstr(xi), ystr(d[i].Yi))
			}
			first = false
		}
		if first {
			fmt.Fprintf(w, "  %s\n", xstr(xi))
		}
	}
}

func Slice(nx, ny int, equal func(ix, iy int) bool) SliceDiffs {
	var dist [][]int = make([][]int, nx+1)
	for i := range dist {
		dist[i] = make([]int, ny+1)
	}

	for i := range nx + 1 {
		dist[i][0] = i
	}

	for j := range ny + 1 {
		dist[0][j] = j
	}

	for i := 1; i <= nx; i++ {
		for j := 1; j <= ny; j++ {
			z := 1
			if equal(i-1, j-1) {
				z = 0
			}
			dist[i][j] = min(dist[i-1][j]+1, dist[i][j-1]+1, dist[i-1][j-1]+z)
		}
	}

	var diffs []SliceDiff
	i := nx
	j := ny
	cost := dist[i][j]
	for cost > 0 {
		switch {
		case i > 0 && cost == dist[i-1][j]+1:
			diffs = append(diffs, SliceDiff{
				Xi: i - 1,
				Yi: j - 1,
				Op: OpDeletion,
			})
			i = i - 1
			cost--
		case j > 0 && cost == dist[i][j-1]+1:
			diffs = append(diffs, SliceDiff{
				Xi: i - 1,
				Yi: j - 1,
				Op: OpInsertion,
			})
			j = j - 1
			cost--
		case i > 0 && j > 0 && cost == dist[i-1][j-1]+1:
			diffs = append(diffs, SliceDiff{
				Xi: i - 1,
				Yi: j - 1,
				Op: OpSubStitution,
			})
			i = i - 1
			j = j - 1
			cost--
		case i > 0 && j > 0 && cost == dist[i-1][j-1]:
			i = i - 1
			j = j - 1
		default:
			panic("unreachable")
		}
	}
	slices.Reverse(diffs)
	return diffs
}
