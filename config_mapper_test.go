package gobatis

import (
	"strings"
	"testing"
)

func TestBuildConfig(t *testing.T)  {
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
	<update id="updateByCond">
        update user 
		<set>
			<if test="Name != nil and Name != ''">name = #{Name},</if>
			<if test="Email != nil and Email != ''">email = #{Email},</if>
		</set>
        where id = #{Id}
    </update>
    <delete id="deleteById">
        delete from user where id=#{id}
    </delete>
</mapper>
`
	r:= strings.NewReader(xmlStr)
	conf := buildMapperConfig(r)
	assertNotNil(conf.getMappedStmt("Mapper.findMapById"), "Mapper.findMapById mapped stmt is nil")
	assertNotNil(conf.getMappedStmt("Mapper.insertStructsBatch"), "Mapper.insertStructsBatch mapped stmt is nil")
	assertNotNil(conf.getMappedStmt("Mapper.updateByStruct"), "Mapper.updateByStruct mapped stmt is nil")
	assertNotNil(conf.getMappedStmt("Mapper.deleteById"), "Mapper.deleteById mapped stmt is nil")
	assertNotNil(conf.getMappedStmt("Mapper.updateByCond"), "Mapper.deleteById mapped stmt is nil")
}
