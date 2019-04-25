package railinfo

import (
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"strconv"
	"strings"
)

const railServiceBaseURL = "https://api.wmata.com/Rail.svc"

type RailInfo interface {
	GetLines() (*GetLinesResponse, error)
	GetParkingInformation(stationCode string) (*GetParkingInformationResponse, error)
	GetPathBetweenStations(fromStation, toStation string) (*GetPathBetweenStationsResponse, error)
	GetStationEntrances(getStationEntranceRequest *GetStationEntrancesRequest) (*GetStationEntrancesResponse, error)
	GetStationInformation(stationCode string) (*GetStationInformationResponse, error)
	GetStationList(lineCode string) (*GetStationListResponse, error)
	GetStationTimings(stationCode string) (*GetStationTimingsResponse, error)
	GetStationToStationInformation(fromStation, toStation string) (*GetStationToStationInformationResponse, error)
}

var _ RailInfo = (*Service)(nil)

// Service provides all API methods for Rail Station Information from WMATA
type Service struct {
	client       *wmata.Client
	responseType wmata.ResponseType
}

// NewService returns a new Service service with a reference to an existing wmata.Client
func NewService(client *wmata.Client, responseType wmata.ResponseType) *Service {
	return &Service{
		client:       client,
		responseType: responseType,
	}
}

type LineResponse struct {
	DisplayName          string `json:"DisplayName" xml:"DisplayName"`
	EndStationCode       string `json:"EndStationCode" xml:"EndStationCode"`
	InternalDestination1 string `json:"InternalDestination1" xml:"InternalDestination1"`
	InternalDestination2 string `json:"InternalDestination2" xml:"InternalDestination2"`
	LineCode             string `json:"LineCode" xml:"LineCode"`
	StartStationCode     string `json:"StartStationCode" xml:"StartStationCode"`
}

type GetLinesResponse struct {
	Lines []LineResponse `json:"Lines" xml:"Lines"`
}

type StationParking struct {
	StationCode string           `json:"Code" xml:"Code"`
	Notes       string           `json:"Notes" xml:"Notes"`
	AllDay      AllDayParking    `json:"AllDayParking" xml:"AllDayParking"`
	ShortTerm   ShortTermParking `json:"ShortTermParking" xml:"ShortTermParking"`
}

type AllDayParking struct {
	TotalCount           int     `json:"TotalCount" xml:"TotalCount"`
	RiderCost            float64 `json:"RiderCost" xml:"RiderCost"`
	NonRiderCost         float64 `json:"NonRiderCost" xml:"NonRiderCost"`
	SaturdayRiderCost    float64 `json:"SaturdayRiderCost" xml:"SaturdayRiderCost"`
	SaturdayNonRiderCost float64 `json:"SaturdayNonRiderCost" xml:"SaturdayNonRiderCost"`
}

type ShortTermParking struct {
	TotalCount int    `json:"TotalCount" xml:"TotalCount"`
	Notes      string `json:"Notes" xml:"Notes"`
}

type GetParkingInformationResponse struct {
	ParkingInformation []StationParking `json:"StationsParking" xml:"StationsParking"`
}

type PathItem struct {
	DistanceToPreviousStation int    `json:"DistanceToPrev" xml:"DistanceToPrev"`
	LineCode                  string `json:"LineCode" xml:"LineCode"`
	SequenceNumber            int    `json:"SeqNum" xml:"SeqNum"`
	StationCode               string `json:"StationCode" xml:"StationCode"`
	StationName               string `json:"StationName" xml:"StationName"`
}

type GetPathBetweenStationsResponse struct {
	Path []PathItem `json:"Path" xml:"Path"`
}

type GetStationEntrancesRequest struct {
	latitude  float64
	longitude float64
	radius    float64
}

type StationEntrance struct {
	Description string `json:"Description" xml:"Description"`
	// Deprecated: ID response field is deprecated
	ID           string  `json:"ID" xml:"ID"`
	Latitude     float64 `json:"Lat" xml:"Lat"`
	Longitude    float64 `json:"Lon" xml:"Lon"`
	Name         string  `json:"Name" xml:"Name"`
	StationCode1 string  `json:"StationCode1" xml:"StationCode1"`
	StationCode2 string  `json:"StationCode2" xml:"StationCode2"`
}

