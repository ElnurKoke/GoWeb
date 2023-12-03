package main

import (
	"fmt"
	"forum/eternals"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	eternals.Initdb()
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", eternals.IndexHandler)
	http.HandleFunc("/register", eternals.RegisterHandler)
	http.HandleFunc("/login", eternals.LoginHandler)
	http.HandleFunc("/taskday/", eternals.TaskHandler)
	http.HandleFunc("/logout", eternals.LogoutHandler)
	http.HandleFunc("/adminControl", eternals.ControlHandler)

	fmt.Println("Сервер запущен на порту http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
