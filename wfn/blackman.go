// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package wfn

import "math"

func Blackman(i float64) float64 {
	return 0.42 - 0.5*math.Cos(math.Pi-i) + 0.08*math.Cos(2*(math.Pi-i))
}
