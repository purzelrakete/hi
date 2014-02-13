package fftw

// #cgo LDFLAGS: -lfftw3 -lm
// #include <fftw3.h>
import "C"

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// Plan for fftw
type Plan struct {
	fftw_p C.fftw_plan
	sync.Mutex
}

// New fftw plan. Destroy() after use.
func New(data [][]complex128) *Plan {
	pData := (unsafe.Pointer)(&data[0][0])
	return &Plan{
		fftw_p: C.fftw_plan_dft_2d(
			(C.int)(len(data)),       // nrows
			(C.int)(len(data[0])),    // ncolumns
			(*C.fftw_complex)(pData), // input slice
			(*C.fftw_complex)(pData), // output slice
			C.int(C.FFTW_FORWARD),
			C.uint(C.FFTW_ESTIMATE)),
	}
}

// Execute the plan
func (p *Plan) Execute() {
	p.Lock()
	defer p.Unlock()
	C.fftw_execute(p.fftw_p)
}

// Destroy the plan
func (p *Plan) Destroy() {
	p.Lock()
	defer p.Unlock()
	C.fftw_destroy_plan(p.fftw_p)
}

// Alloc2D allocates a complex matrix for fftw
func Alloc2D(m, n int) ([][]complex128, error) {
	if m <= 0 || n <= 0 {
		return [][]complex128{}, fmt.Errorf("positive dimensions required")
	}

	samples := m * n
	var slice []complex128
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	*header = reflect.SliceHeader{
		Data: uintptr(C.fftw_malloc((C.size_t)(16 * samples))),
		Len:  samples,
		Cap:  samples,
	}

	r := make([][]complex128, m)
	for i := range r {
		r[i] = slice[i*n : (i+1)*n]
	}

	return r, nil
}
