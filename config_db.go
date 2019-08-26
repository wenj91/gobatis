package gobatis

type DataSource struct {
	DataSource     string `yaml:"datasource"`
	DriverName     string `yaml:"driverName"`
	DataSourceName string `yaml:"dataSourceName"`
	MaxLifeTime    int    `yaml:"maxLifeTime"`
	MaxOpenConns   int    `yaml:"maxOpenConns"`
	MaxIdleConns   int    `yaml:"maxIdleConns"`
}

// NewDataSource new data source
func NewDataSource(dataSource string, driverName string, dataSourceName string) *DataSource {
	return &DataSource{
		DataSource:     dataSource,
		DriverName:     driverName,
		DataSourceName: dataSourceName,
	}
}

// NewDataSource_ new data source
func NewDataSource_(dataSource string, driverName string, dataSourceName string,
	maxLifeTime int, maxOpenConns int, maxIdleConns int) *DataSource {
	return &DataSource{
		DataSource:     dataSource,
		DriverName:     driverName,
		DataSourceName: dataSourceName,
		MaxLifeTime:    maxLifeTime,
		MaxOpenConns:   maxOpenConns,
		MaxIdleConns:   maxIdleConns,
	}
}

type DBConfig struct {
	DB      []*DataSource `yaml:"db"`
	ShowSQL bool          `yaml:"showSQL"`
	Mappers []string      `yaml:"mappers"`
	db      map[string]*GoBatisDB
	dbMap   map[string]*DataSource
}

func (this *DBConfig) getDataSourceByName(datasource string) *DataSource {
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
