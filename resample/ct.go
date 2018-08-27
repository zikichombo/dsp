package resample

import (
	"errors"
	"io"
	"math"

	"zikichombo.org/dsp/wfn"
	"zikichombo.org/sound"
	"zikichombo.org/sound/cil"
	"zikichombo.org/sound/freq"
)

const (
	shiftSize = 64
)

// Type C holds state for giving a continuous time
// representation of a sound.Source.
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

// SampleRateConverter provides an interface to a dynamic resample rate
// conversion.
type SampleRateConverter interface {
	// Convert is called by the resampling methods in this package to determine
	// the result frequency in a conversion.
	//
	// It is called once for every output sample except the first sample, which
	// is taken to be at the same point in time as the first input sample.
	//
	// The return value should provide the ratio of the output rate to the input
	// rate.  It is assumed the input rate is fixed and determined by calling
	// context.
	Convert() float64
}

// ConstSampleRateStretcher is a type which implements SampleRateConverter
// based on a float64 constant sample rate conversion ratio.
type constSampleRateConverter float64

// Stretch implements SampleRateStretcher
func (c constSampleRateConverter) Convert() float64 {
	return float64(c)
}

// DynResampler is used to dynamically resample a source.
// It does not implement sound.Source, since the sample rate is
// fixed.
type DynResampler struct {
	ct    *C
	src   sound.Source
	conv  SampleRateConverter
	lasti float64
	buf   []float64
}

// NewDynResampler creates a new Dynamic Resampler from a continuous
// time representation and a sample rate converter.
func NewDynResampler(c *C, conv SampleRateConverter) *DynResampler {
	return &DynResampler{ct: c, conv: conv, lasti: 0.0, buf: make([]float64, c.Channels())}
}

// DynResampler returns the number of channels.
func (r *DynResampler) Channels() int {
	return r.ct.src.Channels()
}

// Close implements sound.Close
func (r *DynResampler) Close() error {
	return r.ct.Close()
}

// Receive is as in sound.Source.Receive.
func (r *DynResampler) Receive(d []float64) (int, error) {
	nC := r.ct.Channels()
	if len(d)%nC != 0 {
		return 0, sound.ErrChannelAlignment
	}
	nF := len(d) / nC
	for f := 0; f < nF; f++ {
		if err := r.ct.FrameAt(r.buf, r.lasti); err != nil {
			if err == io.EOF {
				cil.Compact(d, nC, f)
				return f, nil
			}
			return 0, err
		}
		r.lasti += r.conv.Convert()
		for c := range r.buf {
			d[c*nF+f] = r.buf[c]
		}
	}
	return nF, nil
}

type constResampler struct {
	*DynResampler
	outRate freq.T
}

// SampleRate returns the output sample rate of c.
func (c *constResampler) SampleRate() freq.T {
	return c.outRate
}

// Resample takes a sound.Source src, a desired samplerate r, and
// an interpolator itp.
//
// If itp is nil, it will default to a high quality interpolator
// (order 10 Blackman windowed sinc interpolation).
//
// Resample returns a sound.Source whose SampleRate() is equal to
// r.
//
// After a call to Resample, either the Receive method of src
// should not be called, or the Receive method of the result
// should not be called.  Clearly, the former is the usual use case.
func Resample(src sound.Source, r freq.T, itp Itper) sound.Source {
	sr := src.SampleRate()
	if sr == r {
		return src
	}
	tr := float64(sr) / float64(r)
	conv := constSampleRateConverter(tr)
	ct := NewC(src, itp)
	dyn := NewDynResampler(ct, conv)
	return &constResampler{DynResampler: dyn, outRate: r}
}

// NewC creates a new continuous time representation
// of the source src using an interpolator itp.
//
// if itp is nil, then a default interpolator of high
// quality will be used (order 10 Blackman windowed sinc interpolation).
//
// NewC calls src.Receive in this process, so src.Receive
// should not be called if the resulting continuous time
// interface is used.
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
//	var buf [1]float64
//	if err := c.FrameAt(buf[:], i); err != nil {
//		return 0.0, err
//	}
//	return buf[0], nil
//
func (c *C) At(i float64) (float64, error) {
	if c.Channels() != 1 {
		return 0.0, errMultiChanAt
	}
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

// Close closes the underlying source and returns the resulting error.
func (c *C) Close() error {
	return c.src.Close()
}

// FrameAt places a continuous time interpolated frame at index i in
// dst.
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
// FrameAt returns sound.ErrChannelAlignment if len(dst) != c.Channels().
func (c *C) FrameAt(dst []float64, i float64) error {
	nC := c.Channels()
	if len(dst) != nC {
		return sound.ErrChannelAlignment
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
