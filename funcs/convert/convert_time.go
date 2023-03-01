package convert

import (
	"time"

	"github.com/perpower/goframe/funcs/judge"
	"github.com/perpower/goframe/funcs/ptime"
)

// Time converts `any` to time.Time.
func Time(any interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v
		}
	}
	if t := GTime(any, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func Duration(any interface{}) time.Duration {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v
	}
	s := String(any)
	if !judge.IsNumeric(s) {
		d, _ := ptime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(any))
}

// GTime converts `any` to *gtime.Time.
// The parameter `format` can be used to specify the format of `any`.
// If no `format` given, it converts `any` using gtime.NewFromTimeStamp if `any` is numeric,
// or using gtime.StrToTime if `any` is string.
func GTime(any interface{}, format ...string) *ptime.Time {
	if any == nil {
		return nil
	}
	if v, ok := any.(iGTime); ok {
		return v.GTime(format...)
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*ptime.Time); ok {
			return v
		}
		if t, ok := any.(time.Time); ok {
			return ptime.New(t)
		}
		if t, ok := any.(*time.Time); ok {
			return ptime.New(t)
		}
	}
	s := String(any)
	if len(s) == 0 {
		return ptime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		t, _ := ptime.StrToTimeFormat(s, format[0])
		return t
	}
	if judge.IsNumeric(s) {
		return ptime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := ptime.StrToTime(s)
		return t
	}
}
