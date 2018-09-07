// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package lpc

import (
	"math"
	"testing"

	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/gen"
)

func TestLpc(t *testing.T) {
	testLpcConst(16, 1, t)
	testLpcConst(16, 2, t)
	testLpcConst(16, 3, t)
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		testLpcSin(512, i, t)
	}
}

func testLpcConst(N, order int, t *testing.T) {
	lpc := New(order)
	d := make([]float64, N)
	for i := range d {
		d[i] = 103.0
	}
	lpc.Model(d)
	lpc.Residue(d)
	for i := order; i < len(d); i++ {
		if d[i] != 0.0 {
			t.Errorf("lpc constant data[%d] prediction order %d failed at %d: gave %f wanted %f\n", N, order, i, d[i], 0.0)
		}
	}
	lpc.Restore(d)
	for i, v := range d {
		if v != 103.0 {
			t.Errorf("%d: residue/restore got %f not 103", i, v)
		}
	}
}

func testLpcSin(N, order int, t *testing.T) {
	lpc := New(order)
	//src, _ := ops.Add(ops.Amplify(gen.Noise(), 0.001), ops.Amplify(gen.Note(440*freq.Hertz), 0.999))
	src := gen.Note(440 * freq.Hertz)
	d := make([]float64, N)
	e := make([]float64, N)
	src.Receive(d)
	copy(e, d)
	lpc.Model(d)
	lpc.Residue(d)
	consumer := lpc.State(d[:order])
	producer := lpc.State(d[:order])
	for i := order; i < N; i++ {
		r := consumer.Consume(e[i])
		v := producer.Produce(r)
		if math.Abs(r-d[i]) > 1e-3 {
			t.Errorf("N=%d, order %d, at %d residue mismatch: got %f wanted %f\n", N, order, i, r, e[i])
		}
		if math.Abs(v-e[i]) > 1e-3 {
			t.Errorf("N=%d, order %d, at %d produce from residue mismatch got %f wanted %f\n", N, order, i, v, d[i])
		}
	}
	lpc.Restore(d)
	for i, v := range d {
		if math.Abs(v-e[i]) > 1e-3 {
			t.Errorf("at %d lpc restore got %f not %f\n", i, e[i], v)
		}
	}
}
