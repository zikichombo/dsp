// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package convol

import (
	"fmt"

	"github.com/zikichombo/dsp/fft"
)

// K describes a convolver object T which has a precomputed argument component
// ("kernel").
type K struct {
	t      *T
	kernel fft.HalfComplex
}

// Conv computes a convolution of the kernel used
// to construct k with arg, placing the results in
// arg and returning them.  Conv returns a non-nil
// error and nil convolution if arg isn't the right
// length.
func (k *K) Conv(arg []float64) ([]float64, error) {
	if len(arg) != k.t.n {
		return nil, fmt.Errorf("arg dimension mismatch, %d != %d", len(arg), k.t.m)
	}
	arg = k.t.pad(k.t.WinB(arg), k.t.PadL())
	hc := k.t.ft.Do(arg)
	hc.MulElems(k.kernel)
	arg = k.t.ft.Inv(hc)
	return arg[:k.t.L()], nil
}

// ConvTo computes convolution of the kernel
func (k *K) ConvTo(dst, arg []float64) ([]float64, error) {
	dst = k.t.WinDst(dst)
	copy(dst, arg)
	dst = dst[:len(arg)]
	return k.Conv(dst)
}

// Win returns a slice containing everything in c with
// length and capacity set so that no copying
// copying takes place if c is used to house argument data
//
//  c := k.Win(nil)
//  for i := range c {
//    c[i] = ...
//  }
//  k.Conv(c)  // no copying
func (k *K) Win(c []float64) []float64 {
	return k.t.WinB(c)
}

// M() returns the length of the kernel.
func (k *K) M() int {
	return k.t.m
}

// N() returns the length of the argument.
func (k *K) N() int {
	return k.t.n
}

// L() returns the length of the result, which is
//
//  M() + N() - 1
//
func (k *K) L() int {
	return k.t.L()
}

// NewK creates a new convolver using "kernel" as
// the kernel.
func NewK(kernel []float64, argLen int) *K {
	k, e := New(len(kernel), argLen).K(kernel)
	if e != nil {
		panic(e)
	}
	return k
}

// K() creates a new kernel-convolver with using "kernel"
// as the kernel.
func (t *T) K(kernel []float64) (*K, error) {
	if len(kernel) != t.m {
		return nil, fmt.Errorf("kernel length wrong %d != %d", len(kernel), t.m)
	}
	krn := t.WinA(nil)
	copy(krn, kernel)
	krn = t.pad(krn, t.PadL())
	t.ft.Scale(false)
	hc := t.ft.Do(krn)
	t.ft.Scale(true)
	res := &K{
		t:      t,
		kernel: hc}
	return res, nil
}
