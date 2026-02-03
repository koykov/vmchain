# Victoria Metrics chains

Эта библиотека предлагает быстрое и удобное управление записью метрик в [VictoriaMetrics](https://github.com/VictoriaMetrics/metrics).
В чем идея: метрики VM неважно как они создаются, через `New*` конструкторы или через `GetOrCreate*` адаптеры принимают
только полное имя метрики в виде строки. Это не проблема, когда метрики не содержат никаких срезов или срезов заранее
известное количество с известными значениями, например:
```
foo
foo{bar="baz"}
foo{bar="baz",aaa="b"}
```
Но как только появляются динамические метки, то не остаётся иного пути кроме как пользоваться конкатенацией
```go
GetOrCreateCounter(`foo{stage="` + stageName + `",groupID="` + strconv.Itoa(groupID) + `"}`)
```
или пакетом `fmt`
```go
GetOrCreateCounter(fmt.Sprintf(`foo{stage="%s",groupID="%d"}`, stageName, groupID))
```

Оба эти способа, помимо того что ими неудобно пользоваться, делают аллокации, что приводит к нагрузке на GC и в целом
достаточно медленные.

Давайте посмотрим что предлагает эта библиотека взамен.

## Концепт решения

Начнём с типичного примера использования:
```go
import "github.com/koykov/vmchain"

func myfunc() {
    stageName := "auth"
    groupID := 123
    vmchain.Gauge("myservice_feature_counter").
        WithLabel("stage", stageName).
        WithAnyLabel("groupID", groupID).
        Add(321)
}
```

В результате выполнения этого кода будет создано имя метрики из стартового имени `"myservice_feature_counter"`, строковой
метки `stage` и числовой метки `groupID` и финальное имя будет таким `myservice_feature_counter{stage="auth",groupID="123"}`.
Далее если метрика по этому имени вызывается впервые, она будет зарегистрирована в VM и значение изменится. Если
метрика ранее была зарегистрирована, она просто поменяет своё значение.

Преимуществами использования библиотеки является:
* полное отсутствие аллокаций
* метод `WithAnyChain` избавляет от необходимости использования `strconv` пакета, т.е. от ещё одной аллокации
* метафора построения цепи - посредством цепочки вызовов методов `WithChain`/`WithAnyChain` нет необходимости прибегать к использованию конкатенаций или `fmt` пакета и метки добавляются в название метрики удобным и очевидным способом

Аллокациям уделяется такое повышенное внимание потому, что сама идея метрик заключается в том, чтобы замерить что происходит
в проекте, а не добавлять в него оверхед. Особенно важно это в высоконагруженных проектах, где каждая лишняя аллокация
в горячем пути может привести к проблемам с GC.

## API

В данный момент поддерживаются четыре основных типа метрик:
* [Gauge](gauge.go)
* [Counter](counter.go)
* [FloatCounter](float_counter.go)
* [Histogram](historgram.go)

Все эти обёртки объединены в сущности [Chain](chain.go), которая является хранилищем самих метрик, внутренних буферов и
прочих вспомогательных механизмов. Библиотека по умолчанию содержит уже инициализированный chain и предоставляет доступ
к нему через удобные функции `Gauge`, `Counter`, `FloatCounter` и `Histogram`, расположенные в [default.go](default.go).

Свой chain можно создать посредством функции `NewChain` и использовать нужным образом.

## Производительность

Проект [versus/vmchain](https://github.com/koykov/versus/tree/master/vmchain) демонстрирует сравнительные бенчмарки
четырёх участников соревнования:
* vmchain
* [strings.Builder](https://pkg.go.dev/strings#Builder)
* конкатенация
* [fmt.Sprintf](https://pkg.go.dev/fmt#Sprintf)

и демонстрирует такие показатели:
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

Более высокая скорость работы vmchain в первую очередь достигается отсутствием аллокаций, что в высоконагруженных
проектах является весьма существенным преимуществом.
