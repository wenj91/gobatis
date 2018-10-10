package gobatis

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

type s struct {
	A string
}

func TestTextSqlNode_build(t *testing.T) {

	ctx := &dynamicContext{
		params: map[string]interface{}{},
	}

	textSqlNode := &textSqlNode{
		content: "select 1 from t_gap",
	}

	textSqlNode.build(ctx)

	expc := "select 1 from t_gap "
	assertEqual(t, ctx.sqlStr, expc, "test failed, actual:" + ctx.sqlStr)
}

func TestIfSqlNode_True_build(t *testing.T) {
	ctx := &dynamicContext{
		params: map[string]interface{}{
			"name": "wenj91",
		},
	}

	ifSqlNode := &ifSqlNode{
		test: "name == 'wenj91'",
		sqlNode: &textSqlNode{
			content: "select 1 from t_gap",
		},
	}

	ifSqlNode.build(ctx)

	expc := "select 1 from t_gap "
	assertEqual(t, ctx.sqlStr, expc, "test failed, actual:" + ctx.sqlStr)
}

func TestIfSqlNode_False_build(t *testing.T) {
	ctx := &dynamicContext{
		params: map[string]interface{}{
			"name": "wenj91",
		},
	}

	ifSqlNode := &ifSqlNode{
		test: "name != 'wenj91'",
		sqlNode: &textSqlNode{
			content: "select 1 from t_gap",
		},
	}

	ifSqlNode.build(ctx)

	expc := ""
	assertEqual(t, ctx.sqlStr, expc, "test failed, actual:" + ctx.sqlStr)
}

func TestForeachSqlNode_build(t *testing.T)  {
	ctx := &dynamicContext{
		params: map[string]interface{}{
			"array": []int{1, 2, 3},
		},
	}

	f := &foreachSqlNode{
		sqlNode: &mixedSqlNode{
			sqlNodes: []iSqlNode{
				&textSqlNode{
					content: "#{ item }",
				},
			},
		},
		item:       "item",
		open:       "select 1 from t_gap where id in (",
		close:      ")",
		separator:  ",",
		collection: "array",
	}

	f.build(ctx)

	expc := "select 1 from t_gap where id in ( #{_ls_item_p_item0}  , #{_ls_item_p_item1}  , #{_ls_item_p_item2}  ) "
	assertEqual(t, ctx.sqlStr, expc, "test failed, actual:" + ctx.sqlStr)
	assertEqual(t, ctx.params["_ls_item_p_item0"], 1, "test failed, actual:" + fmt.Sprintf("%d", ctx.params["_ls_item_p_item0"]))
	assertEqual(t, ctx.params["_ls_item_p_item1"], 2, "test failed, actual:" + fmt.Sprintf("%d", ctx.params["_ls_item_p_item1"]))
	assertEqual(t, ctx.params["_ls_item_p_item2"], 3, "test failed, actual:" + fmt.Sprintf("%d", ctx.params["_ls_item_p_item2"]))
}

func TestMixedSqlNode_build(t *testing.T) {
	params := map[string]interface{}{
		"name":   "wenj91",
		"array":  []map[string]interface{}{{"idea": "11"}, {"idea": "22"}, {"idea": "33"}},
		"array1": []string{"11", "22", "33"},
		"array2": []s{{A: "aa"}, {A: "bb"}, {A: "cc"}},
	}

	mixedSqlNode := mixedSqlNode{
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
						&ifSqlNode{
							test: "item.B == nil",
							sqlNode: &textSqlNode{
								content: "1, ",
							},
						},
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

	ctx := &dynamicContext{
		params: params,
	}

	mixedSqlNode.build(ctx)

	expc := "select 1 from t_gap where 1 = 1 and name = #{name} and id in ( #{_ls_item_p_item0.A}  , #{_ls_item_p_item1.A}  , #{_ls_item_p_item2.A}  ) "
	assertEqual(t, ctx.sqlStr, expc, "test failed, actual:" + ctx.sqlStr)
	assertEqual(t, ctx.params["_ls_item_p_item0.A"], "aa", "test failed, actual:" + fmt.Sprintf("%s", ctx.params["_ls_item_p_item0.A"]))
	assertEqual(t, ctx.params["_ls_item_p_item1.A"], "bb", "test failed, actual:" + fmt.Sprintf("%s", ctx.params["_ls_item_p_item1.A"]))
	assertEqual(t, ctx.params["_ls_item_p_item2.A"], "cc", "test failed, actual:" + fmt.Sprintf("%s", ctx.params["_ls_item_p_item2.A"]))
}
