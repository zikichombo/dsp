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
	F := 1024 * freq.Hertz * 16
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
	ttlErr := 0
	for i := 0; i < L/2; i++ {
		m1, _ := cmplx.Polar(dc[i])
		m2, _ := cmplx.Polar(ddc[i])
		if math.Abs(m2-m1) > 0.001 {
			ttlErr++
		}
	}
	// Many places cite this dilate mechanism as a pitch shift.  But it is not
	// purely a pitch shift, as 1) frequencies in the signal are mapped to sinc shaped
	// functions in the quantized fft frequencies, and 2) edge effects.  For
	// 1), the sinc function determining the magnitude of a frequency bin for
	// single sinusoid is at a different distance from the center frequency
	// after a pitch shift.  For 2), window functions or window size at LCM
	// of sinusoid wavelengths can help.
	//
	// The bound 0.05 was just found manually.
	if float64(ttlErr)/float64(L/2) > 0.05 {
		t.Errorf("too many errors")
	}
}
