// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"testing"
)

type convolveCfg struct {
	a, b []complex128
	m    int
}

func (c *convolveCfg) test(t *testing.T) {
	d := c.dumb()
	dd := make([]complex128, len(d))
	//Convolve(c.a, c.b, dd)
	for i := range d {
		if d[i] != dd[i] {
			t.Errorf("convolve result mismatch at %d: %f != %f", i, dd[i], d[i])
		}
	}
}

func mod(i, j int) int {
	if i >= 0 {
		return i % j
	}
	return (j + (i % j)) % j
}

func rev(x []complex128) {
	n := len(x)
	h := n / 2
	n--

	for i := 0; i < h; i++ {
		x[i], x[n-i] = x[n-i], x[i]
	}
}

func (c *convolveCfg) dumb() []complex128 {
	res := make([]complex128, c.m)
	a, b := c.a, c.b
	if len(a) < len(b) {
		a, b = b, a
	}
	rev(b)
	n := len(a)
	m := len(b)

	fmt.Printf("indices:\n")
	for i := range res {
		v := 0i
		for j := range a {
			k := (j + n) % n
			k -= n - m
			if k < 0 {
				continue
			}
			fmt.Printf("\t*(%d) a[%d] * b[%d] aL%d bL%d\n", i, j, k, len(a), len(b))
			v += a[j] * b[k]
		}
		res[i] = v
		fmt.Printf("result: %f\n", v)
	}
	return res
}

var convolveCfgs = [...]convolveCfg{
	{[]complex128{0i}, []complex128{0i}, 1},
	{[]complex128{0i, 1i, 0i}, []complex128{0i, 1i, 0i}, 3},
	{[]complex128{1i, 11i, 0 + 2i}, []complex128{2i, 3i, 5i}, 3},
	{[]complex128{2i, 3i, 5i}, []complex128{0i, 0i, 1i}, 3},
	{[]complex128{2i, 3i, 5i}, []complex128{2i, 0i, 0i}, 3}}

func testConvolve(t *testing.T) {
	for i := range convolveCfgs {
		c := &convolveCfgs[i]
		c.test(t)
	}
}
