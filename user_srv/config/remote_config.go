package config

// ServerConfig consul上记录的远程配置信息
type ServerConfig struct {
	Port  string      `mapstructure:"port"`
	Mysql MySQLConfig `mapstructure:"mysql"`
}

// MySQLConfig MySQL信息
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"db_name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
