package railinfo

import (
	"encoding/xml"
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"strconv"
	"strings"
)

const railServiceBaseURL = "https://api.wmata.com/Rail.svc"

// RailInfo defines the methods available in the WMATA "Rail Station Information" API
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
	XMLName xml.Name       `json:"-" xml:"http://www.wmata.com LinesResp"`
	Lines   []LineResponse `json:"Lines" xml:"Lines>Line"`
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
	XMLName            xml.Name         `json:"-" xml:"http://www.wmata.com StationParkingResp"`
	ParkingInformation []StationParking `json:"StationsParking" xml:"StationsParking>StationParking"`
}

type PathItem struct {
	DistanceToPreviousStation int    `json:"DistanceToPrev" xml:"DistanceToPrev"`
	LineCode                  string `json:"LineCode" xml:"LineCode"`
	SequenceNumber            int    `json:"SeqNum" xml:"SeqNum"`
	StationCode               string `json:"StationCode" xml:"StationCode"`
	StationName               string `json:"StationName" xml:"StationName"`
}

type GetPathBetweenStationsResponse struct {
	XMLName xml.Name   `json:"-" xml:"http://www.wmata.com PathResp"`
	Path    []PathItem `json:"Path" xml:"Path>MetroPathItem"`
}

type GetStationEntrancesRequest struct {
	Latitude  float64
	Longitude float64
	Radius    float64
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
	XMLName   xml.Name          `json:"-" xml:"http://www.wmata.com StationEntrancesResp"`
	Entrances []StationEntrance `json:"Entrances" xml:"Entrances>StationEntrance"`
}

type StationAddress struct {
	City   string `json:"City" xml:"City"`
	State  string `json:"State" xml:"State"`
	Street string `json:"Street" xml:"Street"`
	Zip    string `json:"Zip" xml:"Zip"`
}

type GetStationInformationResponse struct {
	XMLName          xml.Name       `json:"-" xml:"http://www.wmata.com Station"`
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
	XMLName  xml.Name                     `json:"-" xml:"http://www.wmata.com StationsResp"`
	Stations []GetStationListResponseItem `json:"Stations" xml:"Stations>Station"`
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
	FirstTrains []StationTrainInformation `json:"FirstTrains" xml:"FirstTrains>Train"`
	LastTrains  []StationTrainInformation `json:"LastTrains" xml:"LastTrains>Train"`
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
	XMLName      xml.Name      `json:"-" xml:"http://www.wmata.com StationTimeResp"`
	StationTimes []StationTime `json:"StationTimes" xml:"StationTimes>StationTime"`
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
	XMLName                     xml.Name           `json:"-" xml:"http://www.wmata.com StationToStationInfoResp"`
	StationToStationInformation []StationToStation `json:"StationToStationInfos" xml:"StationToStationInfos>StationToStationInfo"`
}

// GetLines retrieves information about all rail lines
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe330c
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

	return &lines, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), nil, &lines)
}

// GetParkingInformation retrieves parking information for a given station
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe330d?
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

	return &parkingInformation, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), map[string]string{"StationCode": stationCode}, &parkingInformation)
}

// GetPathBetweenStations retrieves an ordered list of stations and distances between two stations
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe330e?
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

	return &path, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), map[string]string{"FromStationCode": fromStation, "ToStationCode": toStation}, &path)

}

// GetStationEntrances retrieves a list of station entrances near the provided coordinates
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe330f?
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
		queryParams = make(map[string]string)

		if getStationEntranceRequest.Latitude != 0 {
			queryParams["Lat"] = strconv.FormatFloat(getStationEntranceRequest.Latitude, 'g', -1, 64)
		}

		if getStationEntranceRequest.Longitude != 0 {
			queryParams["Lon"] = strconv.FormatFloat(getStationEntranceRequest.Longitude, 'g', -1, 64)
		}

		if getStationEntranceRequest.Radius != 0 {
			queryParams["Radius"] = strconv.FormatFloat(getStationEntranceRequest.Radius, 'g', -1, 64)
		}
	}
	entrances := GetStationEntrancesResponse{}

	return &entrances, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), queryParams, &entrances)
}

// GetStationInformation retrieves station location and address information by station code
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe3310?
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

	return &stationInformation, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), map[string]string{"StationCode": stationCode}, &stationInformation)
}

// GetStationList retrieves a list of station location and address information for all stations on a given line
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe3311?
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

	return &stationList, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), map[string]string{"LineCode": lineCode}, &stationList)
}

// GetStationTimings retrieves opening and scheduled first and last train times for a given station
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe3312?
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

	return &stationTimings, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), map[string]string{"StationCode": stationCode}, &stationTimings)
}

// GetStationToStationInformation retrieves distance, fare and estimated travel time between the given two stations
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe3313?
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

	return &stationToStation, railService.client.BuildAndSendGetRequest(railService.responseType, requestUrl.String(), queryParams, &stationToStation)
}
