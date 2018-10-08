package gobatis

type dynamicContext struct {
	sqlStr string
	params map[string]interface{}
}

func (this *dynamicContext) appendSql(sqlStr string)  {
	this.sqlStr += sqlStr + " "
}

type iSqlNode interface {
	build(ctx *dynamicContext) bool
}

// mixed node
type mixedSqlNode struct {
	sqlNodes []iSqlNode
}

func (this *mixedSqlNode) build(ctx *dynamicContext) bool {
	for i := 0; i < len(this.sqlNodes); i++ {
		sqlNode := this.sqlNodes[i]
		sqlNode.build(ctx)
	}

	return true
}

// if node
type ifSqlNode struct {
	test    string
	sqlNode iSqlNode
}

func (this *ifSqlNode) build(ctx *dynamicContext) bool {
	if ok := exprProcess(this.test, ctx.params); ok {
		this.sqlNode.build(ctx)
		return true
	}

	return false
}


// text node
type textSqlNode struct {
	content string
}

func (this *textSqlNode) build(ctx *dynamicContext) bool {
	ctx.appendSql(this.content)
	return true
}

// for node
type foreachSqlNode struct {
	sqlNode    iSqlNode
	collection string
	open       string
	close      string
	separator  string
	item       string
	index      string
}

func (this *foreachSqlNode) build(ctx *dynamicContext) bool {

	return true
}

// set node

// trim node

// where node

// choose node
