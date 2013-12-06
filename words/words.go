package words

import (
	"bufio"
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

// Words contains a word vector definition for each included term.
type Words interface {
	Vector(term string) ([]float32, bool)
	NearestNeighbours(term string, k int) ([]string, bool)
	Len() int
}

// NewWords creates a dictionary given a word vector file.
func NewWords(r *bufio.Reader) (Words, error) {
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
		buf     = make([]float32, words*dims)       // allocate contiguous slice for vectors
		dictmap = make(map[string][]float32, words) // allocate dictionary for corpus
		term    = make([]byte, 50)                  // assume max term length is 50
	)

	for i := 0; i < words; i++ {
		if !lines.Scan() {
			return &dict{}, fmt.Errorf("invalid dictionary file")
		}

		var (
			off    = i * dims
			vector = buf[off : off+dims]
		)

		fields := strings.Fields(lines.Text())

		for j := 1; j <= dims; j++ {
			weight, err := strconv.ParseFloat(fields[j], 32)
			if err != nil {
				return &dict{}, fmt.Errorf("could not parse weight: %s", fields[j])
			}

			vector[j-1] = float32(weight)
		}

		// Copy term name and create a new string. The strings in the slice
		// returned by strings.Fields() are backed by the input string, in this
		// case one line of the input file. This ensures that we don't hold on to
		// the lines we have read.
		copy(term, fields[0])

		key := string(term[:len(fields[0])])
		dictmap[key] = vector
	}

	return &dict{dictmap: dictmap}, nil
}

type dict struct {
	dictmap map[string][]float32
}

// NearestNeighbours returns k nearest tags in vector space.
func (d *dict) NearestNeighbours(term string, k int) ([]string, bool) {
	termVector, ok := (*d).dictmap[term]
	if !ok {
		return []string{}, false
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	for t, v := range d.dictmap {
		if term == t {
			continue
		}

		similarity, err := cosine(termVector, v)
		if err != nil {
			msg := fmt.Sprintf("could not compare %s with %s. should never happen", term, t)
			panic(msg)
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
