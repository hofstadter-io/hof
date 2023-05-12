package chat

import (
	"errors"
	"math"
)

// https://en.wikipedia.org/wiki/Cosine_similarity
// 1 is similare
// 0 is orthogonal
// -1 is opposite
func CosineSimilarity(a []float64, b []float64) (cosine float64, err error) {
	if len(a) != len(b) {
		return 0.0, errors.New("CosineSimilarity: vectors must be the same length.")
	}

	sumN := 0.0
	sumA := 0.0
	sumB := 0.0
	for k := 0; k < len(a); k++ {
		sumN += a[k] * b[k]
		sumA += math.Pow(a[k], 2)
		sumB += math.Pow(b[k], 2)
	}
	if sumA == 0 || sumB == 0 {
		return 0.0, errors.New("Vectors should not sum to zero")
	}
	return sumN / (math.Sqrt(sumA) * math.Sqrt(sumB)), nil
}

// finds the clostes in bs to a, returning the max similarity and position
func ClosestVector(a []float64, bs [][]float64) (max float64, pos int, err error) {
	max = -2.0
	pos = -1
	for i,b := range bs {
		cs, err := CosineSimilarity(a, b)
		if err != nil {
			return max, pos, err
		}
		if cs > max {
			max = cs
			pos = i
		}
	}

	return max, pos, nil
}
