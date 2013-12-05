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

	words, err := strconv.Atoi(meta[0])
	if err != nil {
		return Dictionary{}, fmt.Errorf("words NaN: %s", meta)
	}

	dims, err := strconv.Atoi(meta[1])
	if err != nil {
		return Dictionary{}, fmt.Errorf("dimensions NaN: %s", meta)
	}

	var (
		buf  = make([]float32, words*dims)       // allocate contiguous slice for vectors
		dict = make(map[string][]float32, words) // allocate dictionary for corpus
		term = make([]byte, 50)                  // assume max term length is 50
	)

	for i := 0; i < words; i++ {
		if !lines.Scan() {
			return Dictionary{}, fmt.Errorf("invalid dictionary file")
		}

		var (
			off    = i * dims
			vector = buf[off : off+dims]
		)

		fields := strings.Fields(lines.Text())

		for j := 1; j <= dims; j++ {
			weight, err := strconv.ParseFloat(fields[j], 32)
			if err != nil {
				return Dictionary{}, fmt.Errorf("could not parse weight: %s", fields[j])
			}

			vector[j-1] = float32(weight)
		}

		// Copy term name and create a new string. The strings in the slice
		// returned by strings.Fields() are backed by the input string, in this
		// case one line of the input file. This ensures that we don't hold on to
		// the lines we have read.
		copy(term, fields[0])
		dict[string(term[:len(fields[0])])] = vector
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
