// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"

	"zikichombo.org/dsp/mathutil/qitp"
	"zikichombo.org/sound/sample"
)

// S provides convenience wrappers around a ft spectrum.
type S struct {
	mags     []float64
	phases   []float64
	neg      int
	min, max float64
}

// NewS creates a spectrum from the data in d.  Once
// NewS returns, d is free to be used or gc'd.
func NewS(d []complex128) *S {
	s := NewSN(len(d))
	s.min = math.Inf(1)
	s.max = math.Inf(-1)
	for i, c := range d {
		m, p := cmplx.Polar(c)
		s.mags[i] = m
		s.phases[i] = p
		if m < s.min {
			s.min = m
		}
		if m > s.max {
			s.max = m
		}
	}
	return s
}

// NewSN creates a new S for ft spectrum of size n.
func NewSN(n int) *S {
	mem := make([]float64, n*2)
	mags := mem[:n]
	phases := mem[n:]
	neg := Ny(n)
	min, max := math.Inf(1), math.Inf(-1)
	return &S{mags: mags, phases: phases, neg: neg, min: min, max: max}
}

// At returns the complex representing the spectrum value at
// symmetric index i.
func (s *S) At(i int) complex128 {
	j := s.at(i)
	return cmplx.Rect(s.mags[j], s.phases[j])
}

// Mag returns the magnitude at symmetric index i.
func (s *S) Mag(i int) float64 {
	return s.mags[s.at(i)]
}

// SetMag sets the magnitude at symmetric index i.
func (s *S) SetMag(i int, m float64) {
	s.mags[s.at(i)] = m
	if m < s.min {
		s.min = m
	}
	if m > s.max {
		s.max = m
	}
}

// MagDb returns the magnitude in decibels
func (s *S) MagDb(i int) float64 {
	v := s.Mag(i)
	if v == 0 {
		v += 1e-20
	}
	return 20 * math.Log10(v)
}

// Phase returns the phase at symmetric index i.
func (s *S) Phase(i int) float64 {
	return s.phases[s.at(i)]
}

// SetPhase sets the phase at symmetric index i to p.
func (s *S) SetPhase(i int, p float64) {
	s.phases[s.at(i)] = p
}

// Ny returns the index of the first bin at or above
// the Nyquist from the index of the input from NewS().
func (s *S) Ny() int {
	return s.neg
}

// N returns the number of frequency bins in s.
func (s *S) N() int {
	return len(s.mags)
}

// Power returns the total power of the spectrum,
// the sum of squares of magnitudes.  Power assumes
// s represents real data.
func (s *S) Power() float64 {
	ttl := 0.0
	for i := 0; i < s.neg; i++ {
		ttl += s.mags[i] * s.mags[i]
	}
	ttl *= 2
	return math.Sqrt(ttl)
}

// ItpQMag uses quadratic interpolation to find the magnitude
// at index i.  Linear interpolation is used if s.N() == 2.
// ItpQMag panics if s.N() < 2.
func (s *S) ItpQMag(f float64) float64 {
	return sample.FromDb(qitp.SliceMap(s.mags, f, sample.ToDb))
}

// Peaks returns the indices of the non-negative frequency bins
// which are higher than one of their two neighbors and not
// less than either neighbor.  If there is only one element,
// that element is returned.  Endpoints are treated as though
// they are strictly higher than beyond the endpoint.
func (s *S) Peaks() []int {
	return s.PeaksTo(nil)
}

// PeaksTo places the peaks in dst by appending and returns
// the result.
func (s *S) PeaksTo(dst []int) []int {
	dst = dst[:0]
	n := len(s.mags)
	if n == 0 {
		return dst
	}
	if n == 1 {
		return dst
	}
	if n == 2 {
		dst = append(dst, 1)
		return dst
	}

	m := Ny(n)
	l := 0.0
	c := s.mags[0]
	r := s.mags[1]
	j := 2
	for j < m {
		l, c, r = c, r, s.mags[j]
		if c >= l && c >= r && (c > l || c > r) {
			dst = append(dst, j-1)
		}
		j++
	}
	if r >= c {
		dst = append(dst, m-1)
	}
	return dst
}

// PeakItpQ performs interpolation of spectrum peaks, giving
// a floating point index, magnitude, and phase.  To retrieve the
// frequency at the index, FreqAt is available.  Peaks
// often can correspond to  sinusoidal waves which are off-center
// of the frequency bin.  The peak shape is modelled as a parabola
// from neighboring points.
//
// If i is <= 1 or at the end of the Nyquist limit, then no
// interpolation takes place and the bin information is returned.
func (s *S) PeakItpQ(i int) (idx float64, mag float64, phase float64) {
	if i <= 1 || i >= s.neg-2 {
		return float64(i), s.mags[i], s.phases[i]
	}
	l, c, r := s.mags[i-1], s.mags[i], s.mags[i+1]
	l = sample.ToDb(l)
	c = sample.ToDb(c)
	r = sample.ToDb(r)
	a, b, c := qitp.Abc(l, c, r)
	h, k := qitp.Abc2Hk(a, b, c)
	mag = sample.FromDb(k)
	idx = float64(i) + h
	phase = 0.0
	return
}

// ItpPeaks interpolates all the peaks in the spectrum, returning their
// interpolated indices, magnitudes and phases.  Quadratic interpolation on log
// scale magnitudes is used, as in PeakItpQ, but the returned magnitudes are
// not log scale.
//
// The returned slice contains interpolated indices at n*3 positions and
// interpolated magnitudes at n*3+1 positions, and interpolated phases at n*3+2
// position.
//
// For example, if there are two peaks at 1.23 and 30.91 with magnitudes 10 and
// 100, and phases pi/49, 8pi/9 then the returned slice would be
//
//  {1.23, 10, pi/49, 30.91, 100, 8pi/9}
func (s *S) ItpPeaks(dst []float64) []float64 {
	ps := s.Peaks()
	for _, p := range ps {
		i, m, p := s.PeakItpQ(p)
		dst = append(dst, i)
		dst = append(dst, m)
		dst = append(dst, p)
	}
	return dst
}

