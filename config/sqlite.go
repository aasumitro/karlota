package config

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SQLiteConnection() Option {
	return func(cfg *Config) {
		conn, err := gorm.Open(sqlite.Open(cfg.SQLiteDsnURL), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			//Logger: logger.New(
			//	log.New(os.Stdout, "\r\n", log.LstdFlags),
			//	logger.Config{
			//		SlowThreshold: time.Second,
			//		LogLevel:      logger.Info,
			//		Colorful:      true,
			//	},
			//),
		})
		if err != nil {
			log.Fatalf("GORM_OPEN_ERROR: %s", err.Error())
		}
		db, err := conn.DB()
		if err != nil {
			log.Fatalf("DATABASE_ERROR: %s", err.Error())
		}
		// Set SQLite PRAGMA configurations for optimized RW
		pragmaStatements := []string{
			"PRAGMA synchronous = OFF;",
			"PRAGMA journal_mode = WAL;",
			"PRAGMA temp_store = MEMORY;",
		}
		for _, pragma := range pragmaStatements {
			if _, err := db.Exec(pragma); err != nil {
				log.Fatalf("SQLITE_PRAGMA_ERROR: %v", err)
			}
		}
		// Configure connection pooling
		const dbMaxOpenConnection, dbMaxIdleConnection = 100, 10
		db.SetMaxIdleConns(dbMaxIdleConnection)
		db.SetMaxOpenConns(dbMaxOpenConnection)
		db.SetConnMaxLifetime(time.Hour)
		// Validate the database connection
		if err := db.Ping(); err != nil {
			log.Fatalf("SQLITE_PING_ERROR: %s\n", err.Error())
		}
		// Assign the configured DB connection to the config
		cfg.Infra.GormPool = conn
		cfg.Infra.SQLPool = db
	}
}
