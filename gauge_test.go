package vmchain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func gfn() GaugeChain {
	return Gauge("myservice_feature", nil).
		WithLabel("groupID", "foobar").
		WithLabel("countryID", "123")
}

func TestGauge(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		gfn().Add(10)
		gfn().Add(10)
		v := gfn().Get()
		assert.Equal(t, float64(20), v)
	})
	t.Run("set", func(t *testing.T) {
		gfn().Set(15)
		v := gfn().Get()
		assert.Equal(t, float64(15), v)
	})
	t.Run("inc", func(t *testing.T) {
		gfn().Set(0)
		gfn().Inc()
		gfn().Inc()
		gfn().Inc()
		gfn().Inc()
		gfn().Inc()
		v := gfn().Get()
		assert.Equal(t, float64(5), v)
	})
	t.Run("dec", func(t *testing.T) {
		gfn().Set(10)
		gfn().Dec()
		gfn().Dec()
		gfn().Dec()
		gfn().Dec()
		gfn().Dec()
		v := gfn().Get()
		assert.Equal(t, float64(5), v)
	})
}

func BenchmarkGauge(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			gfn().Add(1)
		}
	})
	b.Run("set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			gfn().Add(float64(i))
		}
	})
	b.Run("inc", func(b *testing.B) {
		b.ReportAllocs()
		gfn().Set(0)
		for i := 0; i < b.N; i++ {
			gfn().Inc()
		}
	})
	b.Run("dec", func(b *testing.B) {
		b.ReportAllocs()
		gfn().Set(math.MaxFloat64)
		for i := 0; i < b.N; i++ {
			gfn().Dec()
		}
	})
}

func BenchmarkGaugeParallel(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				gfn().Add(1)
			}
		})
	})
	b.Run("set", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				gfn().Set(123)
			}
		})
	})
	b.Run("inc", func(b *testing.B) {
		b.ReportAllocs()
		gfn().Set(0)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				gfn().Inc()
			}
		})
	})
	b.Run("dec", func(b *testing.B) {
		b.ReportAllocs()
		gfn().Set(math.MaxFloat64)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				gfn().Dec()
			}
		})
	})
}
