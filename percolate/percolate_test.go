package main

import "testing"

func Test3Percolation(t *testing.T) {
	side := 3
	c := []bool{
		false, true, false,
		true, true, false,
		false, true, false,
	}

	m, _ := NewMaterial(side, side, c)
	connected := Search(m)
	if !isEqualSliceBools(connected, c) {
		t.Fatalf("expected %v but got %v", c, connected)
	}
}

func Test4Percolation(t *testing.T) {
	side := 4
	c := []bool{
		true, true, true, true,
		true, true, false, true,
		true, false, true, true,
		true, true, false, true,
	}

	m, _ := NewMaterial(side, side, c)
	connected := Search(m)
	if !isEqualSliceBools(connected, c) {
		t.Fatalf("expected %v but got %v", c, connected)
	}
}

func isEqualSliceBools(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i, val := range a {
		if b[i] != val {
			return false
		}
	}

	return true
}
