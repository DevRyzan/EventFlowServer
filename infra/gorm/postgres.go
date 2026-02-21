package gorm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Open creates a GORM DB connection from a Postgres DSN.
func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
