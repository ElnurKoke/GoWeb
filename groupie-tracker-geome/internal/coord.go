package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const mapboxAPIURL = "https://maps.googleapis.com/maps/api/geocode/"

// https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key=YOUR_API_KEY
type Locations struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Geometry struct {
	Locations Locations `json:"location"`
}

type Result struct {
	Geometry Geometry `json:"geometry"`
}

type GeocodeResponse struct {
	Results []Result `json:"results"`
}
type Index1 struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Artistt struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
}
type DateLocation struct {
	Date string `json:"date"`
}

type DatesLocations map[string][]string

type IndexItem struct {
	ID             int            `json:"id"`
	DatesLocations DatesLocations `json:"datesLocations"`
}

type Index struct {
	Items []IndexItem `json:"index"`
}

func GetLocationsAndDates() Index {
	jsonData, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatal(err)
	}
	var index Index
	err = json.NewDecoder(jsonData.Body).Decode(&index)

	if err != nil {
		log.Fatal(err)
	}

	return index
}

func GetArtists() []Artistt {
	data, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}

	var artists []Artistt

	err = json.NewDecoder(data.Body).Decode(&artists)

	if err != nil {
		log.Fatal(err)
	}
	return artists
}

func GetData(id int) (Artistt, IndexItem, error) {
	var empty1 Artistt
	var empty2 IndexItem
	err := errors.New("Wrong number")
	artists := GetArtists()
	locADate := GetLocationsAndDates()

	for _, ch := range artists {
		for _, el := range locADate.Items {
			if ch.Id == id && el.ID == id {
				return ch, el, nil
			}
		}
	}

	return empty1, empty2, err
}

func GetLocations() Index1 {
	data, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Fatal(err)
	}
	var index Index1

	err = json.NewDecoder(data.Body).Decode(&index)

	if err != nil {
		log.Fatal(err)
	}
	return index
}

func GetLocationById(id int) ([]string, error) {
	err := errors.New("Wrong id")
	locations := GetLocations()
	for _, ch := range locations.Index {
		if ch.ID == id {
			return ch.Locations, nil
		}
	}
	return nil, err
}

func GetCoordinates(location, apiKey string) []float64 {
	result := make([]float64, 2)

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", location, apiKey)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return result
	}
	defer response.Body.Close()

	var geocodeResponse GeocodeResponse
	err = json.NewDecoder(response.Body).Decode(&geocodeResponse)
	if err != nil {
		log.Fatal(err)
	}

	result[1] = geocodeResponse.Results[0].Geometry.Locations.Lat
	result[0] = geocodeResponse.Results[0].Geometry.Locations.Lng

	return result
}

func GetCoordinatesBatch(locations []string) [][]float64 {
	apiKey := "AIzaSyBXH2PYMwKrL18rjTE-O5OdtEgUZywoIgo"
	coordinateCh := make(chan []float64, len(locations))
	var wg sync.WaitGroup

	for _, location := range locations {
		wg.Add(1)
		go func(loc string) {
			defer wg.Done()
			coordinates := GetCoordinates(loc, apiKey)
			coordinateCh <- coordinates
		}(location)
	}

	go func() {
		wg.Wait()
		close(coordinateCh)
	}()

	var coor [][]float64
	for coords := range coordinateCh {
		coor = append(coor, coords)
	}

	return coor
}
