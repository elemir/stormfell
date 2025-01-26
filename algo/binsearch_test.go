package algo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/elemir/stormfell/algo"
)

func TestBinSearch(t *testing.T) {
	xs := []float64{-10, 4.5, 9}

	require.Equal(t, 0, algo.BinSearch(xs, -12), "smallest value")
	require.Equal(t, 2, algo.BinSearch(xs, 5), "somewhere in middle")
	require.Equal(t, 3, algo.BinSearch(xs, 10), "biggest value")

	ys := []float64{-5}

	require.Equal(t, 0, algo.BinSearch(ys, -10), "small value")
	require.Equal(t, 1, algo.BinSearch(ys, 10), "big value")

	require.Equal(t, 0, algo.BinSearch(nil, 0), "work on empty slices")
}
