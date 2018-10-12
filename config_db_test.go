package gobatis

import (
	"fmt"
	"testing"
)

func TestDbConfig(t *testing.T)  {
	ymlStr := `
db:
  driverName: mysql
  dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
  maxLifeTime: 120
  maxOpenConns: 10
  maxIdleConns: 5
mappers:
  - userMapper.xml
  - orderMapper.xml
`
	dbconf := buildDbConfig(ymlStr)

	assertEqual(dbconf.DB.DriverName, "mysql", "test fail, actual:" + dbconf.DB.DriverName)
	assertEqual(dbconf.DB.DataSourceName, "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", "test fail, actual:" + dbconf.DB.DataSourceName)
	assertEqual(dbconf.DB.MaxLifeTime, 120, "test fail, actual:" + fmt.Sprintf("%d", dbconf.DB.MaxLifeTime))
	assertEqual(dbconf.DB.MaxOpenConns, 10, "test fail, actual:" + fmt.Sprintf("%d", dbconf.DB.MaxOpenConns))
	assertEqual(dbconf.DB.MaxIdleConns, 5, "test fail, actual:" + fmt.Sprintf("%d", dbconf.DB.MaxIdleConns))
	assertTrue(len(dbconf.Mappers) == 2, "len(dbconf.Mappers) != 2")
	assertEqual(dbconf.Mappers[0], "userMapper.xml", "test fail, actual:" + dbconf.Mappers[0])
	assertEqual(dbconf.Mappers[1], "orderMapper.xml", "test fail, actual:" + dbconf.Mappers[1])
}