type GetStationEntrancesResponse struct {
	Entrances []StationEntrance `json:"Entrances" xml:"Entrances"`
}

type StationAddress struct {
	City   string `json:"City" xml:"City"`
	State  string `json:"State" xml:"State"`
	Street string `json:"Street" xml:"Street"`
	Zip    string `json:"Zip" xml:"Zip"`
}

type GetStationInformationResponse struct {
	Address          StationAddress `json:"Address" xml:"Address"`
	Latitude         float64        `json:"Lat" xml:"Lat"`
	LineCode1        string         `json:"LineCode1" xml:"LineCode1"`
	LineCode2        string         `json:"LineCode2" xml:"LineCode2"`
	Longitude        float64        `json:"Lon" xml:"Lon"`
	Name             string         `json:"Name" xml:"Name"`
	StationCode      string         `json:"Code" xml:"Code"`
	StationTogether1 string         `json:"StationTogether1" xml:"StationTogether1"`
	StationTogether2 string         `json:"StationTogether2" xml:"StationTogether2"`
}

type GetStationListResponse struct {
	Stations []GetStationListResponseItem `json:"Stations" xml:"Stations"`
}

type GetStationListResponseItem struct {
	Address          StationAddress `json:"Address" xml:"Address"`
	StationCode      string         `json:"Code" xml:"Code"`
	Latitude         float64        `json:"Lat" xml:"Lat"`
	LineCode1        string         `json:"LineCode1" xml:"LineCode1"`
	LineCode2        string         `json:"LineCode2" xml:"LineCode2"`
	LineCode3        string         `json:"LineCode3" xml:"LineCode3"`
	LineCode4        string         `json:"LineCode4" xml:"LineCode4"`
	Longitude        float64        `json:"Lon" xml:"Lon"`
	Name             string         `json:"Name" xml:"Name"`
	StationTogether1 string         `json:"StationTogether1" xml:"StationTogether1"`
	StationTogether2 string         `json:"StationTogether2" xml:"StationTogether2"`
}

type StationTrainInformation struct {
	Time               string `json:"Time" xml:"Time"`
	DestinationStation string `json:"DestinationStation" xml:"DestinationStation"`
}

type StationDayItem struct {
	OpeningTime string                    `json:"OpeningTime" xml:"OpeningTime"`
	FirstTrains []StationTrainInformation `json:"FirstTrains" xml:"FirstTrains"`
	LastTrains  []StationTrainInformation `json:"LastTrains" xml:"LastTrains"`
}

type StationTime struct {
	StationCode string         `json:"Code" xml:"Code"`
	StationName string         `json:"StationName" xml:"StationName"`
	Monday      StationDayItem `json:"Monday" xml:"Monday"`
	Tuesday     StationDayItem `json:"Tuesday" xml:"Tuesday"`
	Wednesday   StationDayItem `json:"Wednesday" xml:"Wednesday"`
	Thursday    StationDayItem `json:"Thursday" xml:"Thursday"`
	Friday      StationDayItem `json:"Friday" xml:"Friday"`
	Saturday    StationDayItem `json:"Saturday" xml:"Saturday"`
	Sunday      StationDayItem `json:"Sunday" xml:"Sunday"`
}

type GetStationTimingsResponse struct {
	StationTimes []StationTime `json:"StationTimes" xml:"StationTimes"`
}

type RailFare struct {
	OffPeakTime    float64 `json:"OffPeakTime" xml:"OffPeakTime"`
	PeakTime       float64 `json:"PeakTime" xml:"PeakTime"`
	SeniorDisabled float64 `json:"SeniorDisabled" xml:"SeniorDisabled"`
}

type StationToStation struct {
	CompositeMiles     float64  `json:"CompositeMiles" xml:"CompositeMiles"`
	DestinationStation string   `json:"DestinationStation" xml:"DestinationStation"`
	Fare               RailFare `json:"RailFare" xml:"RailFare"`
	Time               int      `json:"RailTime" xml:"RailTime"`
	SourceStation      string   `json:"SourceStation" xml:"SourceStation"`
}

