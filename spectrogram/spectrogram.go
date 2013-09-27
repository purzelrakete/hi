package spectrogram

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
	"github.com/mjibson/go-dsp/window"
	"io"
)

// Spectrogram of an audio file
func Spectrogram(r io.Reader, windowLen, overlap int) ([][]complex128, error) {
	data, err := wav.ReadWav(r)
	if err != nil {
		return [][]complex128{}, err
	}

	if expected, actual := uint16(16), data.BitsPerSample; expected != actual {
		return [][]complex128{}, fmt.Errorf("rate %d, not %d", actual, expected)
	}

	if expected, actual := uint16(1), data.NumChannels; expected != actual {
		return [][]complex128{}, fmt.Errorf("%d chans, not %d", actual, expected)
	}

	sampleData := make([]float64, len(data.Data16[0]))
	for i, s := range data.Data16[0] {
		sampleData[i] = float64(s)
	}

	hamming := window.Hamming(len(data.Data16[0]))
	windows := spectral.Segment(sampleData, windowLen, overlap)
	for i, w := range windows {
		for j, _ := range w {
			windows[i][j] = windows[i][j] * hamming[j]
		}
	}

	return fft.FFT2Real(windows), nil
}
