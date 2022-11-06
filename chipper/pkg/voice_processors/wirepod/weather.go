package wirepod

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/digital-dream-labs/chipper/pkg/logger"
)

type weatherAPIResponseStruct struct {
	Location struct {
		Name      string `json:"name"`
		Localtime string `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
	} `json:"current"`
}
type weatherAPICladStruct []struct {
	APIValue string `json:"APIValue"`
	CladType string `json:"CladType"`
}

type openWeatherMapAPIGeoCodingStruct struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state"`
}

type WeatherStruct struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type openWeatherMapAPIResponseStruct struct {
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Weather []WeatherStruct `json:"weather"`
	Base    string          `json:"base"`
	Main    struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	DT  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		Id      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func getWeather(location string, botUnits string) (string, string, string, string, string, string) {
	var weatherEnabled bool
	var condition string
	var is_forecast string
	var local_datetime string
	var speakable_location_string string
	var temperature string
	var temperature_unit string
	weatherAPIEnabled := os.Getenv("WEATHERAPI_ENABLED")
	weatherAPIKey := os.Getenv("WEATHERAPI_KEY")
	weatherAPIUnit := os.Getenv("WEATHERAPI_UNIT")
	if weatherAPIEnabled == "true" && weatherAPIKey != "" {
		weatherEnabled = true
		logger.Logger("Weather API Enabled")
	} else {
		weatherEnabled = false
		logger.Logger("Weather API not enabled, using placeholder")
		if weatherAPIEnabled == "true" && weatherAPIKey == "" {
			logger.Logger("Weather API enabled, but Weather API key not set")
		}
	}
	if weatherEnabled {
		if botUnits != "" {
			if botUnits == "F" {
				logger.Logger("Weather units set to F")
				weatherAPIUnit = "F"
			} else if botUnits == "C" {
				logger.Logger("Weather units set to C")
				weatherAPIUnit = "C"
			}
		} else if weatherAPIUnit != "F" && weatherAPIUnit != "C" {
			logger.Logger("Weather API unit not set, using F")
			weatherAPIUnit = "F"
		}
	}

	if weatherEnabled {
		// First use geocoding api to convert location into coordinates
		// E.G. http://api.openweathermap.org/geo/1.0/direct?q={city name},{state code},{country code}&limit={limit}&appid={API key}
		url := "http://api.openweathermap.org/geo/1.0/direct?q=" + location + "&limit=1&appid=" + weatherAPIKey
		resp, err := http.Get(url)
		if err != nil {
			logger.Logger(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		geoCodingResponse := string(body)
		logger.Logger(geoCodingResponse)

		var geoCodingInfoStruct []openWeatherMapAPIGeoCodingStruct

		err = json.Unmarshal([]byte(geoCodingResponse), &geoCodingInfoStruct)
		if err != nil {
			panic(err)
		}
		if len(geoCodingInfoStruct) < 0 || len(geoCodingInfoStruct) == 0 {
			condition = "undefined"
			is_forecast = "false"
			local_datetime = "test"              // preferably local time in UTC ISO 8601 format ("2022-06-15 12:21:22.123")
			speakable_location_string = location // preferably the processed location
			temperature = "120"
			temperature_unit = "C"
			return condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit
		}
		Lat := fmt.Sprintf("%f", geoCodingInfoStruct[0].Lat)
		Lon := fmt.Sprintf("%f", geoCodingInfoStruct[0].Lon)

		logger.Logger("Lat: " + Lat + ", Lon: " + Lon)
		logger.Logger("Name: " + geoCodingInfoStruct[0].Name)
		logger.Logger("Country: " + geoCodingInfoStruct[0].Country)

		// Now that we have Lat and Lon, let's query the weather
		units := "metric"
		if weatherAPIUnit == "F" {
			units = "imperial"
		}
		url = "https://api.openweathermap.org/data/2.5/weather?lat=" + Lat + "&lon=" + Lon + "&units=" + units + "&appid=" + weatherAPIKey
		resp, err = http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ = io.ReadAll(resp.Body)
		weatherResponse := string(body)

		logger.Logger(weatherResponse)

		var openWeatherMapAPIResponse openWeatherMapAPIResponseStruct

		err = json.Unmarshal([]byte(weatherResponse), &openWeatherMapAPIResponse)
		if err != nil {
			panic(err)
		}

		conditionCode := openWeatherMapAPIResponse.Weather[0].Id

		logger.Logger(weatherResponse)
		logger.Logger(conditionCode)

		if conditionCode < 300 {
			// Thunderstorm
			condition = "Thunderstorms"
		} else if conditionCode < 400 {
			// Drizzle
			condition = "Rain"
		} else if conditionCode < 600 {
			// Rain
			condition = "Rain"
		} else if conditionCode < 700 {
			// Snow
			condition = "Snow"
		} else if conditionCode < 800 {
			// Athmosphere
			if openWeatherMapAPIResponse.Weather[0].Main == "Mist" ||
				openWeatherMapAPIResponse.Weather[0].Main == "Fog" {
				condition = "Rain"
			} else {
				condition = "Windy"
			}
		} else if conditionCode == 800 {
			// Clear
			if openWeatherMapAPIResponse.DT < openWeatherMapAPIResponse.Sys.Sunset {
				condition = "Sunny"
			} else {
				condition = "Stars"
			}
		} else if conditionCode < 900 {
			// Cloud
			condition = "Cloudy"
		} else {
			condition = openWeatherMapAPIResponse.Weather[0].Main
		}

		is_forecast = "false"
		t := time.Unix(int64(openWeatherMapAPIResponse.DT), 0)
		local_datetime = t.Format(time.RFC850)
		logger.Logger(local_datetime)
		speakable_location_string = openWeatherMapAPIResponse.Name
		temperature = fmt.Sprintf("%f", math.Round(openWeatherMapAPIResponse.Main.Temp))
		if weatherAPIUnit == "C" {
			temperature_unit = "C"
		} else {
			temperature_unit = "F"
		}
	} else {
		condition = "Snow"
		is_forecast = "false"
		local_datetime = "test"              // preferably local time in UTC ISO 8601 format ("2022-06-15 12:21:22.123")
		speakable_location_string = location // preferably the processed location
		temperature = "120"
		temperature_unit = "C"
	}
	return condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit
}

func weatherParser(speechText string, botLocation string, botUnits string) (string, string, string, string, string, string) {
	var specificLocation bool
	var apiLocation string
	var speechLocation string
	if strings.Contains(speechText, " in ") {
		splitPhrase := strings.SplitAfter(speechText, " in ")
		speechLocation = strings.TrimSpace(splitPhrase[1])
		if len(splitPhrase) == 3 {
			speechLocation = speechLocation + " " + strings.TrimSpace(splitPhrase[2])
		} else if len(splitPhrase) == 4 {
			speechLocation = speechLocation + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
		} else if len(splitPhrase) > 4 {
			speechLocation = speechLocation + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
		}
		splitLocation := strings.Split(speechLocation, " ")
		if len(splitLocation) == 2 {
			speechLocation = splitLocation[0] + ", " + splitLocation[1]
		} else if len(splitLocation) == 3 {
			speechLocation = splitLocation[0] + " " + splitLocation[1] + ", " + splitLocation[2]
		}
		logger.Logger("Location parsed from speech: " + "`" + speechLocation + "`")
		specificLocation = true
	} else {
		logger.Logger("No location parsed from speech")
		specificLocation = false
	}
	if specificLocation {
		apiLocation = speechLocation
	} else {
		apiLocation = botLocation
	}
	// call to weather API
	condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := getWeather(apiLocation, botUnits)
	return condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit
}
