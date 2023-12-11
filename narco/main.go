package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// User структура представляет данные о пользователе
type User struct {
	ID       int
	Username string
	Password string
}

func main() {
	// Инициализация базы данных SQLite3
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы пользователей, если её нет
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Обработчики маршрутов
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/authenticate", authenticateHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/createAccount", createAccountHandler)

	// Запуск сервера на порту 8080
	log.Println("Server started on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Отображение HTML формы
	tmpl, err := template.New("index").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Login</title>
	</head>
	<body>
		<h2>Login</h2>
		<form action="/authenticate" method="post">
			<label for="username">Username:</label>
			<input type="text" id="username" name="username" required><br>
			<label for="password">Password:</label>
			<input type="password" id="password" name="password" required><br>
			<input type="submit" value="Login">
		</form>
		<p>Don't have an account? <a href="/register">Register here</a>.</p>
	</body>
	</html>
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// В этой функции вы можете добавить логику, связанную с обработкой страницы входа
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка аутентификации
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := getUser(username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if user.Password == password {
		fmt.Fprintf(w, "Welcome, %s!", user.Username)
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Отображение HTML формы регистрации
	tmpl, err := template.New("register").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Register</title>
	</head>
	<body>
		<h2>Register</h2>
		<form action="/createAccount" method="post">
			<label for="username">Username:</label>
			<input type="text" id="username" name="username" required><br>
			<label for="password">Password:</label>
			<input type="password" id="password" name="password" required><br>
			<input type="submit" value="Register">
		</form>
		<p>Already have an account? <a href="/login">Login here</a>.</p>
	</body>
	</html>
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка создания нового аккаунта
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Проверка, что пользователь с таким именем не существует
	existingUser, err := getUser(username)
	if err == nil && existingUser != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Добавление нового пользователя в базу данных
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Account created successfully. You can now <a href=\"/login\">login</a>.")
}

func getUser(username string) (*User, error) {
	// Запрос пользователя из базы данных по имени пользователя
	row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
