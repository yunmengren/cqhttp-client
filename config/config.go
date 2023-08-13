package config

type Config struct {
	SvcAddress    string `yaml:"svc_address"`
	CqhttpAddress string `yaml:"cqhttp_address"`

	// 消息监控
	Keywords []string `yaml:"keywords"`
	ToID     int64    `yaml:"to_id"`

	// 定时任务
	CrontablesCfg []CrontableConfig `yaml:"crontables"`

	LogFile  string `yaml:"log_file"`
	LogLevel string `yaml:"log_level"`
}

type CrontableConfig struct {
	ToID    int64  `yaml:"to_id"`
	Message string `yaml:"message"`
	Cron    string `yaml:"cron"`
}
