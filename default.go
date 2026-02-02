package vmchain

var defaultChain = NewChain()

// Gauge return existing or create and return new gauge metric.
//
// initName is a base name of a metric (without any labels). It must be valid Prometheus-compatible name.
// Labels can be added separately using WithLabel chain method:
//
// vmchain.Gauge("my_gauge_metric_name").	// prepare and return metric with name "my_gauge_metric_name"
//
//	WithLabel("stage", "area").			// add a label, so metric name became "my_gauge_metric_name{stage="area"}
//	WithLabel("userID", "123).			// add a label, so metric name became "my_gauge_metric_name{stage="area",userID="123"}
//	Inc()								// finally construct full name of underlying gauge metric, register it if necessary,
//										// and call method Inc.
func Gauge(initName string, f func() float64) GaugeChain {
	return defaultChain.Gauge(initName, f)
}

// Counter return existing or create and return new counter metric.
//
// initName is a base name of a metric (without any labels). It must be valid Prometheus-compatible name.
// Labels can be added separately using WithLabel chain method:
//
// vmchain.Counter("my_counter_metric_name").	// prepare and return metric with name "my_counter_metric_name"
//
//	WithLabel("stage", "area").			// add a label, so metric name became "my_counter_metric_name{stage="area"}
//	WithLabel("userID", "123).			// add a label, so metric name became "my_counter_metric_name{stage="area",userID="123"}
//	Inc()								// finally construct full name of underlying counter metric, register it if necessary,
//										// and call method Inc.
func Counter(initName string) CounterChain {
	return defaultChain.Counter(initName)
}

// Histogram return existing or create and return new histogram metric.
//
// initName is a base name of a metric (without any labels). It must be valid Prometheus-compatible name.
// Labels can be added separately using WithLabel chain method:
//
// vmchain.Histogram("my_histogram_metric_name").	// prepare and return metric with name "my_histogram_metric_name"
//
//	WithLabel("stage", "area").			// add a label, so metric name became "my_histogram_metric_name{stage="area"}
//	WithLabel("userID", "123).			// add a label, so metric name became "my_histogram_metric_name{stage="area",userID="123"}
//	Inc()								// finally construct full name of underlying gauge metric, register it if necessary,
//										// and call method Inc.
func Histogram(initName string) HistogramChain {
	return defaultChain.Histogram(initName)
}
