package database

import (
	"database/sql"
	"os"
)

func InitBD(db *sql.DB, schemaFile string) error {
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	return err
}
