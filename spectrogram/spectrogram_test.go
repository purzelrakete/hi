package spectrogram

import (
	"os"
	"testing"
)

func TestDrawSpectrogram(t *testing.T) {
	file, err := os.Open("sweep.wav")
	if err != nil {
		t.Fatalf("could not open wav fixture: %s", err.Error())
	}

	windowLen, overlap := 256, 128
	s, err := Spectrogram(file, windowLen, overlap)
	if err != nil {
		t.Fatalf("could not generate spectrogram: %s", err.Error())
	}

	Draw(s, windowLen, "spectrogram.png")

	// TODO: think of a reasonable assertion.
}
