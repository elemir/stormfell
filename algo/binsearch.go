package algo

import (
	"golang.org/x/exp/constraints"
)

// BinSearch returns a number of interval where val is placed in the ordered slice xs.
func BinSearch[T constraints.Ordered](xs []T, val T) int {
	if len(xs) == 0 {
		return 0
	}

	return binSearch(xs, val, 0)
}

func binSearch[T constraints.Ordered](xs []T, val T, idx int) int {
	if len(xs) == 1 {
		if val > xs[0] {
			return idx + 1
		}

		return idx
	}

	half := len(xs) / 2
	if val > xs[half] {
		return binSearch(xs[half:], val, half+idx)
	}

	return binSearch(xs[:half], val, idx)
}
