package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func tableExists(db *sql.DB, tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := db.QueryRow(query, tableName)
	var name string
	if err := row.Scan(&name); err != nil {
		return false
	}
	return name == tableName
}

func InitializeDatabase() (*sql.DB, error) {

	// basePath, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }

	dbFile := filepath.Join("../database", "scheduler.db")
	//_, err = os.Stat(dbFile)

	var install bool
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
		log.Println("Файл базы данных отсутствует. Требуется инициализация.")
	} else if err != nil {
		return nil, err
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Ошибка при подключении к базе данных: %v", err)
		return nil, err
	}

	if install {
		err := createSchedulerTable(db)
		if err != nil {
			return nil, err
		}
		log.Println("База данных успешно инициализирована.")
	}

	if !tableExists(db, "scheduler") {
		log.Println("Таблица 'scheduler' отсутствует. Создаём...")
		if err := createSchedulerTable(db); err != nil {
			return nil, err
		}
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
	    CREATE INDEX scheduler_date ON scheduler(date);
	    `
	log.Println("Проверяем создание таблицы...")
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Ошибка при создании таблицы: %v", err)
		return err
	}
	log.Println("Таблица успешно создана.")
	return nil
}
