package implementation

import (
	"database/sql"
	"log"
	"main/config"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	sql.Register("sqlite3_unicode", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("lower_unicode", func(s string) string {
				return strings.ToLower(s)
			}, true)
		},
	})
}

type SqliteConnection struct {
	dbfile string
	db     *gorm.DB
}

func NewSqliteConnection(dbfile string, cfg *config.Config) (*SqliteConnection, error) {
	connection := &SqliteConnection{
		dbfile: dbfile,
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,              // Slow SQL threshold
			LogLevel:                  cfg.LogLevel, // Log level
			IgnoreRecordNotFoundError: true,                     // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,                     // Don't include params in the SQL log
			Colorful:                  false,                    // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Dialector{DSN: dbfile, DriverName: "sqlite3_unicode"}, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	connection.db = db

	return connection, nil
}

func (c *SqliteConnection) GetEngine() *gorm.DB {
	return c.db
}

func (c *SqliteConnection) Migrate(models ...any) error {
	return c.db.AutoMigrate(models...)
}
