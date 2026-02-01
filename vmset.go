package vmchain

import (
	"sync"
	"unsafe"

	"github.com/VictoriaMetrics/metrics"
	"github.com/koykov/byteconv"
)

type Set interface {
	Gauge(name string, f func() float64) Gauge
	Counter(name string) Counter
	Histogram(name string) Histogram
}

type vmset struct {
	gpool, cpool, hpool sync.Pool
	gmap, cmap, hmap    sync.Map
	gmux, cmux, hmux    sync.Mutex
}

func NewSet() Set {
	s := &vmset{}
	s.gpool = sync.Pool{New: func() any { return &gauge{} }}
	// s.cpool = sync.Pool{New: func() any { return &counter{} }}
	// s.hpool = sync.Pool{New: func() any { return &histogram{} }}
	return s
}

func (s *vmset) Gauge(initName string, f func() float64) Gauge {
	return s.acquireGauge(initName, f)
}

func (s *vmset) Counter(initName string) Counter {
	// todo implement me
	return nil
}

func (s *vmset) Histogram(initName string) Histogram {
	// todo implement me
	return nil
}

func (s *vmset) acquireGauge(initName string, f func() float64) *gauge {
	g := s.gpool.Get().(*gauge)
	g.sptr = s.ptr()
	g.setName(initName)
	g.f = f
	return g
}

func (s *vmset) releaseGauge(g Gauge) {
	if gg, ok := any(g).(*gauge); ok {
		gg.reset()
		s.gpool.Put(g)
	}
}

func (s *vmset) getGauge(fullName string, f func() float64) *metrics.Gauge {
	// Fast check.
	if raw, ok := s.gmap.Load(fullName); ok {
		return raw.(*metrics.Gauge)
	}

	// Slow path.
	s.gmux.Lock()
	defer s.gmux.Unlock()

	// Double check.
	if raw, ok := s.gmap.Load(fullName); ok {
		// Double check passed.
		return raw.(*metrics.Gauge)
	}

	fullNameCpy := make([]byte, 0, len(fullName))
	fullNameCpy = append(fullNameCpy, fullName...)
	g := metrics.NewGauge(byteconv.B2S(fullNameCpy), f)
	s.gmap.Store(fullName, g)
	return g
}

func (s *vmset) ptr() uintptr {
	return uintptr(unsafe.Pointer(s))
}
