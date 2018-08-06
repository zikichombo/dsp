// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"testing"

	"zikichombo.org/sound/freq"
)

func TestR2(t *testing.T) {
	testBasic(2, 1024*freq.Hertz, t)
	testBasic(64, 1024*freq.Hertz, t)
	testRand(1024, 1024*freq.Hertz, t)
}

func TestLine(t *testing.T) {
	N := 128
	ft := New(N)
	buf := ft.Win(nil)
	h := N / 4
	src := make([]complex128, 10)
	for i := range src {
		src[i] = complex(float64(i)/10, 0.0)
	}
	for i := 0; i < h; i++ {
		buf[i] = src[i%len(src)]
	}
	for i := h; i < N; i++ {
		buf[i] = 0i
	}

	ft.Do(buf)
	sp := NewS(buf)
	for i := 0; i < sp.Ny(); i++ {
		mag := sp.Mag(i)
		sp.SetMag(i, mag*4.0)
	}
	sp.FoldReal()
	sp.Rect(buf)
	ft.Inv(buf)
	for i := range buf {
		t.Logf("%f\n", real(buf[i]))
	}
}

func TestBS(t *testing.T) {
	testBasic(63, 1024*freq.Hertz, t)
	testRand(1019, 1024*freq.Hertz, t)
}

func TestBSEqR2(t *testing.T) {
	F := 1000 * freq.Hertz
	for N := 8; N <= 512; N *= 2 {
		w := generate(F, N)
		ft := New(N)
		bsft := NewT(N, 2*N)
		dr2, dbs := ft.Win(nil), bsft.Win(nil)
		copy(dbs, w)
		ft.To(dr2, w)
		bsft.bluestein(dbs, false)
		for i := range w {
			cmplxCmpErr(dr2[i], dbs[i], 0.00001, t)
		}
	}
}

func TestBS2(t *testing.T) {
	d := []complex128{(0.15168410948190972 + 0.4143031535168929i), (0.5680494254123835 + 0.5073763963738084i)}
	ft := NewT(2, 4)
	d = ft.Win(d)
	ft.bluestein(d, false)
}

func TestAutoCorr(t *testing.T) {
	N := 128
	ft := New(N)
	d := ft.Win(nil)
	for i := range d {
		if i%6 == 0 {
			d[i] = 1 + 0i
			continue
		}
		d[i] = complex((rand.Float64()-0.5)*2, 0)
	}
	ft.AutoCorr(d)
	for i := 0; i < 32; i++ {
		t.Logf("%d: %f\n", i, real(d[i]))
	}
}

func ExampleT() {
	var d = []complex128{(1 + 0i), (1 + 0i), (1 + 0i)}
	tr := New(len(d))
	w := tr.Win(nil)
	copy(w, d) // for correct capacity.
	tr.Do(w)
	tr.Inv(w)
}

func ExampleT_Spike() {
	var d = []complex128{
		0i, 0i, 0i, 0i, 0i, 0i, 0i, (1 + 0i),
		0i, 0i, 0i, 0i, 0i, 0i, 0i, 0i}

	tr := New(len(d))
	w := tr.Win(nil)
	copy(w, d) // for correct capacity.
	tr.Do(w)
	s := NewS(w)
	for i := 0; i < s.Ny(); i++ {
		fmt.Printf("%0.3f ", s.Phase(i))
	}
	fmt.Println()
	// Output: 0.000 -2.749 0.785 -1.963 1.571 -1.178 2.356 -0.393

}

func BenchmarkDoR2(b *testing.B) {
	b.StopTimer()
	w := generate(1024*freq.Hertz, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Do(w)
	}
}

func BenchmarkT_DoR2(b *testing.B) {
	b.StopTimer()
	N := 1024
	w := generate(1024*freq.Hertz, N)
	tr := New(N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Do(w)
	}
}

func testBasic(N int, F freq.T, t *testing.T) {
	Eps := 0.001
	w := generate(F, N)
	ft := New(N)
	dst := ft.Win(nil)
	tmp := ft.Win(nil)
	for i := 0; i < 1; i++ {
		_, e := ft.To(dst, w)
		if e != nil {
			t.Error(e)
		}
		printT(dst, F)
		_, e = ft.InvTo(tmp, dst)
		if e != nil {
			t.Error(e)
			continue
		}
		for j := range tmp {
			cmplxCmpErr(tmp[j], w[j], Eps, t)
		}
	}
}

// runs forward and back on random data and checks
// F^{-1}(F(d)) = d
func testRand(N int, F freq.T, t *testing.T) {
	for rounds := 0; rounds < 10; rounds++ {
		d := make([]complex128, N)
		for i := 0; i < N; i++ {
			d[i] = complex(rand.Float64(), rand.Float64())
		}
		ft := New(N)
		tmp := ft.Win(nil)
		ft.To(tmp, d)
		ft.Inv(tmp)
		ttlOrg := 0i
		ttlTrn := 0i
		for i := 0; i < N; i++ {
			cmplxCmpErr(d[i], tmp[i], 0.00001, t)
			ttlOrg += d[i] * d[i]
			ttlTrn += tmp[i] * tmp[i]
		}
		// check parseval equation: sum of squares of input == sum of squares of inverse
		cmplxCmpErr(ttlOrg, ttlTrn, 0.0001, t)
	}
}

func generate(f freq.T, N int) []complex128 {
	w := make([]complex128, N)
	//fa, fb, fc := 128*freq.Hertz, 256*freq.Hertz, 64*freq.Hertz
	fa, fb, fc := f/3, f/4, f/5
	for i := 0; i < N; i++ {
		fi := float64(i)
		va := math.Sin(fa.RadsPerAt(f) * fi)
		vb := 2 * math.Sin(fb.RadsPerAt(f)*fi)
		vc := 5 * math.Sin(fc.RadsPerAt(f)*fi)
		v := va + vb + vc
		w[i] = complex(v, 0)
	}
	return w
}

func printT(fcs []complex128, f freq.T) {
	N := len(fcs) / 2
	//fmt.Printf("Window:\n")
	for i := 1; i < N; i++ {
		l, u := BinRange(f, len(fcs), i)
		c := fcs[i]
		mag, p := cmplx.Polar(c)
		if mag > 0.00005 {
			//fmt.Printf("\t%d: [%s .. %s) |%.2f| <%.2f>\n", i, l, u, mag, p)
			_, _, _ = l, u, p
		}
	}
	//fmt.Println()
}

func cmplxCmpErr(a, b complex128, eps float64, t *testing.T) error {
	var e error
	if cmplx.Abs(a-b) > eps || real(a) == math.NaN() || real(b) == math.NaN() {
		e = fmt.Errorf("got %f expecting %f", a, b)
		t.Error(e)
	}
	return e
}
