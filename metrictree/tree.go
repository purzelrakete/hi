// A Metric Tree structure as described by Uhlmann [1].
//
// [1] - http://people.cs.missouri.edu/~uhlmannj/ImplementGH.pdf
package metrictree

import (
	"fmt"
	"math/rand"
)

// obj the type of object being indexed, fixed at []uint32 for now
type obj []uint64

type MetricTree struct {
	Object       *obj
	Inner, Outer float64
	Left, Right  *MetricTree
}

// NewMetricTree builds and returns a MetricTree
func NewMetricTree(set []obj, d Distance) *MetricTree {
	if len(set) == 0 {
		return &MetricTree{}
	}

	if len(set) == 1 {
		return &MetricTree{
			Object: &set[0],
		}
	}

	// select a random pivot
	nPivot := rand.Intn(len(set))
	pivot := set[nPivot]

	// bisect such that `right` contains the k/2 elements closest to `pivot`
	var left, right []obj
	var inner, outer float64
	for i, _ := range set {
		if i == nPivot {
			continue
		}

		d(pivot, set[i])
	}

	fmt.Printf("pivot: %v, left %v, right %v", pivot, left, right)

	return &MetricTree{
		Object: &pivot,
		Inner:  inner,
		Outer:  outer,
		Left:   NewMetricTree(left, d),
		Right:  NewMetricTree(left, d),
	}
}

// Find an object in the metric tree. returns obj, ok if found
func (t *MetricTree) Find(query obj) (obj, bool) {
	return obj{}, false
}
