package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func InitializeDatabase() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		err := createSchedulerTable(db)
		if err != nil {
			return nil, err
		}
		log.Println("База данных успешно инициализирована.")
	}
	return db, nil
}

func createSchedulerTable(db *sql.DB) error {
	query := `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL,
		title VARCHAR(256) NOT NULL,
		comment TEXT, 
		repeat VARCHAR(128)
	);
	CREATE INDEX scheduler_date ON scheduler(date)
	`
	_, err := db.Exec(query)
	return err
}
