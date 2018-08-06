// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"
	"math/rand"
	"testing"
)

func TestFactor(t *testing.T) {
	n := 1024 * 7
	a, _, p2 := factor(n)
	if p2 != 1024 {
		t.Errorf("didn't find a factor of 2")
	}
	//p := 5915587277
	p := 101483
	p = 100000000019
	p = 10000169
	a, b, _ := factor(p)
	if a != 1 || b != p {
		t.Errorf("couldn't find big prime: %d\n", p)
	}
	for _, sq := range []int{23 * 23, 37 * 37, 83 * 83, 127 * 127} {
		a, b, _ = factor(sq)
		if a != b || a != int(math.Sqrt(float64(sq))) {
			t.Errorf("square %d not found correctly", sq)
		}
	}
	for i := 0; i < 1024; i++ {
		r := rand.Intn(100000) + 100000
		a, b, p2 := factor(r)
		if a*b*p2 != r {
			t.Errorf("%d * %d * %d != %d\n", a, b, p2, r)
		}
		n, aa, tt := factor(a)
		if n != 1 || aa != a || tt != 1 {
			t.Errorf("should have returned prime for a but got %d = %d * %d * %d while factoring %d (=%d*%d*%d)\n", a, n, aa, tt, r, a, b, p2)
		}
	}
}

func BenchmarkFactor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(100000) + 100000
		factor(r)
	}
}
