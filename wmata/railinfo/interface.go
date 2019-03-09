package railinfo

type RailStationInfoAPI interface {
	GetLines() (*GetLinesResponse, error)
	GetParkingInformation(stationCode string) (*GetParkingInformationResponse, error)
	GetPathBetweenStations(fromStation, toStation string) (*GetPathBetweenStationsResponse, error)
	GetStationEntrances(getStationEntranceRequest *GetStationEntrancesRequest) (*GetStationEntrancesResponse, error)
	GetStationInformation(stationCode string) (*GetStationInformationResponse, error)
	GetStationList(lineCode string) (*GetStationListResponse, error)
	GetStationTimings(stationCode string) (*GetStationTimingsResponse, error)
	GetStationToStationInformation(fromStation, toStation string) (*GetStationToStationInformationResponse, error)
}

var _ RailStationInfoAPI = (*RailStationInfo)(nil)
