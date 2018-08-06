package resample

import (
	"math"
	"testing"

	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/gen"
)

func TestCDefault(t *testing.T) {
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
