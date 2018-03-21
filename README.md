# gobatis

## TODO
0 todo: 单元测试编写  
1 todo: 批量插入修改<for>标签实现-->done  
2 todo: 动态sql生成<if>标签实现  
3 todo: ${xxx}解析实现  
4 todo: 结果集映射<resultMap>标签实现  
5 todo: 公共查询字段<sql>标签实现  
6 todo: 一级缓存实现  
7 todo: 二级缓存实现   
8 todo: 加入连接池   
9 todo: 完善文档  

#### mapper配置文件
```xml
<?xml version="1.0" encoding="utf-8"?>
<mapper namespace="Mapper">
    <select id="findMapById" resultType="Map">
        SELECT id, name FROM user where id=#{id} order by id
    </select>
    <select id="findStructByStruct" resultType="Struct">
        SELECT id Id, name Name, create_time CrtTm FROM user where id=#{Id} order by id
    </select>
    <insert id="insertStruct">
        insert into user (name, email, create_time)
        values (#{Name}, #{Email}, #{CrtTm})
    </insert>
    <delete id="deleteById">
        delete from user where id=#{id}
    </delete>
</mapper>
```

#### 使用方法
```go
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wenj91/gobatis"
	"github.com/wenj91/gobatis/structs"
)

type User struct {
	Id    structs.NullInt64
	Name  structs.NullString
	Email structs.NullString
	CrtTm structs.NullTime
}

func main(){
    mapperPath := []string{"./mapper.xml"}
    gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)
    
    //传入id查询Map
    mapRes := make(map[string]interface{})
    i, err := gobatis.Select("Mapper.findMapById", 1)(mapRes)
    fmt.Println("Mapper.findMapById-->", i, mapRes, err)
    	
    //根据传入实体查询对象
    param := User{
        Id: structs.NullInt64{3, true},
    }
    structRes2 := User{}
    i, err = gobatis.Select("Mapper.findStructByStruct", param)(&structRes2)
    fmt.Println("Mapper.findStructByStruct-->", i, structRes2, err)
}
```
