// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package convol

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func rev(c []float64) {
	n := len(c)
	h := n / 2
	j := 0
	for i := n - 1; i > h; i-- {
		c[i], c[j] = c[j], c[i]
		j++
	}
}

func direct(a, b []float64) []float64 {
	if len(a) < len(b) {
		a, b = b, a
	}
	L := len(a) + len(b) - 1
	dst := make([]float64, L)
	var i, j, k int

	for i = 0; i < len(dst); i++ {
		ttl := 0.0
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

func approxEq(a, b []float64, eps float64) int {
	if len(a) != len(b) {
		panic(fmt.Sprintf("cannot compare equality of diff length vectors: %d, %d\n", len(a), len(b)))
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > eps {
			return i
		}
	}
	return -1
}

func gen() ([]float64, []float64) {
	n := rand.Intn(7) + 1
	m := rand.Intn(4) + 1
	if n%2 == 1 {
		n++
	}
	if m%2 == 1 {
		m++
	}
	bufA := make([]float64, m)
	bufB := make([]float64, n)
	for i := range bufA {
		a := float64(rand.Intn(10))
		bufA[i] = a
	}
	for i := range bufB {
		b := float64(rand.Intn(10))
		bufB[i] = b
	}

	return bufA, bufB
}

func TestDirect(t *testing.T) {
	for i := 0; i < 64; i++ {
		a, b := gen()
		d := direct(a, b)
		o := To(nil, a, b)
		if k := approxEq(d, o, 0.001); k != -1 {
			t.Errorf("[%d] %v (*) %v @%d: %.2f v %.2f\n", i, a, b, k, d[k], o[k])
		}
	}
}
