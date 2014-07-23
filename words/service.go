package words

import (
	"bufio"
	"fmt"
	"net/http"
)

// Service returns a list of terms similar to the given one, bool ok, and
// a vector representation of the query term.
type Service interface {
	Vector(term string) ([]float32, bool)
	NN(term string, k, minfq int, θ float32) ([]Hit, bool)
	NNVector(queryVector []float32, k, minfq int, θ float32) ([]Hit, bool)
}

// NewService is a similarity function backed by word2vec vectors
func NewService(modelPath string) (Service, error) {
	words, err := loadWords(modelPath)
	if err != nil {
		return nil, fmt.Errorf("could not get words: %s", err.Error())
	}

	return words, nil
}

// loadWords retrieves the model from a url and constructs Words
func loadWords(modelPath string) (Words, error) {
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

	return words, nil
}
