package railinfo

import "github.com/awiede/where-is-wmata/wmata"

type RailStationInfoAPI interface {
	GetLines(lines []string) (*GetLinesResponse, error)
	GetParkingInformation(stationCode string) (*GetParkingInformationResponse, error)
	GetPathBetweenStations(startingStation, endingStation string) (*GetPathBetweenStationsResponse, error)
	GetStationEntrances(request *GetStationEntrancesRequest) (*GetStationEntrancesResponse, error)
	GetStationInformation(stationCode string) (*GetStationInformationResponse, error)
	GetStationList(line wmata.LineCode) (*GetStationListResponse, error)
	GetStationTimings(stationCode string) (*GetStationTimingsResponse, error)
	GetStationToStationInformation(startingStation, endingStation string) (*GetStationToStationInformationResponse, error)
}

var _ RailStationInfoAPI = (*RailStationInfo)(nil)