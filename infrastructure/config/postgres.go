package config

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

func InitPostgres() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Errorf("Failed to load .env file: %v", err)
		return nil, err
	}

	host := os.Getenv("POSTGRES_HOST")
	logrus.Infof("POSTGRES_HOST: %s", host)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.Errorf("Failed to open connection to Postgres: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logrus.Errorf("Failed to ping Postgres: %v", err)
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		logrus.Errorf("Failed to run migrations: %v", err)
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	migration := `
    CREATE TABLE IF NOT EXISTS queue (
        id SERIAL PRIMARY KEY,
        file_id VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err := db.Exec(migration)
	return err
}
