package vmchain

import "github.com/koykov/indirect"

type Gauge interface {
	WithLabel(name, value string) Gauge
	Add(value float64)
	Set(value float64)
	Inc()
	Get() float64
	Dec()
}

type gauge struct {
	builder
	sptr uintptr
	f    func() float64
}

func (g *gauge) WithLabel(name, value string) Gauge {
	g.setLabel(name, value)
	return g
}

func (g *gauge) Add(value float64) {
	if s := g.indirectSet(); s != nil {
		defer s.releaseGauge(g)
		s.getGauge(g.commit(), g.f).Add(value)
	}
}

func (g *gauge) Set(value float64) {
	if s := g.indirectSet(); s != nil {
		defer s.releaseGauge(g)
		s.getGauge(g.commit(), g.f).Set(value)
	}
}

func (g *gauge) Inc() {
	if s := g.indirectSet(); s != nil {
		defer s.releaseGauge(g)
		s.getGauge(g.commit(), g.f).Inc()
	}
}

func (g *gauge) Get() float64 {
	if s := g.indirectSet(); s != nil {
		defer s.releaseGauge(g)
		return s.getGauge(g.commit(), g.f).Get()
	}
	return 0
}

func (g *gauge) Dec() {
	if s := g.indirectSet(); s != nil {
		defer s.releaseGauge(g)
		s.getGauge(g.commit(), g.f).Dec()
	}
}

func (g *gauge) indirectSet() *vmset {
	if g.sptr == 0 {
		return nil
	}
	return (*vmset)(indirect.ToUnsafePtr(g.sptr))
}

func (g *gauge) reset() {
	g.builder.reset()
	g.sptr = 0
	g.f = nil
}
