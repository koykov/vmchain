package vmchain

type Histogram interface {
	WithLabel(name, value string) Histogram
	Update(value float64)
	Reset()
}
