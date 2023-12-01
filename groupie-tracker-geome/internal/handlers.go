package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandlerError(w, http.StatusMethodNotAllowed)
		return
	}

	artists, err := ParseJson(w, r)
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	if r.URL.Path != "/" {
		HandlerError(w, 404)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		HandlerError(w, 500)
		return
	}
}

func getArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandlerError(w, http.StatusMethodNotAllowed)
		return
	}

	artisId := r.URL.Path[len("/artists/"):]
	if artisId == "" || artisId[0] == '0' {
		HandlerError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(artisId)
	if err != nil {
		HandlerError(w, http.StatusBadRequest)
		return
	}

	artist, locADate, err := GetData(id)
	if err != nil {
		HandlerError(w, http.StatusBadRequest)
		return
	}

	location, err := GetLocationById(id)

	coor := GetCoordinatesBatch(location)
	if err != nil {
		HandlerError(w, http.StatusBadRequest)
		return
	}

	tmp, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		HandlerError(w, http.StatusInternalServerError)
		return
	}
	// fmt.Println(coor)
	ans := map[string]interface{}{
		"Artist":   artist,
		"LocADate": locADate,
		"Location": coor,
	}

	err = tmp.Execute(w, ans)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func HandlerError(w http.ResponseWriter, code int) {
	pageError := struct {
		ErrorNum     int
		ErrorMessage string
	}{
		ErrorNum:     code,
		ErrorMessage: http.StatusText(code),
	}
	// w.WriteHeader(san)
	temp, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	temp.Execute(w, pageError)
}

func RenderSearch(w http.ResponseWriter, artists Artists) {
	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	err := tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleRequest() {
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/static/", ServeCSSFiles(http.StripPrefix("/static", fs)))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", Home)
	http.HandleFunc("/artists/", getArtist)
	http.HandleFunc("/search", SearchHandler)
	http.HandleFunc("/autocomplete", AutocompleteHandler)
	http.HandleFunc("/filter", FilterHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
