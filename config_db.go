package gobatis

type DataSource struct {
	DataSource     string `yaml:"datasource"`
	DriverName     string `yaml:"driverName"`
	DataSourceName string `yaml:"dataSourceName"`
	MaxLifeTime    int    `yaml:"maxLifeTime"`
	MaxOpenConns   int    `yaml:"maxOpenConns"`
	MaxIdleConns   int    `yaml:"maxIdleConns"`
}

func NewDataSource() *DataSource {
	return &DataSource{}
}

type dbConfig struct {
	DB      []*DataSource `yaml:"db"`
	ShowSql bool          `yaml:"showSql"`
	Mappers []string      `yaml:"mappers"`
	dbMap   map[string]*DataSource
}

func NewDbConfig() *dbConfig {
	return &dbConfig{}
}

func (this *dbConfig) getDataSourceByName(datasource string) *DataSource {
	if this.dbMap == nil {
		this.dbMap = make(map[string]*DataSource)
	}

	if v, ok := this.dbMap[datasource]; ok {
		return v
	}

	for _, v := range this.DB {
		if v.DataSource == datasource {
			this.dbMap[datasource] = v
			return v
		}
	}

	return nil
}
