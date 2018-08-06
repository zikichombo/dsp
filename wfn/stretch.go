// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package wfn

func Stretch(fn func(float64) float64, by float64) func(float64) float64 {
	return func(i float64) float64 {
		return fn(by * i)
	}
}
