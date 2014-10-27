package metrictree

import "testing"

func TestEuclidean(t *testing.T) {
	o1 := []uint64{1, 2, 3}
	o2 := []uint64{1, 2, 3}

	got, _ := Euclidean(o1, o2)
	if want := 0.0; want != got {
		t.Fatalf("want %v but got %v", want, got)
	}

	o1 = []uint64{0, 0}
	o2 = []uint64{0, 1}

	got, _ = Euclidean(o1, o2)
	if want := 1.0; want != got {
		t.Fatalf("want %v but got %v", want, got)
	}
}
