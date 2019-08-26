package gobatis

type DSOption struct {
	dss     []*DataSource
	showSQL bool
	mappers []string
}

var _ IOption = &DSOption{}

func NewDSOption() *DSOption {
	return &DSOption{}
}

func NewDSOption_(dss []*DataSource, showSQL bool, mappers []string) *DSOption {
	return &DSOption{
		dss:     dss,
		showSQL: showSQL,
		mappers: mappers,
	}
}

func (ds *DSOption) DS(dss []*DataSource) *DSOption {
	ds.dss = dss
	return ds
}

func (ds *DSOption) ShowSQL(showSQL bool) *DSOption {
	ds.showSQL = showSQL
	return ds
}

func (ds *DSOption) Mappers(mappers []string) *DSOption {
	ds.mappers = mappers
	return ds
}

func (ds *DSOption) Type() OptionType {
	return OptionTypeDS
}

func (ds *DSOption) ToDBConf() *DBConfig {
	dbconf := NewDBConfigBuilder().
		DS(ds.dss).
		ShowSQL(ds.showSQL).
		Mappers(ds.mappers).
		Build()
	return dbconf
}
