// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"
	"math/cmplx"
	"testing"

	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/sndbuf"
)

func TestDilate(t *testing.T) {
	L := 4096
	F := 1024 * freq.Hertz * 32
	fa := 24 * freq.Hertz
	fb := fa * 3
	fc := fa * 5
	waves := sndbuf.New(44100*freq.Hertz, 1)
	for i := 0; i < L; i++ {
		v := 0.0
		for _, f := range []freq.T{fa, fb, fc} {
			r := f.RadsPerAt(F)
			v += math.Sin(float64(i) * r)
		}
		waves.Send([]float64{v})
	}
	waves.Seek(0)
	d := make([]float64, L)
	waves.Receive(d)
	dc := make([]complex128, L)
	for i := range d {
		dc[i] = complex(d[i], 0)
	}

	ft := New(L)
	ft.Do(dc)
	Dilate(dc, 5, 2)

	waves.Seek(0)
	fa = (fa * 5) / 2
	fb = (fb * 5) / 2
	fc = (fc * 5) / 2
	for i := 0; i < L; i++ {
		v := 0.0
		for _, f := range []freq.T{fa, fb, fc} {
			r := f.RadsPerAt(F)
			v += math.Sin(float64(i) * r)
		}
		waves.Send([]float64{v})
	}
	waves.Seek(0)
	waves.Receive(d)
	ddc := make([]complex128, len(d))
	for i := range d {
		ddc[i] = complex(d[i], 0)
	}

	ft.Do(ddc)
	for i := 0; i < L/2; i++ {
		m1, _ := cmplx.Polar(dc[i])
		m2, _ := cmplx.Polar(ddc[i])
		if math.Abs(m2-m1) > 0.001 {
			t.Errorf("bin %d %f v %f (%f v %f)\n", i, dc[i], ddc[i], m1, m2)
		}
	}
}
