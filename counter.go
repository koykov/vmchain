package vmchain

type CounterChain interface {
	WithLabel(name, value string) CounterChain
	Add(value int)
	AddInt64(value int64)
	Set(value uint64)
	Inc()
	Get() uint64
	Dec()
}
