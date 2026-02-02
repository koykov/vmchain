package vmchain

import (
	"strconv"
	"testing"

	"github.com/koykov/x2bytes"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var b builder

		result := b.commit()
		assert.Equal(t, result, "")

		b.setName("metric")
		result = b.commit()
		assert.Equal(t, result, "metric")
	})
	t.Run("base", func(t *testing.T) {
		tests := []struct {
			name     string
			actions  func(*builder)
			expected string
		}{
			{
				name: "only name",
				actions: func(b *builder) {
					b.setName("metric_name")
				},
				expected: "metric_name",
			},
			{
				name: "one label",
				actions: func(b *builder) {
					b.setName("metric")
					b.setLabel("label1", "value1")
				},
				expected: `metric{label1="value1"}`,
			},
			{
				name: "multiple labels",
				actions: func(b *builder) {
					b.setName("http_requests")
					b.setLabel("method", "GET")
					b.setLabel("status", "200")
					b.setLabel("path", "/api/users")
				},
				expected: `http_requests{method="GET",status="200",path="/api/users"}`,
			},
			{
				name: "mixed labels",
				actions: func(b *builder) {
					b.setName("mixed_metric")
					b.setLabel("string_label", "test")
					b.setAnyLabel("int_label", 123)
					b.setLabel("another_string", "value")
					b.setAnyLabel("float_label", 45.67)
				},
				expected: `mixed_metric{string_label="test",int_label="123",another_string="value",float_label="45.67"}`,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				var b builder
				tc.actions(&b)
				result := b.commit()
				assert.Equal(t, result, tc.expected)
			})
		}
	})
	t.Run("reset", func(t *testing.T) {
		var b builder

		b.setName("metric1")
		b.setLabel("label1", "value1")
		result1 := b.commit()
		expected1 := `metric1{label1="value1"}`
		assert.Equal(t, result1, expected1)

		assert.True(t, len(b.buf) != 0 || b.lc != 0)

		b.setName("metric2")
		b.setLabel("label2", "value2")
		result2 := b.commit()
		expected2 := `metric2{label2="value2"}`
		assert.Equal(t, result2, expected2)
	})
	t.Run("any label", func(t *testing.T) {
		type customType string
		type customInt int

		x2bytes.RegisterToBytesFn(func(dst []byte, val any, _ ...any) ([]byte, error) {
			switch x := val.(type) {
			case customType:
				dst = append(dst, x...)
			case customInt:
				dst = strconv.AppendInt(dst, int64(x), 10)
			default:
				return dst, x2bytes.ErrUnknownType
			}
			return dst, nil
		})

		tests := []struct {
			name     string
			value    any
			expected string
		}{
			{"int", 42, `metric{label="42"}`},
			{"int8", int8(127), `metric{label="127"}`},
			{"int32", int32(1000), `metric{label="1000"}`},
			{"int64", int64(999999), `metric{label="999999"}`},
			{"uint", uint(42), `metric{label="42"}`},
			{"uint64", uint64(18446744073709551615), `metric{label="18446744073709551615"}`},
			{"float32", float32(3.140000104904175), `metric{label="3.140000104904175"}`},
			{"float64", 2.718281828, `metric{label="2.718281828"}`},
			{"bool true", true, `metric{label="true"}`},
			{"bool false", false, `metric{label="false"}`},
			{"string", "hello", `metric{label="hello"}`},
			{"byte slice", []byte("world"), `metric{label="world"}`},
			{"custom string type", customType("custom"), `metric{label="custom"}`},
			{"custom int type", customInt(777), `metric{label="777"}`},
			{"nil", nil, `metric{label="<nil>"}`},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var b builder
				b.setName("metric")
				b.setAnyLabel("label", tt.value)
				result := b.commit()
				assert.Equal(t, result, tt.expected)
			})
		}
	})
}

func BenchmarkBuilder(b *testing.B) {
	b.Run("base", func(b *testing.B) {
		b.ReportAllocs()
		var bb builder
		for i := 0; i < b.N; i++ {
			bb.reset()
			bb.setName("http_requests_total")
			bb.setLabel("method", "GET")
			bb.setLabel("status", "200")
			bb.setLabel("path", "/api/v1/users")
			_ = bb.commit()
		}
	})
	b.Run("any label", func(b *testing.B) {
		b.ReportAllocs()
		var bb builder
		for i := 0; i < b.N; i++ {
			bb.reset()
			bb.setName("custom_metric")
			bb.setLabel("string_label", "value")
			bb.setAnyLabel("int_label", 12345)
			bb.setAnyLabel("float_label", 3.14159)
			bb.setAnyLabel("bool_label", true)
			_ = bb.commit()
		}
	})
}
