// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import "math"

// T encapsulates a DCT and its inverse.
type T struct {
	p      uint
	tmp    []float64
	cosTbl [][]float64
	scf    float64
}

// New creates a new T for transforming data of
// length n.  n must be a power of 2 or New panics.
func New(n int) *T {
	p := uint(0)
	for 1<<p < n {
		p++
	}
	if n != 1<<p {
		panic("not a power of 2")
	}
	var ct [][]float64
	if p >= uint(len(cosTbl)) {
		ct = make([][]float64, p+1)
		for i := range cosTbl {
			ct[i] = cosTbl[i]
		}
		for i := uint(len(cosTbl)); i <= p; i++ {
			ct[i] = genCos(i)
		}
	} else {
		ct = cosTbl
	}
	res := &T{tmp: make([]float64, n), cosTbl: ct, p: p}
	res.scf = math.Sqrt(1.0 / float64(n/2))
	return res
}

// Do performs "the" dct (type II dct) on d in place, with scaling.
//
// If |d| is not the size indicated in the call to New() which created t,
// then Do panics.
func (t *T) Do(d []float64) {
	if len(d) != len(t.tmp) {
		panic("wrong size input")
	}
	t.doRec(d, t.tmp, t.p)
	t.scale(d)
}

func (t *T) doRec(d, e []float64, p uint) {
	if p == 0 {
		return
	}
	n := len(d)
	h := n / 2
	top := n - 1
	var x, y float64
	cs := t.cosTbl[p]
	for i := 0; i < h; i++ {
		x = d[i]
		y = d[top-i]
		e[i] = x + y
		e[h+i] = (x - y) / (2 * cs[i])
	}
	t.doRec(e[:h], d[:h], p-1)
	t.doRec(e[h:], d[:h], p-1)
	var i2, j int
	for i := 0; i < h-1; i++ {
		i2 = 2 * i
		d[i2] = e[i]
		j = h + i
		d[i2+1] = e[j] + e[j+1]
	}
	d[top-1] = e[h-1]
	d[top] = e[top]
}

// Inv inverts a transformed slice using the dct type III transform.
//
// If |d| is not the size indicated in the call to New() which created t,
// then Inv panics.
func (t *T) Inv(d []float64) {
	if len(d) != len(t.tmp) {
		panic("wrong input size")
	}
	d[0] /= 2
	t.invRec(d, t.tmp, t.p)
	t.scale(d)
}

func (t *T) scale(d []float64) {
	for i := range d {
		d[i] *= t.scf
	}
}

func (t *T) invRec(d, e []float64, p uint) {
	if p == 0 {
		return
	}
	n := len(d)
	h := n / 2
	top := n - 1
	e[0] = d[0]
	e[h] = d[1]
	var i2 int
	for i := 1; i < h; i++ {
		i2 = 2 * i
		e[i] = d[i2]
		e[h+i] = d[i2-1] + d[i2+1]
	}
	t.invRec(e[:h], d[:h], p-1)
	t.invRec(e[h:], d[:h], p-1)
	var x, y float64
	cs := t.cosTbl[p]
	for i := 0; i < h; i++ {
		x = e[i]
		y = e[h+i] / (2 * cs[i])
		d[i] = x + y
		d[top-i] = x - y
	}
}
