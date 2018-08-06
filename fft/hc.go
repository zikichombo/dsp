// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

// HalfComplex is a format for storing complex spectra of
// real data of length N in a []float64 of length N.
// Such spectra have Hermitian symmetry.
//
// Since the spectrum is so symmetric, all the information
// will fit, provided we reflect around the points correctly.
//
// The format is (like in fftw):
//
//  r0, r1, ..., r_{n/2}, i_{(n+1)/2 - 1}, ..., i2, i1
//
// for complex values rn + in, n an integer in [0...N/2]
//
// Some things to note
//
//  - Due to the symmetry, i0 is always 0 and i_{n/2} is always 0 if N is even.
//
//  - The number of reals is 2 greater than the number of imaginaries when N is even
//
//  - The number of reals is 1 greater than the number of imaginaries when N is odd
//    because i_{n/2} doesn't exist.
//
type HalfComplex []float64

// Cmplx returns the complex128 representation of element i.
func (h HalfComplex) Cmplx(i int) complex128 {
	N := len(h)
	if i == 0 || 2*i == N {
		return complex(h[i], 0.0)
	}
	return complex(h[i], h[N-i])
}

// Sets the complex number i to c in h.
func (h HalfComplex) SetCmplx(i int, c complex128) {
	N := len(h)
	h[i] = real(c)
	if i == 0 || 2*i == N {
		return
	}
	h[N-i] = imag(c)
}

// Get the real part of the complex number at i.
func (h HalfComplex) Real(i int) float64 {
	return h[i]
}

// SetReal sets the real part of the complex number at i.
func (h HalfComplex) SetReal(i int, v float64) {
	h[i] = v
}

// Get the imaginary part of the complex number at i.
func (h HalfComplex) Imag(i int) float64 {
	N := len(h)
	if i == 0 || 2*i == N {
		return 0
	}
	return h[N-i]
}

// SetImag sets the imaginary part of the complex number at i to v.
// Since all imaginary parts at complex number 0 and h.Len()/2 are
// 0, if i == 0 or h.Len()/2, then SetImag is a no-op.
func (h HalfComplex) SetImag(i int, v float64) {
	N := len(h)
	if i == 0 || 2*i == N {
		return
	}
	h[N-i] = v
}

// Len returns the number of complex numbers in h.
func (h HalfComplex) Len() int {
	n := len(h)
	return n/2 + 1
}

// Dot computes the complex dot-product (elementwise multiplication) of a and b, placing
// the result in a and returning it. Dot panics if a.Len() != b.Len().
func (a HalfComplex) Dot(b HalfComplex) HalfComplex {
	var ca, cb complex128
	N := a.Len()
	if b.Len() != N {
		panic("size mismatch")
	}
	for i := 0; i < N; i++ {
		ca = a.Cmplx(i)
		cb = b.Cmplx(i)
		a.SetCmplx(i, ca*cb)
	}
	return a
}
