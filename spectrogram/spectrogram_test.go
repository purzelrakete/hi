package spectrogram

import (
	"os"
	"reflect"
	"testing"
)

func TestDrawSpectrogram(t *testing.T) {
	file, err := os.Open("sweep.wav")
	if err != nil {
		t.Fatalf("could not open wav fixture: %s", err.Error())
	}

	windowLen, noverlap := 256, 128
	s, err := NewSpectrogram(file, windowLen, noverlap)
	if err != nil {
		t.Fatalf("could not generate spectrogram: %s", err.Error())
	}

	Draw(s, "spectrogram.png")

	// TODO: assert something for the love of god
}

func TestSegmentColumns(t *testing.T) {
	expected := [][]float64{
		[]float64{1.0, 3.0, 5.0},
		[]float64{2.0, 4.0, 6.0},
		[]float64{3.0, 5.0, 7.0},
		[]float64{4.0, 6.0, 8.0},
		[]float64{5.0, 7.0, 9.0},
	}

	actual := SegmentColumns(seqFloat64(1.0, 10.0, 1.0), 5, 2)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v but got %v", expected, actual)
	}

	actual = SegmentColumns(seqFloat64(1.0, 9.0, 1.0), 5, 2)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v but got %v", expected, actual)
	}

	expected = [][]float64{
		[]float64{1.0, 4.0},
		[]float64{2.0, 5.0},
		[]float64{3.0, 6.0},
		[]float64{4.0, 7.0},
		[]float64{5.0, 8.0},
		[]float64{6.0, 9.0},
	}

	actual = SegmentColumns(seqFloat64(1.0, 9.0, 1.0), 6, 3)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v but got %v", expected, actual)
	}
}

func TestApplyHamming(t *testing.T) {
	windows := [][]float64{
		[]float64{1.0, 1.0},
		[]float64{1.0, 1.0},
		[]float64{1.0, 1.0},
		[]float64{1.0, 1.0},
		[]float64{1.0, 1.0},
	}

	expected := [][]float64{
		[]float64{0.08000000000000002, 0.08000000000000002},
		[]float64{0.54, 0.54},
		[]float64{1, 1},
		[]float64{0.5400000000000001, 0.5400000000000001},
		[]float64{0.08000000000000002, 0.08000000000000002},
	}

	ApplyHammingColumns(windows)
	if !reflect.DeepEqual(expected, windows) {
		t.Fatalf("expected %v but got %v", expected, windows)
	}
}

// FIXME(rk): this is a temporary test; remove.
func TestOctaveRendering(ts *testing.T) {
	drawMatrix("octave.csv", "octave.png")
}
