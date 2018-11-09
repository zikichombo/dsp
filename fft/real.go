// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"
	"math/cmplx"
)

// Real computes an FFT for a Real only data.
//
// For even length transforms, the implementation
// uses a complex FFT of size N/2 and some pre/post processing.
// For odd length transforms, the implementtion uses a complex
// FFT of size N.
type Real struct {
	ft     *T           // half sized for even
	n      int          //
	cBuf   []complex128 //
	twidz  []complex128 // only for even
	scaler float64      // only for even
}

// NewReal creates a new FFT transformer for
// float data of length n.
func NewReal(n int) *Real {
	if n&1 == 0 {
		m := n / 2
		res := &Real{ft: New(m)}
		res.n = n
		res.cBuf = res.ft.Win(nil)
		res.ft.Scale(false)
		twidz := make([]complex128, m)
		N := float64(n)
		for i := range twidz {
			s, c := math.Sincos(float64(i) * 2.0 * math.Pi / N)
			twidz[i] = complex(c, -s)
		}
		res.twidz = twidz
		res.scaler = 1.0 / math.Sqrt(N)
		return res
	}
	res := &Real{ft: New(n)}
	res.n = n
	res.cBuf = res.ft.Win(nil)
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
	if r.n&1 == 0 {
		return r.evenDo(d)
	}
	return r.oddDo(d)
}

func (r *Real) evenDo(d []float64) HalfComplex {
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

func (r *Real) oddDo(d []float64) HalfComplex {
	for i, v := range d {
		r.cBuf[i] = complex(v, 0.0)
	}
	r.ft.Do(r.cBuf)
	res := HalfComplex(d)
	res.FromCmplx(r.cBuf)
	return res
}

// Inv computes the inverse transform of a real data
// from a HalfComplex object.
//
// Inv operates in place but returns the same data as hc, cast to
// a []float64.
func (r *Real) Inv(hc HalfComplex) []float64 {
	if r.n&1 == 0 {
		return r.evenInv(hc)
	}
	return r.oddInv(hc)
}

func (r *Real) evenInv(hc HalfComplex) []float64 {
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
		sc := 1.0 / float64(len(r.cBuf))
		for i := range res {
			res[i] *= sc
		}
	}
	return res
}

func (r *Real) oddInv(hc HalfComplex) []float64 {
	hc.ToCmplx(r.cBuf)
	r.ft.Inv(r.cBuf)
	for i, c := range r.cBuf {
		hc[i] = real(c)
	}
	return []float64(hc)
}

// Scale sets whether or not r scales the transform results.
// Scale returns whether or not r was configured to scale
// the transform results prior to calling Scale.
func (r *Real) Scale(v bool) {
	if r.n&1 == 0 {
		if !v {
			r.scaler = 1.0
		} else {
			r.scaler = 1.0 / math.Sqrt(float64(r.n))
		}
		return
	}
	r.ft.Scale(v)
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
	if N != r.n {
		panic("invalid length")
	}
	res := HalfComplex(d)
	if N == 0 {
		return res
	}
	cb := r.cBuf
	h := len(cb)

	a := cb[0]
	f0 := halfR * a
	g0 := -halfI * a
	shift := r.twidz[0]
	res.SetCmplx(0, complex(2.0, 0.0)*(f0+shift*g0))
	if h+h == N {
		res[h] = 2.0 * real(f0-g0)
	}

	for i := 1; i < h; i++ {
		a, b := cb[i], cmplx.Conj(cb[h-i])
		fi := halfR * (a + b)
		gi := halfI * (b - a)
		shift := r.twidz[i]
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
	N := len(hc)
	if N != r.n {
		panic("invalid HalfComplex length")
	}
	if N == 0 {
		return
	}
	h := len(r.cBuf)
	a, b := hc.Cmplx(0), 0i
	f := halfR * (a + b)
	g := cmplx.Conj(r.twidz[0]) * (a - f)
	r.cBuf[0] = halfR * (f/halfR - g/halfI)
	var ny float64
	if h+h == N {
		ny = 0.5 * hc[h]
	} else {
		ny = 0.5 * hc[h-1]
	}
	r.cBuf[0] = complex(real(r.cBuf[0])+ny, imag(r.cBuf[0])-ny)
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
	end := len(cb)
	for i := 0; i < end; i++ {
		cb[i] = complex(d[2*i], d[2*i+1])
	}
}

func (r *Real) unpack(d []float64) {
	cb := r.cBuf
	end := len(cb)
	var re, im float64
	var v complex128
	for i := 0; i < end; i++ {
		v = cb[i]
		re = real(v)
		im = imag(v)
		d[2*i] = re
		d[2*i+1] = im
	}
}
