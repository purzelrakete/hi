package metrictree

import (
	"reflect"
	"testing"
)

func TestTreeFind(t *testing.T) {
	set := []obj{
		[]uint64{1, 2, 3},
		[]uint64{4, 5, 6},
		[]uint64{5, 6, 7},
	}

	tree, in := NewMetricTree(set, Euclidean), set[2]

	got, ok := tree.Find(in)
	if !ok {
		t.Fatalf("could not find expected obj")
	}

	if want := in; !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v but got %v", want, got)
	}
}
