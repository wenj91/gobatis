package gobatis

import (
	"strings"
	"testing"
)

func TestXmlNode_parse(t *testing.T) {
	xmlStr := `
<?xml version="1.0" encoding="utf-8"?>
<mapper namespace="Mapper">
    <select id="findMapById" resultType="Map">
		SELECT id, name FROM user where id=#{id} 
		<if test="name != nil">
			<foreach item="item" open="AND name in (" close=")" separator="," collection="names">
				#{item}
			</foreach>
			AND name = #{name}
		</if>
		ORDER BY id
    </select>
	<update id="updateById">
		UPDATE t_gap SET gap = #{gap} WHERE id = #{id}
	</update>
</mapper>
`
	r := strings.NewReader(xmlStr)
	rn := parse(r)
	assertNotNil(rn, "Parse xml result is nil")
}
