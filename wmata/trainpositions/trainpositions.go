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
}

type GetTrackCircuitsResponse struct {
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
	panic("implement me")
}

func (service *Service) GetTrackCircuits() (*GetTrackCircuitsResponse, error) {
	panic("implement me")
}
