package service

import (
	"net/mail"
)

func CheckUsername(username string) bool {
	for _, w := range username {
		if w >= 33 && w <= 126 {
			return true
		} else {
			return false
		}
	}
	return true
}

func CheckPassword(password string) bool {

	for _, rune := range password {
		if rune > 32 && rune <= 126 && len(password) > 8 {
			return true
		} else {
			return false
		}
	}

	return true
}

func CheckLogin(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
