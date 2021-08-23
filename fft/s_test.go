// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/zikichombo/sound/freq"
)

func TestPeaks(t *testing.T) {
	N := 512
	d := generate(5000*freq.Hertz, N)
	Do(d)
	sp := NewS(d)
	peaks := sp.Peaks()
	for _, p := range peaks {
		i, m, p := sp.PeakItpQ(p)
		fmt.Printf("peak at %f: |%.2f| <%.2f>\n", i, m, p)
	}
	if len(peaks) != 3 {
		t.Errorf("wrong number of peaks: %d != 3\n", len(peaks))
	}
}

func TestFoldReal(t *testing.T) {
	for i := 0; i < 4; i++ {
		n := rand.Intn(100) + 10
		sp := NewSN(n)
		for j := 0; j < n; j++ {
			sp.SetMag(j, rand.Float64())
			sp.SetPhase(j, rand.Float64())
		}
		//sp.SetPhase(sp.Ny(), 0.0)
		sp.FoldReal()
		d := sp.Rect(nil)
		Do(d)
		for i, v := range d {
			if math.Abs(imag(v)) > 0.0000000001 {
				t.Errorf("%d: fold real n=%d gave imag component ift: %f\n", i, n, v)
			}
		}
	}
}

func TestSPlot(t *testing.T) {
	N := 512
	d := make([]complex128, N)
	for i := range d {
		d[i] = complex(math.Sin(float64(i)*15000.0*2*math.Pi/44100.0), 0)
	}
	Do(d)
	sp := NewS(d)
	im := sp.PlotMag(image.Rect(0, 0, 640, 480))
	f, e := os.Create("spect.png")
	if e != nil {
		t.Fatal(e)
	}
	if e := png.Encode(f, im); e != nil {
		t.Fatal(e)
	}
}

/*
func TestSPlot(t *testing.T) {
	N := 512
	d := generate(5000*freq.Hertz, N)
	Do(d)
	sp := NewS(d)
	p := sp.PlotMag(8 * vg.Inch)
	fmt := "png"
	f, e := os.Create("spex.png")
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	wt, e := p.WriterTo(8*vg.Inch, 8*vg.Inch, fmt)
	if e != nil {
		t.Fatal(e)
	}
	if _, e := wt.WriteTo(f); e != nil {
		t.Fatal(e)
	}
}
*/

func TestSSymEvenN(t *testing.T) {
	testSSymN(16, t)
}

func TestSSymOddN(t *testing.T) {
	testSSymN(17, t)
}

func testSSymN(N int, t *testing.T) {
	ft := New(N)
	b := ft.Win(nil)
	for i := 0; i < 10; i++ {
		for j := 0; j < N; j++ {
			c := complex(rand.Float64(), 0)
			b[j] = c
		}
		if e := ft.Do(b); e != nil {
			t.Fatal(e)
		}
		s := NewS(b)
		ny := s.Ny()
		for j := 0; j < ny; j++ {
			p, n := s.Mag(j), s.Mag(-j)
			if math.Abs(p-n) > 0.00000000001 {
				t.Errorf("at +/- %d %f, %f\n", j, p, n)
			}
		}
	}
}
