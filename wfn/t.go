// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package wfn

import "math"

// T describes a concrete window function.
type T []float64

// Given a function f with domain from [-pi..pi), return a window
// dividing [-pi..pi) into n values.
func New(f func(float64) float64, n int) T {
	m := float64(n - 1)
	r := 2 * math.Pi / m
	h := m / 2
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		fi := float64(i)
		res[i] = f(r * (fi - h))
	}
	return T(res)
}

// Apply applies the window t to the data t.
// Apply panics if len(d) > len(t).
func (t T) Apply(d []float64) {
	for i := range d {
		d[i] *= t[i]
	}
}

// DcGain returns the power of constant value 1 signal in
func (t T) DcGain() float64 {
	ttl := 0.0
	for _, v := range t {
		ttl += math.Abs(v)
	}
	return ttl / float64(len(t))
}

// AvGain computes gain which is the ratio of sum of coeficients in vs
// to its length.  Since vs can contain negative and zero coeficients,
// AvGain can give a zero result and normalizing accordingly is risky.
func (t T) AvGain() float64 {
	ttl := 0.0
	for _, v := range t {
		ttl += v
	}
	return ttl / float64(len(t))
}

// DcNorm normalizes vs according to unity DcGain
func (t T) DcNorm() {
	g := t.DcGain() * float64(len(t))
	for i := range t {
		t[i] /= g
	}
}

// AvNorm normalizes vs according to unity AvGain.
func (t T) AvNorm() {
	g := t.AvGain() * float64(len(t))
	for i := range t {
		t[i] /= g
	}
}
