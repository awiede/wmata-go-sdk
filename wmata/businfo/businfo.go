package businfo

import (
	"encoding/xml"
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"strconv"
	"strings"
)

const busInfoBaseUrl = "https://api.wmata.com/Bus.svc"

type BusInfo interface {
	GetPositions(request *GetPositionsRequest) (*GetPositionsResponse, error)
	GetRouteDetails(routeID, date string) (*GetRouteDetailsResponse, error)
	GetRoutes() (*GetRoutesResponse, error)
	GetSchedule(routeID, date string, includeVariations bool) (*GetScheduleResponse, error)
	GetScheduleAtStop(stopID, date string) (*GetScheduleAtStopResponse, error)
	GetStops(request *GetStopsRequest) (*GetStopsResponse, error)
}

var _ BusInfo = (*Service)(nil)

// Service provides all API methods for Bus Route and Stop Information from WMATA
type Service struct {
	client       *wmata.Client
	responseType wmata.ResponseType
}

// NewService returns a new BusInfo Service service with a reference to an existing wmata.Client
func NewService(client *wmata.Client, responseType wmata.ResponseType) *Service {
	return &Service{
		client:       client,
		responseType: responseType,
	}
}

type GetPositionsRequest struct {
	RouteID   string
	Latitude  float64
	Longitude float64
	Radius    float64
}

type GetPositionsResponse struct {
	XMLName      xml.Name      `json:"-" xml:"http://www.wmata.com BusPositionsResp"`
	BusPositions []BusPosition `json:"BusPositions" xml:"BusPositions>BusPosition"`
}

type BusPosition struct {
	BlockNumber string `json:"BlockNumber" xml:"BlockNumber"`
	DateTime    string `json:"DateTime" xml:"DateTime"`
	Deviation   int    `json:"Deviation" xml:"Deviation"`
	// Deprecated: DirectionNumber response field is deprecated, use DirectionText
	DirectionNumber int     `json:"DirectionNum" xml:"DirectionNum"`
	DirectionText   string  `json:"DirectionText" xml:"DirectionText"`
	Latitude        float64 `json:"Lat" xml:"Lat"`
	Longitude       float64 `json:"Lon" xml:"Lon"`
	RouteID         string  `json:"RouteID" xml:"RouteID"`
	TripEndTime     string  `json:"TripEndTime" xml:"TripEndTime"`
	TripDestination string  `json:"TripHeadsign" xml:"TripHeadsign"`
	TripID          string  `json:"TripID" xml:"TripID"`
	TripStartTime   string  `json:"TripStartTime" xml:"TripStartTime"`
	VehicleID       string  `json:"VehicleID" xml:"VehicleID"`
}

type GetRouteDetailsResponse struct {
	XMLName    xml.Name  `json:"-" xml:"http://www.wmata.com RouteDetailsInfo"`
	Direction0 Direction `json:"Direction0" xml:"Direction0"`
	Direction1 Direction `json:"Direction1" xml:"Direction1"`
	Name       string    `json:"Name" xml:"Name"`
	RouteID    string    `json:"RouteID" xml:"RouteID"`
}

type Direction struct {
	// Deprecated: DirectionNumber response field is deprecated, use DirectionText
	DirectionNumber string       `json:"DirectionNum" xml:"DirectionNum"`
	DirectionText   string       `json:"DirectionText" xml:"DirectionText"`
	Shapes          []ShapePoint `json:"Shape" xml:"Shape>ShapePoint"`
	Stops           []Stop       `json:"Stops" xml:"Stops>Stop"`
	TripDestination string       `json:"TripHeadsign" xml:"TripHeadsign"`
}

type ShapePoint struct {
	Latitude       float64 `json:"Lat" xml:"Lat"`
	Longitude      float64 `json:"Lon" xml:"Lon"`
	SequenceNumber int     `json:"SeqNum" xml:"SeqNum"`
}

