package spectrogram

import (
	"math"
	"os"
	"testing"
)

func TestDrawSpectrogram(t *testing.T) {
	file440, err := os.Open("440.wav")
	if err != nil {
		t.Fatalf("could not open wav fixture: %s", err.Error())
	}

	defer file440.Close()

	windowLen, noverlap := 990, 0
	s, err := NewSpectrogram(file440, windowLen, noverlap)
	if err != nil {
		t.Fatalf("could not generate spectrogram: %s", err.Error())
	}

	if err := DrawXY(s, "spectrogram-440.png", 200, 1000); err != nil {
		t.Fatalf("could not draw spectrogram: %s", err.Error())
	}

	fileSweep, err := os.Open("sweep.wav")
	if err != nil {
		t.Fatalf("could not open wav fixture: %s", err.Error())
	}

	defer fileSweep.Close()

	s, err = NewSpectrogram(fileSweep, windowLen, noverlap)
	if err != nil {
		t.Fatalf("could not generate spectrogram: %s", err.Error())
	}

	if err = Draw(s, "spectrogram-sweep.png"); err != nil {
		t.Fatalf("could not draw spectrogram: %s", err.Error())
	}

	// TODO: assert something for the love of god
}

func TestSpectrogram440(t *testing.T) {
	file, err := os.Open("440.wav")
	if err != nil {
		t.Fatalf("could not open wav fixture: %s", err.Error())
	}

	windowLen, noverlap := 400, 100
	s, err := NewSpectrogram(file, windowLen, noverlap)
	if err != nil {
		t.Fatalf("could not generate spectrogram: %s", err.Error())
	}

	max := make([]float64, s.NumTimeSlots())
	maxi := make([]int, s.NumTimeSlots())
	for i, a := range s.data {
		for j, b := range a {
			if j == 0 || b > max[i] {
				max[i] = b
				maxi[i] = j
			}
		}
	}

	expected := 440.0
	for _, f := range maxi {
		got, _ := s.IdxToFreq(f)
		if expected != got {
			t.Fatalf("expected %v but got %v", expected, got)
		}
	}
}

func TestSpectrogramSweep(t *testing.T) {
	file, err := os.Open("sweep.wav")
	if err != nil {
		t.Fatalf("could not open wav fixture: %s", err.Error())
	}

	windowLen, noverlap := 392, 0
	s, err := NewSpectrogram(file, windowLen, noverlap)
	if err != nil {
		t.Fatalf("could not generate spectrogram: %s", err.Error())
	}

	max := make([]float64, s.NumTimeSlots())
	maxi := make([]int, s.NumTimeSlots())
	for i, a := range s.data {
		for j, b := range a {
			if j == 0 || b > max[i] {
				max[i] = b
				maxi[i] = j
			}
		}
	}

	// ignore first and second which are noise
	for i := 2; i < len(maxi); i++ {
		got := maxi[i]
		expected := i
		if expected != got {
			t.Fatalf("%d,expected %v but got %v", i, expected, got)
		}
	}
}

func TestApplyHamming(t *testing.T) {
	windows := [][]float64{
		[]float64{1.0, 1.0, 1.0, 1.0, 1.0},
		[]float64{1.0, 1.0, 1.0, 1.0, 1.0},
	}

	expected := [][]float64{
		[]float64{0.08, 0.54, 1, 0.54, 0.08},
		[]float64{0.08, 0.54, 1, 0.54, 0.08},
	}

	ApplyHamming(windows)
	eps := 0.000001
	for i, ws := range windows {
		for j, w := range ws {
			if e := expected[i][j]; math.Abs(w-expected[i][j]) > eps {
				t.Fatalf("expected %v but got %v", e, w)
			}
		}
	}
}
