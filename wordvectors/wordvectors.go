package wordvectors

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Dictionary contains a word vector definition for each included term.
type Dictionary map[string][]float32

// NewDictionary creates a dictionary given a word vector file.
func NewDictionary(r io.Reader) (Dictionary, error) {
	lines := bufio.NewScanner(r)

	// first line is dictionary size, vector dimensions
	lines.Scan()
	meta := strings.Fields(lines.Text())
	if len(meta) != 2 {
		return Dictionary{}, fmt.Errorf("does not start with 2 fields: %v", meta)
	}

	dims, err := strconv.Atoi(meta[1])
	if err != nil {
		return Dictionary{}, fmt.Errorf("dimensions NaN: %s", meta)
	}

	dict := Dictionary{}
	for lines.Scan() {
		fields := strings.Fields(lines.Text())
		vector := make([]float32, dims)
		for i := 1; i <= dims; i++ {
			weight, err := strconv.ParseFloat(fields[i], 32)
			if err != nil {
				return Dictionary{}, fmt.Errorf("could not parse weight: %s", fields[i])
			}

			vector[i-1] = float32(weight)
		}

		dict[fields[0]] = vector
	}

	return dict, nil
}

// CosineSimilarity returns terms within threshold distance of the given term.
func (d *Dictionary) CosineSimilarity(term string, threshold float32) []string {
	return []string{
		term,
	}
}
