package words

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestNewReaderDictionary(t *testing.T) {
	vectors := []string{
		"1 4",
		"dubstep -0.001058 0.002683 0.000132 0.001072",
	}

	r := bufio.NewReader(strings.NewReader(strings.Join(vectors, "\n")))
	d, err := NewReaderDictionary(r)
	if err != nil {
		t.Fatalf("could not load dictionary: %s", err.Error())
	}

	expectedEntries := 1
	if actual := d.Len(); expectedEntries != actual {
		t.Fatalf("expected %d but got %d", expectedEntries, actual)
	}

	actual, err := d.Vector("dubstep")
	if err != nil {
		t.Fatalf("could not find dubstep in dictionary")
	}

	expectedVector := []float32{-0.001058, 0.002683, 0.000132, 0.001072}
	if !reflect.DeepEqual(expectedVector, actual) {
		t.Fatalf("expected %v but got %v", expectedVector, actual)
	}
}

func TestNearestNeighbours(t *testing.T) {
	vectors := []string{
		"4 2",
		"minimalhouse 1.0, 0.0",
		"opera -1.0, 0.0",
		"house 1.0, 0.1",
		"classical -1.0, 0.1",
	}

	r := bufio.NewReader(strings.NewReader(strings.Join(vectors, "\n")))
	d, err := NewReaderDictionary(r)
	if err != nil {
		t.Fatalf("could not load dictionary: %s", err.Error())
	}

	actual, err := d.NearestNeighbours("minimalhouse", 2)
	if err != nil {
		t.Fatalf("error calculating NearestNeighbours: %s", err.Error())
	}

	expected := []string{"house", "classical"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v but got %v", expected, actual)
	}
}
