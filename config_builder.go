package gobatis

type DataSourceBuilder struct {
	ds *DataSource
}

func NewDataSourceBuilder() *DataSourceBuilder {
	return &DataSourceBuilder{
		ds: &DataSource{},
	}
}

// DataSource
func (d *DataSourceBuilder) DataSource(ds string) *DataSourceBuilder {
	d.ds.DataSource = ds
	return d
}

// DriverName
func (d *DataSourceBuilder) DriverName(dn string) *DataSourceBuilder {
	d.ds.DriverName = dn
	return d
}

// DataSourceName
func (d *DataSourceBuilder) DataSourceName(dsn string) *DataSourceBuilder {
	d.ds.DataSourceName = dsn
	return d
}

// MaxLifeTime
func (d *DataSourceBuilder) MaxLifeTime(mlt int) *DataSourceBuilder {
	d.ds.MaxLifeTime = mlt
	return d
}

// MaxOpenConns
func (d *DataSourceBuilder) MaxOpenConns(moc int) *DataSourceBuilder {
	d.ds.MaxOpenConns = moc
	return d
}

// MaxIdleConns
func (d *DataSourceBuilder) MaxIdleConns(mic int) *DataSourceBuilder {
	d.ds.MaxIdleConns = mic
	return d
}

func (d *DataSourceBuilder) Build() *DataSource {
	if d.ds.DataSource == "" {
		panic("DataSource is nil")
	}

	if d.ds.DataSourceName == "" {
		panic("DataSourceName is nil")
	}

	if d.ds.DriverName == "" {
		panic("DriverName is nil")
	}

	return d.ds
}

type DBConfigBuilder struct {
	d *DBConfig
}

func NewDBConfigBuilder() *DBConfigBuilder {
	return &DBConfigBuilder{
		d: &DBConfig{
			DB: make([]*DataSource, 0),
		},
	}
}

func (d *DBConfigBuilder) Mappers(mappers []string) *DBConfigBuilder {
	d.d.Mappers = mappers
	return d
}

func (d *DBConfigBuilder) DS(dss []*DataSource) *DBConfigBuilder {
	d.d.DB = dss
	return d
}

func (d *DBConfigBuilder) DB(db map[string]*GoBatisDB) *DBConfigBuilder {
	d.d.db = db
	return d
}

func (d *DBConfigBuilder) ShowSQL(showSQL bool) *DBConfigBuilder {
	d.d.ShowSQL = showSQL
	return d
}

func (d *DBConfigBuilder) Build() *DBConfig {
	if len(d.d.DB) <= 0 && d.d.db == nil {
		panic("No config for datasource")
	}

	return d.d
}
