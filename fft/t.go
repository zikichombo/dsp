// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"math/cmplx"
)

// Type T maintains state for efficient
// repeated computation on windows of
// data.
type T struct {
	n, padN int       // padN = n if n is power of 2, least power of 2 >= 2*n + 1 otherwise.
	twids   *twiddles // size padN
	chirpz  *chirpz   // nil if n == padN
	itwids  *twiddles
	ichirpz *chirpz
	scale   bool
}

// New creates a new T which can compute repeated windows more efficiently than
// repeatedly calling Do or To
//
// - n is the size of the data windows to transform
//
// - inv is whether or not to use the "inverse" transform
func New(n int) *T {
	if is2pow(n) {
		return NewT(n, n)
	} else {
		return NewT(n, 1<<log2(2*n-1))
	}
}

// we factor this out so we can test case where
// N is power of 2, and padN = 2N.
func NewT(n, padN int) *T {
	return &T{n: n, padN: padN, scale: true}
}

// Win returns a slice of data with dimensions set so no error occurs if used
// as argument to Do() or Inv(), nor if used as destination in To() or InvTo().
//
// These errors ensure not only that the input size makes sense but also the
// capacity, which might vary from length due to internal zero padding.
//
// Returned windows are of length t.N() and capacity t.Cap() and contain all
// the data in c
func (t *T) Win(c []complex128) []complex128 {
	if cap(c) < t.padN {
		tmp := make([]complex128, len(c), t.padN)
		copy(tmp, c)
		c = tmp
	}
	if len(c) < t.n {
		c = c[:t.n]
		for i := len(c); i < t.n; i++ {
			c[i] = 0i
		}
	}
	return c[:t.n]
}

// N returns the size of the transforms to be
// performed.
func (t *T) N() int {
	return t.n
}

// Cap returns the desired capacity of slices passed into
// Do() and as dst to To().
func (t *T) Cap() int {
	return t.padN
}

// To perform a FFT on src, placing the results in dst
// and leaving src untouched.  To returns dst, e
//
// A non-nil error can occur if src and dst are not aligned
// with t.  A new dst is allocated and returned if dst is nil.
func (t *T) To(dst, src []complex128) ([]complex128, error) {
	if len(src) != t.n {
		return nil, fmt.Errorf("wrong input size, got %d expected %d\n", len(src), t.n)
	}
	if dst == nil {
		dst = t.Win(nil)
	}
	if e := t.ckSrc(dst); e != nil {
		return nil, e
	}
	copy(dst, src)
	e := t.Do(dst)
	return dst, e
}

// Do performs an in-place transform on d.
func (t *T) Do(d []complex128) error {
	if e := t.ckSrc(d); e != nil {
		return e
	}
	if t.n == t.padN {
		r2(d, t.getTwids(false), t.scale)
		return nil
	}
	return t.bluestein(d, false)
}

// Inv performs an in-place inverse transform on d.
func (t *T) Inv(d []complex128) error {
	if e := t.ckSrc(d); e != nil {
		return e
	}
	if t.n == t.padN {
		r2(d, t.getTwids(true), t.scale)
		return nil
	}
	return t.bluestein(d, true)
}

// InvTo perform an inverse FFT on src, placing the results in dst
// and leaving src untouched.  To returns dst, e
//
// A non-nil error can occur if src and dst are not aligned
// with t.  A new dst is allocated and returned if dst is nil.
func (t *T) InvTo(dst, src []complex128) ([]complex128, error) {
	if len(src) != t.n {
		return nil, fmt.Errorf("wrong input size, got %d expected %d\n", len(src), t.n)
	}
	if dst == nil {
		dst = t.Win(nil)
	}
	if e := t.ckSrc(dst); e != nil {
		return nil, e
	}
	copy(dst, src)
	e := t.Inv(dst)
	return dst, e
}

func (t *T) bluestein(d []complex128, inv bool) error {
	// bluestein, other convolution arg stored pre-computed in t.chirpz.tD

	// multiply input by conjugate chirps
	iChirpz := t.getChirpz(!inv)
	for i := range d {
		c := d[i]
		d[i] = c * iChirpz.D[i]
	}

	// perform fwd dft on input zero padded (no scaling since chirps are scaled)
	d = pad(d, t.padN)
	r2(d, t.getTwids(inv), inv)

	// pointwise multiply for convolution (also scales since chirps scaled)
	chirpz := t.getChirpz(inv)
	for i := range d {
		d[i] *= chirpz.tD[i]
	}

	// inverse (scaled, since d was scaled since chirps were scaled)
	r2(d, t.getTwids(!inv), !inv)
	d = d[:t.n]
	if t.scale {
		scale(d) // scale to current len
	}

	// post processing
	for i := range d {
		d[i] *= iChirpz.D[i]
	}
	return nil
}

// Scale can be set to turn on or off scaling.  By default, scaling is on.
//
// Scaling in this fft implementation divides the results of both forward and
// backwards transforms by sqrt(t.N()), which enforces Parsevals power
// equivalence between powers input and output to transforms and enforces
// that an inverse is an inverse, rather than inverse on a different scale.
func (t *T) Scale(v bool) {
	t.scale = v
}

func (t *T) ckSrc(d []complex128) error {
	if len(d) != t.n {
		return fmt.Errorf("wrong length %d != %d", len(d), t.n)
	}
	if cap(d) < t.padN {
		return fmt.Errorf("wrong cap got %d expected %d", cap(d), t.padN)
	}
	return nil
}

func (t *T) ensureSrc(d []complex128) []complex128 {
	if d == nil {
		return t.Win(nil)
	}
	if e := t.ckSrc(d); e != nil {
		return t.Win(nil)
	}
	return d[:t.n]
}

// AutoCorr computes the (circular) autocorrelation of the input
// d, returning an error if d isn't proper size/capacity
// as in T.Do().  To obtain normal (actually tapered) autocorrelation,
// zero padding must be added.
func (t *T) AutoCorr(d []complex128) error {
	sc := t.scale
	t.scale = false
	if e := t.Do(d); e != nil {
		return e
	}
	for i, v := range d {
		d[i] *= cmplx.Conj(v)
	}
	t.Inv(d)
	t.scale = sc
	n := complex(float64(len(d)), 0)
	for i := range d {
		d[i] /= n
	}
	return nil
}

func (t *T) getTwids(inv bool) *twiddles {
	if !inv {
		if t.twids == nil {
			t.twids = newTwiddles(t.padN, false)
		}
		return t.twids
	}
	if t.itwids == nil {
		t.itwids = newTwiddles(t.padN, true)
	}
	return t.itwids
}

func (t *T) getChirpz(inv bool) *chirpz {
	if !inv {
		if t.chirpz == nil {
			t.chirpz = newChirpz(t.n, t.padN, t.getTwids(false))
		}
		return t.chirpz
	}
	if t.ichirpz == nil {
		t.ichirpz = newChirpz(t.n, t.padN, t.getTwids(true))
	}
	return t.ichirpz
}
