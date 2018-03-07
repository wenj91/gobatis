# gobatis

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
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/wenj91/gobatis"
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
   	
   	param := User{
    		Id: structs.NullInt64{3, true},
    }
    structRes2 := User{}
    i, err = gobatis.Select("Mapper#findStructByStruct", param)(&structRes2)
    fmt.Println("Mapper#findStructByStruct-->", i, structRes2, err)
}
```
