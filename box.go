// Copyright 2023 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package box

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync/atomic"
	"unsafe"
)

var primTypes = [...]byte{0, 1, 2, 3, 4, 5}

var (
	boolType     = unsafe.Pointer(&primTypes[0])
	int64Type    = unsafe.Pointer(&primTypes[1])
	uint64Type   = unsafe.Pointer(&primTypes[2])
	float64Type  = unsafe.Pointer(&primTypes[3])
	custBitsType = unsafe.Pointer(&primTypes[4])
)

func isPrim(ptr unsafe.Pointer) bool {
	return ptr == nil || (uintptr(ptr) >= uintptr(boolType) &&
		uintptr(ptr) <= uintptr(custBitsType))
}

// Value is a boxed value
type Value struct {
	ext uint64
	ptr unsafe.Pointer
}

// Nil boxes a nil.
// This is the same as the default `box.Value{}`.
func Nil() Value {
	return Value{0, nil}
}

// Bool boxes a bool
func Bool(t bool) Value {
	return Value{uint64(*(*byte)(unsafe.Pointer(&t))), boolType}
}

// Int64 boxes an int64
func Int64(x int64) Value {
	return Value{uint64(x), int64Type}
}

// Uint64 boxes a uint64
func Uint64(x uint64) Value {
	return Value{uint64(x), uint64Type}
}

// Float64 boxes a float64
func Float64(f float64) Value {
	return Value{uint64(math.Float64bits(f)), float64Type}
}

// CustomBits boxes a custom value.
func CustomBits(x uint64) Value {
	return Value{x, custBitsType}
}

var plocker uint64
var ptable map[unsafe.Pointer]struct{}

func plock() {
	for !atomic.CompareAndSwapUint64(&plocker, 0, 1) {
		runtime.Gosched()
	}
}
func punlock() {
	atomic.StoreUint64(&plocker, 0)
}

func psave(p unsafe.Pointer) {
	plock()
	if _, ok := ptable[p]; !ok {
		if ptable == nil {
			ptable = make(map[unsafe.Pointer]struct{})
		}
		ptable[p] = struct{}{}
	}
	punlock()
}

type (
	booler    interface{ Bool() bool }
	int64er   interface{ Int64() int64 }
	uint64er  interface{ Uint64() uint64 }
	float64er interface{ Float64() float64 }
)

type sface struct {
	ptr unsafe.Pointer
	len int
}

type bface struct {
	ptr unsafe.Pointer
	len int
	cap int
}

// maxLen is the maximum length for strings or byte-slices
const maxLen uint64 = 0x7FFFFFFF // int32 -> 2147483647 bytes

// maxCap is the maximum capacity above the length for byte-slices.
const maxCap uint64 = 0x7FFFFF // int24 -> 8388607 bytes

var forceIfaceStrs = false
var forceIfacePtrs = false

// non-primitive types
const (
	_ = iota
	ptrString
	ptrBytes
	ptrIface
	ptrIfacePtr
)

// String boxes a string value
func String(s string) Value {
	slen := uint64((*sface)(unsafe.Pointer(&s)).len)
	if forceIfaceStrs || slen > maxLen {
		return toIface(s)
	}
	return Value{
		ext: (slen << 32) | ptrString,
		ptr: (*sface)(unsafe.Pointer(&s)).ptr,
	}
}

type taggedString struct {
	tag uint16
	str string
}

func (ts *taggedString) String() string {
	return ts.str
}

func StringWithTag(s string, tag uint16) Value {
	slen := uint64((*sface)(unsafe.Pointer(&s)).len)
	if forceIfaceStrs || slen > maxLen {
		return toIface(&taggedString{tag: tag, str: s})
	}
	return Value{
		ext: (slen << 32) | (uint64(tag) << 8) | ptrString,
		ptr: (*sface)(unsafe.Pointer(&s)).ptr,
	}
}

// Bytes boxes a byte slice
func Bytes(b []byte) Value {
	blen := uint64(len(b))
	bcap := uint64(cap(b))
	if forceIfaceStrs || blen > maxLen || bcap-blen > maxCap {
		return toIface(b)
	}

	return Value{
		ext: (blen << 32) | (bcap-blen)<<8 | ptrBytes,
		ptr: (*bface)(unsafe.Pointer(&b)).ptr,
	}
}

