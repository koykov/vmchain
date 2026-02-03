# Victoria Metrics chains

This library provides fast and convenient management of metrics recording in [VictoriaMetrics](https://github.com/VictoriaMetrics/metrics).
The idea: VM metrics, regardless of how they are created, through `New*` constructors or through `GetOrCreate*` adapters,
only accept the full metric name as a string. This is not a problem when metrics do not contain any labels or labels
have a known number with known values, for example:
```
foo
foo{bar="baz"}
foo{bar="baz",aaa="b"}
```
But as soon as dynamic labels appear, there is no other way than to use concatenation
```go
GetOrCreateCounter(`foo{stage="` + stageName + `",groupID="` + strconv.Itoa(groupID) + `"}`).Inc()
```
or `fmt` package
```go
GetOrCreateCounter(fmt.Sprintf(`foo{stage="%s",groupID="%d"}`, stageName, groupID)).Inc()
```

Both of these methods are inconvenient to use and produce allocations, which causes a GC pressure and are slow in general.

Let's see what this library offers instead.

## Solution Concept

Let's start with a typical usage example:
```go
import "github.com/koykov/vmchain"

func myfunc() {
    stageName := "auth"
    groupID := 123
    vmchain.Counter("myservice_feature_counter").
        WithLabel("stage", stageName).
        WithAnyLabel("groupID", groupID).
        Add(123)
}
```

As a result of this code, a metric name will be created from the starting name `"myservice_feature_counter"`, string
label `stage` and numeric label `groupID`, and the final name will be `myservice_feature_counter{stage="auth",groupID="123"}`.
Then, if the metric with this name is called for the first time, it will be registered in VM and its value will be updated.
If the metric has been registered before, it will simply change its value.

The advantages of using the library are:
* avoid to produce allocations
* the `WithAnyChain` method eliminates the need to use the `strconv` package, i.e., another allocation
* the chain-building approach - through a chain of calls to the `WithChain`/`WithAnyChain` methods there is no need to use concatenation or `fmt` package, and labels are added to the metric name in a convenient and obvious way

Allocations are given such heightened attention because the very idea of metrics is to measure what is happening
in the project, but not to add an overhead to it. This is especially important in high-load projects, where every extra
allocation in the hot path can lead to GC issues.

## API

Currently, four main metric types are supported:
* [Gauge](gauge.go)
* [Counter](counter.go)
* [FloatCounter](float_counter.go)
* [Histogram](historgram.go)

All these wrappers are combined into the [Chain](chain.go) entity, which is a storage for the metrics themselves, internal buffers, and
other auxiliary mechanisms. The library by default already contains an initialized chain and provides access
to it through convenient functions `Gauge`, `Counter`, `FloatCounter`, and `Histogram`, located in [default.go](default.go).

You can create your own chain using the `NewChain` function and use it as needed.

## Performance

The project [versus/vmchain](https://github.com/koykov/versus/tree/master/vmchain) demonstrates comparative benchmarks
of four competitors:
* vmchain
* [strings.Builder](https://pkg.go.dev/strings#Builder)
* concatenation
* [fmt.Sprintf](https://pkg.go.dev/fmt#Sprintf)

and shows the following results:
```
BenchmarkVMChain/normal-8    	        11848255	       101.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkVMChain/parallel-8             29461041	       45.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringsBuilder/normal-8     	 9529771	       143.8 ns/op	     147 B/op	       3 allocs/op
BenchmarkStringsBuilder/parallel-8   	 4162002	       437.9 ns/op	     147 B/op	       3 allocs/op
BenchmarkConcat/normal-8             	 7849113	       130.3 ns/op	      83 B/op	       2 allocs/op
BenchmarkConcat/parallel-8           	 5054778	       305.5 ns/op	      83 B/op	       2 allocs/op
BenchmarkFmt/normal-8                	 3316862	       318.0 ns/op	     112 B/op	       3 allocs/op
BenchmarkFmt/parallel-8              	 3262021	       512.3 ns/op	     112 B/op	       3 allocs/op
```

The higher speed of vmchain is primarily achieved due to allocation-free work, which is a very significant advantage in
high-load projects.
