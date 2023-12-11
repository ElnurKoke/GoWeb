package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Открываем или создаем базу данных SQLite
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицу
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// Добавляем данные в таблицу
	insertData := "INSERT INTO users (name, age) VALUES (?, ?);"
	_, err = db.Exec(insertData, "John Doe", 30)
	if err != nil {
		log.Fatal(err)
	}

	// Выполняем запрос и выводим результаты
	rows, err := db.Query("SELECT id, name, age FROM users;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {
		var id, age int
		var name string
		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

	// Проверяем ошибки выполнения запроса
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
