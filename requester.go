
package main

import (
	"net/http"
	"io"
	"bytes"
	"encoding/json"
	"os"
)

/**
 * Post request to the given address with the given parameters
 * @param host The host to post to (e.g. https://incommodities.io)
 * @param path The path to post to (e.g. /a)
 * @param params The parameters to post
 * @return The response body
 */
func (r *AreaRequester) post(address string, params map[string]string, data []byte) (result []byte, code int, err error) {
	// Declare HTTP Method and Url
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(data))

	if err != nil {
		return []byte{}, 0, err;
	}

	req.Header.Set("Authorization", "Bearer " + r.auth_key);
	if data != nil && 0 < len(data) {
		req.Header.Set("Content-Type", "application/json");
	}
	
	q := req.URL.Query()
	// for each parameter, add it to the query
	for key, value := range params {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return []byte{}, 0, err;
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	

	if err != nil {
		return []byte{}, 0, err;
	}

	return body, resp.StatusCode, nil;
}

type AreaRequester struct {
	from_address string
	to_address string
	auth_key string
	areas []string
}

func (r *AreaRequester) GetDataFromArea(area string) (data []byte, statusCode int, err error) {
	params := make(map[string]string);
	params["area"] = area;

	return r.post(r.from_address, params, nil);
}

func (r *AreaRequester) postData(data []byte) (body []byte, statusCode int, err error) {
	// empty params
	params := make(map[string]string);

	return r.post(r.to_address, params, data);
}

func (r *AreaRequester) PostAreaData(payload WeatherData) ([]byte, int, error) {
	data, err := json.Marshal(payload)
		
	if err != nil {
		return []byte{}, 0, err;
	}

	return r.postData(data);
}

func GetAreaRequester(areas []string) AreaRequester{
	return AreaRequester{
		from_address: os.Getenv("API_ADDRESS_FROM"),
		to_address: os.Getenv("API_ADDRESS_TO"),
		auth_key: os.Getenv("API_KEY"),
		areas: areas,
	}
}


/**
 * Get the area from area requester concurrently
 * @return The responses from the different areas
*/

func (r AreaRequester) RequestAreas() []RequestResponse {
	c := make(chan RequestResponse, 1)

	responses := make([]RequestResponse, len(r.areas))
	
	for _, area := range r.areas {
		go r.RequestArea(c, area)
	}

	for i := 0; i < len(r.areas); i++ {
		responses[i] = <-c
	}

	return responses
}

/**
 * Get the area from area requester
 * @param area The area to get the data for
 * @return The response from the area
*/
func (r AreaRequester) RequestArea(respond chan RequestResponse, area string) {
	data, statusCode, err := r.GetDataFromArea(area);
	respond <- RequestResponse{
		area: area,
		data: data,
		statusCode: statusCode,
		err: err,
	}
}

type RequestResponse struct {
	area string
	data []byte
	statusCode int
	err error
}

func (r RequestResponse) writeToFile () error {
	return writeRawCSVFile(r.data, r.area);
}
