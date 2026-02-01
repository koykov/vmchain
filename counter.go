package vmchain

type Counter interface {
	WithLabel(name, value string) Counter
	Add(value int)
	AddInt64(value int64)
	Set(value uint64)
	Inc()
	Get() uint64
	Dec()
}
