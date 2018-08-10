// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestTNaive(t *testing.T) {
	N := 1024
	Iters := 8
	d := make([]float64, N)
	tmp := make([]float64, N)
	ct := New(N)
	for i := 0; i < Iters; i++ {
		for i := 0; i < N; i++ {
			d[i] = rand.Float64()
		}
		copy(tmp, d)
		ct.Do(d)
		Naive(tmp)
		for i, v := range d {
			if math.Abs(v-tmp[i]) > 1e-12 {
				t.Errorf("%d: Lee %f Naive %f\n", i, v, tmp[i])
			}
		}
	}
}

func TestCmp(t *testing.T) {
	d := []float64{-0.999984, -0.736924, 0.511211, -0.082700}
	dct := New(len(d))
	dct.Do(d)
}

func TestTCos(t *testing.T) {
	N := 16
	d := make([]float64, N)
	rps := math.Pi / 16
	for i := 0; i < N; i++ {
		d[i] = math.Cos(float64(i) * rps)
	}
	fmt.Printf("in:  %v\n", d)
	dct := New(N)
	dct.Do(d)
	fmt.Printf("dct: %v\n", d)
	dct.Inv(d)
	fmt.Printf("inv: %v\n", d)
}
