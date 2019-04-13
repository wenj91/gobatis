package gobatis

import (
	"math/rand"
	"testing"
)

func TestBuildConfig(t *testing.T) {
	params := map[string]interface{}{
		"Name":     "",
		"Age":      -1,
		"Password": 1,
	}

	conf := loadingMapper("./examples/mapper")
	stmt := conf.getMappedStmt("userMapper.queryStructsByCond3")
	s := stmt.sqlSource.getBoundSql(params)
	t.Logf("%s", s.sqlStr)

}

func BenchmarkName(b *testing.B) {

	params := make(map[string]interface{})
	params["Name"] = "Sean"
	conf := loadingMapper("./examples/mapper")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		stmt := conf.getMappedStmt("userMapper.queryStructsByCond3")
		params["Age"] = rand.Int()
		params["Password"] = rand.Int()
		_ = stmt.sqlSource.getBoundSql(params)
	}
	b.StopTimer()

}
