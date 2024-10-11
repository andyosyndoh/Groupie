package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"groupie/internals/renders"
	"groupie/utils"
)

type ArtistDetails struct {
	Art       utils.Artists       `json:"artist"`
	Locations utils.Location      `json:"locations"`
	Dates     utils.Date          `json:"dates"`
	Relations utils.ArtistDetails `json:"relations"`
}

var arts []utils.Artists

// HomeHandler handles the homepage route '/'
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		artists, err := utils.GetArtists()
		arts = artists
		if err != nil {
			ServerErrorHandler(w, r)
			log.Printf("Error retrieving artists: %v", err)
			return
		}
		renders.RenderTemplate(w, "home.page.html", artists)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// NotFoundHandler handles unknown routes; 404 status
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renders.RenderTemplate(w, "notfound.page.html", nil)
}

// BadRequestHandler handles bad requests routes
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	renders.RenderTemplate(w, "badrequest.page.html", nil)
}

// ServerErrorHandler handles server failures that result in status 500
func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	renders.RenderTemplate(w, "serverError.page.html", nil)
}

// Details handles the about page route '/details'
func Details(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		queryParams := r.URL.Query()
		idValue := queryParams.Get("id")
		ID, err := strconv.Atoi(idValue)
		if err != nil {
			BadRequestHandler(w, r)
			// http.Error(w, "Invalid ID", http.StatusBadRequest)
			log.Printf("Error converting id param to int value: %v", err)
			return
		}

		if ID <= 0 || ID > 52 {
			BadRequestHandler(w, r)
			// http.Error(w, "ID out of range", http.StatusBadRequest)
			log.Printf("ID out of range: %d", ID)
			return
		}

		location, err := utils.GetLocations(ID)
		artists, _ := utils.Getsingleartist(ID)
		dates, _ := utils.GetDates(ID)
		relations, _ := utils.GetRelation(ID)
		if err != nil {
			ServerErrorHandler(w, r)
			// http.Error(w, "Failed to retrieve location data", http.StatusInternalServerError)
			log.Printf("Error retrieving location data: %v", err)
			return
		}

		artistDetails := ArtistDetails{
			Art:       artists, // Wrap single artist in a slice
			Locations: location,
			Dates:     dates,
			Relations: relations,
		}

		renders.RenderTemplate(w, "details.page.html", artistDetails)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// SearchArtists handles search requests
func SearchArtists(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	if len(query) < 2 {
		json.NewEncoder(w).Encode([]string{})
		return
	}

	var suggestions []string

	for _, artist := range arts {
		// Search by artist/band name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			suggestions = append(suggestions, fmt.Sprintf("%d&%s - artist/band", artist.ID, artist.Name))
		}

		// Search by members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				suggestions = append(suggestions, fmt.Sprintf("%d&%s - member", artist.ID, member))
			}
		}

		// Search by first album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			suggestions = append(suggestions, fmt.Sprintf("%d&%s - first album: %s", artist.ID, artist.Name, artist.FirstAlbum))
		}

		// Search by creation date
		// creationYear := time.Unix(int64(artist.CreationDate), 0).Year()
		if strings.Contains(strings.ToLower(strconv.Itoa(artist.CreationDate)), query) {
			suggestions = append(suggestions, fmt.Sprintf("%d&%v - created: %d", artist.ID, artist.Name, artist.CreationDate))
		}
	}

	locations, err := utils.GetallLocations()
	if err != nil {
		ServerErrorHandler(w, r)
		// http.Error(w, "Failed to retrieve location data", http.StatusInternalServerError)
		log.Printf("Error retrieving location data: %v", err)
		return
	}

	for _, loc := range locations.Index {
		for _, v := range loc.Locations {
			if strings.Contains(v, query) {
				name := ""
				for _, artist := range arts {
					if artist.ID == loc.ID{
						name = artist.Name
					}
				}
				suggestions = append(suggestions, fmt.Sprintf("%v&%v - Performing at: %v", loc.ID,name , v))
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}
