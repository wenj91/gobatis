package gobatis

import (
	"fmt"
	"testing"
)

func TestStaticSqlSource_getBoundSql(t *testing.T)  {
	sss := &staticSqlSource{
		sqlStr:"select * from t_gap where id = #{id} and gap = #{gap}",
		paramMappings:make([]string, 0),
	}

	bs := sss.getBoundSql(map[string]interface{}{
		"id":1,
		"gap":10,
	})

	expc := "select * from t_gap where id = ? and gap = ?"
	assertEqual(bs.sqlStr, expc, "test failed, actual:" + bs.sqlStr)
	assertEqual(bs.params["id"], 1, "test failed, actual:" + fmt.Sprintf("%d", bs.params["id"]))
	assertEqual(bs.params["gap"], 10, "test failed, actual:" + fmt.Sprintf("%d", bs.params["gap"]))
}

func TestDynamicSqlSource_getBoundSql(t *testing.T)  {
	params := map[string]interface{}{
		"name":   "wenj91",
		"array":  []map[string]interface{}{{"idea": "11"}, {"idea": "22"}, {"idea": "33"}},
		"array1": []string{"11", "22", "33"},
		"array2": []s{{A: "aa"}, {A: "bb"}, {A: "cc"}},
	}

	msn := &mixedSqlNode{
		sqlNodes: []iSqlNode{
			&textSqlNode{
				content: "select 1 from t_gap where 1 = 1",
			},
			&ifSqlNode{
				test: "name == 'wenj91'",
				sqlNode: &textSqlNode{
					content: "and name = #{name}",
				},
			},
			&foreachSqlNode{
				sqlNode: &mixedSqlNode{
					sqlNodes: []iSqlNode{
						&textSqlNode{
							content: "#{ item.A }",
						},
					},
				},
				item:       "item",
				open:       "and id in (",
				close:      ")",
				separator:  ",",
				collection: "array2",
			},
		},
	}

	ds := dynamicSqlSource{
		sqlNode: msn,
	}

	bs := ds.getBoundSql(params)

	expc := "select 1 from t_gap where 1 = 1 and name = ? and id in ( ?  , ?  , ?  )"
	assertEqual(bs.sqlStr, expc, "test failed, actual:" + bs.sqlStr)
	assertEqual(bs.params["name"], "wenj91", "test failed, actual:" + fmt.Sprintf("%d", bs.params["id"]))
	assertEqual(bs.extParams["_ls_item_p_item0.A"], "aa", "test failed, actual:" + fmt.Sprintf("%s", bs.extParams["_ls_item_p_item0.A"]))
	assertEqual(bs.extParams["_ls_item_p_item1.A"], "bb", "test failed, actual:" + fmt.Sprintf("%s", bs.extParams["_ls_item_p_item1.A"]))
	assertEqual(bs.extParams["_ls_item_p_item2.A"], "cc", "test failed, actual:" + fmt.Sprintf("%s", bs.extParams["_ls_item_p_item2.A"]))
}