func toIface(v any) Value {
	typ := (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[0]
	ptr := (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1]
	if !forceIfacePtrs && uint64(uintptr(typ)) < uint64(1)<<56 {
		// The interface type pointer is small enough to fit into 56 bits.
		// Save the type and tag the pointer
		psave(typ)
		return Value{(uint64(uintptr(typ)) << 8) | ptrIface, ptr}
	}
	// The interface type is a pointer in the heap or its pointer is too
	// large to store in 56 bits.
	// Use a pointer to the interface.
	return Value{ptrIfacePtr, unsafe.Pointer(&v)}
}

// Any boxes anything
func Any(v any) Value {
	switch v := v.(type) {
	case nil:
		return Nil()
	case string:
		return String(v)
	case []byte:
		return Bytes(v)
	case bool:
		return Bool(v)
	case int8:
		return Int64(int64(v))
	case int16:
		return Int64(int64(v))
	case int32:
		return Int64(int64(v))
	case int64:
		return Int64(int64(v))
	case uint8:
		return Uint64(uint64(v))
	case uint16:
		return Uint64(uint64(v))
	case uint32:
		return Uint64(uint64(v))
	case uint64:
		return Uint64(uint64(v))
	case int:
		return Int64(int64(v))
	case uint:
		return Uint64(uint64(v))
	case uintptr:
		return Uint64(uint64(v))
	case float32:
		return Float64(float64(v))
	case float64:
		return Float64(v)
	}
	return toIface(v)
}

func (v Value) isPrim() bool {
	return isPrim(v.ptr)
}

func (v Value) assertString() string {
	return *(*string)(unsafe.Pointer(&sface{
		ptr: unsafe.Pointer(v.ptr),
		len: int(v.ext >> 32),
	}))
}

func (v Value) assertBytes() []byte {
	blen := int(v.ext >> 32)
	bcap := int((v.ext >> 8) & maxCap)
	return *(*[]byte)(unsafe.Pointer(&bface{
		ptr: unsafe.Pointer(v.ptr),
		len: blen,
		cap: blen + bcap,
	}))
}

func (v Value) assertIfacePtr() any {
	return *(*any)(v.ptr)
}

func (v Value) assertIface() any {
	return *(*any)(unsafe.Pointer(&[2]uintptr{
		uintptr(v.ext >> 8),
		uintptr(v.ptr),
	}))
}

// String returns the value as a string.
func (v Value) String() string {
	if !v.isPrim() {
		if v.ext&0xFF == ptrString {
			return v.assertString()
		}
		if v.ext&0xFF == ptrBytes {
			return string(v.assertBytes())
		}
		var vf any
		if v.ext&0xFF == ptrIface {
			vf = v.assertIface()
		} else if v.ext&0xFF == ptrIfacePtr {
			vf = v.assertIfacePtr()
		}
		switch vf := vf.(type) {
		case []byte:
			return string(vf)
		case string:
			return vf
		}
		return fmt.Sprint(vf)
	}
	return v.primToString()
}

// Bytes returns the value as a byte slice.
// When the boxed value is a `[]byte` then those original bytes are returned.
// Otherwise, the string representation of the value is returned, which will
// be equivalent to `[]byte(value.String())`.
func (v Value) Bytes() []byte {
	if !v.isPrim() {
		if v.ext&0xFF == ptrBytes {
			return v.assertBytes()
		}
		if v.ext&0xFF == ptrString {
			return []byte(v.assertString())
		}
		var vf any
		if v.ext&0xFF == ptrIface {
			vf = v.assertIface()
		} else if v.ext&0xFF == ptrIfacePtr {
			vf = v.assertIfacePtr()
		}
		switch vf := vf.(type) {
		case []byte:
			return vf
		case string:
			return []byte(vf)
		}
		return []byte(fmt.Sprint(vf))
	}
	return v.primToBytes()
}

func (v Value) assertNonPrimAny() any {
	if v.ext&0xFF == ptrIface {
		return v.assertIface()
	}
	if v.ext&0xFF == ptrIfacePtr {
		return v.assertIfacePtr()
	}
	if v.ext&0xFF == ptrString {
		return v.assertString()
	}
	return v.assertBytes()
}

// Any returns the value as an `any/interface{}` type.
func (v Value) Any() any {
	if !v.isPrim() {
		return v.assertNonPrimAny()
	}
	return v.primToAny()
}

func (v Value) primToBytes() []byte {
	return []byte(v.primToString())
}

func (v Value) primToString() string {
	switch v.ptr {
	case boolType:
		return strconv.FormatBool(v.ext != 0)
	case int64Type:
		return strconv.FormatInt(int64(v.ext), 10)
	case uint64Type:
		return strconv.FormatUint(v.ext, 10)
	case float64Type:
		return strconv.FormatFloat(math.Float64frombits(v.ext), 'f', -1, 64)
	case custBitsType:
		return strconv.FormatUint(v.ext, 10)
	}
	return "" // nil
}

func (v Value) primToAny() any {
	switch v.ptr {
	case boolType:
		return v.ext != 0
	case int64Type:
		return int64(v.ext)
	case uint64Type:
		return uint64(v.ext)
	case float64Type:
		return math.Float64frombits(v.ext)
	case custBitsType:
		return uint64(v.ext)
	}
	return nil // nil
}

// Float64 returns the value as a float64
func (v Value) Float64() float64 {
	if v.ptr == float64Type {
		return math.Float64frombits(v.ext)
	}
	return v.toFloat64()
}

func (v Value) toFloat64() float64 {
	switch {
	case v.ptr == nil:
		return 0
	case v.ptr == boolType:
		if v.ext == 0 {
			return 0.0
		}
		return 1.0
	case v.ptr == int64Type:
		return float64(int64(v.ext))
	case v.ptr == uint64Type:
		return float64(v.ext)
	case v.ptr == float64Type:
		return math.Float64frombits(v.ext)
	case v.ptr == custBitsType:
		return float64(v.ext)
	}
	switch v := v.assertNonPrimAny().(type) {
	case string:
		x, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return x
		}
	case []byte:
		x, err := strconv.ParseFloat(string(v), 64)
		if err == nil {
			return x
		}
	case float64er:
		return v.Float64()
	}
	return math.NaN()
}

