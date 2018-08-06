// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"math"
	"math/cmplx"
)

// Type Real computes an FFT for a Real only data.
type Real struct {
	ft     *T
	n      int
	cBuf   []complex128
	twidz  []complex128
	scaler float64
}

// NewReal creates a new FFT transformer for
// float data of length n.
func NewReal(n int) *Real {
	if n%2 != 0 {
		panic("Real only works for even n")
	}
	res := &Real{
		ft: New(n / 2)}
	res.cBuf = res.ft.Win(nil)
	res.ft.Scale(false)
	twidz := make([]complex128, n/2)
	N := float64(n)
	for i := range twidz {
		s, c := math.Sincos(float64(i) * 2.0 * math.Pi / N)
		twidz[i] = complex(c, -s)
	}
	res.twidz = twidz
	res.scaler = 1.0 / math.Sqrt(N)
	res.n = n
	return res
}

// Do performs a DFT on real data in d.
//
// d must be the size specified in the call to NewReal()
// which created r, or Do will panic.
//
// Do operates in place, overwriting d.  Do returns
// d overwritten as (i.e. cast to) a HalfComplex object.
//
func (r *Real) Do(d []float64) HalfComplex {
	h := len(r.cBuf)
	N := 2 * h

	if len(d) != N {
		panic(fmt.Sprintf("size mismatch got %d not %d", len(d), r.n))
	}
	r.pack(d)
	r.ft.Do(r.cBuf)
	hc := r.toHC(d)
	if r.scaler == 1.0 {
		return hc
	}
	for i := range hc {
		hc[i] *= r.scaler
	}
	return hc
}

// Inv computes the inverse transform of a real data
// from a HalfComplex object.
//
// Inv operates in place but returns the same data as hc, cast to
// a []float64.
func (r *Real) Inv(hc HalfComplex) []float64 {
	if r.scaler != 1.0 {
		for i := range hc {
			hc[i] /= r.scaler
		}
	}
	r.fromHC(hc)
	r.ft.Inv(r.cBuf)
	res := []float64(hc)
	r.unpack(res)
	if r.scaler != 1.0 {
		for i := range res {
			res[i] /= float64(len(r.cBuf))
		}
	}
	return res
}

// Scale sets whether or not r scales the transform results.
// Scale returns whether or not r was configured to scale
// the transform results prior to calling Scale.
func (r *Real) Scale(v bool) bool {
	res := r.scaler != 1.0
	if !v {
		r.scaler = 1.0
	} else {
		r.scaler = 1.0 / math.Sqrt(float64(r.n))
	}
	return res
}

// N returns the length of the arguments to the transform
// implemented by r.
func (r *Real) N() int {
	return r.n
}

// Translate r.cBuf into halfcomplex in d.
//
// Method from The DSP Book.
// some bugs fixed (wrong sign of
// twiddles, lack of appropriate DC scaling, index off by one error).
//
// http://dsp-book.narod.ru/FFTBB/0270_PDF_C14.pdf
//
func (r *Real) toHC(d []float64) HalfComplex {
	const (
		halfR = complex(0.5, 0)
		halfI = complex(0, 0.5)
	)
	N := len(d)
	h := N / 2
	res := HalfComplex(d)
	cb := r.cBuf

	if N != 0 {
		a := cb[0]
		f0 := halfR * a
		g0 := -halfI * a
		shift := r.twidz[0]
		res.SetCmplx(0, complex(2.0, 0.)*(f0+shift*g0))
		res[h] = 2.0 * real(f0-g0)
	}

	twidz := r.twidz
	for i := 1; i < h; i++ {
		a, b := cb[i], cmplx.Conj(cb[h-i])
		fi := halfR * (a + b)
		gi := halfI * (b - a)
		shift := twidz[i]
		xi := fi + shift*gi
		res.SetCmplx(i, xi)
	}
	return res
}

func (r *Real) fromHC(hc HalfComplex) {
	// Method derived in part from literature and in part trial and error.
	// We derive backwards fi, gi from Do, and from that derive
	// the original complex value output from the forward fft.
	const (
		halfR = complex(0.5, 0)
		halfI = complex(0, 0.5)
	)
	h := len(r.cBuf)
	if len(hc) != h*2 {
		panic("invalid HalfComplex length")
	}
	if len(hc) > 0 {
		a, b := hc.Cmplx(0), 0i
		f := halfR * (a + b)
		g := cmplx.Conj(r.twidz[0]) * (a - f)
		r.cBuf[0] = halfR * (f/halfR - g/halfI)
		ny := 0.5 * hc[h]
		r.cBuf[0] = complex(real(r.cBuf[0])+ny, imag(r.cBuf[0])-ny)
	}
	var j int
	for i := 1; i < h; i++ {
		j = h - i
		a, b := hc.Cmplx(i), cmplx.Conj(hc.Cmplx(j))
		fi := halfR * (a + b)
		gi := cmplx.Conj(r.twidz[i]) * (a - fi)
		r.cBuf[i] = halfR * (fi/halfR - gi/halfI)
	}
}

func (r *Real) pack(d []float64) {
	cb := r.cBuf
	for i := range cb {
		cb[i] = complex(d[2*i], d[2*i+1])
	}
}

func (r *Real) unpack(d []float64) {
	cb := r.cBuf
	var re, im float64
	for i, v := range cb {
		re = real(v)
		im = imag(v)
		d[2*i] = re
		d[2*i+1] = im
	}
}
