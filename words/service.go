package words

import (
	"bufio"
	"fmt"
	"net/http"
)

// Service returns a list of terms similar to the given one, bool ok, and
// a vector representation of the query term.
type Service func(term string, k, minfq int, θ float32) ([]Hit, bool, []float32)

// NewService is a similarity function backed by word2vec vectors
func NewService(modelPath string) (Service, error) {
	resp, err := http.Get(modelPath)
	if err != nil {
		return nil, fmt.Errorf("could not get %s: %s", modelPath, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get %s: %d", modelPath, resp.StatusCode)
	}

	words, err := New(bufio.NewReader(resp.Body))
	if err != nil {
		return nil, fmt.Errorf("could not get words: %s", err.Error())
	}

	return func(term string, k, minfq int, θ float32) ([]Hit, bool, []float32) {
		return words.NN(term, k, minfq, θ)
	}, nil
}
