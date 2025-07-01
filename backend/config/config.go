package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"gorm.io/gorm/logger"
)

type Config struct {
	Port         int             `env:"PORT"`
	Host         string          `env:"HOST"`
	GitlabUrl    string          `env:"GITLAB_URL"`
	GitlabToken  string          `env:"GITLAB_TOKEN" json:"-"`
	LogLevel     logger.LogLevel `env:"LOG_LEVEL" envDefault:"1"`
	AllowOrigins []string        `env:"ALLOW_ORIGINS" envSeparator:","`
	Logo         string          `env:"LOGO" envDefault:""`
	Caption      string          `env:"CAPTION" envDefault:""`

	DbType string `env:"DB_TYPE" envDefault:"postgresql"`

	PostgresHost string `env:"POSTGRES_HOST"`
	PostgresPort int    `env:"POSTGRES_PORT"`
	PostgresDb   string `env:"POSTGRES_DB"`
	PostgresUser string `env:"POSTGRES_USER"`
	PostgresPass string `env:"POSTGRES_PASSWORD" json:"-"`

	MysqlHost string `env:"MYSQL_HOST"`
	MysqlPort int    `env:"MYSQL_PORT"`
	MysqlDb   string `env:"MYSQL_DATABASE"`
	MysqlUser string `env:"MYSQL_USER"`
	MysqlPass string `env:"MYSQL_PASSWORD" json:"-"`

	SqliteDbFile string `env:"SQLITE_DB_FILE"`

	JwtTokenLifespanHour uint   `env:"JWT_TOKEN_LIFESPAN_HOUR" envDefault:"24"`
	ApiSecret            string `env:"API_SECRET" json:"-"`

	GitlabSyncEnabled      bool `env:"GITLAB_SYNC_ENABLED" envDefault:"true"`
	GitlabSyncPeriodMin    int  `env:"GITLAB_SYNC_PERIOD_MIN" envDefault:"10"`
	MemoryCacheDurationMin int  `env:"MEMORY_CACHE_DURATION_MIN" envDefault:"15"`
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

func (c *Config) Print() {
	var buf []byte
	buf, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Fatal("Can't marshal config")
	}
	fmt.Printf("Config: %s \n", buf)
}

var Settings = NewConfig()
