package boundsql

type BoundSql struct {
	Sql               string
	ParameterMappings []interface{}
	ResultType        string
}
