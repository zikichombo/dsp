// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import (
	"container/heap"
	"math"
)

// Z provides an interface for compression via
// coeficient selection
type Z struct {
	nz           []int
	pi           []int
	vs           []float64
	M            float64
	op2, p2, lp2 float64
}

// NewZ create a new Z for cosine transforms
// of size N.
func NewZ(N int) *Z {
	res := &Z{}
	res.pi = make([]int, N)
	res.nz = make([]int, 0, N)
	res.M = float64(N)
	return res
}

// Init initializes z for the transform data d
// and sets the current transform data of z to d.
// Init then returns the power of d.
//
func (z *Z) Init(d []float64) float64 {
	pwr := 0.0
	z.pi = z.pi[:len(d)]
	for i, v := range d {
		pwr += v * v
		z.pi[i] = i
	}
	z.op2 = pwr
	z.p2 = pwr
	z.lp2 = 2 * pwr
	pwr /= z.M
	pwr = math.Sqrt(pwr)

	z.vs = d
	zh := (*zhp)(z)
	heap.Init(zh)

	z.nz = z.nz[:0]
	return pwr
}

// Top returns
//
// - the index of the most significant remaining coeficient, ci.
//
// - the ratio of the popped power to the original, pRatio.
//
// - the rate at which the ratio changed w.r.t. to the last Pop(), rate.
func (z *Z) Top() (ci int, pRatio, rate float64) {
	ci = z.pi[0]
	v := z.vs[ci]
	org := math.Sqrt(z.op2 / z.M)
	d := z.p2 - v*v
	pRatio = math.Sqrt((z.op2-d)/z.M) / org
	d = z.lp2 - z.p2
	rate = math.Sqrt(d/z.M) / math.Sqrt(z.op2/z.M)
	return
}

// Pop pops a most significant coeficient from
// the current transform data.
func (z *Z) Pop() {
	zh := (*zhp)(z)
	ci := heap.Pop(zh).(int)
	z.nz = append(z.nz, ci)
	v := z.vs[ci]
	z.lp2 = z.p2
	z.p2 -= v * v
	return
}

// Zero zeros all un-popped coefficients
func (z *Z) Zero() {
	for _, vi := range z.pi {
		z.vs[vi] = 0.0
	}
}

// Coefs returns the indices of the popped coeficients
// storing the results in dst if there is space.
func (z *Z) Coefs(dst []int) []int {
	dst = dst[:0]
	dst = append(dst, z.nz...)
	return dst
}

type zhp Z

func (z *zhp) Less(i, j int) bool {
	return math.Abs(z.vs[z.pi[i]]) > math.Abs(z.vs[z.pi[j]])
}

func (z *zhp) Swap(i, j int) {
	z.pi[i], z.pi[j] = z.pi[j], z.pi[i]
}

func (z *zhp) Len() int {
	return len(z.pi)
}

func (z *zhp) Push(v interface{}) {
	z.pi = append(z.pi, v.(int))
}

func (z *zhp) Pop() interface{} {
	n := len(z.pi) - 1
	r := z.pi[n]
	z.pi = z.pi[:n]
	return r
}
