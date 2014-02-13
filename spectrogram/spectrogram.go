package spectrogram

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
	"io"
	"math"
	"math/cmplx"
)

// Spectrogram of an audio file
type Spectrogram struct {
	data         [][]float64
	sampleRate   uint32
	windowLength int
	overlap      int
	numSamples   int
}

// NewSpectrogram constructs a spectrogram
// windowLen should be a power of 2
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

	// find maximum energy
	max := -math.MaxFloat64
	for _, value := range data.Data {
		if max < float64(value[0]) {
			max = float64(value[0])
		}
	}

	sampleData := make([]float64, data.NumSamples)
	for i := range data.Data {
		sampleData[i] = float64(data.Data[i][0]) / max
	}

	// create hammed windows for the fft
	windows := spectral.Segment(sampleData, windowLen, overlap)
	nWindow := len(windows) // (len(sampleData) - windowLen) / (windowLen - overlap) + 1
	ApplyHamming(windows)

	// do the fft, get positive fq components, and discard phase
	s := Spectrogram{
		data:         make([][]float64, nWindow),
		sampleRate:   data.SampleRate,
		windowLength: windowLen,
		overlap:      overlap,
		numSamples:   data.NumSamples,
	}

	for i, w := range windows {
		s.data[i] = make([]float64, int(windowLen/2))
		for j, v := range fft.FFTReal(w)[0:int(windowLen/2)] {
			s.data[i][j] = 2 * cmplx.Abs(v)
		}
	}

	return &s, nil
}

func (s *Spectrogram) NumFreqSlots() int { return len((*s).data[0]) }
func (s *Spectrogram) NumTimeSlots() int { return len((*s).data) }

func (s *Spectrogram) IdxToFreq(i int) (float64, error) {
	if i < 0 || i >= (*s).NumFreqSlots() {
		return math.NaN(), fmt.Errorf("invalid index for frequencies: %d", i)
	}

	return float64(i) / 2.0 * float64((*s).sampleRate) / float64((*s).NumFreqSlots()), nil
}

func (s *Spectrogram) IdxToTime(i int) (float64, error) {
	if i < 0 || i >= (*s).NumTimeSlots() {
		return math.NaN(), fmt.Errorf("invalid index for time windows: %d", i)
	}

	// corresponding slot in time domain
	j := float64(((*s).windowLength-(*s).overlap)*i) + float64((*s).windowLength)/2.0
	return j / float64((*s).sampleRate), nil
}
