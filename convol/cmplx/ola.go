// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cmplx

import (
	"fmt"
)

// Ola keeps state for implementing overlap-add
// block convolution
type Ola struct {
	k    *K
	over []complex128
	conv []complex128
}

// NewOla creates a new overlap block convolver
// based on the kernel krn with processing block
// size L.
func NewOla(krn []complex128, L int) *Ola {
	k := NewK(krn, L)
	return &Ola{
		k:    k,
		over: make([]complex128, len(krn)-1),
		conv: k.Win(nil)}
}

// M returns the length of the kernel
func (o *Ola) M() int {
	return o.k.M()
}

// N returns the block length of the input
func (o *Ola) N() int {
	return o.k.N()
}

// L returns o.M() + o.N() - 1, the zero
// padding size and size of the underlying fft.
func (o *Ola) L() int {
	return o.M() + o.N() - 1
}

// WinSrc takes a candidate window slice c and returns a slice properly
// proportioned, in terms of both length and capacity for passing as src arg of
// Block().
//
// The returned slice uses the backing store of c if possible and contains the
// elements of c.
func (o *Ola) WinSrc(c []complex128) []complex128 {
	return o.k.Win(c)
}

// WinDst takes a candidate window slice c and returns a slice properly
// proportioned, in terms of both length and capacity, for passing as arg dst
// to Block().
//
// The returned slice uses the backing store of c if possible and contains the
// elements of c.
func (o *Ola) WinDst(c []complex128) []complex128 {
	return o.k.t.ft.Win(c)[:o.k.t.n]
}

// Block processes one block of the convolution
func (o *Ola) Block(src, dst []complex128) error {
	conv, e := o.k.ConvTo(o.conv, src)
	if e != nil {
		return e
	}
	o.conv = conv
	M := o.k.t.m - 1
	fmt.Printf("M %d, len(dst) %d len(over) %d len(conv) %d\n", M, len(dst), len(o.over), len(conv))
	for i := 0; i < M; i++ {
		dst[i] = o.over[i] + conv[i]
	}
	copy(o.over, conv[o.N():])
	copy(dst[M:], conv[M:o.N()])
	return nil
}
