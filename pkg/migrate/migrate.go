package migrate

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
)

var db *sql.DB

func RunMigrations() {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration", // Путь к папке с миграциями
		"postgres",         // Имя базы данных
		driver,
	)
	if err != nil {
		log.Fatalf("Error initializing migration: %v", err)
	}

	// Применяем все миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error applying migration: %v", err)
	}

	log.Println("Migrations applied successfully")
}

// Откат последней миграции
func RollbackLastMigration() {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Error initializing migration: %v", err)
	}

	// Откат последней миграции
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error rolling back migration: %v", err)
	}

	log.Println("Last migration rolled back successfully")
}
