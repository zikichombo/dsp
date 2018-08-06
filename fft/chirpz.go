// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"
	"math/cmplx"
)

// chirps are extra factors used in bluesteins algorithm
// of the form e^{i*pi*k*k/N} with k in range 0...padN
// i sqrt(-1).
type chirpz struct {
	D  []complex128
	tD []complex128
}

func newChirpz(n, padN int, twids *twiddles) *chirpz {
	res := &chirpz{
		D:  make([]complex128, padN),
		tD: make([]complex128, padN)}
	N := float64(n)
	M := padN
	res.D[0] = complex(1, 0) // = cos(0) i*sin(0)
	for i := 1; i < n; i++ {
		r := complex(0, -math.Pi*float64(i*i)/N)
		c := cmplx.Exp(r)
		if !twids.inv {
			c = cmplx.Conj(c)
		}
		res.D[i] = c
		res.D[M-i] = c
	}
	for i := n; i < M-n; i++ {
		res.D[i] = 0.0
	}
	copy(res.tD, res.D)
	r2(res.tD, twids, true)
	return res
}

func (c *chirpz) inv(i int) complex128 {
	return cmplx.Conj(c.D[i])
}
