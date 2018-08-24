package resample

import (
	"errors"
	"math"

	"zikichombo.org/dsp/wfn"
	"zikichombo.org/sound"
)

const (
	shiftSize = 64
)

type C struct {
	src     sound.Source
	shift   int
	bufSize int
	cbufs   [][]float64
	rbuf    []float64
	off     int
	err     error
	eps     float64
	itper   Itper
}

func NewC(src sound.Source, itp Itper) *C {
	order := 10
	if itp != nil {
		order = itp.Order()
	}
	if itp == nil {
		n := 2 * order
		m := float64(n - 1)
		r := 2 * math.Pi / m
		itp = NewWinSinc(order, wfn.Stretch(wfn.Blackman, r))
	}
	nC := src.Channels()
	sz := 2*order + shiftSize
	cbufs := make([][]float64, nC)
	for i := range cbufs {
		cbufs[i] = make([]float64, sz)
	}
	rbuf := make([]float64, shiftSize*nC)
	return &C{
		src:     src,
		shift:   shiftSize,
		bufSize: sz,
		off:     -sz,
		itper:   itp,
		eps:     0.0000000001,
		cbufs:   cbufs,
		rbuf:    rbuf}
}

var errMultiChanAt = errors.New("ErrMultiChanAt")

// At returns a continuous time interpolated sample at index i.
// It is the equivalent of
//
//  var buf [1]float64
//	if err := c.FrameAt(buf[:], i); err != nil {
//		return 0.0, err
//	}
//	return buf[0], nil
//
func (c *C) At(i float64) (float64, error) {
	if c.Channels() != 1 {
		return 0.0, errMultiChanAt
	}
	//return c.at(i)
	var buf [1]float64
	if err := c.FrameAt(buf[:], i); err != nil {
		return 0.0, err
	}
	return buf[0], nil
}

// Channels returns the number of channels of the source to
// which c provides continuous time access.
func (c *C) Channels() int {
	return c.src.Channels()
}

// FrameAt returns a continuous time interpolated frame at index i.
//
// FrameAt should be called with i increasing monotonically to guarantee that c
// does not need to go back in time arbitrarily in its underlying source.
//
// If i is not increasing monotonically, the behavior of FrameAt is undefined.
//
// FrameAt returns a non-nil error if i >= the number of samples available in
// the underlying source without returning an error.  The returned error is
// that returned from the underlying source.
//
// At the edges, where insufficient or no neighbors are available, the
// interpolation is truncated symmetrically.
//
// FrameAt returns sound.ChannelAlignmentError if len(dst) != c.Channels().
func (c *C) FrameAt(dst []float64, i float64) error {
	nC := c.Channels()
	if len(dst) != nC {
		return sound.ChannelAlignmentError
	}
	jf, jr := math.Modf(i)
	j := int(jf)
	order := c.itper.Order()
	for j+order >= c.off+c.bufSize {
		if c.err != nil {
			return c.err
		}
		n, e := c.src.Receive(c.rbuf)
		if e != nil {
			c.err = e
		}
		c.off += n
		for i, cb := range c.cbufs {
			copy(cb, cb[c.shift:])
			copy(cb[len(cb)-c.shift:], c.rbuf[i*n:(i+1)*n])
		}
	}
	itp := c.itper
	for ci := range dst {
		buf := c.cbufs[ci]
		cj := j - c.off
		if jr <= c.eps || (1-jr) <= c.eps {
			dst[ci] = buf[cj]
			continue
		}
		if cj+order >= len(buf) {
			order = len(buf) - 1 - cj
		}
		if cj-order < 0 {
			order = cj
		}
		if order == 0 {
			dst[ci] = buf[cj]
			continue
		}
		r := itp.Itp(buf[cj-order+1:cj+order+1], float64(order-1)+jr)
		dst[ci] = r
	}
	return nil
}