type GetStationToStationInformationResponse struct {
	StationToStationInformation []StationToStation `json:"StationToStationInfos" xml:"StationToStationInfos"`
}

func (railService *Service) GetLines() (*GetLinesResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jLines")
	case wmata.XML:
		requestUrl.WriteString("/Lines")
	}

	lines := GetLinesResponse{}

	return &lines, railService.client.BuildAndSendGetRequest(requestUrl.String(), nil, &lines)
}

func (railService *Service) GetParkingInformation(stationCode string) (*GetParkingInformationResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jStationParking")
	case wmata.XML:
		requestUrl.WriteString("/StationParking")
	}

	parkingInformation := GetParkingInformationResponse{}

	return &parkingInformation, railService.client.BuildAndSendGetRequest(requestUrl.String(), map[string]string{"StationCode": stationCode}, &parkingInformation)
}

func (railService *Service) GetPathBetweenStations(fromStation, toStation string) (*GetPathBetweenStationsResponse, error) {
	if fromStation == "" || toStation == "" {
		return nil, errors.New("fromStation and toStation are required parameters")
	}

	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jPath")
	case wmata.XML:
		requestUrl.WriteString("/Path")
	}

	path := GetPathBetweenStationsResponse{}

	return &path, railService.client.BuildAndSendGetRequest(requestUrl.String(), map[string]string{"FromStationCode": fromStation, "ToStationCode": toStation}, &path)

}

func (railService *Service) GetStationEntrances(getStationEntranceRequest *GetStationEntrancesRequest) (*GetStationEntrancesResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jStationEntrances")
	case wmata.XML:
		requestUrl.WriteString("/StationEntrances")
	}

	var queryParams map[string]string

	if getStationEntranceRequest != nil {
		queryParams = map[string]string{
			"Lat":    strconv.FormatFloat(getStationEntranceRequest.latitude, 'g', -1, 64),
			"Lon":    strconv.FormatFloat(getStationEntranceRequest.longitude, 'g', -1, 64),
			"Radius": strconv.FormatFloat(getStationEntranceRequest.radius, 'g', -1, 64),
		}
	}
	entrances := GetStationEntrancesResponse{}

	return &entrances, railService.client.BuildAndSendGetRequest(requestUrl.String(), queryParams, &entrances)
}

func (railService *Service) GetStationInformation(stationCode string) (*GetStationInformationResponse, error) {
	if stationCode == "" {
		return nil, errors.New("stationCode is a required parameter")
	}

	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jStationInfo")
	case wmata.XML:
		requestUrl.WriteString("/StationInfo")
	}

	stationInformation := GetStationInformationResponse{}

	return &stationInformation, railService.client.BuildAndSendGetRequest(requestUrl.String(), map[string]string{"StationCode": stationCode}, &stationInformation)
}

func (railService *Service) GetStationList(lineCode string) (*GetStationListResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jStations")
	case wmata.XML:
		requestUrl.WriteString("/Stations")
	}

	stationList := GetStationListResponse{}

	return &stationList, railService.client.BuildAndSendGetRequest(requestUrl.String(), map[string]string{"LineCode": lineCode}, &stationList)
}

func (railService *Service) GetStationTimings(stationCode string) (*GetStationTimingsResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jStationTimes")
	case wmata.XML:
		requestUrl.WriteString("/StationTimes")
	}

	stationTimings := GetStationTimingsResponse{}

	return &stationTimings, railService.client.BuildAndSendGetRequest(requestUrl.String(), map[string]string{"StationCode": stationCode}, &stationTimings)
}

func (railService *Service) GetStationToStationInformation(fromStation, toStation string) (*GetStationToStationInformationResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(railServiceBaseURL)

	switch railService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jSrcStationToDstStationInfo")
	case wmata.XML:
		requestUrl.WriteString("/SrcStationToDstStationInfo")
	}

	queryParams := make(map[string]string)

	if fromStation != "" {
		queryParams["FromStationCode"] = fromStation
	}

	if toStation != "" {
		queryParams["ToStationCode"] = toStation
	}

	stationToStation := GetStationToStationInformationResponse{}

	return &stationToStation, railService.client.BuildAndSendGetRequest(requestUrl.String(), queryParams, &stationToStation)
}
