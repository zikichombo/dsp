// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package convol

import "math"

func unscale(d []float64) float64 {
	n := float64(len(d))
	r := math.Sqrt(n)
	return r
}
