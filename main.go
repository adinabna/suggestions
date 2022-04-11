package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

type photoAlbumDetails struct {
	id          int
	photosCount int
	year        int
	days        map[int]int
	months      map[string]int
	addresses   map[string]int
}

func main() {
	var apiKey string
	flag.StringVar(&apiKey, "apiKey", "", "Google Maps API key")
	flag.Parse()

	if apiKey == "" {
		log.Fatalf("Please provide an API key for Google Maps API.")
	}

	if len(os.Args) < 3 {
		log.Fatalf("Please provide at least one input file.")
	}

	albums := createPhotoAlbums(os.Args[2:])

	var wg sync.WaitGroup
	pad := make([]photoAlbumDetails, len(albums))
	for i, album := range albums {
		var pa photoAlbumDetails
		pa.id = i

		addresses := make(map[string]int)
		days := make(map[int]int)
		months := make(map[string]int)

		photosCount := len(album.photos)
		pa.photosCount = photosCount
		wg.Add(photosCount)
		for _, photo := range album.photos {
			go func(photo photoMetadata) {
				defer wg.Done()
				lat := fmt.Sprintf("%v", photo.geoLocation.latitude)
				lon := fmt.Sprintf("%v", photo.geoLocation.longitude)
				response := getPhotoLocation(lat, lon, apiKey)

				switch response.Status {
				case "OK":
					// identify the year
					pa.year = photo.date.Year()

					// identify days
					days[photo.date.Day()] += 1
					pa.days = days

					//identify months
					months[photo.date.Month().String()] += 1
					pa.months = months

					// identify addresses
					addresses[response.Results[0].FormattedAddress] += 1
					pa.addresses = addresses

					pad[i] = pa
				case "ZERO_RESULTS":
					log.Println("We could not find any results for the photo metadata provided...")
					break
				case "OVER_QUERY_LIMIT":
					log.Fatal("Too many requests!")
				case "REQUEST_DENIED":
					log.Fatal("Please provide the API key.")
				case "INVALID_REQUEST":
					log.Fatal("Invalid request...")
				default:
					log.Fatal("Unknow error encountered when requesting photo metadata...")
				}
			}(photo)
		}
		wg.Wait()
	}

	for _, album := range pad {
		getAlbumTitleSuggestions(album)
	}
}
