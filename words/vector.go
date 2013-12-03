package words

import (
	"fmt"
	"math"
)

// cosine similarity between two vectors.
func cosine(v1, v2 []float64) (float64, error) {
	if len(v1) != len(v2) {
		return 0, fmt.Errorf("vectors do not have the same dimensions")
	}

	var dot, sum1, sum2 float64
	for i := range v1 {
		dot += v1[i] * v2[i]
		sum1 += v1[i] * v1[i]
		sum2 += v2[i] * v2[i]
	}

	return dot / math.Sqrt(sum1*sum2), nil
}
