package gobatis

type DBOption struct {
	dbs     map[string]*GoBatisDB
	showSQL bool
	mappers []string
}

var _ IOption = &DBOption{}

func NewDBOption() *DBOption {
	return &DBOption{}
}

func NewDBOption_(dbs map[string]*GoBatisDB, showSQL bool, mappers []string) *DBOption {
	return &DBOption{
		dbs:     dbs,
		showSQL: showSQL,
		mappers: mappers,
	}
}

func (ds *DBOption) DB(dbs map[string]*GoBatisDB) *DBOption {
	ds.dbs = dbs
	return ds
}

func (ds *DBOption) ShowSQL(showSQL bool) *DBOption {
	ds.showSQL = showSQL
	return ds
}

func (ds *DBOption) Mappers(mappers []string) *DBOption {
	ds.mappers = mappers
	return ds
}

func (ds *DBOption) Type() OptionType {
	return OptionTypeDB
}

func (ds *DBOption) ToDBConf() *DBConfig {
	dbconf := NewDBConfigBuilder().
		DB(ds.dbs).
		ShowSQL(ds.showSQL).
		Mappers(ds.mappers).
		Build()
	return dbconf
}
