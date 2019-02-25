package wmata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type StationInformation struct {
	StationCode string `json:"Code"`
	LineCode1   string `json:"LineCode1"`
	LineCode2   string `json:"LineCode2"`
	LineCode3   string `json:"LineCode3"`
	LineCode4   string `json:"LineCode4"`
	StationName string `json:"Name"`
}

type StationInformationListResponse struct {
	Stations []StationInformation `json:"Stations"`
}

type Service struct {
	apiKey string
}

func NewService(apiKey string) *Service {
	return &Service{apiKey: apiKey}
}

func (service *Service) GetStationsByLine(lineCode string) (*StationInformationListResponse, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", "https://api.wmata.com/Rail.svc/json/jStations", nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("api_key", service.apiKey)

	query := request.URL.Query()
	query.Add("LineCode", lineCode)

	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	stations := StationInformationListResponse{}

	err = json.Unmarshal(body, &stations)

	return &stations, err
}
