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

	sampleData := make([]float64, data.NumSamples)
	for i := range data.Data {
		sampleData[i] = float64(data.Data[i][0])
	}

	// create hammed windows for the fft
	windows := SegmentColumns(sampleData, windowLen, overlap)
	ApplyHammingColumns(windows)

	// doit
	t := fft.FFT2Real(windows)

	// get the positive fq components
	t = t[int(windowLen/2):windowLen]

	// find maximum energy
	max := 0.0
	for y := range t {
		for _, point := range t[y] {
			amplitude := cmplx.Abs(point) // discard phase
			if amplitude > max {
				max = amplitude
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
