package vmchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func hfn() HistogramChain {
	return Histogram("myservice_feature_histogram").
		WithLabel("groupID", "foobar").
		WithLabel("countryID", "123")
}

func TestHistogram(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		hfn().Reset()
		hfn().Update(10)
		hfn().Update(20)
		var i int
		hfn().VisitNonZeroBuckets(func(vmrange string, count uint64) {
			switch i {
			case 0:
				assert.Equal(t, vmrange, "8.799e+00...1.000e+01")
				assert.Equal(t, count, uint64(1))
			case 1:
				assert.Equal(t, vmrange, "1.896e+01...2.154e+01")
				assert.Equal(t, count, uint64(1))
			}
			i++
		})
	})
	t.Run("update duration", func(t *testing.T) {
		hfn().Reset()
		tm, _ := time.Parse(time.DateTime, time.DateTime)
		hfn().UpdateDuration(tm)
		hfn().UpdateDuration(tm.Add(time.Hour * 24 * 30))
		hfn().VisitNonZeroBuckets(func(vmrange string, count uint64) {
			assert.Equal(t, vmrange, "5.995e+08...6.813e+08")
			assert.Equal(t, count, uint64(2))
		})
	})
}

func BenchmarkHistogram(b *testing.B) {
	b.Run("update", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hfn().Update(1)
		}
	})
	b.Run("update duration", func(b *testing.B) {
		b.ReportAllocs()
		tm, _ := time.Parse(time.DateTime, time.DateTime)
		for i := 0; i < b.N; i++ {
			hfn().UpdateDuration(tm)
		}
	})
}

func BenchmarkHistogramParallel(b *testing.B) {
	b.Run("update", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				hfn().Update(1)
			}
		})
	})
	b.Run("update duration", func(b *testing.B) {
		b.ReportAllocs()
		tm, _ := time.Parse(time.DateTime, time.DateTime)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				hfn().UpdateDuration(tm)
			}
		})
	})
}
