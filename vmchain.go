package vmchain

var defaultChain = NewChain()

func Gauge(initName string, f func() float64) GaugeChain {
	return defaultChain.Gauge(initName, f)
}

func Counter(initName string) CounterChain {
	return defaultChain.Counter(initName)
}

func Histogram(initName string) HistogramChain {
	return defaultChain.Histogram(initName)
}
