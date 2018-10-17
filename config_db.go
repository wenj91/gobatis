package gobatis

type dbConfig struct {
	DB map[string]struct {
		DriverName     string `yaml:"driverName"`
		DataSourceName string `yaml:"dataSourceName"`
		MaxLifeTime    int    `yaml:"maxLifeTime"`
		MaxOpenConns   int    `yaml:"maxOpenConns"`
		MaxIdleConns   int    `yaml:"maxIdleConns"`
	} `yaml:"db"`
	ShowSql bool     `yaml:"showSql"`
	Mappers []string `yaml:"mappers"`
}
