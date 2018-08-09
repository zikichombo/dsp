package resample

import (
	"math"

	"zikichombo.org/dsp/wfn"
	"zikichombo.org/sound"
	"zikichombo.org/sound/ops"
)

// C provides continuous time representation of sampled signal by means of
// interpolation.
type C struct {
	src     sound.Source
	shift   int
	bufSize int
	buf     []float64
	off     int
	err     error
	eps     float64
	itper   Itper
}

const shiftSize = 64

// NewC creates a new continuous time representation of src.
func NewC(src sound.Source, itp Itper) *C {
	order := 10
	if itp != nil {
		order = itp.Order()
	}
	sz := 2*order + shiftSize
	if itp == nil {
		n := 2 * order
		m := float64(n - 1)
		r := 2 * math.Pi / m
		itp = NewWinSinc(order, wfn.Stretch(wfn.Blackman, r))
	}
	return &C{
		src:     src,
		shift:   shiftSize,
		bufSize: sz,
		off:     -sz,
		itper:   itp,
		eps:     0.0000000001,
		buf:     make([]float64, sz)}
}

// Eps gives the smallest difference such that
// interpolation takes place.  If the float index i
// passed to At has rational part <= eps, then
// a value from the underlying source is returned.
// By default, eps is 0.0000000001.
func (c *C) Eps(v float64) {
	c.eps = v
}

// At returns a continuous time interpolated sample
// at index i.
//
// At should be called with i increasing monotonically
// to guarantee that c does not need to go back in time
// arbitrariyly in its underlying source.
//
// If i is not increasing monotonically, the behavior
// of At is undefined.
//
// At returns a non-nil error if i >= the number of samples
// available in the underlying source without returning
// an error.  The returned error is that returned from
// the underlying source.
//
// At the edges, where insufficient or no neighbors are available,
// the interpolation is truncated symmetrically.
func (c *C) At(i float64) (float64, error) {
	jf, jr := math.Modf(i)
	j := int(jf)
	for j+c.itper.Order() >= c.off+c.bufSize {
		if c.err != nil {
			return 0, c.err
		}
		n, e := ops.Hop(c.src, c.buf, c.shift)
		if e != nil {
			c.err = e
		}
		c.off += n
	}
	j -= c.off
	if jr <= c.eps || (1-jr) <= c.eps {
		return c.buf[j], nil
	}
	order := c.itper.Order()
	if j+order >= len(c.buf) {
		order = len(c.buf) - 1 - j
	}
	if j-order < 0 {
		order = j
	}
	if order == 0 {
		return c.buf[j], nil
	}
	return c.itper.Itp(c.buf[j-order+1:j+order+1], float64(order-1)+jr), nil
}

func (c *C) Channels() int {
	return c.src.Channels()
}
