package main

import (
    "io"
		"net/http"
		"net/url"
		"fmt"
		"os"
)


func check(e error) {
	if e != nil {
			panic(e)
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
	/* d1 := []byte("hello\ngo\n")
	err := os.WriteFile("./dat1.csv", d1, 0644)
	check(err) */

	for _, area := range AREAS {
		bytes, err := fetchForecastData(area);
		if err != nil {
			fmt.Println("Error fetching data for area " + area);
			fmt.Println(err);
			continue;
		}
		err = writeToFile(bytes, area + ".csv");
		if err != nil {
			fmt.Println("Error writing data for area " + area);
			fmt.Println(err);
			continue;
		}
	}
}

const API_ADDRESS = "https://incommodities.io/a"


func writeToFile(data []byte, filename string) (err error) {
	return os.WriteFile(filename, data, 0644);
}

/**
 * Fetch the forecast data for the given area
 * @param area The area to fetch data for
 * @return The forecast data
 */

func fetchForecastData(area string) (data []byte, err error) {
	params := []KeyValuePair{
		KeyValuePair{Key: "area", Value: area},
	}

	return postRequest(API_ADDRESS, params);
}

/**
 * Post request to the given address with the given parameters
 * @param address The address to post to
 * @param params The parameters to post
 * @return The response body
 */

func postRequest(address string, params []KeyValuePair) (data []byte, err error) {
	paramObject := url.Values{};
	for _, param := range params {
		paramObject.Add(param.Key, param.Value);
	}
	resp, err := http.PostForm(address, paramObject);

	if err != nil {
		return []byte{}, err;
	}

	defer resp.Body.Close();

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err;
	}

	return body, nil;
}