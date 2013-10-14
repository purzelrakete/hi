package spectrogram

import "github.com/mjibson/go-dsp/window"

// SegmentColumns x segmented into segments of length size with  noverlap.
// Number of segments returned is (len(x) - size) / (size - noverlap) + 1.
func SegmentColumns(x []float64, size, stride int) [][]float64 {
	lx := len(x)

	var segments int
	if lx == size {
		segments = 1
	} else if lx > size {
		segments = (len(x)-size)/stride + 1
	} else {
		segments = 0
	}

	// set up matrix
	r := make([][]float64, size)
	for i := 0; i < size; i++ {
		r[i] = make([]float64, segments)
	}

	// fill matrix
	for i, offset := 0, 0; i < segments; i++ {
		for j := 0; j < size; j++ {
			r[j][i] = x[offset+j]
		}

		offset += stride
	}

	return r
}

// ApplyHammingColumns applies a hamming filter to each column
func ApplyHammingColumns(x [][]float64) {
	size, segments := len(x), len(x[0])
	hamming := window.Hamming(size)
	for i := 0; i < segments; i++ {
		for j := 0; j < size; j++ {
			x[j][i] = x[j][i] * hamming[j]
		}
	}
}
