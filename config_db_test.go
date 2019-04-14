package gobatis

import (
	"fmt"
	"testing"
)

func TestDbConfig(t *testing.T)  {
	ymlStr := `
db:
  -	datasource: ds1
    driverName: mysql
    dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
    maxLifeTime: 120
    maxOpenConns: 10
    maxIdleConns: 5
  -	datasource: ds2
    driverName: mysql
    dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
    maxLifeTime: 120
    maxOpenConns: 10
    maxIdleConns: 5
showSql: true
mappers:
  - userMapper.xml
  - orderMapper.xml
`
	dbconf := buildDbConfig(ymlStr)

	dbc := dbconf.getDataSourceByName("ds1")
	assertTrue(dbc != nil, "test fail: No datasource1")
	assertTrue(dbconf.ShowSql, "test fail: showSql == false")
	assertEqual(dbc.DriverName, "mysql", "test fail, actual:" + dbc.DriverName)
	assertEqual(dbc.DataSourceName, "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", "test fail, actual:" + dbc.DataSourceName)
	assertEqual(dbc.MaxLifeTime, 120, "test fail, actual:" + fmt.Sprintf("%d", dbc.MaxLifeTime))
	assertEqual(dbc.MaxOpenConns, 10, "test fail, actual:" + fmt.Sprintf("%d", dbc.MaxOpenConns))
	assertEqual(dbc.MaxIdleConns, 5, "test fail, actual:" + fmt.Sprintf("%d", dbc.MaxIdleConns))
	assertTrue(len(dbconf.Mappers) == 2, "len(dbconf.Mappers) != 2")
	assertEqual(dbconf.Mappers[0], "userMapper.xml", "test fail, actual:" + dbconf.Mappers[0])
	assertEqual(dbconf.Mappers[1], "orderMapper.xml", "test fail, actual:" + dbconf.Mappers[1])
}
