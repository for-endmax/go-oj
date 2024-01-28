package config

// ServerConfig consul上记录的远程配置信息
type ServerConfig struct {
	Port            int               `mapstructure:"default_port"`
	Host            string            `mapstructure:"host"`
	CheckHost       string            `mapstructure:"check_host"`
	Tags            []string          `mapstructure:"tags"`
	QuestionSrvInfo QuestionSrvConfig `mapstructure:"question_srv"`
	RabbitMQInfo    RabbitMQConfig    `mapstructure:"rabbitmq"`
	JaegerInfo      JaegerConfig      `mapstructure:"jaeger"`
}
type JaegerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}
type QuestionSrvConfig struct {
	Name string `mapstructure:"name"`
}
type RabbitMQConfig struct {
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
