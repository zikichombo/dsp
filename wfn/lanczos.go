// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package wfn

func LanczosItp(a int, t float64) float64 {
	return Sinc(t) * Sinc(t/float64(a))
}

func LanczosItpFn(a int) func(float64) float64 {
	return func(t float64) float64 {
		return LanczosItp(a, t)
	}
}

func LanczosWin(a int, t float64) float64 {
	return Sinc(t / float64(a))
}
