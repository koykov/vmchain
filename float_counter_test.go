package vmchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ffn() FloatCounterChain {
	return FloatCounter("myservice_feature_float_counter").
		WithLabel("groupID", "foobar").
		WithLabel("countryID", "123")
}

func TestFloatCounter(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		ffn().Set(0)
		ffn().Add(10)
		ffn().Add(10)
		v := ffn().Get()
		assert.Equal(t, float64(20), v)
	})
	t.Run("sub", func(t *testing.T) {
		ffn().Set(100)
		ffn().Sub(10)
		ffn().Sub(10)
		v := ffn().Get()
		assert.Equal(t, float64(80), v)
	})
	t.Run("set", func(t *testing.T) {
		ffn().Set(3.14)
		v := ffn().Get()
		assert.Equal(t, 3.14, v)
	})
}

func BenchmarkFloatCounter(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ffn().Add(1)
		}
	})
	b.Run("set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ffn().Add(float64(i))
		}
	})
}

func BenchmarkFloatCounterParallel(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ffn().Add(1)
			}
		})
	})
	b.Run("set", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ffn().Set(123)
			}
		})
	})
}
