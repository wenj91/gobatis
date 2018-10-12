package gobatis

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

func TestDbConfig(t *testing.T)  {
	ymlStr := `
db:
  driverName: mysql
  dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
  maxLifeTime: 10
  maxIdleConns: 1
mappers:
  - userMapper.xml
  - orderMapper.xml
`
	dbconf := &dbConfig{}
	err := yaml.Unmarshal([]byte(ymlStr), &dbconf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	assertEqual(dbconf.DB.DriverName, "mysql", "test fail, actual:" + dbconf.DB.DriverName)
	assertEqual(dbconf.DB.DataSourceName, "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", "test fail, actual:" + dbconf.DB.DataSourceName)
	assertEqual(dbconf.DB.MaxLifeTime, 10, "test fail, actual:" + fmt.Sprintf("%d", dbconf.DB.MaxLifeTime))
	assertEqual(dbconf.DB.MaxIdleConns, 1, "test fail, actual:" + fmt.Sprintf("%d", dbconf.DB.MaxIdleConns))
	assertTrue(len(dbconf.Mappers) == 2, "len(dbconf.Mappers) != 2")
	assertEqual(dbconf.Mappers[0], "userMapper.xml", "test fail, actual:" + dbconf.Mappers[0])
	assertEqual(dbconf.Mappers[1], "orderMapper.xml", "test fail, actual:" + dbconf.Mappers[1])
}
