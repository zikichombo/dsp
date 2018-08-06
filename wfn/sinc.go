// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package wfn

import "math"

func Sinc(x float64) float64 {
	return SincEps(x, 1e-10)
}

func SincEps(x, eps float64) float64 {
	if math.Abs(x) < eps {
		return 1
	}
	return math.Sin(math.Pi*x) / (math.Pi * x)
}
