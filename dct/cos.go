// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import (
	"math"
)

const cosTblBits = 10

var cosTbl [][]float64

func init() {
	cosTbl = make([][]float64, cosTblBits+1)
	i := uint(0)
	var n int
	for {
		n = 1 << i
		cosTbl[i] = genCos(i)
		if n >= 1<<cosTblBits {
			break
		}
		i++
	}
}

func genCos(i uint) []float64 {
	n := 1 << i
	sl := make([]float64, n)
	pion := math.Pi / float64(n)
	for j := 0; j < n; j++ {
		sl[j] = math.Cos((float64(j) + 0.5) * pion)
	}
	return sl
}
