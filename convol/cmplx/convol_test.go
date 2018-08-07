// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cmplx

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func rev(c []complex128) {
	n := len(c)
	h := n / 2
	j := 0
	for i := n - 1; i > h; i-- {
		c[i], c[j] = c[j], c[i]
		j++
	}
}

func direct(a, b []complex128) []complex128 {
	if len(a) < len(b) {
		a, b = b, a
	}
	L := len(a) + len(b) - 1
	dst := make([]complex128, L)
	var i, j, k int

	for i = 0; i < len(dst); i++ {
		ttl := 0i
		for j = 0; j < len(a); j++ {
			k = i - j
			if k < 0 {
				continue
			}
			if k >= len(b) {
				continue
			}
			ttl += a[j] * b[k]
		}
		dst[i] = ttl
	}
	return dst
}

func cmplxApproxEq(a, b []complex128, eps float64) int {
	if len(a) != len(b) {
		panic(fmt.Sprintf("cannot compare equality of diff length vectors: %d, %d\n", len(a), len(b)))
	}
	for i := range a {
		re, im := real(a[i]), imag(a[i])
		cre, cim := real(b[i]), imag(b[i])
		if math.Abs(re-cre) > eps {
			return i
		}
		if math.Abs(im-cim) > eps {
			return i
		}
	}
	return -1
}

func gen() ([]complex128, []complex128) {
	n := rand.Intn(7) + 1
	m := rand.Intn(4) + 1
	bufA := make([]complex128, m)
	bufB := make([]complex128, n)
	for i := range bufA {
		a := float64(rand.Intn(10))
		b := 0.0 //float64(rand.Intn(10))
		bufA[i] = complex(a, b)
	}
	for i := range bufB {
		a := float64(rand.Intn(10))
		b := 0.0 //float64(rand.Intn(10))
		bufB[i] = complex(a, b)
	}

	return bufA, bufB
}

func TestDirect(t *testing.T) {
	for i := 0; i < 1; i++ {
		a, b := gen()
		d := direct(a, b)
		o := To(a, b, nil)
		if k := cmplxApproxEq(d, o, 0.001); k != -1 {
			t.Errorf("[%d] %v (*) %v @%d: %.2f v %.2f\n", i, a, b, k, d[k], o[k])
		}
	}
}
