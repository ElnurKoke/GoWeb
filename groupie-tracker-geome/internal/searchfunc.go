package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func ContainsWordWithStartingLetter(letter, text string) bool {
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(strings.ToLower(word), strings.ToLower(letter)) {
			return true
		}
	}
	return false
}

func RemoveSubstringAndTrim(word, delword string) string {
	result := strings.Replace(word, delword, "", -1)
	result = strings.TrimSpace(result)
	return result
}

func RemoveParenthesesAndContents(word string) string {
	startIndex := strings.Index(word, "(")
	if startIndex == -1 {
		return word
	}

	endIndex := strings.Index(word, ")")
	if endIndex == -1 || endIndex <= startIndex {
		return word
	}

	return strings.TrimSpace(word[:startIndex] + word[endIndex+1:])
}

func AutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	api, _ := ParseJson(w, r)

	var suggestions []string
	for _, artist := range api {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchQuery)) {
			suggestions = append(suggestions, artist.Name+" - artist/band")
		} else {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(searchQuery)) {
					suggestions = append(suggestions, member+" - member")
					break
				}
			}
		}
		if strings.Contains(strconv.Itoa(artist.CreationDate), searchQuery) {
			suggestions = append(suggestions, strconv.Itoa(artist.CreationDate)+" - creation date"+"("+artist.Name+")")
		}
		if strings.Contains(artist.FirstAlbum, searchQuery) {
			suggestions = append(suggestions, artist.FirstAlbum+" - first album date"+"("+artist.Name+")")
		}
		for locations, _ := range artist.Relation.Relations {
			if strings.Contains(strings.ToLower(string(locations)), strings.ToLower(searchQuery)) {
				suggestions = append(suggestions, string(locations)+" - location"+"("+artist.Name+")")
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
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

func SearchArtist(w http.ResponseWriter, r *http.Request, searchname string) Artists {
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
