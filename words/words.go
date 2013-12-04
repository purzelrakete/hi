package words

import (
	"bufio"
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

// Dictionary contains a word vector definition for each included term.
type Dictionary map[string][]float32

// NewDictionary creates a dictionary given a word vector file.
func NewDictionary(r *bufio.Reader) (Dictionary, error) {
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

// NearestNeighbours returns k nearest tags in vector space.
func (d *Dictionary) NearestNeighbours(term string, k int) ([]string, error) {
	termVector, ok := (*d)[term]
	if !ok {
		return []string{}, fmt.Errorf("%s not in dictionary", term)
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	for t, v := range *d {
		if term == t {
			continue
		}

		similarity, err := cosine(termVector, v)
		if err != nil {
			return []string{}, fmt.Errorf("could not compare %s with %s", term, t)
		}

		// fill the queue up with the first K candidates
		if pq.Len() < k {
			heap.Push(pq, &Item{
				value:    t,
				priority: similarity,
			})

			// check candidate proximity
		} else {
			if similarity > (*pq)[0].priority {
				heap.Pop(pq)
				heap.Push(pq, &Item{
					value:    t,
					priority: similarity,
				})
			}
		}
	}

	terms := make([]string, k)
	for i := 0; pq.Len() > 0; i++ {
		item := heap.Pop(pq).(*Item)
		terms[k-i-1] = item.value
	}

	return terms, nil
}