// Uint64 returns the value as a uint64
func (v Value) Uint64() uint64 {
	if v.ptr == uint64Type {
		return v.ext
	}
	return v.toUint64()
}

func (v Value) toUint64() uint64 {
	switch {
	case v.ptr == nil:
		return 0
	case v.ptr == boolType:
		if v.ext == 0 {
			return 0.0
		}
		return 1.0
	case v.ptr == int64Type:
		return v.ext
	case v.ptr == uint64Type:
		return v.ext
	case v.ptr == float64Type:
		return uint64(math.Float64frombits(v.ext))
	case v.ptr == custBitsType:
		return v.ext
	}
	switch v := v.assertNonPrimAny().(type) {
	case string:
		x, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			return x
		}
	case []byte:
		x, err := strconv.ParseUint(string(v), 10, 64)
		if err == nil {
			return x
		}
	case uint64er:
		return v.Uint64()
	}
	return 0
}

// Int64 returns the value as an int64
func (v Value) Int64() int64 {
	if v.ptr == int64Type {
		return int64(v.ext)
	}
	return v.toInt64()
}

func (v Value) toInt64() int64 {
	switch {
	case v.ptr == nil:
		return 0
	case v.ptr == boolType:
		if v.ext == 0 {
			return 0.0
		}
		return 1.0
	case v.ptr == int64Type:
		return int64(v.ext)
	case v.ptr == uint64Type:
		return int64(v.ext)
	case v.ptr == float64Type:
		return int64(math.Float64frombits(v.ext))
	case v.ptr == custBitsType:
		return int64(v.ext)
	}
	switch v := v.assertNonPrimAny().(type) {
	case string:
		x, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return x
		}
	case []byte:
		x, err := strconv.ParseInt(string(v), 10, 64)
		if err == nil {
			return x
		}
	case int64er:
		return v.Int64()
	}
	return 0
}

// Bool returns the value as a bool
func (v Value) Bool() bool {
	if v.ptr == boolType {
		return *(*bool)(unsafe.Pointer(&v.ext))
	}
	return v.toBool()
}

func (v Value) toBool() bool {
	switch {
	case v.ptr == nil:
		return false
	case v.ptr == boolType:
		return v.ext != 0
	case v.ptr == int64Type:
		return v.ext != 0
	case v.ptr == uint64Type:
		return v.ext != 0
	case v.ptr == float64Type:
		x := math.Float64frombits(v.ext)
		return x > 0 || x < 0
	case v.ptr == custBitsType:
		return v.ext != 0
	}
	switch v := v.assertNonPrimAny().(type) {
	case string:
		x, err := strconv.ParseBool(v)
		if err == nil {
			return x
		}
	case []byte:
		x, err := strconv.ParseBool(string(v))
		if err == nil {
			return x
		}
	case booler:
		return v.Bool()
	}
	return false
}

