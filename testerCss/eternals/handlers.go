package eternals

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	id        int          = 0
	dayOfWeek time.Weekday = time.Now().Weekday()
	answer    string
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		clientWords := r.FormValue("chat-input")
		if Contain(clientWords, "плохо") {
			answer = "Почему вам плохо?"
		}

	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := id

	row := db.QueryRow("SELECT ID, Password, Login, Username, Point FROM user WHERE ID = ?", userID)
	var Userinfo user
	err = row.Scan(&Userinfo.ID, &Userinfo.Password, &Userinfo.Login, &Userinfo.Username, &Userinfo.Point)
	if err != nil {
		fmt.Println("Error scanning row:", err)
	}

	UserinfoWithTime := struct {
		User   user
		Time   time.Weekday
		Answer string
	}{
		User:   Userinfo,
		Time:   dayOfWeek,
		Answer: answer,
	}
	tmpl.Execute(w, UserinfoWithTime)
}
