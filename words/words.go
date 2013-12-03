package words

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Dictionary contains a word vector definition for each included term.
type Dictionary map[string][]float64

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
		vector := make([]float64, dims)
		for i := 1; i <= dims; i++ {
			weight, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				return Dictionary{}, fmt.Errorf("could not parse weight: %s", fields[i])
			}

			vector[i-1] = float64(weight)
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
			if similarity > (*pq)[pq.Len()-1].priority {
				heap.Pop(pq)
				heap.Push(pq, &Item{
					value:    t,
					priority: similarity,
				})
			}
		}
	}

	// FIXME(rk): will always return k-1, since self simiarity is 1.0.
	terms := []string{}
	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		terms = append(terms, item.value)
	}

	return terms, nil
}
