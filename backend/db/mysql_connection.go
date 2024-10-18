package db

import (
	"fmt"
	"log"
	"main/config"
	"main/repository"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlConnection struct {
	config RDBConnectionConfig
	db     *gorm.DB
}

func NewMysqlConnection(host string, port int, dbName string, user string, pass string) (repository.ConnectionInterface, error) {
	connection := &MysqlConnection{
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

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: connection.getDSN(),
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	connection.db = db

	return connection, nil
}

func (c *MysqlConnection) GetEngine() interface{} {
	return c.db
}

func (c *MysqlConnection) Migrate(models ...interface{}) error {
	return c.db.AutoMigrate(models...)
}

func (c *MysqlConnection) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.config.user, c.config.pass, c.config.host, c.config.port, c.config.dbName)
}
