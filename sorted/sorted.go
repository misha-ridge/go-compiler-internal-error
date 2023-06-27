package sorted

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Slice returns a sorted copy of the given slice
func Slice[T constraints.Ordered](s []T) []T {
	res := slices.Clone(s)
	slices.Sort(res)
	return res
}

// Keys returns a sorted slice of keys of the given map
func Keys[M ~map[K]V, K constraints.Ordered, V any](m M) []K {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}