// FromRect resets s to use the complex spectrum d.
func (s *S) FromRect(d []complex128) error {
	if len(d) != len(s.mags) {
		return fmt.Errorf("mismatched spectrum lengths: %d != %d", len(d), len(s.mags))
	}
	s.min, s.max = math.Inf(1), math.Inf(-1)
	for i, c := range d {
		m, p := cmplx.Polar(c)
		s.mags[i] = m
		s.phases[i] = p
		if m < s.min {
			s.min = m
		}
		if m > s.max {
			s.max = m
		}
	}
	return nil
}

// FromHalfComplex makes s contain spectrum from hc.
//
// FromHalfComplex returns a non-nil error if
// s doesn't contain the same number of elements as hc.
func (s *S) FromHalfComplex(hc HalfComplex) error {
	if len(hc) != len(s.mags) {
		return fmt.Errorf("mismatched spectrum lengths: %d != %d", len(hc), len(s.mags))
	}
	if len(hc) == 0 {
		return nil
	}
	hc.ToPolar(s.mags, s.phases)
	return nil
}

// Rect puts the spectrum in rectangular complex (real + imag) form in dst.
// If dst doesn't have capacity for the data, then a new slice is allocated
// and returned.  Otherwise, the results are placed in dst and returned.
func (s *S) Rect(dst []complex128) []complex128 {
	if cap(dst) < len(s.mags) {
		dst = make([]complex128, len(s.mags))
	}
	dst = dst[:len(s.mags)]
	for i := range dst {
		dst[i] = cmplx.Rect(s.mags[i], s.phases[i])
	}
	return dst
}

// ToHalfComplex places the spectrum s in dst.
//
// If dst doesn't have capacity for the data, then a new slice
// is allocated and returned.  Otherwise, the results are placed in dst
// and returned.
func (s *S) ToHalfComplex(dst HalfComplex) HalfComplex {
	if cap(dst) != len(s.mags) {
		dst = HalfComplex(make([]float64, len(s.mags)))
	}
	dst = dst[:len(s.mags)]
	dst.FromPolar(s.mags, s.phases)
	return dst
}

func (s *S) PlotMagTo(b image.Rectangle, p string) error {
	f, e := os.Create(p)
	if e != nil {
		return e
	}
	defer f.Close()
	img := s.PlotMag(b)
	return png.Encode(f, img)
}

// PlotMag plots the magnitudes on an image of
// dimensions b and returns it.
func (s *S) PlotMag(b image.Rectangle) *image.RGBA {
	mAcc := 0.0
	wAcc := 0.0
	rAcc := 0.0
	im := image.NewRGBA(b)
	x := 0
	black := color.RGBA{A: 255}
	for x := b.Min.X; x < b.Max.X; x++ {
		im.Set(x, b.Min.Y, black)
		im.Set(x, b.Max.Y-1, black)
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		im.Set(b.Min.X, y, black)
		im.Set(b.Max.X-1, y, black)
	}
	bb := image.Rect(b.Min.X+1, b.Min.Y+1, b.Max.X-1, b.Max.Y-1)
	subIm := im.SubImage(bb).(*image.RGBA)
	r := float64(bb.Dx()) / float64(len(s.mags))

	color := color.RGBA{
		A: 180,
		R: 0,
		G: 200,
		B: 255}
	if s.min == 0 {
		s.min = 1e-10
	}
	minDb := sample.ToDb(s.min)
	maxDb := sample.ToDb(s.max)
	for j := 0; j < len(s.mags); j++ {
		i := j - s.neg
		if i < 0 {
			i += len(s.mags)
		}
		mdb := s.MagDb(i)
		rAcc += r
		if rAcc < 1.0 {
			mAcc += mdb * r
			wAcc += r
			continue
		}
		_, ar := math.Modf(rAcc)
		mAcc += mdb * (r - ar)
		wAcc += r - ar
		v := mAcc / wAcc
		yr := (v - minDb) / (maxDb - minDb)
		Y := int(math.Floor(yr*float64(bb.Dy()) + 0.5))
		for rAcc >= 1.0 {
			for y := 1; y <= Y; y++ {
				subIm.Set(b.Min.X+x, b.Max.Y-y, color)
			}
			rAcc -= 1.0
			x++
		}
		mAcc = mdb * ar
		rAcc = ar
		wAcc = ar
	}
	return im
}

// FoldReal takes the spectrum and makes all negative frequency bins
// complex conjugates of the corresponding positive frequency bin to
// guarantee real output of inverse fft.
//
func (s *S) FoldReal() {
	N := s.Ny()
	M := len(s.phases)
	s.phases[0] = 0
	if M%2 == 0 {
		s.phases[N] = 0
	}
	for i := 1; i < N; i++ {
		s.phases[M-i] = -s.phases[i]
		s.mags[M-i] = s.mags[i]
	}
}

// CopyFrom makes s a copy of t.
func (s *S) CopyFrom(t *S) {
	if cap(s.mags) < len(t.mags) {
		s.mags = make([]float64, len(t.mags))
		s.phases = make([]float64, len(t.phases))
	}
	s.mags = s.mags[:len(t.mags)]
	s.phases = s.phases[:len(t.mags)]
	copy(s.mags, t.mags)
	copy(s.phases, t.phases)
}

func (s *S) at(i int) int {
	if i < 0 {
		return len(s.mags) + i
	}
	return i
}
