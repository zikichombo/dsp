// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package convol

import "testing"

func TestK(t *testing.T) {
	for i := 0; i < 64; i++ {
		a, b := gen()
		krn := NewK(a, len(b))
		td := To(a, b, nil)
		kd, e := krn.ConvTo(b, nil)
		if e != nil {
			t.Error(e)
			continue
		}
		if k := approxEq(td, kd, 0.001); k != -1 {
			t.Errorf("%d kernel/to mismatch %v (*) %v @%d: %.2f %.2f\n", i, a, b, k, td[k], kd[k])
		}
	}
}
