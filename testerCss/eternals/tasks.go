package eternals

import (
	"net/http"
	"regexp"
	"strconv"
	"text/template"
	"time"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/taskday" {
	// 	http.Error(w, "Error task day", http.StatusUnauthorized)
	// }
	pattern := `/taskday/([1-7])\z`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(r.URL.Path)
	ID, err := strconv.Atoi(matches[1])
	if err != nil {
		http.Error(w, "Error task day", http.StatusUnauthorized)
		return
	}
	Task := GetTaskByID(ID)
	tmpl, err := template.ParseFiles("templates/nameday.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	TaskAndDay := struct {
		Time time.Weekday
		Task task
	}{
		Time: dayOfWeek,
		Task: Task,
	}
	tmpl.Execute(w, TaskAndDay)
}

func GetTaskForDay(day time.Weekday) task {
	switch day {
	case time.Monday:
		return task{1, day, "Volleyball", "Computer Games", "Reading"}
	case time.Tuesday:
		return task{2, day, "Football", "Board Games", "Crafting"}
	case time.Wednesday:
		return task{3, day, "Swimming", "Music", "Handicrafts"}
	case time.Thursday:
		return task{4, day, "Basketball", "English Courses", "Pottery"}
	case time.Friday:
		return task{5, day, "Stretching", "Self-defense Courses", "Blogging"}
	case time.Saturday:
		return task{6, day, "Pilates", "Movie Time", "Animal Therapy"}
	case time.Sunday:
		return task{7, day, "No specific activity", "Chess", "English Courses"}
	default:
		return task{0, day, "No specific activity", "", ""}
	}
}

func GetTaskByID(id int) task {
	tasks := []task{
		GetTaskForDay(time.Monday),
		GetTaskForDay(time.Tuesday),
		GetTaskForDay(time.Wednesday),
		GetTaskForDay(time.Thursday),
		GetTaskForDay(time.Friday),
		GetTaskForDay(time.Saturday),
		GetTaskForDay(time.Sunday),
	}
	for _, t := range tasks {
		if t.ID == id {
			return t
		}
	}
	return task{0, time.Monday, "No specific activity", "", ""}
}
