// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"
	"math/cmplx"
)

// R2Size returns the least power of 2 not less than
// n, the size for which fourier transforms use the
// fast radix-2 algorithm.
func R2Size(n int) int {
	j := 1
	for j < n {
		j *= 2
	}
	return j
}

// Do performs an in-place FFT on d, if possible.
//
// The input samples d must have a cap....blah.
func Do(d []complex128) {
	// XXX cap dimensions can break this.
	t := New(len(d))
	wd := t.Win(d)
	t.Do(wd)
	if &d[0] != &wd[0] {
		copy(d, wd)
	}
}

// To performs a FFT on src placing the results
// in dst if dst has approprioate capacity returning
// either dst or the results in a new slice.
func To(dst, src []complex128) []complex128 {
	t := New(len(src))
	dst = t.ensureSrc(dst)
	t.To(dst, src)
	return dst
}

// Inv is like Do but for inverse transform
func Inv(d []complex128) {
	t := New(len(d))
	t.Inv(d)
}

// InvTo is like To but for inverse transform
func InvTo(dst, src []complex128) []complex128 {
	t := New(len(src))
	dst = t.ensureSrc(dst)
	t.InvTo(dst, src)
	return dst
}

// Dilate changes the frequency basis by a factor of n/m.
// This is a pitch shift relative to the quantized
// frequency domain, hence there is informaion loss,
// and some values may just be clobbered...
func Dilate(d []complex128, p, q int) {
	if p == q {
		return
	}
	n := len(d)
	h := n / 2
	if p > q {
		for i := h; i >= 1; i-- {
			dst := (i * p) / q
			if dst > h {
				d[i] = 0i
				continue
			}
			v := d[i]
			d[i] = 0i
			d[dst] += v
		}
		for i := n - 1; i > h; i-- {
			dst := (i * p) / q
			if dst >= n {
				d[i] = 0i
				continue
			}
			v := d[i]
			d[i] = 0i
			d[dst] += v
		}
		return
	}
	for i := 1; i <= h; i++ {
		dst := (i * p) / q
		if dst > h {
			d[i] = 0i
			continue
		}
		v := d[i]
		d[i] = 0i
		d[dst] += v
	}
	for i := h + 1; i < n; i++ {
		dst := (i * p) / q
		if dst > n {
			d[i] = 0i
			continue
		}
		v := d[i]
		d[i] = 0i
		d[dst] += v
	}
}

// older stuff

// radix 2, in place recursive
func radix2(src []float64, dst []complex128, N, stride int, inv float64) {
	if N == 1 {
		dst[0] = complex(src[0], 0.0)
		return
	}
	h := N / 2
	radix2(src, dst[:h], h, stride*2, inv)
	radix2(src[stride:], dst[h:], h, stride*2, inv)

	// recombine
	for i := 0; i < h; i++ {
		ed := dst[i]
		od := dst[i+h]
		exp := complex(0, inv*math.Pi*2.0*float64(i)/float64(N))
		c := cmplx.Exp(exp)
		dst[i] = ed + c*od
		dst[i+h] = ed - c*od
	}
}

// P is prime.

// TBD(wsc) Rader?
func slowft(src []float64, dst []complex128, P, stride int, inv float64) {
	if P == 1 {
		dst[0] = complex(src[0], 0)
		return
	}
	var a, s, fi float64
	var c complex128

	lim := P * stride
	n := float64(len(dst))
	A := inv * 2.0 * math.Pi / n
	for i := 0; i < len(dst); i++ {
		fi = float64(i)
		c = 0i
		for j := 0; j < lim; j += stride {
			a = A * float64(j/stride) * fi
			s = src[j]
			c += complex(s*math.Cos(a), s*math.Sin(a))
		}
		dst[i] = c
	}
}
