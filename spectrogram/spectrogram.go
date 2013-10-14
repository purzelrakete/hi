package spectrogram

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
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

	// create hammed windows for the fft
	windows := SegmentColumns(sampleData, windowLen, overlap)
	ApplyHammingColumns(windows)

	// doit
	t := fft.FFT2Real(windows)

	// discard above the nyquist frequency

	// find maximum energy
	max := 0.0
	for y := range t {
		for x := range t[y] {
			point := cmplx.Abs(t[y][x]) // discard phase
			if point > max {
				max = point
			}
		}
	}

	// normalized by max, convert to unit8, transpose
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
