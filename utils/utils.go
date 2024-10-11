package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Artists struct {
	ID              int      `json:"id"`
	ImageURL        string   `json:"image"`
	Name            string   `json:"name"`
	Members         []string `json:"members"`
	CreationDate    int      `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	LocationsURL    string   `json:"locations"`
	ConcertDatesURL string   `json:"concertDates"`
	RelationsURL    string   `json:"relations"`
}

type Artist struct {
	Index []Artists
}

type Locations struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Index []Date `json:"index"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Index []ArtistDetails `json:"index"`
}

type ArtistDetails struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var CacheDataMap map[string]json.RawMessage = make(map[string]json.RawMessage)

const api = "https://groupietrackers.herokuapp.com/api"

// GetAndUnmarshalArtists returns a list of artists by fetching or using cached data
func GetArtists() ([]Artists, error) {
	artists := []Artists{}
	err := unmarshalData("/artists", &artists)
	return artists, err
}

func Getsingleartist(ID int) (Artists, error) {
	artists := []Artists{}
	err := unmarshalData("/artists", &artists)
	if err != nil {
		return Artists{}, err
	}

	for _, v := range artists {
		if v.ID == ID {
			return v, nil
		}
	}
	return Artists{}, fmt.Errorf("location with ID %d not found", ID)
}

func GetallLocations() (Locations, error){
	locations := Locations{}
	err := unmarshalData("/locations", &locations)
	if err != nil {
		return Locations{}, err
	}
	return locations, nil
}

func GetLocations(ID int) (Location, error) {
	locations := Locations{}
	err := unmarshalData("/locations", &locations)
	if err != nil {
		return Location{}, err
	}

	for _, v := range locations.Index {
		if v.ID == ID {
			return v, nil
		}
	}
	return Location{}, fmt.Errorf("location with ID %d not found", ID)
}

func GetDates(ID int) (Date, error) {
	dates := Dates{}
	err := unmarshalData("/dates", &dates)
	if err != nil {
		return Date{}, err
	}

	for _, v := range dates.Index {
		if v.ID == ID {
			return v, nil
		}
	}
	return Date{}, fmt.Errorf("date with ID %d not found", ID)
}

func GetallDates() (Dates, error){
	dates := Dates{}
	err := unmarshalData("/dates", &dates)
	if err != nil {
		return Dates{}, err
	}
	return dates , nil
}

func GetRelation(ID int) (ArtistDetails, error) {
	relation := Relation{}
	err := unmarshalData("/relation", &relation)
	if err != nil {
		return ArtistDetails{}, err
	}

	for _, v := range relation.Index {
		if v.ID == ID {
			return v, nil
		}
	}
	return ArtistDetails{}, fmt.Errorf("relation with ID %d not found", ID)
}

func unmarshalData(endpoint string, out interface{}) error {
	if data, ok := CacheDataMap[endpoint]; ok {
		return json.Unmarshal(data, out)
	} else {
		jsonData, err := getJSONData(endpoint)
		if err != nil {
			return err
		}

		CacheDataMap[endpoint] = jsonData

		return json.Unmarshal(jsonData, out)
	}
}

func getJSONData(endpoint string) (json.RawMessage, error) {
	resp, err := http.Get(api + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s json data: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received a non-200 response code for %s: %d", endpoint, resp.StatusCode)
	}

	var jsonString json.RawMessage
	err = json.NewDecoder(resp.Body).Decode(&jsonString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s json data: %w", endpoint, err)
	}

	return jsonString, nil
}

// errors is a map of error output value in ErrorHandler
var errors = map[string]string{
	"web":        "😮 Oops! Something went wrong",
	"restricted": "😣 Oops! this is a restricted path.\nplease use another path.",
}

// ErrorHandler outputs errors and safely exits the program
func ErrorHandler(errType string) {
	fmt.Println(errors[errType])
	os.Exit(0)
}
