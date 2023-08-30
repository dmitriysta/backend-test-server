package database

import (
	"back-end/logs"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func ConnectToDataBase(config DatabaseConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logs.LogError("Failed to open database: %s", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		logs.LogError("Failed to connect to database: %s", err)
		return nil
	}

	return db
}
