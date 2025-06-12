package database

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	globalDB *gorm.DB
	once     sync.Once
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	// Additional configuration options
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// DSN returns the Data Source Name for the PostgreSQL connection
func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

// Connect establishes a connection to the database
func (p *PostgresConfig) Connect() error {
	if p.Host == "" || p.Port == "" || p.User == "" || p.DBName == "" {
		return errors.New("missing required database configuration")
	}

	var connErr error
	once.Do(func() {
		config := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}

		db, err := gorm.Open(postgres.Open(p.DSN()), config)
		if err != nil {
			connErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		// Configure connection pool
		sqlDB, err := db.DB()
		if err != nil {
			connErr = fmt.Errorf("failed to get database instance: %w", err)
			return
		}

		// Set connection pool settings
		if p.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		}
		if p.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		}
		if p.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(p.ConnMaxLifetime)
		}

		globalDB = db
		fmt.Println("Connected to the database successfully")
	})

	if connErr != nil {
		return connErr
	}

	return nil
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return globalDB
}

// Migrate applies the database migrations
func (p *PostgresConfig) Migrate(models ...any) error {
	if globalDB == nil {
		return errors.New("database connection is not initialized")
	}

	err := globalDB.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	fmt.Println("Database migration completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if globalDB == nil {
		return nil
	}

	sqlDB, err := globalDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	return sqlDB.Close()
}
