// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cmplx

import "fmt"

// Type K describes a convolver object T which
// has a precomputed argument component ("kernel").
type K struct {
	t      *T
	kernel []complex128
}

// Conv computes a convolution of the kernel used
// to construct k with arg, placing the results in
// arg and returning them.  Conv returns a non-nil
// error and nil convolution if arg isn't the right
// length.
func (k *K) Conv(arg []complex128) ([]complex128, error) {
	if len(arg) != k.t.n {
		return nil, fmt.Errorf("arg dimension mismatch, %d != %d\n", len(arg), k.t.m)
	}
	arg = k.t.pad(k.t.WinB(arg), k.t.PadL())
	k.t.ft.Do(arg)
	for i := range arg {
		arg[i] *= k.kernel[i]
	}
	k.t.ft.Inv(arg)
	return arg[:k.t.L()], nil
}

func (k *K) ConvTo(arg, dst []complex128) ([]complex128, error) {
	dst = k.t.WinDst(dst)
	copy(dst, arg)
	dst = dst[:len(arg)]
	return k.Conv(dst)
}

func (k *K) Win(c []complex128) []complex128 {
	return k.t.WinB(c)
}

func (k *K) M() int {
	return k.t.m
}

func (k *K) N() int {
	return k.t.n
}

func (k *K) L() int {
	return k.t.L()
}

func NewK(kernel []complex128, argLen int) *K {
	k, e := New(len(kernel), argLen).K(kernel)
	if e != nil {
		panic(fmt.Sprintf("%s", e))
	}
	return k
}

func (t *T) K(kernel []complex128) (*K, error) {
	if len(kernel) != t.m {
		return nil, fmt.Errorf("kernel length wrong %d != %d\n", len(kernel), t.m)
	}
	krn := t.WinA(nil)
	copy(krn, kernel)
	krn = t.pad(krn, t.PadL())
	t.ft.Do(krn)
	r := unscale(krn)
	for i := range krn {
		krn[i] *= r
	}
	res := &K{
		t:      t,
		kernel: krn}
	return res, nil
}
