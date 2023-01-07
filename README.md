# box

[![GoDoc](https://godoc.org/github.com/tidwall/box?status.svg)](https://godoc.org/github.com/tidwall/box)

**experimental**

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
goos: darwin
goarch: arm64
pkg: github.com/tidwall/box
BenchmarkIfaceInt/to-10        10000000    8.921 ns/op    7 B/op   0 allocs/op
BenchmarkIfaceInt/from-10      10000000    0.6289 ns/op   0 B/op   0 allocs/op
BenchmarkBoxInt/to-10          10000000    1.334 ns/op    0 B/op   0 allocs/op
BenchmarkBoxInt/from-10        10000000    0.6823 ns/op   0 B/op   0 allocs/op
BenchmarkIfaceString/to-10     10000000   18.17 ns/op    16 B/op   1 allocs/op
BenchmarkIfaceString/from-10   10000000    0.8010 ns/op   0 B/op   0 allocs/op
BenchmarkBoxString/to-10       10000000    3.705 ns/op    0 B/op   0 allocs/op
BenchmarkBoxString/from-10     10000000    2.421 ns/op    0 B/op   0 allocs/op
BenchmarkIfaceBytes/to-10      10000000   21.62 ns/op    24 B/op   1 allocs/op
BenchmarkIfaceBytes/from-10    10000000    0.8104 ns/op   0 B/op   0 allocs/op
BenchmarkBoxBytes/to-10        10000000    2.881 ns/op    0 B/op   0 allocs/op
BenchmarkBoxBytes/from-10      10000000    2.366 ns/op    0 B/op   0 allocs/op
```
