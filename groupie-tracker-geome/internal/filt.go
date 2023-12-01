package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		creationFrom := r.FormValue("creation-date-from")
		creationTo := r.FormValue("creation-date-to")

		albumFrom := r.FormValue("first-album-date-from")
		albumTo := r.FormValue("first-album-date-to")

		membersFrom := r.FormValue("members-from")
		membersTo := r.FormValue("members-to")

		pickedlocations := r.Form["pickedlocations[]"]

		membersFromInt, err := strconv.Atoi(membersFrom)
		if err != nil {
			membersFromInt = 0
		}
		membersToInt, err := strconv.Atoi(membersTo)
		if err != nil {
			membersToInt = 0
		}

		creationFromInt, err := strconv.Atoi(creationFrom)
		if err != nil {
			creationFromInt = 0
		}
		creationToInt, err := strconv.Atoi(creationTo)
		if err != nil {
			creationToInt = 0
		}

		if creationFrom == "" {
			creationFrom = "default"
			creationFromInt = 1960
		}
		if creationTo == "" {
			creationTo = "default"
			creationToInt = 2020
		}
		if membersFrom == "" {
			membersFrom = "default"
			membersFromInt = 1
		}
		if membersTo == "" {
			membersTo = "default"
			membersToInt = 10
		}

		if albumFrom == "" {
			albumFrom = "1950-01-02"
		}
		if albumTo == "" {
			albumTo = "2050-01-02"
		}

		var filteredArtists Artists
		FullArtists, err := ParseJson(w, r)
		if err != nil {
			HandlerError(w, 400)
		}

		for _, artist := range FullArtists {
			addartist := true
			if membersFrom != "" && membersTo != "" {
				if len(artist.Members) >= membersFromInt && len(artist.Members) <= membersToInt {
				} else {
					addartist = false
				}
			}
			if creationFrom != "" && creationTo != "" {
				if artist.CreationDate >= creationFromInt && artist.CreationDate <= creationToInt {
				} else {
					addartist = false
				}
			}
			if albumFrom != "" && albumTo != "" {
				if CheckByDate(albumFrom, albumTo, artist.FirstAlbum) {
				} else {
					addartist = false
				}
			}
			if len(pickedlocations) != 0 {
				if CheckByCountries(pickedlocations, artist.Relation.Relations) {
				} else {
					addartist = false
				}
			}
			if addartist {
				filteredArtists = append(filteredArtists, artist)
			}
		}

		RenderFilter(w, filteredArtists)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This function does not support " + r.Method + " method."))
	}
}

func CheckByCountries(pickedcountries []string, artistcountries map[string][]string) bool {
	truenumber := 0
	for _, pickedcountry := range pickedcountries {
		for artistcountry, _ := range artistcountries {
			if strings.HasSuffix(artistcountry, pickedcountry) {
				truenumber++
				break
			}
		}
	}
	if len(pickedcountries) == truenumber {
		return true
	}
	return false
}

func CheckByDate(albumFrom, albumTo string, artistAlbomDate string) bool {
	dateFormat := "2006-01-02"
	ArtistdateFormat := "02-01-2006"

	albumFromDate, err := time.Parse(dateFormat, albumFrom)
	if err != nil {
		fmt.Println("Ошибка при парсинге даты:", err)
		return false
	}
	albumToDate, err := time.Parse(dateFormat, albumTo)
	if err != nil {
		fmt.Println("Ошибка при парсинге даты:", err)
		return false
	}
	artistDate, err := time.Parse(ArtistdateFormat, artistAlbomDate)
	if err != nil {
		fmt.Println("Ошибка при парсинге даты:", err)
		return false
	}
	return artistDate.After(albumFromDate) && artistDate.Before(albumToDate)
}

func FSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		search := r.URL.Query().Get("search")
		search = RemoveSubstringAndTrim(search, "- artist/band")
		search = RemoveSubstringAndTrim(search, "- member")
		search = RemoveSubstringAndTrim(RemoveParenthesesAndContents(search), "- location")
		search = RemoveSubstringAndTrim(search, "- first album date")
		search = RemoveSubstringAndTrim(search, "- creation date")

		if search != "" {
			filteredArtists := SearchArtist(w, r, search)
			RenderSearch(w, filteredArtists)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FSearchArtist(w http.ResponseWriter, r *http.Request, searchname string) Artists {
	var results Artists
	artists, _ := ParseJson(w, r)
	if len(searchname) == 1 && unicode.IsLetter([]rune(searchname)[0]) {
		for _, artist := range artists {
			if ContainsWordWithStartingLetter(strings.ToLower(searchname), strings.ToLower(artist.Name)) {
				results = append(results, artist)
			}
		}
		return results
	}
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchname)) {
			results = append(results, artist)
		} else if strings.Contains(strconv.Itoa(artist.CreationDate), searchname) || strings.Contains(artist.FirstAlbum, searchname) {
			results = append(results, artist)
		} else {
			ok := false
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(searchname)) {
					results = append(results, artist)
					ok = true
					break
				}
			}
			if !ok {
				for locations, _ := range artist.Relation.Relations {
					if strings.Contains(strings.ToLower(string(locations)), strings.ToLower(searchname)) {
						results = append(results, artist)
						break
					}
				}
			}
		}
	}
	return results
}

func RenderFilter(w http.ResponseWriter, artists Artists) {
	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	err := tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
