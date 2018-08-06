// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import (
	"fmt"
	"math"
	"testing"
)

func TestNaive(t *testing.T) {
	N := 16
	d := make([]float64, N)
	rps := math.Pi / 16
	for i := 0; i < N; i++ {
		d[i] = math.Cos(float64(i) * rps)
	}
	p := 0.0
	for _, v := range d {
		p += v * v
	}
	fmt.Printf("in: power %f %v\n", p, d)
	Naive(d)
	p = 0.0
	for _, v := range d {
		p += v * v
	}
	fmt.Printf("dct: power %f %v\n", p, d)
	NaiveInv(d)
	fmt.Printf("inv: %v\n", d)
}

func TestNaiveSinInv(t *testing.T) {
	N := 16
	d := make([]float64, N)
	rps := 2.0 * math.Pi / float64(N)
	for i := 0; i < N; i++ {
		d[i] = math.Cos(float64(i) * rps)
	}
	p := 0.0
	for _, v := range d {
		p += v * v
	}
	fmt.Printf("in: power %f %v\n", p, d)
	Naive(d)
	p = 0.0
	for _, v := range d {
		p += v * v
	}
	fmt.Printf("dct: power %f %v\n", p, d)
	NaiveSinInv(d)
	fmt.Printf("inv: %v\n", d)
}
