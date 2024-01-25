package config

// LocalConfig 本地配置信息
type LocalConfig struct {
	Name   string       `mapstructure:"name"` //服务名称
	Consul ConsulConfig `mapstructure:"consul"`
}

// ConsulConfig consul的地址和端口
type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
