# box

[![GoDoc](https://godoc.org/github.com/tidwall/box?status.svg)](https://godoc.org/github.com/tidwall/box)

Box is a Go package for wrapping value types.
Works similar to an `interface{}` and is optimized for primitives, strings, and byte slices.

## Features

- Zero new allocations for wrapping primitives, strings, and []byte slices.
- Uses a 128 bit structure. Same as `interface{}` on 64-bit architectures.
- Allows for auto convertions between value types. No panics on assertions.
- Pretty decent [performance](#performance).

## Examples

```go
// A boxed value can hold various types, just like an interface{}.
var v box.Value

// box an int
v = box.Int(123)

// unbox the value
println(v.Int())
println(v.String())
println(v.Bool())

// box a string
v = box.String("hello")
println(v.String())

// Auto conversions between types
println(box.String("123.45").Float64())
println(box.Bool(false).String())
println(box.String("hello").IsString())

// output
// 123
// 123
// true
// hello
// +1.234500e+002
// false
// true
```

## Performance

Below are some benchmarks comparing `interface{}` to `box.Value`.

- `Iface*/to`: Convert a value to an `interface{}`.
- `Iface*/from`: Convert an `interface{}` back to its original value.
- `Box*/to`: Convert a value to `box.Value`.
- `Box*/from`: Convert a `box.Value` back to its original value.

```
BenchmarkIfaceInt/to-10         128508648    8.632 ns/op    7 B/op   0 allocs/op
BenchmarkIfaceInt/from-10      1000000000    0.652 ns/op    0 B/op   0 allocs/op
BenchmarkBoxInt/to-10          1000000000    1.152 ns/op    0 B/op   0 allocs/op
BenchmarkBoxInt/from-10        1000000000    0.640 ns/op    0 B/op   0 allocs/op
BenchmarkIfaceString/to-10       58581472   17.120 ns/op   16 B/op   1 allocs/op
BenchmarkIfaceString/from-10   1000000000    6.781 ns/op    0 B/op   0 allocs/op
BenchmarkBoxString/to-10        474026008    3.736 ns/op    0 B/op   0 allocs/op
BenchmarkBoxString/from-10      492863490    2.416 ns/op    0 B/op   0 allocs/op
BenchmarkIfaceBytes/to-10        53937030   20.870 ns/op   24 B/op   1 allocs/op
BenchmarkIfaceBytes/from-10    1000000000    5.353 ns/op    0 B/op   0 allocs/op
BenchmarkBoxBytes/to-10         438596022    6.365 ns/op    0 B/op   0 allocs/op
BenchmarkBoxBytes/from-10       489872066    2.453 ns/op    0 B/op   0 allocs/op
```

