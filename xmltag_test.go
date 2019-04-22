package gobatis

import (
	"fmt"
	"strings"
	"testing"
)

type s struct {
	A string
	B string
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
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
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
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
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
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
}

func TestForeachSqlNode_build(t *testing.T) {
	ctx := newDynamicContext(map[string]interface{}{
		"array": []int{1, 2, 3},
	})

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
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
	assertEqual(ctx.params["_ls_item_p_item0"], 1, "test failed, actual:"+fmt.Sprintf("%d", ctx.params["_ls_item_p_item0"]))
	assertEqual(ctx.params["_ls_item_p_item1"], 2, "test failed, actual:"+fmt.Sprintf("%d", ctx.params["_ls_item_p_item1"]))
	assertEqual(ctx.params["_ls_item_p_item2"], 3, "test failed, actual:"+fmt.Sprintf("%d", ctx.params["_ls_item_p_item2"]))
}

func TestMixedSqlNode_build(t *testing.T) {
	params := map[string]interface{}{
		"name":   "wenj91",
		"array":  []map[string]interface{}{{"idea": "11"}, {"idea": "22"}, {"idea": "33"}},
		"array1": []string{"11", "22", "33"},
		"array2": []s{{A: "aa"}, {A: "bb"}, {A: "cc"}},
	}

	mixedSqlNode := &mixedSqlNode{
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

	ctx := newDynamicContext(params)

	mixedSqlNode.build(ctx)

	expc := "select 1 from t_gap where 1 = 1 and name = #{name} and id in ( #{_ls_item_p_item0.A}  , #{_ls_item_p_item1.A}  , #{_ls_item_p_item2.A}  )"
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
	assertEqual(ctx.params["_ls_item_p_item0.A"], "aa", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["_ls_item_p_item0.A"]))
	assertEqual(ctx.params["_ls_item_p_item1.A"], "bb", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["_ls_item_p_item1.A"]))
	assertEqual(ctx.params["_ls_item_p_item2.A"], "cc", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["_ls_item_p_item2.A"]))
}

func TestSetSqlNode_build(t *testing.T) {
	params := map[string]interface{}{
		"name":  "wenj91",
		"name2": "wenj91",
	}

	setSqlNode := &setSqlNode{
		sqlNodes: []iSqlNode{
			&ifSqlNode{
				test: "name == 'wenj91'",
				sqlNode: &textSqlNode{
					content: "name = #{name}",
				},
			},
			&ifSqlNode{
				test: "name2 == 'wenj91'",
				sqlNode: &textSqlNode{
					content: "name2 = #{name2}",
				},
			},
		},
	}

	ctx := newDynamicContext(params)

	setSqlNode.build(ctx)

	expc := " set  name = #{name}  , name2 = #{name2} "
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
	assertEqual(ctx.params["name"], "wenj91", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["name"]))
	assertEqual(ctx.params["name2"], "wenj91", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["name2"]))
}

func TestTrimSqlNode_build(t *testing.T) {
	params := map[string]interface{}{
		"name":  "wenj91",
		"name2": "wenj91",
	}

	trimSqlNode := &trimSqlNode{
		prefixOverrides: "and",
		suffixOverrides: ",",
		sqlNodes: []iSqlNode{
			&ifSqlNode{
				test: "name == 'wenj91'",
				sqlNode: &textSqlNode{
					content: "and name = #{name}",
				},
			},
			&ifSqlNode{
				test: "name2 == 'wenj91'",
				sqlNode: &textSqlNode{
					content: "and name2 = #{name2}",
				},
			},
		},
	}

	ctx := newDynamicContext(params)

	trimSqlNode.build(ctx)

	expc := "name = #{name}  and name2 = #{name2} "
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
	assertEqual(ctx.params["name"], "wenj91", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["name"]))
	assertEqual(ctx.params["name2"], "wenj91", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["name2"]))
}

func TestWhereSqlNode_build(t *testing.T) {
	params := map[string]interface{}{
		"name":  "wenj91",
		"name2": "wenj91",
	}

	whereSqlNode := &whereSqlNode{
		sqlNodes: []iSqlNode{
			&ifSqlNode{
				test: "name == 'wenj91'",
				sqlNode: &textSqlNode{
					content: "and name = #{name}",
				},
			},
			&ifSqlNode{
				test: "name2 == 'wenj91'",
				sqlNode: &textSqlNode{
					content: "and name2 = #{name2}",
				},
			},
		},
	}

	ctx := newDynamicContext(params)

	whereSqlNode.build(ctx)

	expc := "where name = #{name}  and name2 = #{name2} "
	assertEqual(ctx.toSql(), expc, "test failed, actual:"+ctx.toSql())
	assertEqual(ctx.params["name"], "wenj91", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["name"]))
	assertEqual(ctx.params["name2"], "wenj91", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["name2"]))
}

func TestChooseSqlNode_build(t *testing.T) {
	params := map[string]interface{}{
		"name": "aa",
	}

	choose := chooseNode{
		sqlNodes: []iSqlNode{
			&ifSqlNode{
				test: "name == 'sean'",
				sqlNode: &textSqlNode{
					content: "and name = 'sean' ",
				},
			},
			&ifSqlNode{
				test: "name == 'Sean'",
				sqlNode: &textSqlNode{
					content: "and name = #{name} ",
				},
			},
		},
		otherwise: &mixedSqlNode{
			sqlNodes: []iSqlNode{
				&ifSqlNode{
					test: "name == 'aa'",
					sqlNode: &textSqlNode{
						content: "and name = 'aa' ",
					},
				},
			},
		},
	}

	ctx := newDynamicContext(params)
	choose.build(ctx)
	expc := "and name = 'aa'"
	assertEqual(strings.Trim(ctx.toSql(), " "), expc, "test failed, actual:"+ctx.toSql())
	assertEqual(ctx.params["name"], "aa", "test failed, actual:"+fmt.Sprintf("%s", ctx.params["name"]))
}
