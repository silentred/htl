package ast

import (
	"math"
)

// Kind constants.
const (
	KindNull          byte = 0
	KindInt64         byte = 1
	KindUint64        byte = 2
	KindFloat32       byte = 3
	KindFloat64       byte = 4
	KindString        byte = 5
	KindBytes         byte = 6
	KindMysqlBit      byte = 7
	KindMysqlDecimal  byte = 8
	KindMysqlDuration byte = 9
	KindMysqlEnum     byte = 10
	KindMysqlHex      byte = 11
	KindMysqlSet      byte = 12
	KindMysqlTime     byte = 13
	KindRow           byte = 14
	KindInterface     byte = 15
	KindMinNotNull    byte = 16
	KindMaxValue      byte = 17
	KindRaw           byte = 18
)

// Datum is a data box holds different kind of data.
// It has better performance and is easier to use than `interface{}`.
type Datum struct {
	k         byte        // datum kind.
	collation uint8       // collation can hold uint8 values.
	decimal   uint16      // decimal can hold uint16 values.
	length    uint32      // length can hold uint32 values.
	i         int64       // i can hold int64 uint64 float64 values.
	b         []byte      // b can hold string or []byte values.
	x         interface{} // x hold all other types.
}

// Kind gets the kind of the datum.
func (d *Datum) Kind() byte {
	return d.k
}

// Collation gets the collation of the datum.
func (d *Datum) Collation() byte {
	return d.collation
}

// SetCollation sets the collation of the datum.
func (d *Datum) SetCollation(collation byte) {
	d.collation = collation
}

// Frac gets the frac of the datum.
func (d *Datum) Frac() int {
	return int(d.decimal)
}

// SetFrac sets the frac of the datum.
func (d *Datum) SetFrac(frac int) {
	d.decimal = uint16(frac)
}

// Length gets the length of the datum.
func (d *Datum) Length() int {
	return int(d.length)
}

// SetLength sets the length of the datum
func (d *Datum) SetLength(l int) {
	d.length = uint32(l)
}

// IsNull checks if datum is null.
func (d *Datum) IsNull() bool {
	return d.k == KindNull
}

// GetInt64 gets int64 value.
func (d *Datum) GetInt64() int64 {
	return d.i
}

// SetInt64 sets int64 value.
func (d *Datum) SetInt64(i int64) {
	d.k = KindInt64
	d.i = i
}

// GetUint64 gets uint64 value.
func (d *Datum) GetUint64() uint64 {
	return uint64(d.i)
}

// SetUint64 sets uint64 value.
func (d *Datum) SetUint64(i uint64) {
	d.k = KindUint64
	d.i = int64(i)
}

// GetFloat64 gets float64 value.
func (d *Datum) GetFloat64() float64 {
	return math.Float64frombits(uint64(d.i))
}

// SetFloat64 sets float64 value.
func (d *Datum) SetFloat64(f float64) {
	d.k = KindFloat64
	d.i = int64(math.Float64bits(f))
}

// GetFloat32 gets float32 value.
func (d *Datum) GetFloat32() float32 {
	return float32(math.Float64frombits(uint64(d.i)))
}

// SetFloat32 sets float32 value.
func (d *Datum) SetFloat32(f float32) {
	d.k = KindFloat32
	d.i = int64(math.Float64bits(float64(f)))
}

// GetString gets string value.
func (d *Datum) GetString() string {
	return String(d.b)
}

// SetString sets string value.
func (d *Datum) SetString(s string) {
	d.k = KindString
	sink(s)
	d.b = Slice(s)
}

// sink prevents s from being allocated on the stack.
var sink = func(s string) {
}

// GetBytes gets bytes value.
func (d *Datum) GetBytes() []byte {
	return d.b
}

// SetBytes sets bytes value to datum.
func (d *Datum) SetBytes(b []byte) {
	d.k = KindBytes
	d.b = b
}

