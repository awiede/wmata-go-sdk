package businfo

import (
	"encoding/xml"
	"github.com/awiede/wmata-go-sdk/wmata"
)

const busInfoBaseUrl = "https://api.wmata.com/Bus.svc"

type BusInfo interface {
	GetPositions(request *GetPositionsRequest) (*GetPositionsResponse, error)
	GetPathDetails(routeID, date string) (*GetPathDetailsResponse, error)
	GetRoutes() (*GetRoutesResponse, error)
	GetSchedule(routeID, date string, includeVariations bool) (*GetScheduleResponse, error)
	GetScheduleAtStop(stopID, date string) (*GetScheduleAtStopResponse, error)
	GetStopsResponse(request *GetStopsRequest) (*GetStopsResponse, error)
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
	routeID   string
	latitude  float64
	longitude float64
	radius    float64
}

type GetPositionsResponse struct {
	XMLName      xml.Name      `xml:"http://www.wmata.com BusPositionResp"`
	BusPositions []BusPosition `json:"BusPositions" xml:"BusPositions>BusPosition"`
}

type BusPosition struct {
	DateTime  string `json:"DateTime" xml:"DateTime"`
	Deviation string `json:"Deviation" xml:"Deviation"`
	// Deprecated: DirectionNumber response field is deprecated, use DirectionText
	DirectionNumber string `json:"DirectionNum" xml:"DirectionNum"`
	DirectionText string `json:"DirectionText" xml:"DirectionText"`
	Latitude

}

type GetPathDetailsResponse struct {
}

type GetRoutesResponse struct {
}

type GetScheduleResponse struct {
}

type GetScheduleAtStopResponse struct {
}

type GetStopsRequest struct {
}

type GetStopsResponse struct {
}

func (busService *Service) GetPositions(request *GetPositionsRequest) (*GetPositionsResponse, error) {
	panic("implement me")
}

func (busService *Service) GetPathDetails(routeID, date string) (*GetPathDetailsResponse, error) {
	panic("implement me")
}

func (busService *Service) GetRoutes() (*GetRoutesResponse, error) {
	panic("implement me")
}

func (busService *Service) GetSchedule(routeID, date string, includeVariations bool) (*GetScheduleResponse, error) {
	panic("implement me")
}

func (busService *Service) GetScheduleAtStop(stopID, date string) (*GetScheduleAtStopResponse, error) {
	panic("implement me")
}

func (busService *Service) GetStopsResponse(request *GetStopsRequest) (*GetStopsResponse, error) {
	panic("implement me")
}
