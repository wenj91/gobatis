package gobatis

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestStruct struct {
}

func TestVal(t *testing.T) {
	paramsInt := 1
	v := reflect.ValueOf(paramsInt)
	assert.True(t, v.Kind() == reflect.Int, "test fail: params is not int")

	paramsInt64 := int64(1)
	v = reflect.ValueOf(paramsInt64)
	assert.True(t, v.Kind() == reflect.Int64, "test fail: params is not int64")

	paramsString := ""
	v = reflect.ValueOf(paramsString)
	assert.True(t, v.Kind() == reflect.String, "test fail: params is not string")

	paramsSlice := []int{1, 2, 3}
	v = reflect.ValueOf(paramsSlice)
	assert.True(t, v.Kind() == reflect.Slice, "test fail: params is not slice")

	paramsStruct := TestStruct{}
	v = reflect.ValueOf(paramsStruct)
	assert.True(t, v.Kind() == reflect.Struct, "test fail: params is not struct")

	paramsPtr := &TestStruct{}
	v = reflect.ValueOf(paramsPtr)
	assert.True(t, v.Kind() == reflect.Ptr, "test fail: params is not ptr")
	v = v.Elem()
	assert.True(t, v.Kind() == reflect.Struct, "test fail: params is not struct")

	paramsStructs := []*TestStruct{{}, {}, {}}
	v = reflect.ValueOf(paramsStructs)
	assert.True(t, v.Kind() == reflect.Slice, "test fail: params is not slice")
	assert.True(t, v.Len() == 3, "test fail: params len != 3")
	v0 := v.Index(0)
	assert.True(t, v0.Kind() == reflect.Ptr, "test fail: ele is not ptr")
	v0 = v0.Elem()
	assert.True(t, v0.Kind() == reflect.Struct, "test fail: ele is not struct")
}
