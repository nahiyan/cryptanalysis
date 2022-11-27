package utils

import (
	"math"
	"math/rand"

	"github.com/samber/lo"
)

func RandomCubes(cubesCount, selectionSize int) []int {
	cubes := lo.Map(rand.Perm(cubesCount), func(index, _ int) int {
		return index + 1
	})

	randomCubeSelectionCount := int(math.Min(float64(cubesCount), float64(selectionSize)))

	return cubes[:randomCubeSelectionCount]
}
