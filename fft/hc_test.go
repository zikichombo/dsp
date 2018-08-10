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
