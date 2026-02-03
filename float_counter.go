package vmchain

import "github.com/koykov/indirect"

type FloatCounterChain interface {
	WithLabel(name, value string) FloatCounterChain
	WithAnyLabel(name string, value any) FloatCounterChain
	Add(value float64)
	Sub(value float64)
	Set(value float64)
	Get() float64
}

type fcounter struct {
	builder
	sptr uintptr
}

func (c *fcounter) WithLabel(name, value string) FloatCounterChain {
	c.setLabel(name, value)
	return c
}

func (c *fcounter) WithAnyLabel(name string, value any) FloatCounterChain {
	c.setAnyLabel(name, value)
	return c
}

func (c *fcounter) Add(value float64) {
	if s := c.indirectSet(); s != nil {
		defer s.releaseFCounter(c)
		s.getFCounter(c.commit()).Add(value)
	}
}

func (c *fcounter) Sub(value float64) {
	if s := c.indirectSet(); s != nil {
		defer s.releaseFCounter(c)
		s.getFCounter(c.commit()).Sub(value)
	}
}

func (c *fcounter) Set(value float64) {
	if s := c.indirectSet(); s != nil {
		defer s.releaseFCounter(c)
		s.getFCounter(c.commit()).Set(value)
	}
}

func (c *fcounter) Get() float64 {
	if s := c.indirectSet(); s != nil {
		defer s.releaseFCounter(c)
		return s.getFCounter(c.commit()).Get()
	}
	return 0
}

func (c *fcounter) indirectSet() *chain {
	if c.sptr == 0 {
		return nil
	}
	return (*chain)(indirect.ToUnsafePtr(c.sptr))
}
