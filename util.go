package jsonexp

import "iter"

func collectKey[K comparable, V any](seq iter.Seq[K], value V) map[K]V {
	m := make(map[K]V)
	for k := range seq {
		m[k] = value
	}
	return m
}
