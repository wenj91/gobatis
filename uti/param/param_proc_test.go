package param

import (
	"github.com/stretchr/testify/assert"
	"github.com/wenj91/gobatis/na"
	"testing"
	"time"
)

type TestStruct2 struct {
	T  na.NullTime
	S  na.NullString
	Id int64
}

func TestParams(t *testing.T) {
	paramB := ""
	res := Process(paramB)
	assert.NotNil(t, res["0"], "test fail: res[0] == nil")
	assert.Equal(t, res["0"], "", "test fail: res[0] != ''")

	paramM := map[string]interface{}{
		"id":   nil,
		"name": "wenj91",
	}
	res = Process(paramM)
	assert.Nil(t, res["id"], "test fail: res['id'] != nil")
	assert.NotNil(t, res["name"], "test fail: res['name'] == nil")
	assert.Equal(t, res["name"], "wenj91", "test fail: res['name'] != 'wenj91'")

	paramNil := na.NullString{"str", true}
	res = Process(paramNil)
	assert.NotNil(t, res["0"], "test fail: res['0'] == nil")
	assert.Equal(t, res["0"], "str", "test fail: res['0'] != 'str'")

	tt, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	paramS := &TestStruct2{
		T: na.NullTime{tt, true},
	}
	res = Process(paramS)
	assert.NotNil(t, res["T"], "test fail: res['T'] == nil")
	assert.Equal(t, res["T"], "2006-01-02 15:04:05", "test fail: res['T'] != '2006-01-02 15:04:05'")
	assert.Nil(t, res["S"], "test fail: res['S'] != nil")
	assert.Equal(t, res["Id"], int64(0), "test fail: res['Id'] != 0")
}
