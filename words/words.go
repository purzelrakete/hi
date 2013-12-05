package words

import (
	"bufio"
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

// Dictionary contains a word vector definition for each included term.
type Dictionary struct {
	size     int            // current dictionary size
	capacity int            // dictionary capacity
	dims     int            // vector dimensions
	vectors  []float32      // all vectors combined
	keys     []string       // list of keys into vectors[offset * dims]
	index    map[string]int // quick lookup of key indices
}

// NewDictionary initializes internal data structures
func NewDictionary(words, dims int) Dictionary {
	return Dictionary{
		size:     0,
		capacity: words,
		dims:     dims,
		vectors:  make([]float32, words*dims),
		keys:     make([]string, words),
	}
}

// AddTerm adds a term to the dictionary and returns an ordinal key for it
func (d *Dictionary) AddTerm(term string) (int, error) {
	if d.size > d.capacity {
		return 0, fmt.Errorf("exceeded capcity of %d", d.capacity)
	}

	d.keys[d.size] = term
	d.index[term] = d.size
	d.size++

	return d.size, nil
}

// AddWeight adds a weight for a given term at a given dimension of the vector
func (d *Dictionary) AddWeight(termOrdinal, dimension int, weight float32) error {
	if termOrdinal > d.size {
		return fmt.Errorf("invalid ordinal %d", termOrdinal)
	}

	d.vectors[(termOrdinal*d.dims)+dimension] = weight

	return nil
}

// Vector returns the vector associated with the term
func (d *Dictionary) Vector(term string) ([]float32, error) {
	ordinal, ok := (*d).index[term]
	if !ok {
		return []float32{}, fmt.Errorf("%s not in dictionary", term)
	}

	offset := ordinal * d.dims
	return d.vectors[offset : offset+d.dims], nil
}

// Len of the dictionary
func (d *Dictionary) Len() int { return d.size }

// NewReaderDictionary creates a dictionary given a word vector file.
func NewReaderDictionary(r *bufio.Reader) (Dictionary, error) {
	lines := bufio.NewScanner(r)

	// first line is dictionary size, vector dimensions
	lines.Scan()
	meta := strings.Fields(lines.Text())
	if len(meta) != 2 {
		return Dictionary{}, fmt.Errorf("does not start with 2 fields: %v", meta)
	}

	words, err := strconv.Atoi(meta[0])
	if err != nil {
		return Dictionary{}, fmt.Errorf("words? %s", err.Error())
	}

	dims, err := strconv.Atoi(meta[1])
	if err != nil {
		return Dictionary{}, fmt.Errorf("dimensions NaN: %s", meta)
	}

	dict := NewDictionary(words, dims)
	for lines.Scan() {
		fields := strings.Fields(lines.Text())
		term := fields[0]
		ordinal, err := dict.AddTerm(term)
		if err != nil {
			return Dictionary{}, fmt.Errorf("could not add term %s: %s", term, err.Error())
		}

		// skip fields[0], this is the term name.
		for i := 1; i <= len(fields); i++ {
			weight, err := strconv.ParseFloat(fields[i], 32)
			if err != nil {
				return Dictionary{}, fmt.Errorf("could not parse weight: %s", fields[i])
			}

			dimension := i - 1
			dict.AddWeight(ordinal, dimension, float32(weight))
		}
	}

	return dict, nil
}

// NearestNeighbours returns k nearest tags in vector space.
func (d *Dictionary) NearestNeighbours(t string, k int) ([]string, error) {
	termVector, err := d.Vector(t)
	if err != nil {
		return []string{}, fmt.Errorf("%s not in dictionary", t)
	}

	pq := &PriorityQueue{}
	heap.Init(pq)

	var term string
	for i := 0; i < d.size; i++ {
		if i%d.dims == 0 { // new term vector starts here
			if term == t {
				continue
			}
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
