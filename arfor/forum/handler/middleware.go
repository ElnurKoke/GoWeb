package handler

import (
	"context"
	"errors"
	"fmt"
	"foruum/models"
	"net/http"
	"time"
)

type ctxKey int8

const ctxUserKey ctxKey = iota

func (h *Handler) MiddleWare(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		t, err := r.Cookie("session_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxUserKey, models.User{})))
			case errors.Is(err, t.Valid()):
				fmt.Println(w, http.StatusBadRequest, "invalid cookie value")
			}
			fmt.Println(w, http.StatusBadRequest, "failed to get cookie")
			return
		}
		user, err = h.service.Authorization.GetUserByToken(t.Value)
		if err != nil {
			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxUserKey, models.User{})))
			return
		}
		if user.TokenDuration.Before(time.Now()) {
			if err := h.service.DeleteToken(user.Token); err != nil {
				fmt.Println(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxUserKey, user)))
	}
}
