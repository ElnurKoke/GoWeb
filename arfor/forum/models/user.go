package models

import "time"

type User struct {
	Id            int
	Password      string
	Login         string
	Username      string
	Token         string
	TokenDuration time.Time
}
