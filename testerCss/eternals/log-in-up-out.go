package eternals

import (
	"net/http"
	"text/template"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")

		// Проверяем существование пользователя и правильность пароля
		var user user
		row := db.QueryRow("SELECT id, password FROM user WHERE login=?", login)
		err := row.Scan(&user.ID, &user.Password)
		if err != nil {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		if user.Password != password {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}
		id = user.ID
		// fmt.Fprintln(w, "Вход выполнен успешно!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "Error logout", http.StatusUnauthorized)
	}
	id = 0
	// fmt.Fprintln(w, "Выход выполнен успешно!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")
		username := r.FormValue("username")
		point := 0

		// Вставляем данные пользователя в таблицу
		_, err := db.Exec("INSERT INTO user (login, password, username,point) VALUES (?, ?, ?,?)", login, password, username, point)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fmt.Fprintln(w, "Регистрация успешно завершена!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
