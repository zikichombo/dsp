// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"math"
	"testing"

	"zikichombo.org/sound/freq"
)

func gnr() []complex128 {
	N := 256
	res := make([]complex128, N)
	f := 1024 * freq.Hertz
	fa, fb, fc, fd := 128*freq.Hertz, 256*freq.Hertz, 64*freq.Hertz, 512*freq.Hertz
	for i := 0; i < N; i++ {
		fi := float64(i)
		va := math.Sin(fa.RadsPerAt(f) * fi)
		vb := 2 * math.Sin(fb.RadsPerAt(f)*fi)
		vc := 5 * math.Sin(fc.RadsPerAt(f)*fi)
		vd := 7 * math.Sin(fd.RadsPerAt(f)*fi)
		res[i] = complex(va+vb+vc+vd, 0)
	}
	return res
}

func TestRadix2(t *testing.T) {
	d := gnr()
	dc := make([]complex128, len(d))
	copy(dc, d)

	r2(d, newTwiddles(len(d), false), true)
	r2(d, newTwiddles(len(d), true), true)
	eps := 0.001
	for i := range d {
		re := math.Abs(real(d[i]) - real(dc[i]))
		ie := math.Abs(imag(d[i]) - imag(dc[i]))
		if re+ie > eps {
			l, u := BinRange(1024*freq.Hertz, len(d), i)
			fmt.Printf("bin %d: [%s..%s) %f\n", i, l, u, d[i])
			t.Errorf("transform and back wasn't identity: %f != %f\n", d[i], dc[i])
		}
	}
}
