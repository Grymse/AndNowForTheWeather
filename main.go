package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func loadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
}

var AREAS = []string{
	"Eagle River",
	"Kincaid Park",
	"Far North Bicentennial Park",
	"Bear Valley",
	"Fire Island",
}

func main() {
	loadEnvVariables()
	requester := GetAreaRequester(AREAS)
	requester.passInformationIfFreshData()
	for {
		time.Sleep(time.Second * 5)
		if requester.passInformationIfFreshData() {
			break
		}
	}

	for {
		timeBefore := time.Now().UnixMilli()
		requester.passInformationIfFreshData()
		timeAfter := time.Now().UnixMilli()
		timeTaken := timeAfter - timeBefore
		timeToSleep := 1800*1000 - timeTaken
		time.Sleep(time.Millisecond * time.Duration(timeToSleep))
	}
}

func (responses RequestResponse) IsSuccessful() bool {
	return responses.err == nil && responses.statusCode == 200
}

func (requester AreaRequester) passInformationIfFreshData() bool {
	// FETCH DATA FROM API
	responses := requester.RequestAreas()

	fmt.Println("Attempt to pass information at " + time.Now().Format("2006-01-02 15:04:05"))

	var containsAnyDifferentData = false
	// COMPARE DATA TO DATA IN FILES
	for _, response := range responses {
		if !response.IsSuccessful() {
			continue
		}

		// COMPARE DATA TO DATA IN FILES
		data, err := os.ReadFile(getFilePath(response.area, "csv"))
		if err != nil {
			containsAnyDifferentData = true
			break
		}

		if !bytes.Equal(data, response.data) {
			fmt.Println("diff")
			containsAnyDifferentData = true
			break
		}
	}
	if !containsAnyDifferentData {
		return false
	}

	fmt.Println("Pass information to API")

	// IF DATA IS NEW, DO REQUEST
	for _, response := range responses {
		go requester.attemptInformationPass(response)
	}

	return true
}

func (requester AreaRequester) attemptInformationPass(response RequestResponse) {
	area := response.area
	data := response.data

	writeRawCSVFile(data, area)
	reader := csv.NewReader(bytes.NewBuffer(data))
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error converting to csv file for " + area)
		return
	}

	weatherData := CSVWeatherData(rows, area)
	weatherDataRequest, statusCode, err := requester.PostAreaData(weatherData)
	if err != nil {
		fmt.Println("Error posting data for " + area)
		return
	}

	fmt.Print("completed ")
	fmt.Print(string(area))
	fmt.Print(" request with status code ")
	fmt.Print(statusCode)
	fmt.Print(" and data ")
	fmt.Println(string(weatherDataRequest))
}
