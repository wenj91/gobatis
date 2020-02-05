# gobatis

目前代码都是基于mysql编写测试的,其他数据库暂时还未做兼容处理

## gobatis接口

```go
type GoBatis interface {
	// Select 查询数据
	Select(stmt string, param interface{}) func(res interface{}) error
	// SelectContext 查询数据with context
	SelectContext(ctx context.Context, stmt string, param interface{}) func(res interface{}) error
	// Insert 插入数据
	Insert(stmt string, param interface{}) (int64, int64, error)
	// InsertContext 插入数据with context
	InsertContext(ctx context.Context, stmt string, param interface{}) (int64, int64, error)
	// Update 更新数据
	Update(stmt string, param interface{}) (int64, error)
	// UpdateContext 更新数据with context
	UpdateContext(ctx context.Context, stmt string, param interface{}) (int64, error)
	// Delete 刪除数据
	Delete(stmt string, param interface{}) (int64, error)
	// DeleteContext 刪除数据with context
	DeleteContext(ctx context.Context, stmt string, param interface{}) (int64, error)
}
```

## db数据源配置
- 支持多数据源配置
- db子级配置为一个map，map的key即为数据源名称标识  
- map的value为数据源具体配置，具体配置项如下表

| 配置 | 是否必填配置 | 默认值 | 说明 |
|:---|:----:|:----:|----|
| driverName | 是 | | 数据源驱动名，必填配置项
| dataSourceName | 是 | | 数据源名称，必填配置项，例如: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
| maxLifeTime | 否 | 120(单位: s)| 连接最大存活时间，默认值为: 120 单位为: s
| maxOpenConns | 否 | 10 | 最大打开连接数，默认值为: 10
| maxIdleConns | 否 | 5 | 最大挂起连接数，默认值为: 5

### 示例
* db配置示例(配置较之前的有所调整)  
以下为多数据源配置示例: db.yml
```yaml
# 数据库配置
db:
  # 数据源名称1
  - datasource: ds1
    # 驱动名
    driverName: mysql
    # 数据源
    dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
    # 连接最大存活时间（单位: s）
    maxLifeTime: 120
    # 最大open连接数
    maxOpenConns: 10
    # 最大挂起连接数
    maxIdleConns: 5
  # 数据源名称2
  - datasource: ds2
    # 驱动名
    driverName: mysql
    # 数据源
    dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
    # 连接最大存活时间（单位: s）
    maxLifeTime: 120
    # 最大open连接数
    maxOpenConns: 10
    # 最大挂起连接数
    maxIdleConns: 5
# 是否显示SQL语句
showSql: true
# 数据表映射文件路径配置
mappers:
  # 映射文件路径， 可以为绝对路径，如: /usr/local/mapper/userMapper.xml
  - mapper/userMapper.xml
```

* mapper配置  
1. mapper可以配置namespace属性  
1. mapper可以包含: select, insert, update, delete标签  
1. mapper子标签id属性则为标签唯一标识, 必须配置属性
1. 其中select标签必须包含resultType属性，resultType可以是: map, maps, array, arrays, struct, structs, value
  
* 标签说明  
select: 用于查询操作   
insert: 用于插入sql操作  
update: 用于更新sql操作  
delete: 用于删除sql操作

* resultType说明  
map: 则数据库查询结果为map  
maps: 则数据库查询结果为map数组  
array: 则数据库查询结果为值数组  
arrays: 则数据库查询结果为多个值数组  
struct: 则数据库查询结果为单个结构体  
structs: 则数据库查询结果为结构体数组  
value: 则数据库查询结果为单个数值  
 
以下是mapper配置示例: mapper/userMapper.xml
```xml
<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE mapper PUBLIC "gobatis"
        "https://raw.githubusercontent.com/wenj91/gobatis/master/gobatis.dtd">
<mapper namespace="userMapper">
    <select id="findMapById" resultType="map">
        SELECT id, name FROM user where id=#{id} order by id
    </select>
    <select id="findMapByValue" resultType="map">
            SELECT id, name FROM user where id=#{0} order by id
    </select>
    <select id="findStructByStruct" resultType="struct">
        SELECT id, name, crtTm FROM user where id=#{Id} order by id
    </select>
    <select id="queryStructs" resultType="structs">
        SELECT id, name, crtTm FROM user order by id
    </select>
    <select id="queryStructsByOrder" resultType="structs">
        SELECT id, name, crtTm FROM user order by ${id} desc
    </select>
    <insert id="insertStruct">
        insert into user (name, email, crtTm)
        values (#{Name}, #{Email}, #{CrtTm})
    </insert>
    <delete id="deleteById">
        delete from user where id=#{id}
    </delete>
    <select id="queryStructsByCond" resultType="structs">
         SELECT id, name, crtTm, pwd, email FROM user
         <where>
             <if test="Name != nil and Name != ''">and name = #{Name}</if>
         </where>
         order by id
    </select>
     <select id="queryStructsByCond2" resultType="structs">
         SELECT id, name, crtTm, pwd, email FROM user
         <trim prefixOverrides="and" prefix="where" suffixOverrides="," suffix="and 1=1">
              <if test="Name != nil and Name != ''">and name = #{Name}</if>
         </trim>
         order by id
    </select>
    <update id="updateByCond">
        update user
        <set>
            <if test="Name != nil and Name2 != ''">name = #{Name},</if>
        </set>
        where id = #{Id}
    </update>
</mapper>
```

## 使用方法