type Stop struct {
	Latitude  float64  `json:"Lat" xml:"Lat"`
	Longitude float64  `json:"Lon" xml:"Lon"`
	Name      string   `json:"Name" xml:"Name"`
	Routes    []string `json:"Routes" xml:"Routes>string"`
	StopID    string   `json:"StopID" xml:"StopID"`
}

type GetRoutesResponse struct {
	XMLName xml.Name `json:"-" xml:"http://www.wmata.com RoutesResp"`
	Routes  []Route  `json:"Routes" xml:"Routes>Route"`
}

type Route struct {
	Name            string `json:"Name" xml:"Name"`
	RouteID         string `json:"RouteID" xml:"RouteID"`
	LineDescription string `json:"LineDescription" xml:"LineDescription"`
}

type GetScheduleResponse struct {
	XMLName    xml.Name `json:"-" xml:"http://www.wmata.com RouteScheduleInfo"`
	Direction0 []Trip   `json:"Direction0" xml:"Direction0>Trip"`
	Direction1 []Trip   `json:"Direction1" xml:"Direction1>Trip"`
	Name       string   `json:"Name" xml:"Name"`
}

type Trip struct {
	DirectionNumber string     `json:"DirectionNum" xml:"DirectionNum"`
	EndTime         string     `json:"EndTime" xml:"EndTime"`
	RouteID         string     `json:"RouteID" xml:"RouteID"`
	StartTime       string     `json:"StartTime" xml:"StartTime"`
	StopTimes       []StopTime `json:"StopTimes" xml:"StopTimes>StopTime"`
	TripDirection   string     `json:"TripDirectionText" xml:"TripDirectionText"`
	TripDestination string     `json:"TripHeadsign" xml:"TripHeadsign"`
	TripID          string     `json:"TripID" xml:"TripID"`
}

type StopTime struct {
	StopID       string `json:"StopID" xml:"StopID"`
	StopName     string `json:"StopName" xml:"StopName"`
	StopSequence int    `json:"StopSeq" xml:"StopSeq"`
	Time         string `json:"Time" xml:"Time"`
}

type GetScheduleAtStopResponse struct {
	XMLName          xml.Name          `json:"-" xml:"http://www.wmata.com StopScheduleInfo"`
	ScheduleArrivals []ScheduleArrival `json:"ScheduleArrivals" xml:"ScheduleArrivals>StopScheduleArrival"`
	StopInfo         Stop              `json:"Stop" xml:"Stop"`
}

type ScheduleArrival struct {
	DirectionNumber string `json:"DirectionNum" xml:"DirectionNum"`
	EndTime         string `json:"EndTime" xml:"EndTime"`
	RouteID         string `json:"RouteID" xml:"RouteID"`
	ScheduleTime    string `json:"ScheduleTime" xml:"ScheduleTime"`
	StartTime       string `json:"StartTime" xml:"StartTime"`
	TripDirection   string `json:"TripDirectionText" xml:"TripDirectionText"`
	TripDestination string `json:"TripHeadsign" xml:"TripHeadsign"`
	TripID          string `json:"TripID" xml:"TripID"`
}

type GetStopsRequest struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}

type GetStopsResponse struct {
	XMLName xml.Name `json:"-" xml:"http://www.wmata.com StopsResp"`
	Stops   []Stop   `json:"Stops" xml:"Stops>Stop"`
}

func (busService *Service) GetPositions(request *GetPositionsRequest) (*GetPositionsResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(busInfoBaseUrl)

	switch busService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jBusPositions")
	case wmata.XML:
		requestUrl.WriteString("/BusPositions")
	}

	var queryParams map[string]string
	if request != nil {
		queryParams = make(map[string]string)

		if request.RouteID != "" {
			queryParams["RouteID"] = request.RouteID
		}

		if request.Latitude != 0 {
			queryParams["Lat"] = strconv.FormatFloat(request.Latitude, 'g', -1, 64)
		}

		if request.Longitude != 0 {
			queryParams["Lon"] = strconv.FormatFloat(request.Longitude, 'g', -1, 64)
		}

		if request.Radius != 0 {
			queryParams["Radius"] = strconv.FormatFloat(request.Radius, 'g', -1, 64)
		}
	}

	positions := GetPositionsResponse{}

	return &positions, busService.client.BuildAndSendGetRequest(busService.responseType, requestUrl.String(), queryParams, &positions)
}