// IsString returns true if the boxed value is a string.
func (v Value) IsString() bool {
	if v.isPrim() {
		return false
	}
	switch v.ext & 0xFF {
	case ptrString:
		return true
	case ptrBytes:
		return false
	}
	_, ok := v.assertNonPrimAny().(string)
	return ok
}

// IsBytes returns true if the boxed value is a []byte.
func (v Value) IsBytes() bool {
	if v.isPrim() {
		return false
	}
	switch v.ext & 0xFF {
	case ptrBytes:
		return true
	case ptrString:
		return false
	}
	_, ok := v.assertNonPrimAny().([]byte)
	return ok
}

// IsNil returns true if the boxed value is nil.
func (v Value) IsNil() bool { return v.ptr == nil }

// IsCustomBits returns true if the boxed value was created using
// box.CustomBits.
func (v Value) IsCustomBits() bool { return v.ptr == custBitsType }

// IsInt returns true if the boxed value is an int-like primitive:
// int, int8, int16, int32, int64, byte
func (v Value) IsInt() bool { return v.ptr == int64Type }

// IsUint returns true if the boxed value is an uint-like primitive:
// uint, uint8, uint16, uint32, uint64
func (v Value) IsUint() bool { return v.ptr == uint64Type }

// IsFloat returns true if the boxed value is an float-like primitive:
// float32, float64
func (v Value) IsFloat() bool { return v.ptr == float64Type }

// IsNumber returns true if the boxed value is an numeric-like primitive:
// int, int8, int16, int32, int64, byte,
// uint, uint8, uint16, uint32, uint64,
// float32, float64
func (v Value) IsNumber() bool {
	return v.IsInt() || v.IsUint() || v.IsFloat()
}

// IsBool returns true if the boxed value is a bool primitive.
func (v Value) IsBool() bool { return v.ptr == boolType }

// Byte boxes an byte
func Byte(x byte) Value { return Int64(int64(x)) }

// Int8 boxes an int8
func Int8(x int8) Value { return Int64(int64(x)) }

// Int16 boxes an int16
func Int16(x int16) Value { return Int64(int64(x)) }

// Int32 boxes an int32
func Int32(x int32) Value { return Int64(int64(x)) }

// Int boxes an int
func Int(x int) Value { return Int64(int64(x)) }

// Uint8 boxes a uint8
func Uint8(x uint8) Value { return Uint64(uint64(x)) }

// Uint16 boxes a uint16
func Uint16(x uint16) Value { return Uint64(uint64(x)) }

// Uint32 boxes a uint32
func Uint32(x uint32) Value { return Uint64(uint64(x)) }

// Uint boxes a uint
func Uint(x uint) Value { return Uint64(uint64(x)) }

// Float32 boxes a float32
func Float32(x float32) Value { return Float64(float64(x)) }

// Byte returns the value as a byte
func (v Value) Byte() byte { return byte(v.Int64()) }

// Int8 returns the value as an int8
func (v Value) Int8() int8 { return int8(v.Int64()) }

// Int16 returns the value as an int16
func (v Value) Int16() int16 { return int16(v.Int64()) }

// Int32 returns the value as an int32
func (v Value) Int32() int32 { return int32(v.Int64()) }

// Int returns the value as an int
func (v Value) Int() int { return int(v.Int64()) }

// Uint8 returns the value as a uint8
func (v Value) Uint8() uint8 { return uint8(v.Uint64()) }

// Uint16 returns the value as a uint16
func (v Value) Uint16() uint16 { return uint16(v.Uint64()) }

// Uint32 returns the value as a uint32
func (v Value) Uint32() uint32 { return uint32(v.Uint64()) }

// Uint returns the value as a uint
func (v Value) Uint() uint { return uint(v.Uint64()) }

// Float32 returns the value as a float32
func (v Value) Float32() float32 { return float32(v.Float64()) }

// Tag returns the tag from a value created by box.StringWithTag
func (v Value) Tag() uint16 {
	if v.isPrim() {
		return 0
	}
	switch v.ext & 0xFF {
	case ptrString:
		return uint16(v.ext >> 8)
	case ptrBytes:
		return 0
	default:
		if s, ok := v.assertNonPrimAny().(*taggedString); ok {
			return s.tag
		}
		return 0
	}
}
