package eternals

import (
	"strings"
	"time"
)

type user struct {
	ID       int
	Password string
	Login    string
	Username string
	Point    int
}

type task struct {
	ID      int
	WeekDay time.Weekday
	First   string
	Second  string
	Third   string
}

func Contain(str string, flag string) bool {
	return strings.Contains(str, flag)
}
