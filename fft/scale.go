// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import "math"

// this scaling function, applied to forward and inverse
// transforms makes the fwd and inv reciprocal and
// makes parsevals equation hold (sum of squares of
// input is equal to sum of squares of transform)
//
// this latter condition facilitates using this package
// in a way where frequency domain values are proportional
// to input values.
//
func scale(d []complex128) {
	n := len(d)
	if n <= 1 {
		return
	}
	m := 1 / math.Sqrt(float64(len(d)))
	a := complex(m, 0)
	for i := range d {
		d[i] *= a
	}
}
