package vmchain

type HistogramChain interface {
	WithLabel(name, value string) HistogramChain
	Update(value float64)
	Reset()
}
