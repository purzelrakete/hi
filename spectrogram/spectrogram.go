package spectrogram

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
	"github.com/mjibson/go-dsp/window"
	"io"
	"math/cmplx"
)

// Spectrogram of an audio file
type Spectrogram [][]uint8

// NewSpectrogram constructs a spectrogram
func NewSpectrogram(r io.Reader, windowLen, overlap int) (*Spectrogram, error) {
	data, err := wav.ReadWav(r)
	if err != nil {
		return &Spectrogram{}, err
	}

	if expected, actual := uint16(16), data.BitsPerSample; expected != actual {
		return &Spectrogram{}, fmt.Errorf("rate %d, not %d", actual, expected)
	}

	if expected, actual := uint16(1), data.NumChannels; expected != actual {
		return &Spectrogram{}, fmt.Errorf("%d chans, not %d", actual, expected)
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

	t := fft.FFT2Real(windows)

	max := 0.0
	for y := range t {
		for x := range t[y] {
			point := t[y][x]
			if cmplx.Abs(point) > max {
				max = cmplx.Abs(point)
			}
		}
	}

	// normalized by max, unit8
	s := make(Spectrogram, len(t))
	for y := range t {
		s[y] = make([]uint8, len(t[y]))
		for x := range t[y] {
			s[y][x] = uint8(255 * cmplx.Abs(t[y][x]) / max)
		}
	}

	return &s, nil
}

func (s *Spectrogram) Height() int { return len(*s) }
func (s *Spectrogram) Width() int  { return len((*s)[0]) }