// SetBytesAsString sets bytes value to datum as string type.
func (d *Datum) SetBytesAsString(b []byte) {
	d.k = KindString
	d.b = b
}

// GetInterface gets interface value.
func (d *Datum) GetInterface() interface{} {
	return d.x
}

// SetInterface sets interface to datum.
func (d *Datum) SetInterface(x interface{}) {
	d.k = KindInterface
	d.x = x
}

// GetRow gets row value.
func (d *Datum) GetRow() []Datum {
	return d.x.([]Datum)
}

// SetRow sets row value.
func (d *Datum) SetRow(ds []Datum) {
	d.k = KindRow
	d.x = ds
}

// SetNull sets datum to nil.
func (d *Datum) SetNull() {
	d.k = KindNull
	d.x = nil
}

// SetRaw sets raw value.
func (d *Datum) SetRaw(b []byte) {
	d.k = KindRaw
	d.b = b
}

// GetRaw gets raw value.
func (d *Datum) GetRaw() []byte {
	return d.b
}

// GetValue gets the value of the datum of any kind.
func (d *Datum) GetValue() interface{} {
	switch d.k {
	case KindInt64:
		return d.GetInt64()
	case KindUint64:
		return d.GetUint64()
	case KindFloat32:
		return d.GetFloat32()
	case KindFloat64:
		return d.GetFloat64()
	case KindString:
		return d.GetString()
	case KindBytes:
		return d.GetBytes()
	default:
		return d.GetInterface()
	}
}

// SetValue sets any kind of value.
func (d *Datum) SetValue(val interface{}) {
	switch x := val.(type) {
	case nil:
		d.SetNull()
	case bool:
		if x {
			d.SetInt64(1)
		} else {
			d.SetInt64(0)
		}
	case int:
		d.SetInt64(int64(x))
	case int64:
		d.SetInt64(x)
	case uint64:
		d.SetUint64(x)
	case float32:
		d.SetFloat32(x)
	case float64:
		d.SetFloat64(x)
	case string:
		d.SetString(x)
	case []byte:
		d.SetBytes(x)
	case []Datum:
		d.SetRow(x)
	case []interface{}:
		ds := MakeDatums(x...)
		d.SetRow(ds)
	default:
		d.SetInterface(x)
	}
}

// NewDatum creates a new Datum from an interface{}.
func NewDatum(in interface{}) (d Datum) {
	switch x := in.(type) {
	case []interface{}:
		d.SetValue(MakeDatums(x...))
	default:
		d.SetValue(in)
	}
	return d
}

// NewIntDatum creates a new Datum from an int64 value.
func NewIntDatum(i int64) (d Datum) {
	d.SetInt64(i)
	return d
}

// NewUintDatum creates a new Datum from an uint64 value.
func NewUintDatum(i uint64) (d Datum) {
	d.SetUint64(i)
	return d
}

// NewBytesDatum creates a new Datum from a byte slice.
func NewBytesDatum(b []byte) (d Datum) {
	d.SetBytes(b)
	return d
}

// NewStringDatum creates a new Datum from a string.
func NewStringDatum(s string) (d Datum) {
	d.SetString(s)
	return d
}

// NewFloat64Datum creates a new Datum from a float64 value.
func NewFloat64Datum(f float64) (d Datum) {
	d.SetFloat64(f)
	return d
}

// NewFloat32Datum creates a new Datum from a float32 value.
func NewFloat32Datum(f float32) (d Datum) {
	d.SetFloat32(f)
	return d
}

// NewRawDatum create a new Datum that is not decoded.
func NewRawDatum(b []byte) (d Datum) {
	d.k = KindRaw
	d.b = b
	return d
}

// MakeDatums creates datum slice from interfaces.
func MakeDatums(args ...interface{}) []Datum {
	datums := make([]Datum, len(args))
	for i, v := range args {
		datums[i] = NewDatum(v)
	}
	return datums
}
