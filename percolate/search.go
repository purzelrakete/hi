package main

// Search uses union find to search for a percolation path thorough porous
// material. Returns an array of bools indicating if the cell index is open
// and connected to the top.

func Search(m *Material) []bool {
	// add virtual nodes top, bottom.
	cs := make([]int, m.L+2)
	for i, _ := range cs {
		cs[i] = i
	}

	// connect top row
	for i := 0; i < m.W; i++ {
		if m.Vacant(i) {
			cs[i+1] = cs[0]
		}
	}

	// connect bottom row
	for i := m.L - m.W; i < m.L; i++ {
		if m.Vacant(i) {
			cs[i+1] = cs[len(cs)-1]
		}
	}

	union := func(insert, erase int) {
		source, sink := cs[insert+1], cs[erase+1]
		for i, v := range cs {
			if v == sink {
				cs[i] = source
			}
		}
	}

	// connect cells
	for i, _ := range m.Cells {
		if m.Vacant(i) {
			if j := i - 1; i%m.W != 0 && m.Vacant(j) {
				union(i, j) // left
			}

			if j := i - m.W; j > -1 && m.Vacant(j) {
				union(i, j) // top
			}

			if j := i + 1; j%m.W != 0 && m.Vacant(j) {
				union(i, j) // right
			}

			if j := i + m.W; j < m.L && m.Vacant(j) {
				union(i, j) // bottom
			}
		}
	}

	// find top component
	p := make([]bool, m.L)
	for i, _ := range p {
		if cs[i+1] == cs[0] {
			p[i] = true
		}
	}

	return p
}
