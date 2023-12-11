package handler

import (
	"fmt"
	"foruum/models"
	"html/template"
	"log"
	"net/http"
)

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxUserKey).(models.User)
	if user.Token != "" {
		ErrorHandler(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	temp, err := template.ParseFiles("ui/signin.html")
	if err != nil {
		log.Fatal(err)
	}
	switch r.Method {
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("psw")
		user, err := h.service.GenerateToken(email, password)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:  "session_token",
				Value: user.Token,
				Path:  "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
	temp.Execute(w, temp)
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxUserKey).(models.User)
	if user.Token != "" {
		ErrorHandler(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	temp, err := template.ParseFiles("ui/signup.html")
	if err != nil {
		log.Fatal(err)
	}
	type Alert struct {
		Mess string
	}
	alert := Alert{
		Mess: "",
	}
	switch r.Method {
	case http.MethodPost:
		var signUp models.User
		email := r.FormValue("email")
		password := r.FormValue("psw")
		username := r.FormValue("username")
		repeatpass := r.FormValue("repeatspw")

		if repeatpass != password {
			ErrorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else {
			signUp = models.User{
				Username: username,
				Login:    email,
				Password: password,
			}
		}
		if err := h.service.CreateUser(signUp); err != nil {
			alert.Mess = err.Error()
			fmt.Println(alert)
			temp.Execute(w, alert)
			return

		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	temp.Execute(w, alert)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: "",
		Path:  "/",
	})
	/*user := r.Context().Value(ctxUserKey).(models.User)
	if err := h.service.DeleteToken(user.Token); err != nil {
		log.Fatal("delete token in logout", err)
	}*/
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
