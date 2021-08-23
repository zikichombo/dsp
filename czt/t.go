// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package czt

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/zikichombo/dsp/fft"
	"github.com/zikichombo/sound/freq"
)

type T struct {
	nS, nB, nPad     int
	start, step      float64
	kern, aTab, wTab []complex128
	ft               *fft.T
}

// New creates a new chirp-z transformer object.
//
// The transformer takes as input nS samples and
// returns a freqrency domain picture of those
// samples with nB bins, focusing on frequencies
// from start to end.  Start and end are in radians
// per sample.
func New(nS, nB int, start, end float64) *T {
	nPad := findL(nS, nB)
	step := (end - start) / float64(nB)
	res := &T{
		nS:    nS,
		nB:    nB,
		nPad:  nPad,
		start: start,
		step:  step,
		ft:    fft.New(nPad),
		kern:  make([]complex128, nPad),
		aTab:  make([]complex128, nS),
		wTab:  make([]complex128, nS)}
	res.initWK()
	res.initA()
	return res
}

// Do performs the transform as configured in
// the call to New() which created t.
func (t *T) Do(src []complex128) []complex128 {
	src = t.Win(src)
	src = pad(src, t.nPad)

	// weight inputs
	for i := range src[:t.nS] {
		src[i] *= t.aTab[i] * t.wTab[i]
	}

	// re scaling: kernel is scaled, so we don't scale here
	// and scaling factor of kernel is incorporated in m
	// multiplication below.
	t.ft.Scale(false)
	e := t.ft.Do(src)
	t.ft.Scale(true)
	if e != nil {
		panic(fmt.Sprintf("%s", e))
	}

	for i := 0; i < t.nPad; i++ {
		src[i] *= t.kern[i]
	}

	t.ft.Inv(src)

	// output step
	r := complex(1/math.Sqrt(float64(t.nS)), 0)
	for i := 0; i < t.nB; i++ {
		src[i] *= t.wTab[i] * r
	}
	return src[:t.nB]
}

// NB returns the number of frequency bins produced by the transform
func (t *T) NB() int {
	return t.nB
}

// NS returns the number of samples expected by the transform.
func (t *T) NS() int {
	return t.nS
}

// PadN returns the underlying fft pad size.
func (t *T) PadN() int {
	return t.nPad
}

// FreqRange returns the range of frequencies for which
// coefficients are computed.
func (t *T) FreqRange(sf freq.T) (l freq.T, u freq.T) {
	l = sf.FreqOf(t.start)
	d := sf.FreqOf(t.step * float64(t.nB))
	return l, l + d
}

// FreqStep returns the difference between the
// center frequencies of any two adjacent frequency bins.
func (t *T) FreqStep(sf freq.T) freq.T {
	l, u := t.FreqRange(sf)
	return (u - l) / freq.T(t.nB)
}

// BinRange returns the frequency range of a given frequency bin
// for which a coefficient can be computed.
func (t *T) BinRange(sf freq.T, i int) (l freq.T, u freq.T) {
	L, _ := t.FreqRange(sf)
	step := t.FreqStep(sf)
	return L + step*freq.T(i), L + step*freq.T(i+1)
}

// Win returns a window with appropriate capacity and length
// with the contents of c.
func (t *T) Win(c []complex128) []complex128 {
	if cap(c) < t.nPad {
		tmp := make([]complex128, t.nPad)
		copy(tmp, c)
		c = tmp
	}
	return c[:t.nS]
}

func (t *T) initA() {
	for i := 0; i < t.nS; i++ {
		t.aTab[i] = cmplx.Exp(complex(0, -float64(i)*t.start))
	}
}

func (t *T) initWK() {
	for i := 0; i < t.nS; i++ {
		c := cmplx.Exp(complex(0, t.step*(float64(i*i)/2)))
		t.kern[i] = c
		t.wTab[i] = cmplx.Conj(c)
	}
	// wrap it for special properties of kernel def of convolution
	// (k[-n] = k[n] unlike linear convolution)
	for i := 1; i < t.nS; i++ {
		t.kern[t.nPad-i] = t.kern[i]
	}
	// nb: scaled.
	t.ft.Do(t.kern)
}

func pad(d []complex128, n int) []complex128 {
	m := len(d)
	d = d[:n]
	for i := m; i < n; i++ {
		d[i] = 0i
	}
	return d
}

func findL(nS, nB int) int {
	//L := nS + nB - 1
	L := 2*nS - 1
	res := 1
	for res < L {
		res *= 2
	}
	return res
}
