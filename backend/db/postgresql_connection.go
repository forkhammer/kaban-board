package db

import (
	"fmt"
	"log"
	"main/config"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresqlConnection struct {
	config RDBConnectionConfig
	db     *gorm.DB
}

func NewPostgresqlConnection(host string, port int, dbName string, user string, pass string) (*PostgresqlConnection, error) {
	connection := &PostgresqlConnection{
		config: RDBConnectionConfig{
			host:   host,
			port:   port,
			dbName: dbName,
			user:   user,
			pass:   pass,
		},
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,              // Slow SQL threshold
			LogLevel:                  config.Settings.LogLevel, // Log level
			IgnoreRecordNotFoundError: true,                     // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,                     // Don't include params in the SQL log
			Colorful:                  false,                    // Disable color
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connection.getDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	connection.db = db

	return connection, nil
}

func (c *PostgresqlConnection) GetEngine() interface{} {
	return c.db
}

func (c *PostgresqlConnection) Migrate(models ...interface{}) error {
	return c.db.AutoMigrate(models...)
}

func (c *PostgresqlConnection) getDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC", c.config.host, c.config.port, c.config.user, c.config.pass, c.config.dbName)
}
