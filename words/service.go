package words

import (
	"bufio"
	"fmt"
	"net/http"
)

// WordsService returns a list of terms similar to the given one.
type WordsService func(term string, k, minfq int, θ float32) ([]Hit, bool)

// NewWordsService is a similarity function backed by word2vec vectors
func NewWordsService(modelPath string) (WordsService, error) {
	resp, err := http.Get(modelPath)
	if err != nil {
		return nil, fmt.Errorf("could not get %s: %s", modelPath, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get %s: %d", modelPath, resp.StatusCode)
	}

	words, err := NewWords(bufio.NewReader(resp.Body))
	if err != nil {
		return nil, fmt.Errorf("could not get words: %s", err.Error())
	}

	return func(term string, k, minfq int, θ float32) ([]Hit, bool) {
		return words.NN(term, k, minfq, θ)
	}, nil
}
