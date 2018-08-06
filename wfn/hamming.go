// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package wfn

import "math"

func Hamming(i float64) float64 {
	return 0.53836 - 0.46164*math.Cos(math.Pi-i)
}
