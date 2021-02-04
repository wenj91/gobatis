package gobatis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDbConfig(t *testing.T) {
	ymlStr := `
db:
  - datasource: ds1
    driverName: mysql
    dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
    maxLifeTime: 120
    maxOpenConns: 10
    maxIdleConns: 5
  - datasource: ds2
    driverName: mysql
    dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
    maxLifeTime: 120
    maxOpenConns: 10
    maxIdleConns: 5
showSQL: true
mappers:
  - userMapper.xml
  - orderMapper.xml
`
	dbconf := buildDbConfig(ymlStr)

	dbc := dbconf.getDataSourceByName("ds1")
	assert.True(t, dbc != nil, "test fail: No datasource1")
	assert.True(t, dbconf.ShowSQL, "test fail: showSql == false")
	assert.Equal(t, dbc.DriverName, "mysql", "test fail, actual:"+dbc.DriverName)
	assert.Equal(t, dbc.DataSourceName, "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", "test fail, actual:"+dbc.DataSourceName)
	assert.Equal(t, dbc.MaxLifeTime, 120, "test fail, actual:"+fmt.Sprintf("%d", dbc.MaxLifeTime))
	assert.Equal(t, dbc.MaxOpenConns, 10, "test fail, actual:"+fmt.Sprintf("%d", dbc.MaxOpenConns))
	assert.Equal(t, dbc.MaxIdleConns, 5, "test fail, actual:"+fmt.Sprintf("%d", dbc.MaxIdleConns))
	assert.True(t, len(dbconf.Mappers) == 2, "len(dbconf.Mappers) != 2")
	assert.Equal(t, dbconf.Mappers[0], "userMapper.xml", "test fail, actual:"+dbconf.Mappers[0])
	assert.Equal(t, dbconf.Mappers[1], "orderMapper.xml", "test fail, actual:"+dbconf.Mappers[1])
}

func TestDbConfigCodeInit(t *testing.T) {
	ds1 := NewDataSourceBuilder().
		DataSource("ds1").
		DriverName("mysql").
		DataSourceName("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8").
		MaxLifeTime(120).
		MaxOpenConns(10).
		MaxIdleConns(5).
		Build()

	dbconf := NewDBConfigBuilder().
		DS([]*DataSource{ds1}).
		ShowSQL(true).
		Mappers([]string{"userMapper.xml", "orderMapper.xml"}).
		Build()

	dbc := dbconf.getDataSourceByName("ds1")
	assert.True(t, dbc != nil, "test fail: No datasource1")
	assert.True(t, dbconf.ShowSQL, "test fail: showSql == false")
	assert.Equal(t, dbc.DriverName, "mysql", "test fail, actual:"+dbc.DriverName)
	assert.Equal(t, dbc.DataSourceName, "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", "test fail, actual:"+dbc.DataSourceName)
	assert.Equal(t, dbc.MaxLifeTime, 120, "test fail, actual:"+fmt.Sprintf("%d", dbc.MaxLifeTime))
	assert.Equal(t, dbc.MaxOpenConns, 10, "test fail, actual:"+fmt.Sprintf("%d", dbc.MaxOpenConns))
	assert.Equal(t, dbc.MaxIdleConns, 5, "test fail, actual:"+fmt.Sprintf("%d", dbc.MaxIdleConns))
	assert.True(t, len(dbconf.Mappers) == 2, "len(dbconf.Mappers) != 2")
	assert.Equal(t, dbconf.Mappers[0], "userMapper.xml", "test fail, actual:"+dbconf.Mappers[0])
	assert.Equal(t, dbconf.Mappers[1], "orderMapper.xml", "test fail, actual:"+dbconf.Mappers[1])
}
