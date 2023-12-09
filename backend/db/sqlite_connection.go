package db

import (
	"log"
	"main/config"
	"main/tools"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqliteConnection struct {
	dbfile string
	db     *gorm.DB
}

func NewSqliteConnection(dbfile string) (tools.ConnectionInterface, error) {
	connection := &SqliteConnection{
		dbfile: dbfile,
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

	db, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	connection.db = db

	return connection, nil
}

func (c *SqliteConnection) GetEngine() interface{} {
	return c.db
}

func (c *SqliteConnection) Migrate(models ...interface{}) error {
	return c.db.AutoMigrate(models...)
}
