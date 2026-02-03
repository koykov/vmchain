package vmchain

import (
	"sync"
	"unsafe"

	"github.com/VictoriaMetrics/metrics"
	"github.com/koykov/byteconv"
)

// Chain represents a set for chains of metrics.
type Chain interface {
	// Gauge initialize with initName a gauge chain and return it.
	Gauge(initName string, f func() float64) GaugeChain
	// Counter initialize with initName a counter chain and return it.
	Counter(initName string) CounterChain
	FloatCounter(initName string) FloatCounterChain
	// Histogram initialize with initName a histogram chain and return it.
	Histogram(initName string) HistogramChain
}

type chain struct {
	gpool, cpool, fpool, hpool sync.Pool
	gmux, cmux, fmux, hmux     sync.Mutex
	gmap, cmap, fmap, hmap     sync.Map

	vmset *metrics.Set
	gnew  func(string, func() float64) *metrics.Gauge
	cnew  func(string) *metrics.Counter
	fnew  func(string) *metrics.FloatCounter
	hnew  func(string) *metrics.Histogram
}

// NewChain makes a new chain set.
func NewChain(options ...Option) Chain {
	c := &chain{
		gnew: metrics.NewGauge,
		cnew: metrics.NewCounter,
		fnew: metrics.NewFloatCounter,
		hnew: metrics.NewHistogram,
	}
	c.gpool = sync.Pool{New: func() any { return &gauge{} }}
	c.cpool = sync.Pool{New: func() any { return &counter{} }}
	c.fpool = sync.Pool{New: func() any { return &fcounter{} }}
	c.hpool = sync.Pool{New: func() any { return &histogram{} }}
	for _, fn := range options {
		fn(c)
	}
	if c.vmset != nil {
		c.gnew = c.vmset.NewGauge
		c.cnew = c.vmset.NewCounter
		c.fnew = c.vmset.NewFloatCounter
		c.hnew = c.vmset.NewHistogram
	}
	return c
}

func (c *chain) Gauge(initName string, f func() float64) GaugeChain {
	return c.acquireGauge(initName, f)
}

func (c *chain) Counter(initName string) CounterChain {
	return c.acquireCounter(initName)
}

func (c *chain) FloatCounter(initName string) FloatCounterChain {
	return c.acquireFCounter(initName)
}

func (c *chain) Histogram(initName string) HistogramChain {
	return c.acquireHistogram(initName)
}

func (c *chain) acquireGauge(initName string, f func() float64) *gauge {
	g := c.gpool.Get().(*gauge)
	g.sptr = c.ptr()
	g.setName(initName)
	g.f = f
	return g
}

func (c *chain) releaseGauge(g GaugeChain) {
	if gg, ok := any(g).(*gauge); ok {
		gg.reset()
		c.gpool.Put(g)
	}
}

func (c *chain) getGauge(fullName string, f func() float64) *metrics.Gauge {
	// Fast check.
	if raw, ok := c.gmap.Load(fullName); ok {
		return raw.(*metrics.Gauge)
	}

	// Slow path.
	c.gmux.Lock()
	defer c.gmux.Unlock()

	// Double check.
	if raw, ok := c.gmap.Load(fullName); ok {
		// Double check passed.
		return raw.(*metrics.Gauge)
	}

	g := c.gnew(scopy(fullName), f)
	c.gmap.Store(fullName, g)
	return g
}

func (c *chain) acquireCounter(initName string) *counter {
	cc := c.cpool.Get().(*counter)
	cc.sptr = c.ptr()
	cc.setName(initName)
	return cc
}

func (c *chain) releaseCounter(cc CounterChain) {
	if cc_, ok := any(cc).(*counter); ok {
		cc_.reset()
		c.cpool.Put(cc)
	}
}

func (c *chain) getCounter(fullName string) *metrics.Counter {
	// Fast check.
	if raw, ok := c.cmap.Load(fullName); ok {
		return raw.(*metrics.Counter)
	}

	// Slow path.
	c.cmux.Lock()
	defer c.cmux.Unlock()

	// Double check.
	if raw, ok := c.cmap.Load(fullName); ok {
		// Double check passed.
		return raw.(*metrics.Counter)
	}

	cc := c.cnew(scopy(fullName))
	c.cmap.Store(fullName, cc)
	return cc
}

func (c *chain) acquireFCounter(initName string) *fcounter {
	cc := c.fpool.Get().(*fcounter)
	cc.sptr = c.ptr()
	cc.setName(initName)
	return cc
}

func (c *chain) releaseFCounter(cc FloatCounterChain) {
	if cc_, ok := any(cc).(*fcounter); ok {
		cc_.reset()
		c.fpool.Put(cc)
	}
}

func (c *chain) getFCounter(fullName string) *metrics.FloatCounter {
	// Fast check.
	if raw, ok := c.fmap.Load(fullName); ok {
		return raw.(*metrics.FloatCounter)
	}

	// Slow path.
	c.fmux.Lock()
	defer c.fmux.Unlock()

	// Double check.
	if raw, ok := c.fmap.Load(fullName); ok {
		// Double check passed.
		return raw.(*metrics.FloatCounter)
	}

	cc := c.fnew(scopy(fullName))
	c.fmap.Store(fullName, cc)
	return cc
}

func (c *chain) acquireHistogram(initName string) *histogram {
	h := c.hpool.Get().(*histogram)
	h.sptr = c.ptr()
	h.setName(initName)
	return h
}

func (c *chain) releaseHistogram(h HistogramChain) {
	if hh, ok := any(h).(*histogram); ok {
		hh.reset()
		c.hpool.Put(h)
	}
}

func (c *chain) getHistogram(fullName string) *metrics.Histogram {
	// Fast check.
	if raw, ok := c.hmap.Load(fullName); ok {
		return raw.(*metrics.Histogram)
	}

	// Slow path.
	c.hmux.Lock()
	defer c.hmux.Unlock()

	// Double check.
	if raw, ok := c.hmap.Load(fullName); ok {
		// Double check passed.
		return raw.(*metrics.Histogram)
	}

	cc := c.hnew(scopy(fullName))
	c.hmap.Store(fullName, cc)
	return cc
}

func (c *chain) ptr() uintptr {
	return uintptr(unsafe.Pointer(c))
}

func scopy(s string) string {
	buf := make([]byte, len(s))
	copy(buf, s)
	return byteconv.B2S(buf)
}
