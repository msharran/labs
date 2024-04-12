package server

import (
	"fmt"
	"go-htmx-kvstore/internal/web/data"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(name string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// enable sqlite foreign key support
	result := db.Exec("PRAGMA foreign_keys = ON")
	if result.Error != nil {
		return nil, fmt.Errorf("error enabling foreign key support: %w", result.Error)
	}

	err = db.AutoMigrate(&data.User{}, &data.KeyValue{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
