package resample

import (
	"math"

	"zikichombo.org/dsp/wfn"
)

// Itper provides an interface for an interpolator.
type Itper interface {
	// Order returns the maximum number of discrete neighbors on either
	// side of the point to be interpolated.
	Order() int

	// Itp performs the interpolation for point i.  i must be in
	// the range [0..2*Itper.Order()).  Interpolation neighborhood
	// is truncated according to the bounds of neighbors.
	Itp(neighbors []float64, i float64) float64

	// ItpCirc interpolates for point i in neighbors supposing
	// neighbors is circularly wrapped.  This is useful for
	// when neighbors is in a circular buffer and the context
	// guarantees there are sufficient neighbors around index i.
	CircItp(neighbors []float64, i float64) float64
}

type linItp struct{}

func (l *linItp) Order() int {
	return 1
}

func (l *linItp) Itp(nbrs []float64, p float64) float64 {
	pi, pf := math.Modf(p)
	q := int(pi)
	return (1-pf)*nbrs[q] + pf*nbrs[q+1]
}

func (l *linItp) CircItp(nbrs []float64, p float64) float64 {
	pi, pf := math.Modf(p)
	q := int(pi) % len(nbrs)
	r := q + 1
	if r == len(nbrs) {
		r = 0
	}
	return (1-pf)*nbrs[q] + pf*nbrs[r]
}

// LinItp returns a linear interpolator.
func LinItp() Itper {
	return &linItp{}
}

type fItp struct {
	order int
	fn    func(float64) float64
}

// NewFnItp returns a new interpolator from a weighting function fn.
//
// The function fn should accept values in the range (-o..o)
// giving the (signed) distance to the point to be interpolated.
// It should return an appropriate weight for input point at the specified
// distance.
func NewFnItp(o int, fn func(dist float64) float64) Itper {
	return &fItp{order: o, fn: fn}
}

func (i *fItp) Order() int {
	return i.order
}

func (i *fItp) Itp(nbrs []float64, p float64) float64 {
	acc := 0.0
	order := i.order
	qf, qr := math.Modf(p)
	q := int(qf)
	fn := i.fn
	for o := 0; o < order; o++ {
		l, r := q-o, q+o+1
		if l < 0 || r >= len(nbrs) {
			break
		}
		fo := float64(o)
		acc += fn(-(fo + qr)) * nbrs[l]
		acc += fn(fo+(1-qr)) * nbrs[r]
	}
	return acc
}

func (i *fItp) CircItp(nbrs []float64, p float64) float64 {
	acc := 0.0
	order := i.order
	qf, qr := math.Modf(p)
	q := int(qf) % len(nbrs)
	fn := i.fn
	for o := 0; o < order; o++ {
		l, r := q-o, q+o+1
		if l < 0 {
			l += len(nbrs)
		}
		if r >= len(nbrs) {
			r -= len(nbrs)
		}
		fo := float64(o)
		acc += fn(-(fo + qr)) * nbrs[l]
		acc += fn(fo+(1-qr)) * nbrs[r]
	}
	return acc
}

// NewSincItp returns a new Sinc interpolator from the Shannon
// interpolation theorem.
func NewSincItp(o int) Itper {
	return &fItp{order: o, fn: wfn.Sinc}
}

// NewWinSinc returns a new windowed sinc interpolator where
// the interpolation weighting function is a windowed sinc
// windowed by wf.
func NewWinSinc(o int, wf func(float64) float64) Itper {
	ws := func(d float64) float64 {
		return wfn.Sinc(d) * wf(d)
	}
	return &fItp{order: o, fn: ws}
}

// NewLanczos returns a new Lanczos interpolator with
// stretch "stretch" of order order.
func NewLanczos(order, stretch int) Itper {
	sf := func(d float64) float64 {
		return wfn.LanczosItp(stretch, d)
	}
	return &fItp{order: order, fn: sf}
}
