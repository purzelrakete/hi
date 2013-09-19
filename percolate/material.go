package main

import (
	"fmt"
	"math/rand"
	"time"
)

func NewMaterial(height, width int, cells []bool) (*Material, error) {
	return &Material{
		Cells: cells,
		H:     height,
		W:     width,
		L:     height * width,
	}, nil
}

func NewRandomMaterial(height, width int, p float64) (*Material, error) {
	if p < 0 || p > 1 {
		return &Material{}, fmt.Errorf("p ∉ [0,1]")
	}

	cells := make([]bool, width*height)
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, _ := range cells {
		if r.Float64() < p {
			cells[i] = true
		}
	}

	return NewMaterial(height, width, cells)
}

type Material struct {
	Cells []bool
	H     int // height
	W     int // width
	L     int // length
}

func (m *Material) Vacant(i int) bool {
	return m.Cells[i]
}
