// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package cmplx

import "fmt"

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

func To(a, b, dst []complex128) []complex128 {
	t := New(len(a), len(b))
	a = t.WinA(a)
	b = t.WinB(b)
	res, e := t.ConvTo(a, b, dst)
	if e != nil {
		panic(fmt.Sprintf("error %s\n", e))
	}
	return res
}
