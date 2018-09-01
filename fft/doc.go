// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Package fft provides support for 1 dimensional fast fourier transform and
// spectra.
//
// Package fft is part of http://zikichombo.org
//
// Fast fourier transforms provide frequency domain representation of time
// domain signals.  Package fft supports a fairly efficient fourier transform
// implementation with a lot of simplicity and convenience.
//
//
// Features
//
// Package fft is designed primarily for repeated usage on a data stream.  As
// such, the incremental interface supports the creation of buffers which
// automatically have size and capacity which guarantee no copying or
// allocations during executation without the user needing to worry about the
// details of the sizes and capacities.
//
// Package fft provides an efficient real-only interface for even transform
// sizes with spectra represented in a half complex format like fftw.  An odd
// transform size real-only interface is also supported, but is not as efficient
// as the even transform size real-only interface.
//
// The interface guarantees the Parseval equation, which states the sum of
// squares of the amplitudes in the time domain equals the sum of squares of
// the frequency coeficients in the frequency domain.  It also guarantees that
// the inverse transform is an inverse without the user needing to worry about
// scaling.
//
// Non Features
//
// Package fft is designed exclusively for 1d data.  Support for matrices is a
// non-goal of this package.
//
// Algorithms
//
// The implementation uses in place decimation in time binary radix 2 Cooley
// Tuckey for sizes of powers of 2, and Bluestein algorithm otherwise.  Twiddle
// factors are pre-computed, and attention is paid to minimize allocations and
// copying, both internally and in the interface.  Package fft provides O(N Log
// N) transforms for all inputs and allows for in place transforms.  The Real
// only interface uses the complex interface for half-sized inputs together
// with some O(N) pre/post processing.
//
package fft /* import "zikichombo.org/dsp/fft" */
