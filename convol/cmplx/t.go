// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cmplx

import (
	"fmt"

	"zikichombo.org/dsp/fft"
)

// Type T holds state for performing n by m sized convolutions.
type T struct {
	n, m int
	winB []complex128
	ft   *fft.T
}

// New creates a new convolver object
func New(m, n int) *T {
	L := n + m - 1
	res := &T{
		n: n,
		m: m}
	res.ft = fft.New(res.PadL())
	res.winB = res.ft.Win(nil)[:L]
	return res
}

// N returns the length of the second argument.
func (t *T) N() int {
	return t.n
}

// M returns the length of the first argument
func (t *T) M() int {
	return t.m
}

// L returns the length of the result.
func (t *T) L() int {
	return t.n + t.m - 1
}

// PadL returns the fft padded length (can be
// > L).
func (t *T) PadL() int {
	L := t.L()
	if L&(L-1) == 0 {
		return L
	}
	res := 1
	for res < L {
		res *= 2
	}
	return res
}

// Conv performs a linear convolution of a and b, placing
// the results in a and returning them.
func (t *T) Conv(a, b []complex128) ([]complex128, error) {
	if len(a) != t.m {
		return nil, fmt.Errorf("operand dimension mismatch: %d != %d\n", len(a), t.m)
	}
	if len(b) != t.n {
		return nil, fmt.Errorf("kernel dimension mismatch: %d != %d\n", len(b), t.n)
	}
	copy(t.winB, b)
	return t.conv(a, t.winB[:len(b)])
}

// clobbers b with ifft.
func (t *T) conv(a, b []complex128) ([]complex128, error) {
	L := t.PadL()
	a = t.pad(a, L)
	b = t.pad(b, L)
	fmt.Printf("b %v\n", b)
	if e := t.ft.Do(b); e != nil {
		return nil, e
	}
	t.ft.Scale(false)
	if e := t.ft.Do(a); e != nil {
		return nil, e
	}
	t.ft.Scale(true)
	for i := range a {
		fmt.Printf("a %f b %f\n", a[i], b[i])
		a[i] *= b[i]
	}
	if e := t.ft.Inv(a); e != nil {
		return nil, e
	}
	return a[:t.L()], nil
}

// ConvTo performs a linear convolution of a and b,
// placing the result in dst.  If dst is nil, an
// appropriately len- and cap- dimensioned slice
// is allocated and returned.
func (t *T) ConvTo(a, b, dst []complex128) ([]complex128, error) {
	dst = t.WinDst(dst)
	copy(dst, a)
	return t.Conv(dst[:t.m], t.WinB(b))
}

// WinA returns a slice with len and cap dimensions
// set so that if the returned slice a is passed
// as the first argument to Conv(), then no
// a-argument related copying or allocations take place during the
// execution of Conv.
func (t *T) WinA(c []complex128) []complex128 {
	return t.ft.Win(c)[:t.m]
}

// WinB returns a slice with len and cap dimensions
// set so that if the returned slice b is passed
// as the second argument to Conv(), then no
// b-argument related copying or allocations take place during the
// execution of Conv.
func (t *T) WinB(c []complex128) []complex128 {
	return t.ft.Win(c)[:t.n]
}

// WinDst returns a slice with len and cap dimensions
// set so that if the returned slice is passed as
// the dst argument to Conv, then no dst related
// allocations or copying takes place during the
// execution of Conv()
func (t *T) WinDst(c []complex128) []complex128 {
	return t.ft.Win(c)[:t.L()]
}

func (t *T) pad(sl []complex128, L int) []complex128 {
	n := len(sl)
	if cap(sl) < L {
		tmp := make([]complex128, n, L)
		copy(tmp, sl)
		sl = tmp
	}
	sl = sl[:L]
	for i := n; i < L; i++ {
		sl[i] = 0i
	}
	return sl
}
