package resample

import (
	"math"
	"testing"

	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/gen"
	"zikichombo.org/sound/ops"
)

func TestCDefaultMonoChan(t *testing.T) {
	gnr := gen.New(44100 * freq.Hertz)
	rps := (44100 * freq.Hertz).RadsPer(800 * freq.Hertz)
	rps /= 10
	d := 0.0
	src := gnr.Sin(800 * freq.Hertz)
	c := NewC(src, nil)
	d = 0.0
	err := 0.0
	for i := 0; i < 10000; i++ {
		fi := float64(i) / 10.0
		v, e := c.At(fi)
		if e != nil {
			t.Fatal(e)
		}
		//fmt.Printf("%d: order %d itp %f org %f err %f\n", i, o, v, math.Sin(d), math.Abs(v-math.Sin(d)))
		err += math.Abs(v - math.Sin(d))
		d += rps
	}
	eps := err / 10000
	if eps > 0.1 {
		t.Errorf("default error per sample too large: %f\n", eps)
	}
	//fmt.Printf("default error per sample %f:\n", err/10000)
}

func TestCDefaultMultiChan(t *testing.T) {
	gnr := gen.New(44100 * freq.Hertz)
	rps := (44100 * freq.Hertz).RadsPer(800 * freq.Hertz)
	rps /= 10
	d := 0.0
	src0, src1 := gnr.Sin(800*freq.Hertz), gnr.Sin(800*freq.Hertz)
	c := NewC(ops.MustJoin(src0, src1), nil)
	d = 0.0
	err := 0.0
	frame := make([]float64, 2)
	for i := 0; i < 10000; i++ {
		fi := float64(i) / 10.0
		if e := c.FrameAt(frame, fi); e != nil {
			t.Fatal(e)
		}
		//fmt.Printf("%d: order %d itp %f org %f err %f\n", i, o, v, math.Sin(d), math.Abs(v-math.Sin(d)))
		for _, v := range frame {
			err += math.Abs(v - math.Sin(d))
		}
		d += rps
	}
	eps := err / 20000
	if eps > 0.1 {
		t.Errorf("default error per sample too large: %f\n", eps)
	}
}

func TestCSinc(t *testing.T) {
	gnr := gen.New(44100 * freq.Hertz)
	rps := (44100 * freq.Hertz).RadsPer(800 * freq.Hertz)
	rps /= 10
	d := 0.0
	for o := 1; o <= 30; o++ {
		src := gnr.Sin(800 * freq.Hertz)
		itper := NewSincItp(o)
		c := NewC(src, itper)
		d = 0.0
		err := 0.0
		for i := 0; i < 10000; i++ {
			fi := float64(i) / 10.0
			v, e := c.At(fi)
			if e != nil {
				t.Fatal(e)
			}
			//fmt.Printf("%d: order %d itp %f org %f err %f\n", i, o, v, math.Sin(d), math.Abs(v-math.Sin(d)))
			err += math.Abs(v - math.Sin(d))
			d += rps
		}
		eps := err / 10000
		if eps > 0.5/float64(o) {
			t.Errorf("sinc %d error per sample too large: %f\n", o, eps)
		}
		//fmt.Printf("sinc %d error per sample %f:\n", o, err/10000)
	}
}

func TestCLanczos(t *testing.T) {
	gnr := gen.New(44100 * freq.Hertz)
	rps := (44100 * freq.Hertz).RadsPer(800 * freq.Hertz)
	rps /= 10
	d := 0.0
	for a := 1; a < 5; a++ {
		for o := 1; o <= 20; o++ {
			src := gnr.Sin(800 * freq.Hertz)
			itper := NewLanczos(o, a)
			c := NewC(src, itper)
			d = 0.0
			err := 0.0
			for i := 0; i < 10000; i++ {
				fi := float64(i) / 10.0
				v, e := c.At(fi)
				if e != nil {
					t.Fatal(e)
				}
				//fmt.Printf("%d: order %d itp %f org %f err %f\n", i, o, v, math.Sin(d), math.Abs(v-math.Sin(d)))
				err += math.Abs(v - math.Sin(d))
				d += rps
			}
			eps := err / 10000
			if eps > 0.2/float64(o) && a != 1 {
				t.Errorf("lanczos order %d stretch %d error per sample too large: %f\n", o, a, eps)
			}
			//fmt.Printf("lanczos order %d stretch %d error per sample %f:\n", o, a, err/10000)
		}
	}
}

func TestResampleConst(t *testing.T) {
	fa := 800 * freq.Hertz
	src := gen.Sin(fa)
	r := 48000 * freq.Hertz
	rez := Resample(src, r, nil)
	rps := fa.RadsPerAt(rez.SampleRate())
	rads := 0.0
	N := 1000
	C := 4
	d := make([]float64, N)
	err := 0.0
	for c := 0; c < C; c++ {
		n, e := rez.Receive(d)
		if e != nil {
			t.Fatal(e)
		}
		if n != N {
			t.Fatalf("got %d frames not %d\n", n, N)
		}
		for i := range d {
			ref := math.Sin(rads)
			rads += rps
			got := d[i]
			err += math.Abs(got - ref)
		}
	}
	if err > 0.001*float64(N)*float64(C) {
		t.Errorf("resample error too large %f\n", err)
	}
}
