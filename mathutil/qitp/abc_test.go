// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package qitp_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/zikichombo/dsp/mathutil/qitp"
)

func TestAbc(t *testing.T) {
	for i := 0; i < 100; i++ {
		p, q, r := rand.Float64(), rand.Float64(), rand.Float64()
		a, b, c := qitp.Abc(p, q, r)
		if math.Abs(p-(a*-1.0*-1.0+b*-1.0+c)) > 1e-10 {
			t.Errorf("p didn't interpolate: Abc(%f,%f,%f) = (%f,%f,%f)", p, q, r, a, b, c)
		}
		if math.Abs(q-(a*0*0+b*0+c)) > 1e-10 {
			t.Errorf("q didn't interpolate: Abc(%f,%f,%f) = (%f,%f,%f)", p, q, r, a, b, c)
		}
		if math.Abs(r-(a*1.0*1.0+b*1.0+c)) > 1e-10 {
			t.Errorf("r didn't interpolate: Abc(%f,%f,%f) = (%f,%f,%f)", p, q, r, a, b, c)
		}
	}
}

func TestSlice2(t *testing.T) {
	sl := make([]float64, 2)
	defer func() {
		if e := recover(); e != nil {
			t.Error(e)
		}
	}()
	for i := 0; i < 16; i++ {
		sl[0] = rand.Float64()
		sl[1] = rand.Float64()
		x := rand.Float64()
		qitp.Slice(sl, x)
	}
}

func TestAbcV(t *testing.T) {
	for i := 0; i < 16; i++ {
		a, b, c := rand.Float64(), rand.Float64(), rand.Float64()
		h, k := qitp.Abc2Hk(a, b, c)
		for i := 0; i < 10; i++ {
			x := 2 * (rand.Float64() - 0.5)
			u := a*x*x + b*x + c
			v := a*(x-h)*(x-h) + k
			if math.Abs(u-v) > 1e-10 {
				t.Errorf("abchk")
			}
		}
	}
}

func TestSliceN(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Error(e)
		}
	}()
	for i := 0; i < 16; i++ {
		n := rand.Intn(10) + 2
		fn := float64(n - 1)

		sl := make([]float64, n)
		for i := range sl {
			sl[i] = rand.Float64()
		}
		for j := 0; j < n; j++ {
			x := rand.Float64() * fn
			qitp.Slice(sl, x)
		}
	}
}
