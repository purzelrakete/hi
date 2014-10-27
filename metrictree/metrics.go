package metrictree

import (
	"fmt"
	"math"
)

// Distance between two objects
type Distance func(o1, o2 obj) (float64, error)

// Euclidean is the L2 norm distance metric
func Euclidean(o1, o2 obj) (float64, error) {
	if l1, l2 := len(o1), len(o2); l1 != l2 {
		return 0.0, fmt.Errorf("incompatible sizes: (%d, %d)", l1, l2)
	}

	var ret float64
	for i := 0; i < len(o1); i++ {
		ret += math.Pow(float64(o1[i])-float64(o2[i]), 2)
	}

	return math.Sqrt(ret), nil
}
