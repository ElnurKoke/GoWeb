package main

import (
	"fmt"
	"forum/eternals"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	eternals.Initdb()
	http.HandleFunc("/", eternals.IndexHandler)
	http.HandleFunc("/register", eternals.RegisterHandler)
	http.HandleFunc("/login", eternals.LoginHandler)
	fmt.Println("Сервер запущен на порту http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
