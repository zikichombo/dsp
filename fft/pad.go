// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

func pad(d []complex128, toLen int) []complex128 {
	n := len(d)
	c := cap(d)
	if n >= toLen {
		return d[:toLen]
	}
	if c < toLen {
		t := make([]complex128, toLen)
		copy(t, d)
		return t
	}
	d = d[:toLen]
	for i := n; i < toLen; i++ {
		d[i] = 0i
	}
	return d
}
