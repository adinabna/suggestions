package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/enescakir/emoji"
)

type PeriodSuggestion int
type WeatherSuggestion int

const (
	Day PeriodSuggestion = iota
	Weekend
	Trip
	Adventure
)

const (
	Sunny WeatherSuggestion = iota
	Cloudy
	Foggy
	Windy
	Rainy
)

type titleSuggestionDetails struct {
	weather string
	period  PeriodSuggestion
	year    int
	month   string
	address string
}

func getAlbumTitleSuggestions(album photoAlbumDetails) {
	var tsd titleSuggestionDetails
	tsd.year = album.year

	// randomize weather for now
	// TODO: fetch weather for a day or a specified period of time using https://openweathermap.org/api/one-call-api
	// e.g. GET https://api.openweathermap.org/data/2.5/onecall?lat=33.44&lon=-94.04&exclude=minutely,hourly,daily,alerts&appid={WEATHER_API_KEY}
	tsd.weather = getWeather(WeatherSuggestion(rand.Intn(4)))

	tsd.period = getTripPeriod(len(album.days))
	tsd.month = getTripMonth(album.months)
	tsd.address = getFinalAddress(album.addresses)

	fmt.Println("Album", album.id, "with", album.photosCount, "photos has the following title suggestions:")
	generateTitleSuggestionDetails(tsd)
	fmt.Println("")
}

func generateTitleSuggestionDetails(tsd titleSuggestionDetails) {
	generateMonthSuggestions(tsd)
	preposition := getPreposition(tsd.period)
	article := getArticle(tsd.period)

	// if the address contains City, County, Country, generate suggestions based on the City and Country
	if strings.Contains(tsd.address, ",") {
		addressDetails := strings.Split(tsd.address, ",")

		city := strings.TrimSpace(addressDetails[0])
		fmt.Println("A", tsd.weather, tsd.period, preposition, city)
		getPeriodSuggestions(article, tsd, preposition, city)

		country := strings.TrimSpace(addressDetails[2])
		getCountrySuggestions(article, country, preposition, tsd)
	} else {
		// otherwise generate suggestions based on the address itself, which should contain just the country name
		getCountrySuggestions(article, tsd.address, preposition, tsd)
	}
}

func getCountrySuggestions(article string, country string, preposition string, tsd titleSuggestionDetails) {
	e := getWeatherEmoji(tsd.weather)

	switch country {
	case "USA", "US", "United States":
		getPeriodSuggestions(article, tsd, preposition, "the United States")
		fmt.Println("A", tsd.weather, tsd.period, preposition, "the United States", e)
	case "UK", "United Kingdom":
		getPeriodSuggestions(article, tsd, preposition, "the United Kingdom")
		fmt.Println("A", tsd.weather, tsd.period, preposition, "the United Kingdom", e)
	default:
		getPeriodSuggestions(article, tsd, preposition, country)
		fmt.Println("A", tsd.weather, tsd.period, preposition, country, e)
	}
}

func getWeatherEmoji(weather string) string {
	switch weather {
	case Sunny.String():
		return emoji.Parse(":sun:")
	case Cloudy.String():
		return emoji.Parse(":cloud:")
	case Foggy.String():
		return emoji.Parse(":fog:")
	default:
		return ""
	}

	return ""
}

func getPeriodSuggestions(article string, tsd titleSuggestionDetails, preposition string, cityOrCountry string) {
	fmt.Println(article, tsd.period, preposition, cityOrCountry)
	fmt.Println("A brief", tsd.period, preposition, cityOrCountry, emoji.OkHand)
	fmt.Println("A short", tsd.period, preposition, cityOrCountry, emoji.OkHand)
	fmt.Println("A quick", tsd.period, preposition, cityOrCountry, emoji.OkHand)
	fmt.Println("A fun", tsd.period, preposition, cityOrCountry, emoji.Parse(":sunglasses:"))
	fmt.Println("An amazing", tsd.period, preposition, cityOrCountry, emoji.Parse(":sunglasses:"))
}

func getArticle(period PeriodSuggestion) string {
	article := "A"
	if period == Adventure {
		article = "An"
	}

	return article
}

func generateMonthSuggestions(t titleSuggestionDetails) {
	if t.month != "" {
		if strings.Contains(t.address, ",") {
			addressDetails := strings.Split(t.address, ",")
			city := strings.TrimSpace(addressDetails[0])
			fmt.Println(city, "in", t.month)
			fmt.Println(city, "in", t.month, t.year)
			fmt.Println(t.month, t.year, "as seen in", city)
		}
		fmt.Println(t.address, "in", t.month)
		fmt.Println(t.address, "in", t.month, t.year)
		fmt.Println(t.month, t.year, "as seen in", t.address)
	}
}

func getPreposition(period PeriodSuggestion) string {
	preposition := "to"
	if period == Weekend || period == Adventure || period == Day {
		preposition = "in"
	}

	return preposition
}

// figure out the best final address out of all addresses for a photo album
func getFinalAddress(addresses map[string]int) string {
	var finalAddress string
	if len(addresses) == 1 {
		for k, _ := range addresses {
			finalAddress = k
		}
	} else if len(addresses) > 1 {
		for k, _ := range addresses {
			addressDetails := strings.Split(k, ",")
			finalAddress = addressDetails[2]
			break
		}
	}

	return strings.TrimSpace(finalAddress)
}

// figure out the month of the trip in the photo album
func getTripMonth(months map[string]int) string {
	// leave blank if the trip spans across more than one month
	month := ""
	if len(months) == 1 {
		for k, _ := range months {
			month = k
		}
	}

	return month
}

// figure out the duration of the trip in the photo album
func getTripPeriod(daysCounter int) PeriodSuggestion {
	var period PeriodSuggestion
	switch daysCounter {
	case 0:
		period = Day
	case 1:
		period = Weekend
	default:
		if daysCounter%2 == 0 {
			period = Trip
		} else {
			period = Adventure
		}
	}

	return period
}

func (ps PeriodSuggestion) String() string {
	switch ps {
	case Day:
		return "day"
	case Weekend:
		return "weekend"
	case Trip:
		return "trip"
	case Adventure:
		return "adventure"
	}
	return "unknown"
}

func (ws WeatherSuggestion) String() string {
	switch ws {
	case Sunny:
		return "sunny"
	case Cloudy:
		return "cloudy"
	case Foggy:
		return "foggy"
	case Windy:
		return "windy"
	case Rainy:
		return "rainy"
	}
	return "unknown"
}

func getWeather(ws WeatherSuggestion) string {
	return fmt.Sprintf("%v", ws)
}
