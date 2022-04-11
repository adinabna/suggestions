package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type location struct {
	latitude  float64
	longitude float64
}

type photoMetadata struct {
	date        time.Time
	geoLocation location
}

type photoAlbum struct {
	id     int
	photos []photoMetadata
}

func createPhotoAlbums(filenames []string) []photoAlbum {
	files := make(map[int]fileContent, len(filenames))
	for i, filename := range filenames {
		data := readFile(filename)
		files[i] = data
	}

	albums := make([]photoAlbum, len(filenames))
	for i, file := range files {
		album := createPhotoAlbum(i, file)
		albums[i] = album
	}

	return albums
}

func createPhotoAlbum(id int, file fileContent) photoAlbum {
	var album photoAlbum
	album.id = id
	for _, line := range file {
		if len(line) != 3 {
			log.Fatalf("Invalid input file!")
		}

		lat, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			log.Fatalf("Invalid latitude: %s", err)
		}
		if lat < -90 || lat > 90 {
			log.Fatalf("Invalid latitude. A latitude should be a value between -90 and 90.")
		}

		lon, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			log.Fatalf("Invalid longitude: %s", err)
		}
		if lon < -180 || lon > 180 {
			log.Fatalf("Invalid longitute. A longitute should be a value between -180 and 180.")
		}

		var record photoMetadata
		record.date = parseDate(line[0])
		record.geoLocation.latitude = lat
		record.geoLocation.longitude = lon
		album.photos = append(album.photos, record)
	}

	return album
}

func parseDate(date string) time.Time {
	// parse dates equivalent to this date format: "2020-03-30T14:12:19Z"
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		// fix date format "2020-03-30 14:12:19" to be compliant to time.RFC3339,
		// equivalent to "2020-03-30T14:12:19Z"
		re := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
		if re.MatchString(date) {
			dateString := strings.Replace(date, " ", "T", 1) + "Z"
			t, err = time.Parse(time.RFC3339, dateString)
			if err != nil {
				log.Fatalf("Invalid date: %s", err)
			}
		} else {
			log.Fatalf("Invalid date: %s", err)
		}
	}

	return t
}