func (busService *Service) GetRouteDetails(routeID, date string) (*GetRouteDetailsResponse, error) {
	if routeID == "" {
		return nil, errors.New("routeID is required")
	}

	var requestUrl strings.Builder
	requestUrl.WriteString(busInfoBaseUrl)

	switch busService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jRouteDetails")
	case wmata.XML:
		requestUrl.WriteString("/RouteDetails")
	}

	queryParams := map[string]string{
		"RouteID": routeID,
	}

	if date != "" {
		queryParams["Date"] = date
	}

	path := GetRouteDetailsResponse{}

	return &path, busService.client.BuildAndSendGetRequest(busService.responseType, requestUrl.String(), queryParams, &path)
}

func (busService *Service) GetRoutes() (*GetRoutesResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(busInfoBaseUrl)

	switch busService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jRoutes")
	case wmata.XML:
		requestUrl.WriteString("/Routes")
	}

	routes := GetRoutesResponse{}

	return &routes, busService.client.BuildAndSendGetRequest(busService.responseType, requestUrl.String(), nil, &routes)

}

func (busService *Service) GetSchedule(routeID, date string, includeVariations bool) (*GetScheduleResponse, error) {
	if routeID == "" {
		return nil, errors.New("routeID is required")
	}

	var requestUrl strings.Builder
	requestUrl.WriteString(busInfoBaseUrl)

	switch busService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jRouteSchedule")
	case wmata.XML:
		requestUrl.WriteString("/RouteSchedule")
	}

	queryParams := map[string]string{
		"RouteID":             routeID,
		"IncludingVariations": strconv.FormatBool(includeVariations),
	}

	if date != "" {
		queryParams["Date"] = date
	}

	schedule := GetScheduleResponse{}

	return &schedule, busService.client.BuildAndSendGetRequest(busService.responseType, requestUrl.String(), queryParams, &schedule)

}

func (busService *Service) GetScheduleAtStop(stopID, date string) (*GetScheduleAtStopResponse, error) {
	if stopID == "" {
		return nil, errors.New("stopID is required")
	}

	var requestUrl strings.Builder
	requestUrl.WriteString(busInfoBaseUrl)

	switch busService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jStopSchedule")
	case wmata.XML:
		requestUrl.WriteString("/StopSchedule")
	}

	queryParams := map[string]string{
		"StopID": stopID,
	}

	if date != "" {
		queryParams["Date"] = date
	}

	stopSchedule := GetScheduleAtStopResponse{}

	return &stopSchedule, busService.client.BuildAndSendGetRequest(busService.responseType, requestUrl.String(), queryParams, &stopSchedule)
}

func (busService *Service) GetStops(request *GetStopsRequest) (*GetStopsResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(busInfoBaseUrl)

	switch busService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jStops")
	case wmata.XML:
		requestUrl.WriteString("/Stops")
	}

	var queryParams map[string]string
	if request != nil {
		queryParams = make(map[string]string)

		if request.Latitude != 0 {
			queryParams["Lat"] = strconv.FormatFloat(request.Latitude, 'g', -1, 64)
		}

		if request.Longitude != 0 {
			queryParams["Lon"] = strconv.FormatFloat(request.Longitude, 'g', -1, 64)
		}

		if request.Radius != 0 {
			queryParams["Radius"] = strconv.FormatFloat(request.Radius, 'g', -1, 64)
		}
	}

	stops := GetStopsResponse{}

	return &stops, busService.client.BuildAndSendGetRequest(busService.responseType, requestUrl.String(), queryParams, &stops)
}
