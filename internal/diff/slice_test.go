package diff

import (
	"strconv"
	"strings"
	"testing"
)

func TestSlice(t *testing.T) {
	x := []int{0, 9, 2, 3, 4, 5}
	y := []int{10, 0, 8, 2, 2, 3, 4}
	diffs := Slice(len(x), len(y), func(i, j int) bool { return x[i] == y[j] })
	t.Logf("%+v", diffs)
	b := strings.Builder{}
	diffs.Text(&b, len(x), len(y),
		func(i int) string { return strconv.Itoa(x[i]) },
		func(i int) string { return strconv.Itoa(y[i]) })
	t.Log(b.String())
}
