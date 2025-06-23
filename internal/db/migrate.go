package db

import (
	"log"
	"subgen/internal/admin"
	"subgen/internal/config"
	"subgen/internal/userlink"
)

func Migrate() {
	if DB == nil {
		log.Fatal("DB not initialized")
	}
	if err := DB.AutoMigrate(&admin.Admin{}, &config.Config{}, &userlink.UUID{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
