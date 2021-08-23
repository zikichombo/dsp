package resample

import (
	"math"
	"math/rand"
	"testing"

	"github.com/zikichombo/dsp/wfn"
)

var d = []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2, 0.1, 0,
	-0.1, -0.2, -0.3, -0.4, -0.5, -0.6, -0.7, -0.8, -0.9, -0.8, -0.7, -0.6, -0.5, -0.4, -0.3, -0.2, -0.1}

func TestItpLin(t *testing.T) {
	itper := LinItp()
	n := float64(len(d) - 1)
	p := rand.Float64() * n
	v := itper.Itp(d, p)
	pi, _ := math.Modf(p)
	l := int(pi)
	r := l + 1
	if !between(d[l], d[r], v) {
		t.Errorf("bad linear interpolation: %f not in {%f..%f}", v, d[l], d[r])
	}
}

func TestItpSinc(t *testing.T) {
	N := len(d) / 2
	N--
	ct := 0
	for o := 1; o < N; o++ {
		itper := NewSincItp(o)
		p := rand.Float64() * float64(N)
		v := itper.Itp(d, p)
		pi, _ := math.Modf(p)
		l := int(pi)
		r := l + 1
		if !between(d[l], d[r], v) {
			ct++
		}
		//fmt.Printf("sinc order %d {%f %f} of %f itp %f\n", o, d[l], d[r], pr, v)
	}
	if ct > 5 {
		t.Errorf("too many out of whack sinc interpolations: %d/%d\n", ct, N)
	}
}

func TestItpWinSinc(t *testing.T) {
	N := len(d) / 2
	N--
	ct := 0
	for o := 1; o < N; o++ {
		itper := NewWinSinc(o, wfn.Stretch(wfn.Blackman, math.Pi/float64(o)))
		p := rand.Float64() * float64(N)
		v := itper.Itp(d, p)
		pi, _ := math.Modf(p)
		l := int(pi)
		r := l + 1
		if !between(d[l], d[r], v) {
			ct++
		}
		//fmt.Printf("win sinc order %d {%f %f} of %f itp %f\n", o, d[l], d[r], pr, v)
	}
	if ct > 5 {
		t.Errorf("windowed sinc interpolation too many not between: %d/%d\n", ct, N)
	}
}

func TestItpLanczos(t *testing.T) {
	N := len(d) / 2
	N--
	ct := 0
	for a := 1; a < 5; a++ {
		for o := 1; o < N; o++ {
			itper := NewLanczos(o, a)
			p := rand.Float64() * float64(N)
			v := itper.Itp(d, p)
			pi, _ := math.Modf(p)
			l := int(pi)
			r := l + 1
			if !between(d[l], d[r], v) {
				ct++
			}
			//fmt.Printf("lanczos order %d stretch %d {%f %f} of %f itp %f\n", o, a, d[l], d[r], pr, v)
		}
	}
	if ct > 5 {
		t.Errorf("too many outliers: %d/%d\n", ct, N)
	}
}

func between(l, r, v float64) bool {
	if l > r {
		l, r = r, l
	}
	l -= 0.01
	r += 0.01
	return v >= l && v <= r
}
