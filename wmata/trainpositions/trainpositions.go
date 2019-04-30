package trainpositions

import (
	"encoding/xml"
	"github.com/awiede/wmata-go-sdk/wmata"
	"strings"
)

const trainPositionsServiceBaseURL = "https://api.wmata.com/TrainPositions"

type GetLiveTrainPositionsResponse struct {
	XMLName   xml.Name        `json:"-" xml:"http://www.wmata.com TrainPositionResp"`
	Positions []TrainPosition `json:"TrainPositions" xml:"TrainPositions>TrainPosition"`
}

type TrainPosition struct {
	CarCount               int    `json:"CarCount" xml:"CarCount"`
	CircuitID              int    `json:"CircuitId" xml:"CircuitId"`
	DestinationStationCode string `json:"DestinationStationCode" xml:"DestinationStationCode"`
	DirectionNumber        int    `json:"DirectionNum" xml:"DirectionNum"`
	LineCode               string `json:"LineCode" xml:"LineCode"`
	SecondsAtLocation      int    `json:"SecondsAtLocation" xml:"SecondsAtLocation"`
	ServiceType            string `json:"ServiceType" xml:"ServiceType"`
	TrainID                string `json:"TrainId" xml:"TrainId"`
	TrainNumber            string `json:"TrainNumber" xml:"TrainNumber"`
}

type GetStandardRoutesResponse struct {
	XMLName xml.Name `json:"-" xml:"http://www.wmata.com StandardRouteResp"`
	Routes  []Route  `json:"StandardRoutes" xml:"StandardRoutes>StandardRoute"`
}

type Route struct {
	LineCode      string                 `json:"LineCode" xml:"LineCode"`
	TrackNumber   int                    `json:"TrackNum" xml:"TrackNum"`
	TrackCircuits []StandardTrackCircuit `json:"TrackCircuits" xml:"TrackCircuits>TrackCircuit"`
}

type StandardTrackCircuit struct {
	CircuitID      int    `json:"CircuitId" xml:"CircuitId"`
	SequenceNumber int    `json:"SeqNum" xml:"SeqNum"`
	StationCode    string `json:"StationCode" xml:"StationCode"`
}

type GetTrackCircuitsResponse struct {
	XMLName       xml.Name       `json:"-" xml:"http://www.wmata.com TrackCircuitResp"`
	TrackCircuits []TrackCircuit `json:"TrackCircuits" xml:"TrackCircuits>TrackCircuit"`
}

type TrackCircuit struct {
	CircuitID int        `json:"CircuitId" xml:"CircuitId"`
	Track     int        `json:"Track" xml:"Track"`
	Neighbors []Neighbor `json:"Neighbors" xml:"Neighbors>TrackCircuitNeighbor"`
}

type Neighbor struct {
	CircuitIDs   []int  `json:"CircuitIds" xml:"CircuitIds>int"`
	NeighborType string `json:"NeighborType" xml:"NeighborType"`
}

type TrainPositions interface {
	GetLiveTrainPositions() (*GetLiveTrainPositionsResponse, error)
	GetStandardRoutes() (*GetStandardRoutesResponse, error)
	GetTrackCircuits() (*GetTrackCircuitsResponse, error)
}

var _ TrainPositions = (*Service)(nil)

// NewService returns a new Incidents service with a reference to an existing wmata.Client
func NewService(client *wmata.Client, responseType wmata.ResponseType) *Service {
	return &Service{
		client:       client,
		responseType: responseType,
	}
}

type Service struct {
	client       *wmata.Client
	responseType wmata.ResponseType
}

func (service *Service) GetLiveTrainPositions() (*GetLiveTrainPositionsResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(trainPositionsServiceBaseURL)
	requestUrl.WriteString("/TrainPositions")

	queryParams := map[string]string{}
	switch service.responseType {
	case wmata.JSON:
		queryParams["contentType"] = "json"
	case wmata.XML:
		queryParams["contentType"] = "xml"
	}

	livePositions := GetLiveTrainPositionsResponse{}

	return &livePositions, service.client.BuildAndSendGetRequest(service.responseType, requestUrl.String(), queryParams, &livePositions)
}

func (service *Service) GetStandardRoutes() (*GetStandardRoutesResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(trainPositionsServiceBaseURL)
	requestUrl.WriteString("/StandardRoutes")

	queryParams := map[string]string{}
	switch service.responseType {
	case wmata.JSON:
		queryParams["contentType"] = "json"
	case wmata.XML:
		queryParams["contentType"] = "xml"
	}

	routes := GetStandardRoutesResponse{}

	return &routes, service.client.BuildAndSendGetRequest(service.responseType, requestUrl.String(), queryParams, &routes)
}

func (service *Service) GetTrackCircuits() (*GetTrackCircuitsResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(trainPositionsServiceBaseURL)
	requestUrl.WriteString("/TrackCircuits")

	queryParams := map[string]string{}
	switch service.responseType {
	case wmata.JSON:
		queryParams["contentType"] = "json"
	case wmata.XML:
		queryParams["contentType"] = "xml"
	}

	circuits := GetTrackCircuitsResponse{}

	return &circuits, service.client.BuildAndSendGetRequest(service.responseType, requestUrl.String(), queryParams, &circuits)
}
