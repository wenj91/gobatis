package gobatis

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"testing"
	"time"
)

func TestGoBatis(t *testing.T) {

	xmlStr := `
<?xml version="1.0" encoding="utf-8"?>
<mapper namespace="Mapper">
    <select id="findMapById" resultType="Map">
        SELECT id, name FROM user where id=#{id} order by id
    </select>
    <insert id="insertStructsBatch">
        insert into user (name, email, create_time)
        values
        <foreach item="item" collection="list" open="(" close=")" separator=",">
            #{Name}, #{Email}, #{CrtTm}
        </foreach>
    </insert>
    <update id="updateByStruct">
        update user set name = #{Name}, email = #{Email}
        where id = #{Id}
    </update>
    <delete id="deleteById">
        delete from user where id=#{id}
    </delete>
</mapper>
`
	r:= strings.NewReader(xmlStr)
	mapperConf := buildMapperConfig(r)

	ymlStr := `
db:
  driverName: mysql
  dataSourceName: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8
  maxLifeTime: 10
  maxOpenConns: 10
  maxIdleConns: 1
  showSql: true
mappers:
  - userMapper.xml
  - orderMapper.xml
`
	dbconf := buildDbConfig(ymlStr)

	conf := &config{
		dbConf:dbconf,
		mapperConf:mapperConf,
	}

	db, err := sql.Open(conf.dbConf.DB.DriverName, conf.dbConf.DB.DataSourceName)
	if nil != err {
		log.Println(err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Println(err)
		panic(err)
	}


	if conf.dbConf.DB.MaxLifeTime == 0 {
		db.SetConnMaxLifetime(120 * time.Second)
	}else{
		db.SetConnMaxLifetime(time.Duration(conf.dbConf.DB.MaxLifeTime) * time.Second)
	}

	if conf.dbConf.DB.MaxOpenConns == 0 {
		db.SetMaxOpenConns(10)
	}else{
		db.SetMaxOpenConns(conf.dbConf.DB.MaxOpenConns)
	}

	if conf.dbConf.DB.MaxOpenConns == 0 {
		db.SetMaxIdleConns(5)
	}else{
		db.SetMaxIdleConns(conf.dbConf.DB.MaxIdleConns)
	}

	gb := &gobatis{
		gbBase{
			db:db,
			mapperConfig:conf.mapperConf,
		},
	}

	res, _ := gb.db.Query("select 1")

	cols, _ := res.Columns()

	fmt.Println(cols)
}
