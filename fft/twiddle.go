// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"math"
	"strings"
)

// implementation of twiddle caching for FFT algorithms

// Reminder:  e^{ix} = cos(x) + i sin(x), i == sqrt(-1)
// and e^{-ix} = cos(x) - i sin(x)
//
// so e^{ix} is complex conjugate of e^{-ix}
//

type twiddles struct {
	cosTbl []float64
	sinTbl []float64
	twoPi  int
	inv    bool
	invSin float64
}

func newTwiddles(n int, inv bool) *twiddles {
	tbl := make([]float64, 2*n)
	res := &twiddles{
		cosTbl: tbl[:n],
		sinTbl: tbl[n:],
		twoPi:  n}
	w := 2.0 * math.Pi / float64(n)
	for i := 0; i < n; i++ {
		s, c := math.Sincos(float64(i) * w)
		res.cosTbl[i] = c
		res.sinTbl[i] = s
	}
	res.inv = inv
	if !inv {
		res.invSin = -1.0
	} else {
		res.invSin = 1.0
	}
	return res
}

func (t *twiddles) sincos(i int) (float64, float64) {
	return t.sin(i), t.cos(i)
}

func (t *twiddles) sincosQ(i, q int) (float64, float64) {
	j := i * t.twoPi / q
	return t.sin(j), t.cos(j)
}

func (t *twiddles) cmplx(i int) complex128 {
	im, re := t.sincos(i)
	return complex(re, im)
}

func (t *twiddles) cmplxQ(i, q int) complex128 {
	j := i * t.twoPi / q
	return t.cmplx(j)
}

func (t *twiddles) sin(i int) float64 {
	return t.invSin * t.sinTbl[i%t.twoPi]
}

// nb this is a hot spot on profiling.
func (t *twiddles) cos(i int) float64 {
	return t.cosTbl[i%t.twoPi]
}

func (t *twiddles) String() string {
	res := make([]string, 0, len(t.cosTbl))
	for _, d := range t.cosTbl {
		res = append(res, fmt.Sprintf("%.2f", d))
	}
	return strings.Join(res, " ")
}
