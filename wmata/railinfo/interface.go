package railinfo

import "github.com/awiede/wmata-go-sdk/wmata"

type RailStationInfoAPI interface {
	GetLines() (*GetLinesResponse, error)
	GetParkingInformation(stationCode string) (*GetParkingInformationResponse, error)
	GetPathBetweenStations(fromStation, toStation string) (*GetPathBetweenStationsResponse, error)
	GetStationEntrances(request *GetStationEntrancesRequest) (*GetStationEntrancesResponse, error)
	GetStationInformation(stationCode string) (*GetStationInformationResponse, error)
	GetStationList(line wmata.LineCode) (*GetStationListResponse, error)
	GetStationTimings(stationCode string) (*GetStationTimingsResponse, error)
	GetStationToStationInformation(fromStation, toStation string) (*GetStationToStationInformationResponse, error)
}

var _ RailStationInfoAPI = (*RailStationInfo)(nil)