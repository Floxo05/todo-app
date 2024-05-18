package tools

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		dbname := os.Getenv("DB_NAME")

		dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

		db, err = sql.Open("mysql", dbURL)
		if err != nil {
			err = fmt.Errorf("error opening database: %v", err)
			return
		}

		err = db.Ping()
		if err != nil {
			err = fmt.Errorf("error connecting to the database: %v", err)
			return
		}
	})

	return db, err
}
