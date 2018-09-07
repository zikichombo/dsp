// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cmplx

import "fmt"

// Do performs linear convolution of a and b, placing
// the result in a and returning it if there is space,
// otherwise a new slice is allocated.
//
// To avoid the allocation, it should be that
//
//  cap(a) >= len(a) + len(b) - 1
//
func Do(a, b []complex128) []complex128 {
	t := New(len(a), len(b))
	a = t.WinA(a)
	b = t.WinB(b)
	res, e := t.Conv(a, b)
	if e != nil {
		panic(fmt.Sprintf("error %s\n", e))
	}
	return res
}

// To performs linear convolution of a and b, placing
// the result in dst and returning it.
//
// if dst does not have sufficient capacity, a new slice
// is allocated and returned in its place.
func To(dst, a, b []complex128) []complex128 {
	t := New(len(a), len(b))
	a = t.WinA(a)
	b = t.WinB(b)
	res, e := t.ConvTo(dst, a, b)
	if e != nil {
		panic(fmt.Sprintf("error %s\n", e))
	}
	return res
}
