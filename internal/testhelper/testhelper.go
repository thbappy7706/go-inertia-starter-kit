package testhelper

import (
	"context"
	"testing"
	"time"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=starter_kit_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	db.Exec("TRUNCATE TABLE users, products RESTART IDENTITY CASCADE")

	cleanup := func() {
		db.Exec("TRUNCATE TABLE users, products RESTART IDENTITY CASCADE")
	}

	return db, cleanup
}
