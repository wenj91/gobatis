package gobatis

type dbConfig struct {
	DB struct {
		DriverName     string `yaml:"driverName"`
		DataSourceName string `yaml:"dataSourceName"`
		MaxLifeTime    int    `yaml:"maxLifeTime"`
		MaxOpenConns   int    `yaml:"maxOpenConns"`
		MaxIdleConns   int    `yaml:"maxIdleConns"`
	} `yaml:"db"`
	Mappers []string `yaml:"mappers"`
}
