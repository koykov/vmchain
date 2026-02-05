package vmchain

import "github.com/VictoriaMetrics/metrics"

type Option func(c *chain)

// WithVMSet sets Victoria Metrics Set to use for metrics storage.
// Caution! You must register set using metrics.RegisterSet yourself.
func WithVMSet(vmset *metrics.Set) Option {
	return func(c *chain) {
		c.vmset = vmset
	}
}

func WithHasher(hasher Hasher) Option {
	return func(c *chain) {
		c.hash = hasher
	}
}
