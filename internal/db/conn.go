package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // Blank import registers the driver
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func Connect(cfg Config) (*sql.DB, error) {
	// DSN format: username:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening db connection: %w", err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(25)                 // Max active connections
	db.SetMaxIdleConns(25)                 // Max idle connections in pool
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum age of a connection
	db.SetConnMaxIdleTime(5 * time.Minute) // Max time a connection can remain idle

	// Ping to verify the credentials and network connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	return db, nil
}
