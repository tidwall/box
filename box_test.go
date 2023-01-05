// Copyright 2023 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package box

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

type Jello struct {
	Neat int
	Feet int
}

type Pudding struct {
	Neat int
	Feet int
}

func (p Pudding) Float64() float64 {
	return float64(p.Neat) * float64(p.Feet)
}
func (p Pudding) Uint64() uint64 {
	return uint64(p.Neat) * uint64(p.Feet)
}
func (p Pudding) Int64() int64 {
	return int64(p.Neat) * int64(p.Feet)
}
func (p Pudding) Bool() bool {
	return true
}

func (p Pudding) String() string {
	return fmt.Sprintf("Yum{%d %d}", p.Neat, p.Feet)
}

func assert(cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestValue(t *testing.T) {
	assert(Nil().Int() == 0)
	assert(Nil().IsNil() == true)
	assert(Nil().IsCustomBits() == false)
	assert(CustomBits(0).IsNil() == false)
	assert(CustomBits(0).Int() == 0)
	assert(CustomBits(1).Int() == 1)
	assert(CustomBits(1).Uint64() == 1)
	assert(Bool(true).Int() == 1)
	assert(Bool(false).Int() == 0)
	assert(Int(0).Bool() == false)
	assert(Int(1).Bool() == true)
	assert(Float64(math.NaN()).Bool() == false)
	assert(CustomBits(1).Bool() == true)
	assert(CustomBits(0).Bool() == false)
	assert(Int(1).Int() == 1)
	assert(Float64(1.0).Int() == 1)
	assert(Float64(math.NaN()).Int() == 0)
	assert(Uint64(99).Int() == 99)
	assert(String("hello world").String() == "hello world")
	assert(String("hello world").Int() == 0)
	assert(String("hello world").IsNil() == false)
	assert(string(String("hello world").Bytes()) == "hello world")
	assert(Bytes([]byte("hello world")).String() == "hello world")
	assert(Bytes([]byte("hello world")).Int() == 0)
	assert(Bytes([]byte("hello world")).IsNil() == false)
	assert(string(Bytes([]byte("hello world")).Bytes()) == "hello world")
	forceIfaceStrs = true
	assert(String("hello world").String() == "hello world")
	assert(String("hello world").Int() == 0)
	assert(String("hello world").IsNil() == false)
	assert(string(String("hello world").Bytes()) == "hello world")
	assert(Bytes([]byte("hello world")).String() == "hello world")
	assert(Bytes([]byte("hello world")).Int() == 0)
	assert(Bytes([]byte("hello world")).IsNil() == false)
	assert(string(Bytes([]byte("hello world")).Bytes()) == "hello world")
	forceIfaceStrs = false
	assert(Any(Jello{1, 2}).IsNil() == false)
	assert(Any(Jello{1, 2}).String() == "{1 2}")
	assert(Any(Pudding{1, 2}).String() == "Yum{1 2}")
	assert(string(Any(Pudding{1, 2}).Bytes()) == "Yum{1 2}")
	forceIfacePtrs = true
	assert(Any(Jello{1, 2}).IsNil() == false)
	assert(Any(Jello{1, 2}).String() == "{1 2}")
	assert(Any(Jello{1, 2}).Any().(Jello).Feet == 2)
	assert(Any(Pudding{1, 2}).String() == "Yum{1 2}")
	assert(string(Any(Pudding{1, 2}).Bytes()) == "Yum{1 2}")
	forceIfacePtrs = false
	assert(Any(nil).IsNil())
	assert(Any("hello").String() == "hello")
	assert(Any([]byte("hello")).String() == "hello")
	assert(Any(true).Bool() == true)
	assert(Any(false).Bool() == false)
	assert(Any(int8(-1)).Int8() == -1)
	assert(Any(int16(-2)).Int16() == -2)
	assert(Any(int32(-3)).Int32() == -3)
	assert(Any(int64(-4)).Int64() == -4)
	assert(Any(uint8(1)).Int8() == 1)
	assert(Any(uint16(2)).Int16() == 2)
	assert(Any(uint32(3)).Int32() == 3)
	assert(Any(uint64(4)).Int64() == 4)
	assert(Any(int(1)).Int8() == 1)
	assert(Any(uint(2)).Int16() == 2)
	assert(Any(uintptr(3)).Int32() == 3)
	assert(Any(float32(4)).Float32() == 4)
	assert(Any(float64(5)).Float64() == 5)
	assert(Int(123).String() == "123")
	assert(string(Int(123).Bytes()) == "123")
	assert(Int(123).Any().(int64) == 123)
	assert(Any(Jello{1, 2}).Any().(Jello).Neat == 1)

	assert(CustomBits(99).String() == "99")
	assert(Bool(true).String() == "true")
	assert(Bool(false).String() == "false")
	assert(Uint64(99).String() == "99")
	assert(Int64(-99).String() == "-99")
	assert(Float64(-998).String() == "-998")
	assert(Nil().String() == "")

	assert(Any(CustomBits(99).Any()).String() == "99")
	assert(Any(Bool(true).Any()).String() == "true")
	assert(Any(Bool(false).Any()).String() == "false")
	assert(Any(Uint64(99).Any()).String() == "99")
	assert(Any(Int64(-99).Any()).String() == "-99")
	assert(Any(Float64(-998).Any()).String() == "-998")
	assert(Any(Nil().Any()).String() == "")

	assert(Int(99).Float64() == 99.0)
	assert(Nil().Float64() == 0)
	assert(CustomBits(1).Float64() == 1)
	assert(Bool(true).Float64() == 1)
	assert(Bool(false).Float64() == 0)
	assert(Uint64(98).Float64() == 98)
	assert(Int64(-98).Float64() == -98)
	assert(Float64(99).Float64() == 99.0)
	assert(Any("-99").Float64() == -99)
	assert(Any([]byte("-99")).Float64() == -99)
	assert(Any(interface{}(nil)).Float64() == 0)
	assert(Any(nil).Float64() == 0)
	assert(math.IsNaN(Any("hello").Float64()))
	assert(math.IsNaN(Any(Jello{10, 20}).Float64()))
	assert(Any(Pudding{10, 20}).Float64() == 200)

	assert(Int(99).Uint64() == 99.0)
	assert(Nil().Uint64() == 0)
	assert(CustomBits(1).Uint64() == 1)
	assert(Bool(true).Uint64() == 1)
	assert(Bool(false).Uint64() == 0)
	assert(Uint64(98).Uint64() == 98)
	assert(Int64(980).Uint64() == 980)
	assert(Float64(99).Uint64() == 99)
	assert(Any("990").Uint64() == 990)
	assert(Any([]byte("990")).Uint64() == 990)
	assert(Any(interface{}(nil)).Uint64() == 0)
	assert(Any(nil).Uint64() == 0)
	assert(Any("hello").Uint64() == 0)
	assert(Any(Jello{10, 20}).Uint64() == 0)
	assert(Any(Pudding{10, 20}).Uint64() == 200)

	assert(Int(99).Int64() == 99.0)
	assert(Nil().Int64() == 0)
	assert(CustomBits(1).Int64() == 1)
	assert(Bool(true).Int64() == 1)
	assert(Bool(false).Int64() == 0)
	assert(Uint64(98).Int64() == 98)
	assert(Int64(-98).Int64() == -98)
	assert(Float64(99).Int64() == 99.0)
	assert(Any("-99").Int64() == -99)
	assert(Any([]byte("-99")).Int64() == -99)
	assert(Any(interface{}(nil)).Int64() == 0)
	assert(Any(nil).Int64() == 0)
	assert(Any("hello").Int64() == 0)
	assert(Any(Jello{10, 20}).Int64() == 0)
	assert(Any(Pudding{10, 20}).Int64() == 200)

	assert(Int(99).Bool() == true)
	assert(Nil().Bool() == false)
	assert(CustomBits(1).Bool() == true)
	assert(Bool(true).Bool() == true)
	assert(Bool(false).Bool() == false)
	assert(Uint64(98).Bool() == true)
	assert(Int64(-98).Bool() == true)
	assert(Float64(99).Bool() == true)
	assert(Any("-99").Bool() == false)
	assert(Any("true").Bool() == true)
	assert(Any([]byte("-99")).Bool() == false)
	assert(Any([]byte("true")).Bool() == true)
	assert(Any(interface{}(nil)).Bool() == false)
	assert(Any(nil).Bool() == false)
	assert(Any("hello").Bool() == false)
	assert(Any(Jello{10, 20}).Bool() == false)
	assert(Any(Pudding{10, 20}).Bool() == true)

	assert(Any(nil).IsString() == false)
	assert(Any(123).IsString() == false)
	assert(Any("hello").IsString() == true)
	assert(Any([]byte("hello")).IsString() == false)
	forceIfaceStrs = true
	assert(Any("hello").IsString() == true)
	assert(Any([]byte("hello")).IsString() == false)
	forceIfaceStrs = false

	assert(Any(nil).IsBytes() == false)
	assert(Any(123).IsBytes() == false)
	assert(Any("hello").IsBytes() == false)
	assert(Any([]byte("hello")).IsBytes() == true)
	forceIfaceStrs = true
	assert(Any("hello").IsBytes() == false)
	assert(Any([]byte("hello")).IsBytes() == true)
	forceIfaceStrs = false

	assert(Int8(-10).Int8() == -10)
	assert(Int(500).Int8() == -12)
	assert(Int16(-10).Int16() == -10)
	assert(Int32(-10).Int32() == -10)
	assert(Int64(-10).Int64() == -10)
	assert(Int64(-10).Float32() == -10)
	assert(Float32(10.1239123).Float32() == 10.1239123)

	assert(Uint8(10).Uint8() == 10)
	assert(Uint(500).Uint8() == 500&0xFF)
	assert(Uint16(10).Uint16() == 10)
	assert(Uint32(11).Uint32() == 11)
	assert(Uint64(12).Uint64() == 12)
	assert(Uint64(12).Uint() == 12)

	assert(Uint64(10).IsUint() == true)
	assert(Uint8(10).IsUint() == true)
	assert(Int64(10).IsUint() == false)

	assert(Int64(10).IsInt() == true)
	assert(Int8(10).IsInt() == true)
	assert(Uint64(10).IsInt() == false)

	assert(Float64(10).IsFloat() == true)
	assert(Float32(10).IsFloat() == true)
	assert(Uint64(10).IsFloat() == false)

	assert(Bool(true).IsBool() == true)
	assert(Bool(false).IsBool() == true)
	assert(Uint64(10).IsBool() == false)

	assert(Byte(10).Byte() == 10)
	assert(Uint64(257).Byte() == 1)

	assert(String("hello").IsNumber() == false)
	assert(Int(10).IsNumber() == true)
	assert(Uint(10).IsNumber() == true)
	assert(Float64(10).IsNumber() == true)
	assert(Any(10).IsNumber() == true)

	assert(Uint64(10).Tag() == 0)
	assert(Bytes(nil).Tag() == 0)
	assert(Bytes([]byte{}).Tag() == 0)

	assert(String("hello").Tag() == 0)
	assert(StringWithTag("hello", 999).Tag() == 999)
	assert(StringWithTag("hello", 999).String() == "hello")
	forceIfaceStrs = true
	assert(String("hello").Tag() == 0)
	assert(StringWithTag("hello", 999).Tag() == 999)
	assert(StringWithTag("hello", 999).String() == "hello")
	forceIfaceStrs = false

}

func TestUnits(t *testing.T) {
	assert(Float64(-98).toFloat64() == -98)
	assert(Uint64(98).toUint64() == 98)
	assert(Int64(-98).toInt64() == -98)
	assert(Bool(true).toBool() == true)
	assert(Bool(false).toBool() == false)
}

func TestPLocks(t *testing.T) {
	// Tests the psave() with plock/punlock using multiple goroutines.
	// Best if used with -race
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			start := time.Now()
			for time.Since(start) < time.Second/10 {
				Any(&Jello{1, 2})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkIfaceInt(b *testing.B) {
	gen := func(b *testing.B, reset bool) []interface{} {
		arr := make([]interface{}, b.N)
		if reset {
			b.ResetTimer()
			b.ReportAllocs()
		}
		for i := 0; i < b.N; i++ {
			arr[i] = i
		}
		return arr
	}
	b.Run("to", func(b *testing.B) {
		gen(b, true)
	})
	b.Run("from", func(b *testing.B) {
		arr := gen(b, false)
		b.ReportAllocs()
		b.ResetTimer()
		var res int
		for i := 0; i < b.N; i++ {
			res += arr[i].(int)
		}
	})
}

func BenchmarkBoxInt(b *testing.B) {
	gen := func(b *testing.B, reset bool) []Value {
		arr := make([]Value, b.N)
		if reset {
			b.ResetTimer()
			b.ReportAllocs()
		}
		for i := 0; i < b.N; i++ {
			arr[i] = Int(i)
		}
		return arr
	}
	b.Run("to", func(b *testing.B) {
		gen(b, true)
	})
	b.Run("from", func(b *testing.B) {
		arr := gen(b, false)
		b.ReportAllocs()
		b.ResetTimer()
		var res int
		for i := 0; i < b.N; i++ {
			res += arr[i].Int()
		}
	})
}

func BenchmarkIfaceString(b *testing.B) {
	gen := func(b *testing.B, reset bool) []interface{} {
		strs := make([]string, b.N)
		for i := 0; i < b.N; i++ {
			strs[i] = fmt.Sprint(i)
		}
		arr := make([]interface{}, b.N)
		if reset {
			b.ResetTimer()
			b.ReportAllocs()
		}
		for i := 0; i < b.N; i++ {
			arr[i] = strs[i]
		}
		return arr
	}
	b.Run("to", func(b *testing.B) {
		gen(b, true)
	})
	b.Run("from", func(b *testing.B) {
		arr := gen(b, false)
		b.ResetTimer()
		b.ReportAllocs()
		var n int
		for i := 0; i < b.N; i++ {
			s := arr[i].(string)
			n += int(s[0]) + int(s[len(s)-1])
		}
	})
}

func BenchmarkBoxString(b *testing.B) {
	gen := func(b *testing.B, reset bool) []Value {
		strs := make([]string, b.N)
		for i := 0; i < b.N; i++ {
			strs[i] = fmt.Sprint(i)
		}
		arr := make([]Value, b.N)
		if reset {
			b.ResetTimer()
			b.ReportAllocs()
		}
		for i := 0; i < b.N; i++ {
			arr[i] = String(strs[i])
		}
		return arr
	}
	b.Run("to", func(b *testing.B) {
		gen(b, true)
	})
	b.Run("from", func(b *testing.B) {
		arr := gen(b, false)
		b.ResetTimer()
		b.ReportAllocs()
		var n int
		for i := 0; i < b.N; i++ {
			s := arr[i].String()
			n += int(s[0]) + int(s[len(s)-1])
		}
	})
}

func BenchmarkIfaceBytes(b *testing.B) {
	gen := func(b *testing.B, reset bool) []interface{} {
		strs := make([][]byte, b.N)
		for i := 0; i < b.N; i++ {
			strs[i] = []byte(fmt.Sprint(i))
		}
		arr := make([]interface{}, b.N)
		if reset {
			b.ResetTimer()
			b.ReportAllocs()
		}
		for i := 0; i < b.N; i++ {
			arr[i] = strs[i]
		}
		return arr
	}
	b.Run("to", func(b *testing.B) {
		gen(b, true)
	})
	b.Run("from", func(b *testing.B) {
		arr := gen(b, false)
		b.ResetTimer()
		b.ReportAllocs()
		var n int
		for i := 0; i < b.N; i++ {
			n += len(arr[i].([]byte))
		}
	})
}

func BenchmarkBoxBytes(b *testing.B) {
	gen := func(b *testing.B, reset bool) []Value {
		strs := make([][]byte, b.N)
		for i := 0; i < b.N; i++ {
			strs[i] = []byte(fmt.Sprint(i))
		}
		arr := make([]Value, b.N)
		if reset {
			b.ResetTimer()
			b.ReportAllocs()
		}
		for i := 0; i < b.N; i++ {
			arr[i] = Bytes(strs[i])
		}
		return arr
	}
	b.Run("to", func(b *testing.B) {
		gen(b, true)
	})
	b.Run("from", func(b *testing.B) {
		arr := gen(b, false)
		b.ResetTimer()
		b.ReportAllocs()
		var n int
		for i := 0; i < b.N; i++ {
			n += len(arr[i].Bytes())
		}
	})
}
