package words

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestNewDictionary(t *testing.T) {
	vectors := []string{
		"1 4",
		"dubstep -0.001058 0.002683 0.000132 0.001072",
	}

	r := bufio.NewReader(strings.NewReader(strings.Join(vectors, "\n")))
	d, err := NewDictionary(r)
	if err != nil {
		t.Fatalf("could not load dictionary: %s", err.Error())
	}

	expectedEntries := 1
	if actual := len(d); expectedEntries != actual {
		t.Fatalf("expected %d but got %d", expectedEntries, actual)
	}

	actual, ok := d["dubstep"]
	if !ok {
		t.Fatalf("could not find dubstep in dictionary")
	}

	expectedVector := []float32{-0.001058, 0.002683, 0.000132, 0.001072}
	if !reflect.DeepEqual(expectedVector, actual) {
		t.Fatalf("expected %v but got %v", expectedVector, actual)
	}
}

func TestNearestNeighbours(t *testing.T) {
	d := Dictionary{}
	d["minimalhouse"] = []float32{1.0, 0.0}
	d["opera"] = []float32{-1.0, 0.0}
	d["house"] = []float32{1.0, 0.1}
	d["classical"] = []float32{-1.0, 0.1}

	actual, err := d.NearestNeighbours("minimalhouse", 2)
	if err != nil {
		t.Fatalf("error calculating NearestNeighbours: %s", err.Error())
	}

	expected := []string{"house", "classical"}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v but got %v", expected, actual)
	}
}
