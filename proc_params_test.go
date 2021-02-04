package gobatis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestStruct2 struct {
	T  NullTime
	S  NullString
	Id int64
}

func TestParams(t *testing.T) {
	paramB := ""
	res := paramProcess(paramB)
	assert.NotNil(t, res["0"], "test fail: res[0] == nil")
	assert.Equal(t, res["0"], "", "test fail: res[0] != ''")

	paramM := map[string]interface{}{
		"id":   nil,
		"name": "wenj91",
	}
	res = paramProcess(paramM)
	assert.Nil(t, res["id"], "test fail: res['id'] != nil")
	assert.NotNil(t, res["name"], "test fail: res['name'] == nil")
	assert.Equal(t, res["name"], "wenj91", "test fail: res['name'] != 'wenj91'")

	paramNil := NullString{"str", true}
	res = paramProcess(paramNil)
	assert.NotNil(t, res["0"], "test fail: res['0'] == nil")
	assert.Equal(t, res["0"], "str", "test fail: res['0'] != 'str'")

	tt, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	paramS := &TestStruct2{
		T: NullTime{tt, true},
	}
	res = paramProcess(paramS)
	assert.NotNil(t, res["T"], "test fail: res['T'] == nil")
	assert.Equal(t, res["T"], "2006-01-02 15:04:05", "test fail: res['T'] != '2006-01-02 15:04:05'")
	assert.Nil(t, res["S"], "test fail: res['S'] != nil")
	assert.Equal(t, res["Id"], int64(0), "test fail: res['Id'] != 0")
}
