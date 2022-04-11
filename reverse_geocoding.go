package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Results struct {
	Results []Result `json:"results"`
	Status  string   `json:"status"`
}

type Result struct {
	AddressComponents []Address `json:"address_components"`
	FormattedAddress  string    `json:"formatted_address"`
	Geometry          Geometry  `json:"geometry"`
	PlaceId           string    `json:"place_id"`
	Types             []string  `json:"types"`
}

type Address struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Bounds       Bounds `json:"bounds"`
	Location     LatLng `json:"location"`
	LocationType string `json:"location_type"`
	Viewport     Bounds `json:"viewport"`
}

type Bounds struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// reverse geocode the photo location
func getPhotoLocation(lat string, lon string, apiKey string) Results {
	url := "https://maps.googleapis.com/maps/api/geocode/json?latlng=" + lat + "," + lon + "&result_type=locality&key=" + apiKey

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return Results{}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return Results{}
	}
	defer resp.Body.Close()

	var results Results
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&results)

	return results
}
