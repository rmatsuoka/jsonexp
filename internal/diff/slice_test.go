package diff

import (
	"testing"
)

func Test_diffArray(t *testing.T) {
	x := []int{0, 9, 2, 3, 4}
	y := []int{9, 2}
	diffs := Slice(len(x), len(y), func(i, j int) bool { return x[i] == y[j] })
	t.Logf("%+v", diffs)
	// b := strings.Builder{}
	// diffs.Text(&b, len(x), len(y),
	// 	func(i int) string { return strconv.Itoa(x[i]) },
	// 	func(i int) string { return strconv.Itoa(y[i]) })
	// t.Log(b.String())
}
