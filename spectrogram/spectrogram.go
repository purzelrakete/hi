package spectrogram

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
	"github.com/mjibson/go-dsp/window"
	"io"
)

// Spectrogram of audio content
type Spectrogram [][]complex128

// MakeSpectrogram returns a spectrogram of the audio file
func MakeSpectrogram(r io.Reader) (Spectrogram, error) {
	data, err := wav.ReadWav(r)
	if err != nil {
		return Spectrogram{}, err
	}

	if expected, actual := uint16(16), data.BitsPerSample; expected != actual {
		return Spectrogram{}, fmt.Errorf("bitrate %d, not %d", actual, expected)
	}

	if expected, actual := uint16(1), data.NumChannels; expected != actual {
		return Spectrogram{}, fmt.Errorf("%d channels, not %d", actual, expected)
	}

	sampleData := make([]float64, len(data.Data16[0]))
	for i, s := range data.Data16[0] {
		sampleData[i] = float64(s)
	}

	sampleLen, windowLen, overlap := len(data.Data16[0]), 256, 10
	hamming := window.Hamming(sampleLen)

	windows := spectral.Segment(sampleData, windowLen, overlap)
	for i, w := range windows {
		for j, _ := range w {
			windows[i][j] = windows[i][j] * hamming[j]
		}
	}

	return fft.FFT2Real(windows), nil
}
