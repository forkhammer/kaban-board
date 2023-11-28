package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"gorm.io/gorm/logger"
	"log"
)

type Config struct {
	Port         int             `env:"PORT"`
	Host         string          `env:"HOST"`
	GitlabUrl    string          `env:"GITLAB_URL"`
	GitlabToken  string          `env:"GITLAB_TOKEN"`
	LogLevel     logger.LogLevel `env:"LOG_LEVEL" envDefault:"1"`
	AllowOrigins []string        `env:"ALLOW_ORIGINS" envSeparator:","`

	PostgresHost string `env:"POSTGRES_HOST"`
	PostgresPort int    `env:"POSTGRES_PORT"`
	PostgresDb   string `env:"POSTGRES_DB"`
	PostgresUser string `env:"POSTGRES_USER"`
	PostgresPass string `env:"POSTGRES_PASSWORD"`

	JwtTokenLifespanHour uint   `env:"JWT_TOKEN_LIFESPAN_HOUR" envDefault:"24"`
	ApiSecret            string `env:"API_SECRET"`

	GitlabSyncPeriodMin    int `env:"GITLAB_SYNC_PERIOD_MIN" envDefault:"10"`
	MemoryCacheDurationMin int `env:"MEMORY_CACHE_DURATION_MIN" envDefault:"15"`
}

func NewConfig() *Config {
	config := Config{}
	config.parseConfig()
	return &config
}

func (c *Config) parseConfig() {
	if err := env.Parse(c); err != nil {
		log.Fatal("Can't parse config")
	}
}

func (c *Config) GetHostPort() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

var Settings = NewConfig()
