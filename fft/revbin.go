// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

func revBin(i int, N uint) int {
	if i > (1 << (N + 1)) {
		panic("i too big")
	}
	r := 0
	N--
	for i != 0 {
		r |= (i % 2) << N
		N--
		i /= 2
	}
	return r
}

func revBinLim(i int) uint {
	N := uint(0)
	for i > (1 << N) {
		N++
	}
	return N
}

func revBinPermute(d []complex128) {
	m := len(d)
	N := revBinLim(m)
	for i := range d {
		j := revBin(i, N)
		if j < i {
			d[i], d[j] = d[j], d[i]
		}
	}
}
