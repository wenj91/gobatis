package gobatis

type mappedStmt struct {
	dbType     DBType
	sqlSource  iSqlSource
	resultType ResultType
}
