package spectrogram

// seqFloat64(1, 5, 1) -> [1 2 3 4 5]
func seqFloat64(lower, upper, step float64) []float64 {
	steps := int((upper - lower) / step)
	ret := make([]float64, steps+1)

	for i := 0; i <= steps; i++ {
		ret[i] = float64(i)*step + lower
	}

	return ret
}
