// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cmplx

import "math"

func unscale(d []complex128) complex128 {
	n := float64(len(d))
	r := math.Sqrt(n)
	return complex(r, 0)
}
