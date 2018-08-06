// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package qitp

import "math"

// Abc takes three y values denoting points (-1, p), (0, q), (1, r)
// and returns the coeficients a, b, c such that
//
//  y = f(x) = ax^2 + bx + c
func Abc(p, q, r float64) (a, b, c float64) {
	// c = q by assumption of x values of p,q,r.
	c = q
	// 1) r = a + b + q <=> 1a) r - q = a + b
	// 2) p = a - b + q <=> 2a) p - q = a - b
	//
	// add 1a) + 2a) gives r + p - 2*q = 2*a
	a = 0.5*(r+p) - q

	// subst a in 1)
	// r =  a + b + q  <=> b = r - a - q
	// rearrange to p - 2q + b + q = 0 <=> b = q - p
	b = r - a - q
	return
}

// translates standard to vertex form y = a(x - h)^2 + k
func Abc2Hk(a, b, c float64) (h, k float64) {
	h = -b / (2 * a)
	k = a*h*h + b*h + c
	return
}

func Ahk2Bc(a, h, k float64) (b, c float64) {
	b = -2 * h * a
	c = a*h*h + k
	return
}

// AbcX returns the point interpolated with x in [-1,1]
func AbcX(p, q, r, x float64) float64 {
	if x < -1+1e-10 || x > 1-1e-10 {
		panic("x oob")
	}
	a, b, c := Abc(p, q, r)
	return a*x*x + b*x + c
}

// ParabX0 returns a function which can be queried for
// values in the range [-1,1] returning interpolated values
// from p,q,r
func ParabX0(p, q, r float64) func(float64) float64 {
	a, b, c := Abc(p, q, r)
	return func(x float64) float64 {
		return a*x*x + b*x + c
	}
}

// Slice interpolates the index x in vs quadratically
// if possible and linearly if len(vs) == 2.
//
// Slice panics if
//
// - len(vs) <= 1
//
// - x is out of bounds
func Slice(vs []float64, x float64) float64 {
	return SliceMap(vs, x, func(v float64) float64 { return v })
}

// SliceMap is like slice but maps the values in vs using m(v)
// before interpolation.  The returned value is by extension also mapped.
//
// Mapping is provided since the hard part is finding the indices.
func SliceMap(vs []float64, x float64, m func(float64) float64) float64 {
	if len(vs) <= 1 {
		panic("not big nuf")
	}
	xi, xf := math.Modf(x)
	c := int(xi)
	if c < 0 || c >= len(vs) || (c == len(vs)-1 && xf > 1e-10) {
		panic("oob")
	}
	if c == len(vs)-1 {
		return m(vs[c])
	}
	if len(vs) == 2 { // back of to linear.
		return (1-xf)*m(vs[0]) + xf*m(vs[1])
	}
	if xf < 1e-10 {
		return m(vs[c])
	}
	if 1-xf < 1e-10 {
		return m(vs[c+1])
	}
	r := c + 1
	l := c - 1
	if r+1 < len(vs) && (l < 0 || xf >= 0.5) {
		l, c, r = c, r, r+1
		// x is in [l..c)
		return AbcX(m(vs[l]), m(vs[c]), m(vs[r]), -1+xf)
	}
	// x is in (c,r)
	return AbcX(m(vs[l]), m(vs[c]), m(vs[r]), xf)
}
