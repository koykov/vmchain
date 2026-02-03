package vmchain

import "github.com/koykov/indirect"

type CounterChain interface {
	WithLabel(name, value string) CounterChain
	WithAnyLabel(name string, value any) CounterChain
	Add(value int)
	AddInt64(value int64)
	Set(value uint64)
	Inc()
	Get() uint64
	Dec()
}

type counter struct {
	builder
	sptr uintptr
}

func (c *counter) WithLabel(name, value string) CounterChain {
	c.setLabel(name, value)
	return c
}

func (c *counter) WithAnyLabel(name string, value any) CounterChain {
	c.setAnyLabel(name, value)
	return c
}

func (c *counter) Add(value int) {
	if s := c.indirectSet(); s != nil {
		defer s.releaseCounter(c)
		s.getCounter(c.commit()).Add(value)
	}
}

func (c *counter) AddInt64(value int64) {
	if s := c.indirectSet(); s != nil {
		defer s.releaseCounter(c)
		s.getCounter(c.commit()).AddInt64(value)
	}
}

func (c *counter) Set(value uint64) {
	if s := c.indirectSet(); s != nil {
		defer s.releaseCounter(c)
		s.getCounter(c.commit()).Set(value)
	}
}

func (c *counter) Inc() {
	if s := c.indirectSet(); s != nil {
		defer s.releaseCounter(c)
		s.getCounter(c.commit()).Inc()
	}
}

func (c *counter) Get() uint64 {
	if s := c.indirectSet(); s != nil {
		defer s.releaseCounter(c)
		return s.getCounter(c.commit()).Get()
	}
	return 0
}

func (c *counter) Dec() {
	if s := c.indirectSet(); s != nil {
		defer s.releaseCounter(c)
		s.getCounter(c.commit()).Dec()
	}
}

func (c *counter) indirectSet() *chain {
	if c.sptr == 0 {
		return nil
	}
	return (*chain)(indirect.ToUnsafePtr(c.sptr))
}
