// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import "testing"

func TestChirpz(t *testing.T) {
	c := newChirpz(10, 32, newTwiddles(32, false))
	_ = c
}
