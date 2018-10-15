package gobatis

type mappedStmt struct {
	dbType     DbType
	sqlSource  iSqlSource
	resultType ResultType
}
