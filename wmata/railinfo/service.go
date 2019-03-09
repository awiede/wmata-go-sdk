package railinfo

import "github.com/awiede/wmata-go-sdk/wmata"

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

type GetLinesResponse struct {
}

type GetParkingInformationResponse struct {
}

type GetPathBetweenStationsResponse struct {
}

type GetStationEntrancesRequest struct {
	latitude  int
	longitude int
	radius    int
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

func (railService *RailStationInfo) GetLines(lines []string) (*GetLinesResponse, error) {
	return nil, nil
}

func (railService *RailStationInfo) GetParkingInformation(stationCode string) (*GetParkingInformationResponse, error) {
	return nil, nil
}

func (railService *RailStationInfo) GetPathBetweenStations(startingStation, endingStation string) (*GetPathBetweenStationsResponse, error) {
	return nil, nil
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

func (railService *RailStationInfo) GetStationToStationInformation(startingStation, endingStation string) (*GetStationToStationInformationResponse, error) {
	return nil, nil
}
