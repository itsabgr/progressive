package span

import "golang.org/x/exp/constraints"

func compare[T constraints.Ordered](t1, t2 T) int {
	if t1 < t2 {
		return -1
	}
	if t1 == t2 {
		return 0
	}
	return 1
}
