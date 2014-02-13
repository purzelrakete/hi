package fftw

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"
)

func TestFFTW(t *testing.T) {
	data, _ := Alloc2D(100, 1)
	fwd := New(data)
	defer fwd.Destroy()

	for i := 0; i < len(data); i++ {
		data[i][0] = complex(math.Sin(float64(i)), 0)
	}

	fwd.Execute()

	max, imax := 0.0, 0
	for i := 0; i < len(data)/2; i++ {
		if modulus := cmplx.Abs(data[i][0]); modulus > max {
			max, imax = modulus, i
		}
	}

	// FIXME(rk): calcuate if 16 is the correct fq band
	if expectedBand, actualBand := 16, imax; actualBand != expectedBand {
		fmt.Errorf("expected band %d but got %d", expectedBand, actualBand)
	}
}
