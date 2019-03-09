package railinfo

import (
	"encoding/json"
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"io/ioutil"
	"net/http"
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

type GetStationEntrancesResponse struct {
}

type GetStationInformationResponse struct {
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

func (railService *RailStationInfo) GetStationEntrances(request *GetStationEntrancesRequest) (*GetStationEntrancesResponse, error) {
	return nil, nil
}

func (railService *RailStationInfo) GetStationInformation(stationCode string) (*GetStationInformationResponse, error) {
	return nil, nil
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
