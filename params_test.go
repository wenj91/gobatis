package gobatis

import (
	"reflect"
	"testing"
)

type TestStruct struct {

}

func TestVal(t *testing.T)  {
	paramsInt := 1
	v := reflect.ValueOf(paramsInt)
	assertTrue(v.Kind() == reflect.Int, "test fail: params is not int")

	paramsInt64 := int64(1)
	v = reflect.ValueOf(paramsInt64)
	assertTrue(v.Kind() == reflect.Int64, "test fail: params is not int64")

	paramsString := ""
	v = reflect.ValueOf(paramsString)
	assertTrue(v.Kind() == reflect.String, "test fail: params is not string")

	paramsSlice := []int{1, 2, 3}
	v = reflect.ValueOf(paramsSlice)
	assertTrue(v.Kind() == reflect.Slice, "test fail: params is not slice")

	paramsStruct := TestStruct{}
	v = reflect.ValueOf(paramsStruct)
	assertTrue(v.Kind() == reflect.Struct, "test fail: params is not struct")

	paramsPtr := &TestStruct{}
	v = reflect.ValueOf(paramsPtr)
	assertTrue(v.Kind() == reflect.Ptr, "test fail: params is not ptr")
	v = v.Elem()
	assertTrue(v.Kind() == reflect.Struct, "test fail: params is not struct")

	paramsStructs := []*TestStruct{{},{},{}}
	v = reflect.ValueOf(paramsStructs)
	assertTrue(v.Kind() == reflect.Slice, "test fail: params is not slice")
	assertTrue(v.Len() == 3, "test fail: params len != 3")
	v0 := v.Index(0)
	assertTrue(v0.Kind() == reflect.Ptr, "test fail: ele is not ptr")
	v0 = v0.Elem()
	assertTrue(v0.Kind() == reflect.Struct, "test fail: ele is not struct")
}
