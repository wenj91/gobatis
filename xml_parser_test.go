package gobatis

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestXmlNode_parse(t *testing.T)  {
	xmlStr := `
<?xml version="1.0" encoding="utf-8"?>
<mapper namespace="Mapper">
    <select id="findMapById" resultType="Map">
		SELECT id, name FROM user where id=#{id} order by id
    </select>
</mapper>
`
	r := strings.NewReader(xmlStr)

	node := parse(r)
	res, _ := json.Marshal(node)
	fmt.Println(string(res))
}