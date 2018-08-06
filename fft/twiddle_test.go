// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"
	"testing"
)

func TestTwiddle(t *testing.T) {
	N := 1024
	tw := newTwiddles(1024, false)
	itw := newTwiddles(1024, true)
	eps := 1e-8
	for i := 0; i < N; i++ {
		s, c := tw.sincos(i)
		r := (float64(i) / float64(N)) * math.Pi * 2
		ms, mc := math.Sin(r), math.Cos(r)
		es, ec := math.Abs(s+ms), math.Abs(c-mc)
		if es > eps {
			t.Errorf("sin error too large: %f", es)
		}
		if ec > eps {
			t.Errorf("cos error too large: %f", ec)
		}
		twid := tw.cmplx(i)
		itwid := itw.cmplx(i)
		if math.Abs(imag(twid)+imag(itwid)) > eps {
			t.Errorf("imag doesn't add to 0")
		}
	}
}
