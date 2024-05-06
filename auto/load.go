package auto

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	util "github.com/tae2089/go-sql-script-util"
)

func init() {
	filePath := os.Getenv("SQLITE3_FILE_PATH")
	if filePath == "" {
		log.Println("SQLITE3_FILE_PATH is not set")
		return
	}
	dbPath := os.Getenv("SQLITE3_DB_PATH")
	if dbPath == "" {
		log.Println("SQLITE3_DB_PATH is not set")
		return
	}
	db, err := getDB(filePath)
	if err != nil {
		log.Printf("Failed to get db: %v", err)
		return
	}
	if err := util.ExecuteSqlDir(db, dbPath); err != nil {
		log.Printf("Failed to execute sql file: %v", err)
		return
	}
}

func getDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}
