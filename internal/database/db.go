package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
)

//const (
//	host     = "localhost"
//	port     = "5432"
//	user     = "user"
//	password = "password"
//	dbname   = "authdb"
//)

func Initialize(logger *zap.SugaredLogger) (*sql.DB, error) {
	// Чтение переменных окружения
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Проверка, что все переменные окружения установлены
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		logger.Error("Database configuration error: one or more environment variables are missing")
		return nil, fmt.Errorf("database configuration error")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Infow("Successfully connected!")

	return db, nil
}

func SetupSchema(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil

}
