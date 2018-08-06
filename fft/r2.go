// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

// in place, radix 2 decimation in time Cooley-Tuckey
//
// len(d) should be power of 2
//
func r2(d []complex128, tw *twiddles, sc bool) {
	N := len(d)
	M := 2
	if !is2pow(N) {
		panic("length not power of 2")
	}
	if N == 1 {
		return
	}
	if N == 2 {
		// problem with twiddles or something, treat explicitly.
		e := d[0]
		co := d[1]
		d[0] = e + co
		d[1] = e - co
		if sc {
			scale(d)
		}
		return
	}
	revBinPermute(d)
	var q, r, a, b, H int
	var c, e, co complex128
	for M <= N {
		H = M / 2
		for q = 0; q < N; q += M {
			for r = 0; r < H; r++ {
				a = q + r
				b = a + H
				c = tw.cmplxQ(r, M)
				e, co = d[a], d[b]*c
				d[a], d[b] = e+co, e-co
			}
		}
		M *= 2
	}
	if sc {
		scale(d)
	}
}
