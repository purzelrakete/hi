package words

import (
	"bufio"
	"container/heap"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Words contains a word vector definition for each included term
type Words interface {
	Vector(term string) ([]float32, bool)
	NN(term string, k, fq int, θ float32) ([]Hit, bool)
	NNVector(vector []float32, k, fq int, θ float32) ([]Hit, bool)
	Len() int
}

// Hit is a similar term
type Hit struct {
	Term       string    `json:"term"`
	Frequency  int       `json:"frequency"`
	Similarity float32   `json:"similarity"`
}

// New creates a dictionary given a word vector file.
func New(r *bufio.Reader) (Words, error) {
	lines := bufio.NewScanner(r)

	// first line is dictionary size, vector dimensions
	lines.Scan()
	meta := strings.Fields(lines.Text())
	if len(meta) != 2 {
		return &dict{}, fmt.Errorf("does not start with 2 fields: %v", meta)
	}

	words, err := strconv.Atoi(meta[0])
	if err != nil {
		return &dict{}, fmt.Errorf("words NaN: %s", meta)
	}

	dims, err := strconv.Atoi(meta[1])
	if err != nil {
		return &dict{}, fmt.Errorf("dimensions NaN: %s", meta)
	}

	var (
		buf         = make([]float32, words*dims)       // allocate contiguous slice for vectors
		dictmap     = make(map[string][]float32, words) // allocate dictionary for corpus
		terms       = make([]string, words)             // allocate term ordinal map
		frequencies = make([]int, words*dims)           // term frequencies
		term        = make([]byte, 0, 128)              // temp buffer for copying term
		vecOffset   = 2                                 // index at which vectors start
	)

	for i := 0; i < words; i++ {
		if !lines.Scan() {
			if err := lines.Err(); err != nil {
				return &dict{}, fmt.Errorf("error scanning dictionary file: %s", err)
			}

			return &dict{}, fmt.Errorf("eof. %d words in header, %d processed", words, i)
		}

		var (
			off    = i * dims
			vector = buf[off : off+dims]
			fields = strings.Fields(lines.Text())
		)

		for j := 0; j < dims; j++ {
			weight, err := strconv.ParseFloat(fields[j+vecOffset], 32)
			if err != nil {
				return &dict{}, fmt.Errorf("could not parse weight: %s", fields[j])
			}

			vector[j] = float32(weight)
		}

		fq, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			return &dict{}, fmt.Errorf("could not parse frequency: %s", fields[1])
		}

		frequencies[i] = int(fq)

		// Copy term name and create a new string. The strings in the slice
		// returned by strings.Fields() are backed by the input string, in this
		// case one line of the input file. This ensures that we don't hold on to
		// the lines we have read.
		term = append(term[:0], fields[0]...)
		terms[i] = string(term)
		dictmap[string(term)] = vector
	}

	return &dict{
		dictmap:     dictmap,
		buf:         buf,
		terms:       terms,
		frequencies: frequencies,
		dims:        dims,
		words:       words,
	}, nil
}

type dict struct {
	dictmap     map[string][]float32
	buf         []float32
	terms       []string
	frequencies []int
	dims        int
	words       int
}

// NN returns k nearest neighbours in vector space.
func (d *dict) NN(term string, k, minFq int, θ float32) ([]Hit, bool) {
	termVector, ok := d.dictmap[term]
	if !ok {
		return []Hit{}, false
	}

	return d.NNVector(termVector, k, minFq, θ)
}

// NN returns k nearest neighbours in vector space for a given vector.
func (d *dict) NNVector(termVector []float32, k, minFq int, θ float32) ([]Hit, bool) {
	pq := &PriorityQueue{}
	heap.Init(pq)

	for i := 0; i < d.words; i++ {
		var (
			off = i * d.dims
			v   = d.buf[off : off+d.dims]
			t   = d.terms[i]
			fq  = d.frequencies[i]
		)

 		if fq < minFq {
 			continue
 		}
 
 		if reflect.DeepEqual(v, termVector) {
 			continue
 		}

		similarity, err := cosine(termVector, v)
		if err != nil {
			panic(fmt.Sprintf("could not compare %l with %l.", termVector, t))
		}

		if similarity < θ {
			continue
		}

		// fill the queue up with the first K candidates
		if pq.Len() < k {
			heap.Push(pq, &Item{
				ordinal:  i,
				priority: similarity,
			})

			// check candidate proximity
		} else {
			if similarity > (*pq)[0].priority {
				heap.Pop(pq)
				heap.Push(pq, &Item{
					ordinal:  i,
					priority: similarity,
				})
			}
		}
	}

	length := pq.Len()
	terms := make([]Hit, length)
	for i := 0; i < length; i++ {
		item := heap.Pop(pq).(*Item)
		terms[length-i-1] = Hit{
			Term:       d.terms[item.ordinal],
			Similarity: item.priority,
			Frequency:  d.frequencies[item.ordinal],
		}
	}

	return terms, true
}

// Vector returns vector for a given term
func (d *dict) Vector(term string) ([]float32, bool) {
	v, ok := d.dictmap[term]
	return v, ok
}

// Len of current dictionary
func (d *dict) Len() int {
	return len(d.dictmap)
}
