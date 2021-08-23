// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package czt

import (
	"math"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/zikichombo/dsp/fft"
	"github.com/zikichombo/sound"
	"github.com/zikichombo/sound/freq"
	"github.com/zikichombo/sound/gen"
	"github.com/zikichombo/sound/ops"
)

func TestCztFftEq(t *testing.T) {
	L := 2
	U := 129
	N := 128
	for i := 0; i < N; i++ {
		d := genCmplx(L, U)
		czt := New(len(d), len(d), 0, 2*math.Pi)
		ft := fft.New(len(d))
		d = czt.Win(d)
		e := ft.Win(nil)

		copy(e, d)
		czt.Do(d)
		ft.Do(e)
		if b := cmplxApproxEq(d, e, 0.0001); b != -1 {
			t.Errorf("bs run %d, sz %d bin %d, %.3f v %.3f\n", i, len(d), b, d[b], e[b])
		}
	}
}

func TestCztFftZoomStartEq(t *testing.T) {
	SF := 1000 * freq.Hertz
	f1 := 53 * freq.Hertz
	f2 := f1 + 4*freq.Hertz
	sins, _ := ops.Add([]sound.Source{
		gen.Sin(f1),
		gen.Sin(f2)}...)
	//sins = gen.Constant(1.0)
	M := 10
	N := M
	L := 0 * freq.Hertz
	U := 100 * freq.Hertz
	for i := 1; i < 8; i++ {
		//fmt.Printf("zoom %d\n", i)
		ft := fft.New(N)
		ftw := ft.Win(nil)
		ops.SlurpCmplx(sins, ftw)
		ct := New(N, N/M, SF.RadsPer(L), SF.RadsPer(U))
		ctw := ct.Win(nil)
		copy(ctw, ftw)
		ft.Do(ftw)
		ct.Do(ctw)
		d := N / M
		k := 0
		for k < d {
			if !cmplxClose(ftw[k], ctw[k], 0.0001) {
				fl, fu := fft.BinRange(SF, N, k)
				cl, cu := ct.BinRange(SF, k)
				fm, fp := cmplx.Polar(ftw[k])
				cm, cp := cmplx.Polar(ctw[k])
				t.Errorf("[%s..%s] fft |%.3f| <%.3f> %.2f [%s..%s] czt |%.3f| <%.3f> %.2f\n", fl, fu, fm, fp, ftw[k], cl, cu, cm, cp, ctw[k])
			}
			k++
		}
		N *= 2
	}
}

func TestCztFftZoomMidEq(t *testing.T) {
	SF := 1000 * freq.Hertz
	f1 := 253 * freq.Hertz
	f2 := f1 + 4*freq.Hertz
	sins, _ := ops.Add([]sound.Source{
		gen.Sin(f1),
		gen.Sin(f2)}...)
	//sins = gen.Constant(1.0)
	M := 10
	N := M
	L := 200 * freq.Hertz
	U := L + 100*freq.Hertz
	for i := 1; i < 8; i++ {
		//fmt.Printf("zoom %d\n", i)
		ft := fft.New(N)
		ftw := ft.Win(nil)
		ops.SlurpCmplx(sins, ftw)
		ct := New(N, N/M, SF.RadsPer(L), SF.RadsPer(U))
		ctw := ct.Win(nil)
		copy(ctw, ftw)
		ft.Do(ftw)
		ct.Do(ctw)
		d := N / M
		j := N / 5
		k := 0
		for k < d {
			if !cmplxClose(ftw[j], ctw[k], 0.00001) {
				fl, fu := fft.BinRange(SF, N, j)
				cl, cu := ct.BinRange(SF, k)
				fm, fp := cmplx.Polar(ftw[j])
				cm, cp := cmplx.Polar(ctw[k])
				t.Errorf("[%s..%s] fft |%.3f| <%.3f> %.2f [%s..%s] czt |%.3f| <%.3f> %.2f\n", fl, fu, fm, fp, ftw[j], cl, cu, cm, cp, ctw[k])
			}
			j++
			k++
		}
		N *= 2
	}
}

func genCmplx(l, u int) []complex128 {
	r := rand.Intn(u - l)
	s := l + r
	res := make([]complex128, s)
	for i := range res {
		res[i] = complex(rand.Float64(), rand.Float64())
	}
	return res
}

func cmplxClose(a, b complex128, eps float64) bool {
	ra, ia := real(a), imag(a)
	rb, ib := real(b), imag(b)
	return math.Abs(ra-rb) < eps && math.Abs(ia-ib) < eps
}

func cmplxApproxEq(a, b []complex128, eps float64) int {
	for i := range a {
		ca := a[i]
		cb := b[i]
		ra := real(ca)
		rb := real(cb)
		ia := imag(ca)
		ib := imag(cb)
		if math.Abs(ra-rb) >= eps {
			return i
		}
		if math.Abs(ia-ib) >= eps {
			return i
		}
	}
	return -1
}
