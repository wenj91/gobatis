package gobatis

type boundSql struct {
	sqlStr string
	paramMappings []string
	params map[string]interface{}
	extParams map[string]interface{}
}

type iSqlSource interface {
	getBoundSql() boundSql
}

type dynamicSqlSource struct {
	sqlNode iSqlNode
}

func (*dynamicSqlSource) getBoundSql() boundSql {
	panic("implement me")
}

type staticSqlSource struct {
	sqlStr string
	paramMappings []string
}

func (*staticSqlSource) getBoundSql() boundSql {
	panic("implement me")
}


