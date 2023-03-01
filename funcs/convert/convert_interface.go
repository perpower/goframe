package convert

import "github.com/perpower/goframe/funcs/ptime"

// iString is used for type assert api for String().
type iString interface {
	String() string
}

// iBool is used for type assert api for Bool().
type iBool interface {
	Bool() bool
}

// iInt64 is used for type assert api for Int64().
type iInt64 interface {
	Int64() int64
}

// iUint64 is used for type assert api for Uint64().
type iUint64 interface {
	Uint64() uint64
}

// iFloat32 is used for type assert api for Float32().
type iFloat32 interface {
	Float32() float32
}

// iFloat64 is used for type assert api for Float64().
type iFloat64 interface {
	Float64() float64
}

// iError is used for type assert api for Error().
type iError interface {
	Error() string
}

// iBytes is used for type assert api for Bytes().
type iBytes interface {
	Bytes() []byte
}

// iInterface is used for type assert api for Interface().
type iInterface interface {
	Interface() interface{}
}

// iInterfaces is used for type assert api for Interfaces().
type iInterfaces interface {
	Interfaces() []interface{}
}

// iFloats is used for type assert api for Floats().
type iFloats interface {
	Floats() []float64
}

// iInts is used for type assert api for Ints().
type iInts interface {
	Ints() []int
}

// iStrings is used for type assert api for Strings().
type iStrings interface {
	Strings() []string
}

// iUints is used for type assert api for Uints().
type iUints interface {
	Uints() []uint
}

// iGTime is the interface for gtime.Time converting.
type iGTime interface {
	GTime(format ...string) *ptime.Time
}
