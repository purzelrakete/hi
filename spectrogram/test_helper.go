package spectrogram

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// seqFloat64(1, 5, 1) -> [1 2 3 4 5]
func seqFloat64(lower, upper, step float64) []float64 {
	steps := int((upper - lower) / step)
	ret := make([]float64, steps+1)

	for i := 0; i <= steps; i++ {
		ret[i] = float64(i)*step + lower
	}

	return ret
}

func readMatrix(filename string) [][]float64 {
	file, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}

	matrix := [][]float64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")

		row := []float64{}
		for _, field := range fields {
			n, err := strconv.ParseFloat(field, 64)
			if err != nil {
				panic(err.Error())
			}

			row = append(row, n)
		}

		matrix = append(matrix, row)
	}

	return matrix
}

func drawMatrix(in, out string) {
	t := readMatrix(in)

	// find maximum energy
	max := 0.0
	for y := range t {
		for _, amplitude := range t[y] {
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
			s[y][x] = uint8(255 * t[y][x] / max)
		}
	}

	s = s[int(len(s)/2):len(s)]
	Draw(&s, out)
}
