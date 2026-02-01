package vmchain

import (
	"sync"
	"unsafe"

	"github.com/VictoriaMetrics/metrics"
	"github.com/koykov/byteconv"
)

type Chain interface {
	Gauge(name string, f func() float64) GaugeChain
	Counter(name string) CounterChain
	Histogram(name string) HistogramChain
}

type chain struct {
	gpool, cpool, hpool sync.Pool
	gmap, cmap, hmap    sync.Map
	gmux, cmux, hmux    sync.Mutex
	vmset               *metrics.Set

	gnew func(string, func() float64) *metrics.Gauge
	cnew func(string) *metrics.Counter
}

func NewChain(options ...Option) Chain {
	c := &chain{
		gnew: metrics.NewGauge,
		cnew: metrics.NewCounter,
	}
	c.gpool = sync.Pool{New: func() any { return &gauge{} }}
	c.cpool = sync.Pool{New: func() any { return &counter{} }}
	// c.hpool = sync.Pool{New: func() any { return &histogram{} }}
	for _, fn := range options {
		fn(c)
	}
	if c.vmset != nil {
		c.gnew = c.vmset.NewGauge
		c.cnew = c.vmset.NewCounter
	}
	return c
}

func (c *chain) Gauge(initName string, f func() float64) GaugeChain {
	return c.acquireGauge(initName, f)
}

func (c *chain) Counter(initName string) CounterChain {
	return c.acquireCounter(initName)
}

func (c *chain) Histogram(initName string) HistogramChain {
	// todo implement me
	return nil
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
	if gg, ok := any(cc).(*counter); ok {
		gg.reset()
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

func (c *chain) ptr() uintptr {
	return uintptr(unsafe.Pointer(c))
}

func scopy(s string) string {
	buf := make([]byte, len(s))
	copy(buf, s)
	return byteconv.B2S(buf)
}