###  使用配置文件配置
example1.go
```go
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 引入驱动
	"github.com/wenj91/gobatis"        // 引入gobatis
)

// 实体结构示例， tag：field为数据库对应字段名称
type User struct {
	Id    gobatis.NullInt64  `field:"id"`
	Name  gobatis.NullString `field:"name"`
	Email gobatis.NullString `field:"email"`
	CrtTm gobatis.NullTime   `field:"crtTm"`
}

func main() {
	// 初始化db，参数为db.yml路径，如：db.yml	
	gobatis.Init(gobatis.NewFileOption("db.yml"))

	// 获取数据源，参数为数据源名称，如：datasource1
	gb := gobatis.Get("ds1")

	//传入id查询Map
	mapRes := make(map[string]interface{})
	// stmt标识为：namespace + '.' + id, 如：userMapper.findMapById
	// 查询参数可以是map，也可以是数组，也可以是实体结构
	err := gb.Select("userMapper.findMapById", map[string]interface{}{"id": 1})(mapRes)
	fmt.Println("userMapper.findMapById-->", mapRes, err)

	// 根据传入实体查询对象
	param := User{Id: gobatis.NullInt64{Int64: 1, Valid: true}}
	var structRes *User
	err = gb.Select("userMapper.findStructByStruct", param)(&structRes)
	fmt.Println("userMapper.findStructByStruct-->", structRes, err)

	// 查询实体列表
	structsRes := make([]*User, 0)
	err = gb.Select("userMapper.queryStructs", map[string]interface{}{})(&structsRes)
	fmt.Println("userMapper.queryStructs-->", structsRes, err)

	param = User{
		Id:   gobatis.NullInt64{Int64: 1, Valid: true},
		Name: gobatis.NullString{String: "wenj1993", Valid: true},
	}

	// set tag
	affected, err := gb.Update("userMapper.updateByCond", param)
	fmt.Println("updateByCond:", affected, err)

	param = User{Name: gobatis.NullString{String: "wenj1993", Valid: true}}
	// where tag
	res := make([]*User, 0)
	err = gb.Select("userMapper.queryStructsByCond", param)(&res)
	fmt.Println("queryStructsByCond", res, err)

	// trim tag
	res = make([]*User, 0)
	err = gb.Select("userMapper.queryStructsByCond2", param)(&res)
	fmt.Println("queryStructsByCond2", res, err)
	
	// ${id}
	res = make([]*User, 0)
	err = gb.Select("userMapper.queryStructsByOrder", map[string]interface{}{
		"id":"id",
	})(&res)
	fmt.Println("queryStructsByCond", res, err)

	// 开启事务示例
	tx, _ := gb.Begin()
	defer tx.Rollback()
	tx.Select("userMapper.findMapById", map[string]interface{}{"id": 1,})(mapRes)
	fmt.Println("tx userMapper.findMapById-->", mapRes, err)
	tx.Commit()
}
```

### 代码配置方式

example2.go

```go
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 引入驱动
	"github.com/wenj91/gobatis"        // 引入gobatis
)

// 实体结构示例， tag：field为数据库对应字段名称
type User struct {
	Id    gobatis.NullInt64  `field:"id"`
	Name  gobatis.NullString `field:"name"`
	Email gobatis.NullString `field:"email"`
	CrtTm gobatis.NullTime   `field:"crtTm"`
}

func main() {
	// 初始化db
	ds1 := gobatis.NewDataSourceBuilder().
		DataSource("ds1").
		DriverName("mysql").
		DataSourceName("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8").
		MaxLifeTime(120).
		MaxOpenConns(10).
		MaxIdleConns(5).
		Build()

	option := gobatis.NewDSOption().
		DS([]*gobatis.DataSource{ds1}).
		Mappers([]string{"examples/mapper/userMapper.xml"}).
		ShowSQL(true)

	gobatis.Init(option)

	// 获取数据源，参数为数据源名称，如：ds1
	gb := gobatis.Get("ds1")

	//传入id查询Map
	mapRes := make(map[string]interface{})
	// stmt标识为：namespace + '.' + id, 如：userMapper.findMapById
	// 查询参数可以是map，也可以是数组，也可以是实体结构
	err := gb.Select("userMapper.findMapById", map[string]interface{}{"id": 1})(mapRes)
	fmt.Println("userMapper.findMapById-->", mapRes, err)
}
```

example3.go

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 引入驱动
	"github.com/wenj91/gobatis"        // 引入gobatis
)

// 实体结构示例， tag：field为数据库对应字段名称
type User struct {
	Id    gobatis.NullInt64  `field:"id"`
	Name  gobatis.NullString `field:"name"`
	Email gobatis.NullString `field:"email"`
	CrtTm gobatis.NullTime   `field:"crtTm"`
}

func main() {
	// 初始化db
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	dbs := make(map[string]*gobatis.GoBatisDB)
	dbs["ds1"] = gobatis.NewGoBatisDB(gobatis.DBTypeMySQL, db)

	option := gobatis.NewDBOption().
		DB(dbs).
		ShowSQL(true).
		Mappers([]string{"examples/mapper/userMapper.xml"})

	gobatis.Init(option)

	// 获取数据源，参数为数据源名称，如：ds1
	gb := gobatis.Get("ds1")

	//传入id查询Map
	mapRes := make(map[string]interface{})
	// stmt标识为：namespace + '.' + id, 如：userMapper.findMapById
	// 查询参数可以是map，也可以是数组，也可以是实体结构
	err := gb.Select("userMapper.findMapById", map[string]interface{}{"id": 1})(mapRes)
	fmt.Println("userMapper.findMapById-->", mapRes, err)
}
```