// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

// ZeroPad returns a zero-padding of the Fourier Transform d for time interpolation.
// If we divide a transform into a 0(dc) component, components not exceeding the
// Nyquist limit, and higher frequencies as "negative" frequencies, then the zero
// padding goes in-between the non-negative frequencies and the negative frequencies.
// The inverse FT will then generate  time-interpolation of the data for which
// d is a FT.
func ZeroPad(d []complex128, n int) []complex128 {
	return ZeroPadTo(nil, d, n)
}

// ZeroPadTo is like ZeroPad but allows specifying a destination vector.
// If dst doesn't have sufficient capacity, then a new one is created
// and returned.
func ZeroPadTo(dst, src []complex128, n int) []complex128 {
	if cap(dst) < len(src)+n {
		dst = make([]complex128, len(src)+n)
	} else {
		dst = dst[:len(src)+n]
	}
	l := len(src)
	m := l / 2
	if n%2 == 0 {
		m++
	}
	for i := l; i < l+n; i++ {
		dst[i] = 0i
	}
	copy(dst[:m], src[:m])
	copy(dst[m+n:], src[m:])
	return dst
}
