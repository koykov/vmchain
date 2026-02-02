package vmchain

import (
	"time"

	"github.com/koykov/indirect"
)

type HistogramChain interface {
	WithLabel(name, value string) HistogramChain
	Update(value float64)
	UpdateDuration(startTime time.Time)
	Reset()
}

type histogram struct {
	builder
	sptr uintptr
}

func (h *histogram) WithLabel(name, value string) HistogramChain {
	h.setLabel(name, value)
	return h
}

func (h *histogram) Update(value float64) {
	if s := h.indirectSet(); s != nil {
		defer s.releaseHistogram(h)
		s.getHistogram(h.commit()).Update(value)
	}
}

func (h *histogram) UpdateDuration(startTime time.Time) {
	if s := h.indirectSet(); s != nil {
		defer s.releaseHistogram(h)
		s.getHistogram(h.commit()).UpdateDuration(startTime)
	}
}

func (h *histogram) Reset() {
	if s := h.indirectSet(); s != nil {
		defer s.releaseHistogram(h)
		s.getHistogram(h.commit()).Reset()
	}
}

func (h *histogram) indirectSet() *chain {
	if h.sptr == 0 {
		return nil
	}
	return (*chain)(indirect.ToUnsafePtr(h.sptr))
}
