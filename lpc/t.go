// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package lpc

import (
	"log"
	"math"
)

const eps = 1e-12

type T struct {
	rs    []float64
	k     []float64
	alpha []float64
}

// New returns a new linear predictive coder.
func New(order int) *T {
	return &T{
		rs:    make([]float64, order+1),
		k:     make([]float64, order+1),
		alpha: make([]float64, order+1)}
}

// Order returns the order of p.
func (p *T) Order() int {
	return len(p.rs) - 1
}

// Model causes p to learn coeficients for d, returning
// the model error.
func (p *T) Model(d []float64) float64 {
	err := p.levDurb(d)
	return err
}

// State returns an lpc state for incrementally predicting
// and synthesizing values according to the model in p.
func (p *T) State(seed []float64) *State {
	order := p.Order()
	st := &State{
		hist:  make([]float64, order),
		alpha: make([]float64, order)}
	copy(st.hist, seed)
	copy(st.alpha, p.alpha[1:])
	half := order / 2
	end := len(st.alpha) - 1
	for i := 0; i < half; i++ {
		st.alpha[i], st.alpha[end-i] = st.alpha[end-i], st.alpha[i]
	}
	return st
}

// Residue applies the model to d and for d[p.Order():]
// replaces the value of the input with the error (aka residue)
// of the model.
func (p *T) Residue(d []float64) {
	order := p.Order()
	for i := len(d) - 1; i >= order; i-- {
		iModel := 0.0
		for o := 1; o <= order; o++ {
			iModel += p.alpha[o] * d[i-o]
		}
		d[i] -= iModel
	}
}

// Restore restores d if d[p.Order():] is a residue generated
// from d[:p.Order()].
func (p *T) Restore(d []float64) {
	order := p.Order()
	N := len(d)
	for i := order; i < N; i++ {
		acc := 0.0
		for o := 1; o <= order; o++ {
			acc += p.alpha[o] * d[i-o]
		}
		d[i] += acc
	}
}

// R0 returns the zero autocorrelation value, useful
// for normalizing error.
func (p *T) R0() float64 {
	return p.rs[0]
}

func (p *T) ld2(d []float64) float64 {
	p.autoCorr(d)
	err := p.rs[0]
	r := 0.0
	order := p.Order()
	for i := 1; i <= order; i++ {
		r = -p.rs[i]
		for j := 1; j < i; j++ {
			r -= p.alpha[j] * p.rs[i-j]
		}
		r /= err
		p.alpha[i] = r
		err *= (1.0 - r*r)
		for j := 1; j < i/2; j++ {
			t := p.alpha[j]
			p.alpha[j] += r * p.alpha[i-j]
			p.alpha[i-j] += r * t
		}
		if i%2 == 1 {
			p.alpha[i/2] += p.alpha[i/2] * r
		}
		if err == 0.0 {
			log.Printf("need to limit order...")
		}
	}
	for i := 1; i <= order; i++ {
		p.alpha[i] = -p.alpha[i]
	}
	return err
}

func (p *T) levDurb(d []float64) float64 {
	p.autoCorr(d)
	err := p.rs[0]
	if math.Abs(err) < eps {
		err = 1.0 / eps
	}
	order := p.Order()
	alphaTmp := make([]float64, len(p.alpha))
	i := 1
	for i <= order {
		k := p.rs[i]
		for j := 1; j < i; j++ {
			fub := p.alpha[j] * p.rs[i-j]
			k -= fub
		}
		k /= err
		if math.Abs(k) > 1.0 {
			k = 1.0 / k
		}
		alphaTmp[i] = k
		for j := 1; j < i; j++ {
			alphaTmp[j] -= k * p.alpha[i-j]
		}
		copy(p.alpha, alphaTmp[:i+1])
		err *= (1.0 - k*k)
		i++
		if math.Abs(err) < eps {
			break
		}
	}
	p.rs = p.rs[:i]
	return err
}

func (p *T) autoCorr(d []float64) {
	N := len(d) - p.Order()
	for i := 0; i < N; i++ {
		u := d[i]
		for j := 0; j < len(p.rs); j++ {
			v := d[i+j]
			p.rs[j] += u * v
		}
	}
	for i := range p.rs {
		p.rs[i] /= float64(N)
	}
}
