# Suggesting Story Titles

Every good story needs a title and, with simplicity in mind, we would like you to suggest a
title or titles for a given group of photos.

For this task we have prepared 3 CSV files, each file contains a list of photo-metadata
and the file represents one “album” of photos. The metadata consists of:
● A timestamp (representing the date the photo was taken)
● Geo-coordinates (latitude and longitude)

Please design and implement a component that examines each file and suggests
appropriate titles or a title.

## Requirements

* Go 1.18
* Google Maps API key (used for reverse-geocoding requests)

## Makefile commands

### Linter:

```
make lint
```

### Testing:

```
make test
```

### Build:

```
make build
```

and then run:

```
.\suggestions -apiKey={GOOGLE_API_KEY} .\input\1.csv

or:

.\suggestions -apiKey={GOOGLE_API_KEY} .\input\1.csv .\input\2.csv

or:

.\suggestions -apiKey={GOOGLE_API_KEY} .\input\1.csv .\input\2.csv input\3.csv
```

## How to use

```
go run . -apiKey={GOOGLE_API_KEY} input\1.csv 

or: 

go run . -apiKey={GOOGLE_API_KEY} input\1.csv .\input\2.csv 

or:

go run . -apiKey={GOOGLE_API_KEY} input\1.csv .\input\2.csv input\3.csv
```

## Project description

This component examines each files received and suggests appropriate titles or a title, depending on the information contained in the file/s. The component can receive one or more files. At least one file should be provided.

A file of type CSV is expected.
Each file should contain a list of comma separated value containing photo metadata, such as:
- A timestamp (representing the date the photo was taken)
- Geo-coordinates (latitude and longitude)

Example Data:

Date Photo Was Taken | Latitude  | Longitude
---------------------|-----------|-----------
2019-03-30 14:12:19  | 40.703717 | -74.016094
2019-03-30 15:34:49  | 40.782222 | -73.965278
2019-03-31 12:18:04  | 40.748433 | -73.985656

Example Output
● “A weekend in New York”
● “A trip to New York”
● “A weekend in Manhattan”
● “A rainy trip to New York”
● “A trip to the United States”
● “New York in March”

## Assumptions

I've assumed the component can parse more than one file at once. As such, the output would contain the following:

```
Album 0 with 27 photos has the following title suggestions:
New York in March
New York in March 2020
March 2020 as seen in New York
New York, NY, USA in March
New York, NY, USA in March 2020
March 2020 as seen in New York, NY, USA
A cloudy trip to New York
A trip to New York
A brief trip to New York 👌
A quick trip to New York 👌
A fun trip to New York 😎
An amazing trip to New York 😎
A brief trip to the United States 👌
A short trip to the United States 👌
A quick trip to the United States 👌
An amazing trip to the United States 😎
A cloudy trip to the United States ☁️
```

Where album "0" reffers to the first entry received, which in this case would be input/1.csv.

For reverse-geocoding the location, I integrated with a third-party API: Google Maps.
https://developers.google.com/maps/documentation/geocoding/requests-reverse-geocoding

The API key is expected as the `apiKey` flag input. 
Ideally, when making the service production ready, this key would be kept in HashiCorp Vault (https://www.vaultproject.io/).

To reverse-geocode the location, I made requests like the below, retrieving just the `locality` result type:
https://maps.googleapis.com/maps/api/geocode/json?latlng=40.728808,-73.996106&result_type=locality&key={GOOGLE_API_KEY}

From the results otained, I looked at the `formatted_address` to identify at which location the picture was taken.
If the formatted address contained City, Country, Country (i.e.: `New York, NY, USA`), I would use just the city and the country to create title suggestions for each of these. 
The mistake I made was to only parse one final address, instead of looking at all possible addresses (see `title_suggestions.getFinalAddress()`), which limits the number of responses (e.g. I am missing `Manhattan` examples).

I did not have enough time to integrate with a weather third-party API, so that I could collect temperature stats and add those as suggestions as well. I would have opted for fetching weather for a day or a specified period of time using https://openweathermap.org/api/one-call-api
(e.g. GET https://api.openweathermap.org/data/2.5/onecall?lat=33.44&lon=-94.04&exclude=minutely,hourly,daily,alerts&appid={WEATHER_API_KEY}).
What I did instead was to randomise weather conditions.

I was uncertain whether I could use other external dependencies, but I wanted to be more creative with my title suggestions, so I added emojis to some of them using this library: `github.com/enescakir/emoji`.

If the photo album (aka the CSV file) contained 2 dates worth of photos, the suggestions would refer to a `weekend` as the period when the photos were taken, otherwise, if it would just be one `day`, it would state day, alternatively it would suggest `trip` or `adventure` (whether the number of days was an even number - trip, or an odd number - adventure).

## Future considerations

The following are points I have not had time to address:
- add more tests
I tested the application by adding more invalid input files, with which I ran the service to make sure some scenarios that I thought of as invalid for the input files were handled accordingly. I used this library `github.com/stretchr/testify` to assert results in my tests. 
- fix the addresses collected when reverse-geocoding and make use of all results => more title suggestions
- potentially add a title suggestions limit (too many results might just confuse the user)
- cleanup code structure (I had problems with the GOROOT path + I started hitting the Google Maps quota limit, which limitted me in testing/running the application). The project strcuture I would have gone for would have been the following:

```
- suggestions
	- main.go
	- README.go
	- Makefile
	- go.sum
	- go.mod
	- src
		- domain.go
		- reader
			file_reader.go
			file_reader_test.go
		- storage
			- in_memory.go
			- in_memory_test.go
		- geocoding
			- reverse_geocoding.go
			- reverse_geocoding_test.go
		- writer
			- title_suggestions.go
			- title_suggestions_test.go
```

The following are points I could address to make this component more production ready:
- cleanup/add more logs
- handle internationsalisation
- add more emojis into the title suggestions
- consider more meteo conditions
- add caching -> fetch data from cache when Google API unavailable
- potentially display events as suggestions for the periods the photo albums were taken
- add RESTful/gRPC API implementation for fetching title suggestions
- add rate limitting
- add HashiCorp key storage
- add metrics/monitoring
- Dockerize the solution
- connect to a DB to fetch the data for albums and/or to store the title suggestions per albums and various other information that we might need to use elsewhere (other services potentially)
- if we connect to a DB, we might want to consider adding exponential backoff with retries for recoverable errors
