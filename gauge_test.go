package vmchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGauge(t *testing.T) {
	fn := func() GaugeChain {
		return Gauge("myservice_feature", nil).
			WithLabel("groupID", "foobar").
			WithLabel("countryID", "123")
	}
	t.Run("add", func(t *testing.T) {
		fn().Add(10)
		fn().Add(10)
		v := fn().Get()
		assert.Equal(t, float64(20), v)
	})
	t.Run("set", func(t *testing.T) {
		fn().Set(15)
		v := fn().Get()
		assert.Equal(t, float64(15), v)
	})
	t.Run("inc", func(t *testing.T) {
		fn().Set(0)
		fn().Inc()
		fn().Inc()
		fn().Inc()
		fn().Inc()
		fn().Inc()
		v := fn().Get()
		assert.Equal(t, float64(5), v)
	})
	t.Run("dec", func(t *testing.T) {
		fn().Set(10)
		fn().Dec()
		fn().Dec()
		fn().Dec()
		fn().Dec()
		fn().Dec()
		v := fn().Get()
		assert.Equal(t, float64(5), v)
	})
}
