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
}

func NewChain(options ...Option) Chain {
	c := &chain{}
	c.gpool = sync.Pool{New: func() any { return &gauge{} }}
	// c.cpool = sync.Pool{New: func() any { return &counter{} }}
	// c.hpool = sync.Pool{New: func() any { return &histogram{} }}
	for _, fn := range options {
		fn(c)
	}
	return c
}

func (c *chain) Gauge(initName string, f func() float64) GaugeChain {
	return c.acquireGauge(initName, f)
}

func (c *chain) Counter(initName string) CounterChain {
	// todo implement me
	return nil
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

	var g *metrics.Gauge
	if c.vmset != nil {
		g = c.vmset.NewGauge(scopy(fullName), f)
	} else {
		g = metrics.NewGauge(scopy(fullName), f)
	}
	c.gmap.Store(fullName, g)
	return g
}

func (c *chain) ptr() uintptr {
	return uintptr(unsafe.Pointer(c))
}

func scopy(s string) string {
	buf := make([]byte, len(s))
	copy(buf, s)
	return byteconv.B2S(buf)
}
