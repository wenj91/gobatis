package uti

import (
	"github.com/wenj91/gobatis/na"
	"reflect"
	"time"
)

// PI64 to NullInt64
func PI64(i int64) *int64 {
	return &i
}

// PS to NullString
func PS(s string) *string {
	return &s
}

// PF64 to NullFloat64
func PF64(f float64) *float64 {
	return &f
}

// PT to NullTime
func PT(t time.Time) *time.Time {
	return &t
}

// NB to NullBool
func PB(b bool) *bool {
	return &b
}

// NI64 to NullInt64
func NI64(i int64) na.NullInt64 {
	return na.NullInt64{Int64: i, Valid: true}
}

// NS to NullString
func NS(s string) na.NullString {
	return na.NullString{String: s, Valid: true}
}

// NF64 to NullFloat64
func NF64(f float64) na.NullFloat64 {
	return na.NullFloat64{Float64: f, Valid: true}
}

// NT to NullTime
func NT(t time.Time) na.NullTime {
	return na.NullTime{Time: t, Valid: true}
}

// NB to NullBool
func NB(b bool) na.NullBool {
	return na.NullBool{Bool: b, Valid: true}
}

func IsNil(i interface{}) (bool, isBaseType bool) {
	objVal := reflect.ValueOf(i)
	k := objVal.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if objVal.IsNil() {
			return true, false
		} else {
			return false, false
		}
	}

	return false, true
}

func IsBaseType(i interface{}) bool {
	objVal := reflect.ValueOf(i)
	k := objVal.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return false
	}

	return true
}
