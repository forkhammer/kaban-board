package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"main/config"
	"os"
	"time"
)

type Connection struct {
	Db     *gorm.DB
	host   string
	port   int
	dbName string
	user   string
	pass   string
}

func NewConnection(host string, port int, dbName string, user string, pass string) (*Connection, error) {
	connection := &Connection{
		host:   host,
		port:   port,
		dbName: dbName,
		user:   user,
		pass:   pass,
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

	connection.Db = db

	return connection, nil
}

func (c *Connection) getDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC", c.host, c.port, c.user, c.pass, c.dbName)
}

var DefaultConnection, _ = NewConnection(
	config.Settings.PostgresHost,
	config.Settings.PostgresPort,
	config.Settings.PostgresDb,
	config.Settings.PostgresUser,
	config.Settings.PostgresPass,
)
