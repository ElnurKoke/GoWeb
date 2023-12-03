package eternals

import (
	"database/sql"
)

var db *sql.DB

func Initdb() {
	// Открываем или создаем базу данных SQLite
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	// Создаем таблицу "user", если она еще не существует
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		password TEXT,
		login TEXT UNIQUE,
		username TEXT,
		point INTEGER DEFAULT 0
	)`)
	if err != nil {
		panic(err)
	}
}
