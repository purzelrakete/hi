package words

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	vectors := []string{
		"1 4",
		"dubstep 123 -0.001058 0.002683 0.000132 0.001072",
	}

	r := bufio.NewReader(strings.NewReader(strings.Join(vectors, "\n")))
	d, err := New(r)
	if err != nil {
		t.Fatalf("could not load Words: %s", err.Error())
	}

	expectedEntries := 1
	if actual := d.Len(); expectedEntries != actual {
		t.Fatalf("expected %d but got %d", expectedEntries, actual)
	}

	actual, ok := d.Vector("dubstep")
	if !ok {
		t.Fatalf("could not find dubstep in Words")
	}

	expectedVector := []float32{-0.001058, 0.002683, 0.000132, 0.001072}
	if !reflect.DeepEqual(expectedVector, actual) {
		t.Fatalf("expected %v but got %v", expectedVector, actual)
	}
}

func TestNearestNeighbours(t *testing.T) {
	vectors := []string{
		"4 2",
		"minimalhouse 10 1.0 0.0",
		"opera 11 -1.0 0.0",
		"house 12 1.0 0.1",
		"classical 13 -1.0 0.1",
	}

	r := bufio.NewReader(strings.NewReader(strings.Join(vectors, "\n")))
	d, err := New(r)
	if err != nil {
		t.Fatalf("could not load dictionary: %s", err.Error())
	}

	actual, ok := d.NN("minimalhouse", 2, 0, -1)
	if !ok {
		t.Fatalf("could not find neighbours")
	}

	if len(actual) != 2 {
		t.Fatalf("expected 2 results but got %d", len(actual))
	}

	expected := []string{"house", "classical"}
	for i, hit := range actual {
		if term := hit.Term; term != expected[i] {
			t.Fatalf("expected %v but got %v", expected[i], term)
		}
	}
}

