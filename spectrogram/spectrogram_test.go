package spectrogram

import (
	"fmt"
	"os"
	"testing"
)

func TestSpectrogram(t *testing.T) {
	fmt.Println("hi")
	file, err := os.Open("sweep.wav")
	if err != nil {
		t.Fatalf("could not open wav fixture: %s", err.Error())
	}

	s, err := MakeSpectrogram(file)
	if err != nil {
		t.Fatalf("could not generate spectrogram: %s", err.Error())
	}

	fmt.Println(s)
}
