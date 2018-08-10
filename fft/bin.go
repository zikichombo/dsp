// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"

	"zikichombo.org/sound/freq"
)

// Ny gives the index of the first frequency bin at or above the Nyquist limit
// of a FT of size n.
func Ny(n int) int {
	m := n / 2
	if n%2 == 1 {
		m++
	}
	return m
}

// FreqBin gives the DFT frequency bin of frequency a in
// a window of n-sample at rate sFreq
func FreqBin(sFreq, aFreq freq.T, n int) int {
	fn := freq.T(n)
	df := sFreq / fn
	r, mr := aFreq/df, aFreq%df
	//fmt.Printf("FreqBin(s=%s f=%s, n=%d, b*=%d, mb=%d, d=%d mdf=%d)\n", sFreq, aFreq, n, r, mr, df, r*mdf)
	if mr >= df/2 {
		return int(r) + 1
	}
	return int(r)
}

// BinRange gives the range of frequencies in a DFT
// frequency bin b where the DFT is based on an n-sample
// window at rate sFreq.
func BinRange(sFreq freq.T, n, b int) (l, u freq.T) {
	df := sFreq / freq.T(n)
	l = df * freq.T(b)
	u = l + df
	return l, u
}

// FreqAt gives the frequency at a floating point index i for
// A FT for data of length n with sample frequency sFreq.
func FreqAt(sFreq freq.T, n int, i float64) freq.T {
	ii, r := math.Modf(i)
	l, u := BinRange(sFreq, n, int(ii))
	off := float64(u-l) * r
	return l + freq.T(int64(off))
}

// WinSize gives the smallest positive window size (n, c) (n samples covering c
// cycles of a)  of a DFT applied to input signal with sample rate fS for frequency fA such
// that a cycle appears after n samples.  The cycle is approximate for period functions of frequency
// fA but exact for some frequency within an error range around that (at worst within [fA .. fA + fE]).
//
// Special case:
//  if it overflows, it returns -1, -1
//
// For example, to find a good window size for detecting a=440 Hz in an s=17.6Khz sample,
// we'd call WinSize(440Hz, 17.6KHz, 0.001Hz), giving n, and c.  The window size n
// in any signal of sampling rate s has c full cycles of any signal at frequency a, with
// very little error because (sc % a) is small.  As a result, using a DFT over
// n samples will often have less error and noise in it, since the DFT assumes a cyclic
// signal (actually stationary,  cyclic applied infinitely satisfies).
func WinSize(fS, fA, fE freq.T) (n, c int) {
	fc := freq.T(1)
	for {
		sc := fS * fc
		if sc < fS {
			return -1, -1
		}
		if sc%fA <= fE {
			return int(sc / fA), int(fc)
		}
		fc++
	}
}
