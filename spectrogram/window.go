package spectrogram

import "github.com/mjibson/go-dsp/window"

// ApplyHamming applies a hamming filter to each column
func ApplyHamming(x [][]float64) {
	segments, size := len(x), len(x[0])
	hamming := window.Hamming(size)
	for i := 0; i < segments; i++ {
		for j := 0; j < size; j++ {
			x[i][j] = x[i][j] * hamming[j]
		}
	}
}
