package gobatis

type mappedStmt struct {
	dbid       string
	sqlSource  iSqlSource
	resultType ResultType
}
