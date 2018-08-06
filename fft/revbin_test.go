// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math/rand"
	"testing"
)

func TestRevBin(t *testing.T) {

	L := uint(11)
	for lim := uint(0); lim < L; lim++ {
		for i := 0; i < (1 << lim); i++ {
			r := revBin(i, lim)
			j := revBin(r, lim)
			if j != i {
				t.Errorf("revBin non symmetric: %d, %d\n", i, r)
			}
			//fmt.Printf("%s r %s\n", strconv.FormatInt(int64(i), 2), strconv.FormatInt(int64(r), 2))
		}
	}
}

func testRevBinPermute(t *testing.T) {
	N := 128
	L := revBinLim(N)
	d := make([]complex128, N)
	for i := range d {
		d[i] = complex(rand.Float64(), rand.Float64())
	}
	c := make([]complex128, N)
	copy(d, c)
	revBinPermute(d)
	for i := range d {
		if c[i] != d[revBin(i, L)] {
			t.Errorf("%d %d not reversed\n", i, revBin(i, L))
		}
	}
}
