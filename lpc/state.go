// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package lpc

// State encapsulates linear prediction state for
// incremental usage in synthesis and prediction.
type State struct {
	alpha []float64
	hist  []float64
	i     int
}

// Predict returns the current prediction for the
// next value.
func (s *State) Predict() float64 {
	n := len(s.hist)
	k := 0
	i := s.i
	ttl := 0.0
	for j := 0; j < n; j++ {
		k = (j + i) % n
		ttl += s.hist[k] * s.alpha[j]
	}
	return ttl
}

// Consume advances the state one element (d), and
// returns the residue of the model for d.
func (s *State) Consume(d float64) float64 {
	p := s.Predict()
	s.hist[s.i] = d
	s.i++
	if s.i == len(s.hist) {
		s.i = 0
	}
	return d - p
}

// Produce synthesizes the next state from the residue r.
func (s *State) Produce(r float64) float64 {
	m := s.Predict()
	v := r + m
	s.Consume(v)
	return v
}
