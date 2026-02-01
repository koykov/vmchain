package vmchain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func cfn() CounterChain {
	return Counter("myservice_feature_counter").
		WithLabel("groupID", "foobar").
		WithLabel("countryID", "123")
}

func TestCounter(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		cfn().Set(0)
		cfn().Add(10)
		cfn().Add(10)
		v := cfn().Get()
		assert.Equal(t, uint64(20), v)
	})
	t.Run("addInt64", func(t *testing.T) {
		cfn().Set(0)
		cfn().AddInt64(10)
		cfn().AddInt64(10)
		v := cfn().Get()
		assert.Equal(t, uint64(20), v)
	})
	t.Run("set", func(t *testing.T) {
		cfn().Set(15)
		v := cfn().Get()
		assert.Equal(t, uint64(15), v)
	})
	t.Run("inc", func(t *testing.T) {
		cfn().Set(0)
		cfn().Inc()
		cfn().Inc()
		cfn().Inc()
		cfn().Inc()
		cfn().Inc()
		v := cfn().Get()
		assert.Equal(t, uint64(5), v)
	})
	t.Run("dec", func(t *testing.T) {
		cfn().Set(10)
		cfn().Dec()
		cfn().Dec()
		cfn().Dec()
		cfn().Dec()
		cfn().Dec()
		v := cfn().Get()
		assert.Equal(t, uint64(5), v)
	})
}

func BenchmarkCounter(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cfn().Add(1)
		}
	})
	b.Run("set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cfn().Add(i)
		}
	})
	b.Run("inc", func(b *testing.B) {
		b.ReportAllocs()
		cfn().Set(0)
		for i := 0; i < b.N; i++ {
			cfn().Inc()
		}
	})
	b.Run("dec", func(b *testing.B) {
		b.ReportAllocs()
		cfn().Set(math.MaxUint64)
		for i := 0; i < b.N; i++ {
			cfn().Dec()
		}
	})
}

func BenchmarkCounterParallel(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				cfn().Add(1)
			}
		})
	})
	b.Run("set", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				cfn().Set(123)
			}
		})
	})
	b.Run("inc", func(b *testing.B) {
		b.ReportAllocs()
		cfn().Set(0)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				cfn().Inc()
			}
		})
	})
	b.Run("dec", func(b *testing.B) {
		b.ReportAllocs()
		cfn().Set(math.MaxUint64)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				cfn().Dec()
			}
		})
	})
}
