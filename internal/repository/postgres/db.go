package postgres

import (
	"strings"

	"github.com/yourorg/go-user-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB opens a connection pool to Postgres using GORM.
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)

	return db, nil
}

// buildDSN constructs a libpq-style connection string. Fields with empty
// values (notably password, common for local trust-auth setups) are omitted
// entirely rather than written as "key=" — an empty unquoted value can break
// parsing of every key that follows it.
func buildDSN(cfg *config.Config) string {
	parts := []string{
		"host=" + cfg.DBHost,
		"user=" + cfg.DBUser,
		"dbname=" + cfg.DBName,
		"port=" + cfg.DBPort,
		"sslmode=" + cfg.DBSSLMode,
		"TimeZone=UTC",
	}

	if cfg.DBPassword != "" {
		parts = append(parts, "password="+cfg.DBPassword)
	}

	return strings.Join(parts, " ")
}
