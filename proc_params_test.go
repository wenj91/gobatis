package gobatis

import (
	"testing"
	"time"
)

type TestStruct2 struct {
	T NullTime
	S NullString
	Id int64
}

func TestParams(t *testing.T)  {
	paramB := ""
	res := paramProcess(paramB)
	assertNotNil(res["0"], "test fail: res[0] == nil")
	assertEqual(res["0"], "", "test fail: res[0] != ''")


	paramM := map[string]interface{}{
		"id": nil,
		"name":"wenj91",
	}
	res = paramProcess(paramM)
	assertNil(res["id"], "test fail: res['id'] != nil")
	assertNotNil(res["name"], "test fail: res['name'] == nil")
	assertEqual(res["name"], "wenj91", "test fail: res['name'] != 'wenj91'")

	paramNil := NullString{"str", true}
	res = paramProcess(paramNil)
	assertNotNil(res["0"], "test fail: res['0'] == nil")
	assertEqual(res["0"], "str", "test fail: res['0'] != 'str'")

	tt, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	paramS := &TestStruct2{
		T: NullTime{tt, true},
	}
	res = paramProcess(paramS)
	assertNotNil(res["T"], "test fail: res['T'] == nil")
	assertEqual(res["T"], "2006-01-02 15:04:05", "test fail: res['T'] != '2006-01-02 15:04:05'")
	assertNil(res["S"], "test fail: res['S'] != nil")
	assertEqual(res["Id"], int64(0), "test fail: res['Id'] != 0")
}
