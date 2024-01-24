package config

// ServerConfig consul上记录的远程配置信息
type ServerConfig struct {
	Port          int             `mapstructure:"default_port"`
	Host          string          `mapstructure:"host"`
	CheckHost     string          `mapstructure:"check_host"`
	Tags          []string        `mapstructure:"tags"`
	RecordSrvInfo RecordSrvConfig `mapstructure:"record_srv"`
	RedisInfo     RedisConfig     `mapstructure:"redis"`
	RabbitMQInfo  RabbitMQConfig  `mapstructure:"rabbitmq"`
}

type RecordSrvConfig struct {
	Name string `mapstructure:"name"`
}

type RedisConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type RabbitMQConfig struct {
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
