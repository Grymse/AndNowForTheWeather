package main

import (
	"strconv"
	"github.com/relvacode/iso8601"
	"time"
)

type WeatherDataPoint struct {
	Time string `json:"time"`
	forecast string
	Temperature float64 `json:"temperature"`
	Humidity float64 `json:"humidity"`
	Wind float64 `json:"wind"`
	Pressure float64 `json:"pressure"`
}

type WeatherData struct {
	Area string `json:"area"`
	Forecast []WeatherDataPoint `json:"forecast"`
}

func ParseIso(inp string) string {
	timestamp, err := iso8601.ParseString(inp)

	if err != nil {
		return ""
	}
	hourOffset := 8;
	nanoOffset := 1000000000 * 3600 * hourOffset

	timestamp = timestamp.Add(time.Duration(nanoOffset))

	timestampString := timestamp.Format(time.RFC3339)

	return timestampString
}

func CSVWeatherData(CSVData [][]string, area string) WeatherData {
	amountOfForecast := len(CSVData) - 1
	
	weatherData := WeatherData{
		Area: area,
		Forecast: make([]WeatherDataPoint, amountOfForecast),
	}
	
	for index, row := range CSVData[1:] {
		weatherData.Forecast[amountOfForecast - index - 1] = WeatherDataPoint{
			Time: ParseIso(row[0]),
			forecast: ParseIso(row[1]),
			Temperature: parseFloat(row[2]),
			Humidity: parseFloat(row[3]),
			Wind: parseFloat(row[4]),
			Pressure: parseFloat(row[5]),
		}
	}

	return weatherData;
}


func parseFloat(str string) float64 {
	float, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -10000
	}
	return float
}
