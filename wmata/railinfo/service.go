package railinfo

import (
	"encoding/json"
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"io/ioutil"
	"net/http"
	"strconv"
)

// RailStationInfo provides all API methods for Rail Station Information from WMATA
type RailStationInfo struct {
	client *wmata.Client
}

// NewService returns a new RailStationInfo service with a reference to an existing wmata.Client
func NewService(client *wmata.Client) *RailStationInfo {
	return &RailStationInfo{
		client: client,
	}
}

type LineResponse struct {
	DisplayName          string         `json:"DisplayName"`
	EndStationCode       string         `json:"EndStationCode"`
	InternalDestination1 string         `json:"InternalDestination1"`
	InternalDestination2 string         `json:"InternalDestination2"`
	LineCode             wmata.LineCode `json:"LineCode"`
	StartStationCode     string         `json:"StartStationCode"`
}

type GetLinesResponse struct {
	Lines []LineResponse `json:"Lines"`
}

type StationParking struct {
	StationCode      string `json:"Code"`
	Notes            string `json:"Notes"`
	AllDayParking    `json:"AllDayParking"`
	ShortTermParking `json:"ShortTermParking"`
}

type AllDayParking struct {
	TotalCount           int     `json:"TotalCount"`
	RiderCost            float64 `json:"RiderCost"`
	NonRiderCost         float64 `json:"NonRiderCost"`
	SaturdayRiderCost    float64 `json:"SaturdayRiderCost"`
	SaturdayNonRiderCost float64 `json:"SaturdayNonRiderCost"`
}

type ShortTermParking struct {
	TotalCount int    `json:"TotalCount"`
	Notes      string `json:"Notes"`
}

type GetParkingInformationResponse struct {
	ParkingInformation []StationParking `json:"StationsParking"`
}

type PathItem struct {
	DistanceToPreviousStation int            `json:"DistanceToPrev"`
	LineCode                  wmata.LineCode `json:"LineCode"`
	SequenceNumber            int            `json:"SeqNum"`
	StationCode               string         `json:"StationCode"`
	StationName               string         `json:"StationName"`
}

type GetPathBetweenStationsResponse struct {
	Path []PathItem `json:"Path"`
}

type GetStationEntrancesRequest struct {
	latitude  float64
	longitude float64
	radius    float64
}

type StationEntrance struct {
	Description string `json:"Description"`
	// Deprecated: ID response field is deprecated
	ID           string  `json:"ID"`
	Latitude     float64 `json:"Lat"`
	Longitude    float64 `json:"Lon"`
	Name         string  `json:"Name"`
	StationCode1 string  `json:"StationCode1"`
	StationCode2 string  `json:"StationCode2"`
}

type GetStationEntrancesResponse struct {
	Entrances []StationEntrance `json:"Entrances"`
}

type StationAddress struct {
	City   string `json:"City"`
	State  string `json:"State"`
	Street string `json:"Street"`
	Zip    string `json:"Zip"`
}

type GetStationInformationResponse struct {
	Address          StationAddress `json:"Address"`
	Latitude         float64        `json:"Lat"`
	LineCode1        wmata.LineCode `json:"LineCode1"`
	LineCode2        wmata.LineCode `json:"LineCode2"`
	Longitude        float64        `json:"Lon"`
	Name             string         `json:"Name"`
	StationCode      string         `json:"Code"`
	StationTogether1 string         `json:"StationTogether1"`
	StationTogether2 string         `json:"StationTogether2"`
}

type GetStationListResponse struct {
}

type GetStationTimingsResponse struct {
}

type GetStationToStationInformationResponse struct {
}

func (railService *RailStationInfo) GetLines() (*GetLinesResponse, error) {
	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Rail.svc/json/jLines", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, railService.client.APIKey)

	response, responseErr := railService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	lines := GetLinesResponse{}

	unmarshalErr := json.Unmarshal(body, &lines)

	return &lines, unmarshalErr
}

func (railService *RailStationInfo) GetParkingInformation(stationCode string) (*GetParkingInformationResponse, error) {
	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Rail.svc/json/jStationParking", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, railService.client.APIKey)

	if stationCode != "" {
		query := request.URL.Query()
		query.Add("StationCode", stationCode)

		request.URL.RawQuery = query.Encode()
	}

	response, responseErr := railService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	parkingInformation := GetParkingInformationResponse{}

	unmarshalErr := json.Unmarshal(body, &parkingInformation)

	return &parkingInformation, unmarshalErr
}

func (railService *RailStationInfo) GetPathBetweenStations(fromStation, toStation string) (*GetPathBetweenStationsResponse, error) {
	if fromStation == "" || toStation == "" {
		return nil, errors.New("fromStation and toStation are required parameters")
	}

	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Rail.svc/json/jPath", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, railService.client.APIKey)

	query := request.URL.Query()
	query.Add("FromStationCode", fromStation)
	query.Add("ToStationCode", toStation)

	request.URL.RawQuery = query.Encode()

	response, responseErr := railService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	path := GetPathBetweenStationsResponse{}

	unmarshalErr := json.Unmarshal(body, &path)

	return &path, unmarshalErr
}

func (railService *RailStationInfo) GetStationEntrances(getStationEntranceRequest *GetStationEntrancesRequest) (*GetStationEntrancesResponse, error) {
	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Rail.svc/json/jStationEntrances", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, railService.client.APIKey)

	if getStationEntranceRequest != nil {
		query := request.URL.Query()
		query.Add("Lat", strconv.FormatFloat(getStationEntranceRequest.latitude, 'g', -1, 64))
		query.Add("Lon", strconv.FormatFloat(getStationEntranceRequest.longitude, 'g', -1, 64))
		query.Add("Radius", strconv.FormatFloat(getStationEntranceRequest.radius, 'g', -1, 64))

		request.URL.RawQuery = query.Encode()
	}

	response, responseErr := railService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	entrances := GetStationEntrancesResponse{}

	unmarshalErr := json.Unmarshal(body, &entrances)

	return &entrances, unmarshalErr
}

func (railService *RailStationInfo) GetStationInformation(stationCode string) (*GetStationInformationResponse, error) {
	if stationCode == "" {
		return nil, errors.New("stationCode is a required parameter")
	}

	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Rail.svc/json/jStationInfo", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, railService.client.APIKey)

	query := request.URL.Query()
	query.Add("StationCode", stationCode)

	request.URL.RawQuery = query.Encode()

	response, responseErr := railService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	stationInformation := GetStationInformationResponse{}

	unmarshalErr := json.Unmarshal(body, &stationInformation)

	return &stationInformation, unmarshalErr
}

func (railService *RailStationInfo) GetStationList(line wmata.LineCode) (*GetStationListResponse, error) {
	return nil, nil
}

func (railService *RailStationInfo) GetStationTimings(stationCode string) (*GetStationTimingsResponse, error) {
	return nil, nil
}

func (railService *RailStationInfo) GetStationToStationInformation(fromStation, toStation string) (*GetStationToStationInformationResponse, error) {
	return nil, nil
}
