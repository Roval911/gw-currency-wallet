package repository

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func InitDb() {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
}

// Закрытие соединения с базой данных
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

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

func SetDB(mockDB *sql.DB) {
	db = mockDB
}
