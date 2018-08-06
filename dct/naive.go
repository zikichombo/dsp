// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import "math"

func Naive(d []float64) {
	tmp := make([]float64, len(d))
	N := float64(len(d))
	pion := math.Pi / N
	for i := range d {
		ttl := 0.0
		fi := float64(i)
		for j := range d {
			ttl += d[j] * math.Cos((float64(j)+0.5)*fi*pion)
		}
		tmp[i] = ttl
	}
	copy(d, tmp)
	scale := 1.0 / math.Sqrt(float64(len(d)/2))
	for i := range d {
		d[i] *= scale
	}
}

func NaiveInv(d []float64) {
	tmp := make([]float64, len(d))
	N := float64(len(d))
	hz := d[0] / 2
	pion := math.Pi / N
	for i := range d {
		ttl := hz
		fi := float64(i)
		for j := 1; j < len(d); j++ {
			ttl += d[j] * math.Cos((fi+0.5)*float64(j)*pion)
		}
		tmp[i] = ttl
	}
	scale := 1.0 / math.Sqrt(float64(len(d)/2))
	h := float64(len(d) / 2)
	_ = h
	for i := range tmp {
		tmp[i] *= scale
	}
	copy(d, tmp)
}

func NaiveSinInv(d []float64) {
	tmp := make([]float64, len(d))
	N := float64(len(d))
	hz := d[0] / 2
	pion := math.Pi / N
	for i := range d {
		ttl := hz
		fi := float64(i)
		for j := 1; j < len(d); j++ {
			ttl += d[j] * math.Sin((fi+0.5)*float64(j)*pion)
		}
		tmp[i] = ttl
	}
	scale := 1.0 / math.Sqrt(float64(len(d)/2))
	h := float64(len(d) / 2)
	_ = h
	for i := range tmp {
		tmp[i] *= scale
	}
	copy(d, tmp)
}
