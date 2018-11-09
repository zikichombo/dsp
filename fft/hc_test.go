// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math/cmplx"
	"math/rand"
	"testing"
)

func TestHalfComplex(t *testing.T) {
	d := make([]float64, 64)
	for i := range d {
		d[i] = rand.Float64()
	}
	hc := HalfComplex(d)
	for i := 0; i < hc.Len(); i++ {
		c := hc.Cmplx(i)
		hc.SetCmplx(i, cmplx.Conj(c))
		hc.SetCmplx(i, cmplx.Conj(hc.Cmplx(i)))
		if hc.Cmplx(i) != c {
			t.Errorf("get/set")
		}
	}
}

func TestHCToFromCmplx(t *testing.T) {
	for _, N := range []int{64, 65} {
		d := make([]float64, N)
		e := make([]float64, N)
		c := make([]complex128, N)
		for i := range d {
			d[i] = rand.Float64()
			e[i] = d[i]
		}
		hc := HalfComplex(d)
		hc.ToCmplx(c)
		for i := range hc {
			hc[i] = 0.0
		}
		hc.FromCmplx(c)
		for i, v := range hc {
			if v != e[i] {
				t.Errorf("N=%d i=%d after to/from cmplx got %f not %f\n", N, i, v, e[i])
			}
		}
	}
}
