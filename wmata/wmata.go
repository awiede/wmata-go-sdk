package wmata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// StationInformation holds high level information about metro stations
type StationInformation struct {
	StationCode string `json:"Code"`
	LineCode1   string `json:"LineCode1"`
	LineCode2   string `json:"LineCode2"`
	LineCode3   string `json:"LineCode3"`
	LineCode4   string `json:"LineCode4"`
	StationName string `json:"Name"`
}

// StationInformationListResponse is a response wrapper for API requests for multiple stations
type StationInformationListResponse struct {
	Stations []StationInformation `json:"Stations"`
}

// Service encapsulates all dependencies needed to run the WMATA service
type Service struct {
	APIKey string
}

// GetStationsByLine retrieves all the metro stations by line code. If no line code specified all stations will be retrieved
func (service *Service) GetStationsByLine(lineCode string) (*StationInformationListResponse, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", "https://api.wmata.com/Rail.svc/json/jStations", nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("api_key", service.APIKey)

	query := request.URL.Query()
	query.Add("LineCode", lineCode)

	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("error closing response body: %s", closeErr)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	stations := StationInformationListResponse{}

	err = json.Unmarshal(body, &stations)

	return &stations, err
}
