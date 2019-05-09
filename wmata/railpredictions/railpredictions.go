package railpredictions

import (
	"encoding/xml"
	"github.com/awiede/wmata-go-sdk/wmata"
	"strings"
)

const railPredictionsServiceBaseURL = "https://api.wmata.com/StationPrediction.svc"

type GetNextTrainResponse struct {
	XMLName xml.Name `json:"-" xml:"http://www.wmata.com AIMPredictionResp"`
	Trains  []Train  `json:"Trains" xml:"Trains>AIMPredictionTrainInfo"`
}

type Train struct {
	Car             string `json:"Car" xml:"Car"`
	Destination     string `json:"Destination" xml:"Destination"`
	DestinationCode string `json:"DestinationCode" xml:"DestinationCode"`
	DestinationName string `json:"DestinationName" xml:"DestinationName"`
	Group           string `json:"Group" xml:"Group"`
	Line            string `json:"Line" xml:"Line"`
	LocationCode    string `json:"LocationCode" xml:"LocationCode"`
	LocationName    string `json:"LocationName" xml:"LocationName"`
	Minutes         string `json:"Min" xml:"Min"`
}

// RailPredictions defines the method available in the WMATA "Real-Time Rail Predictions" API
type RailPredictions interface {
	GetNextTrains(stationCodes []string) (*GetNextTrainResponse, error)
}

var _ RailPredictions = (*Service)(nil)

// NewService returns a new Incidents service with a reference to an existing wmata.Client
func NewService(client *wmata.Client, responseType wmata.ResponseType) *Service {
	return &Service{
		client:       client,
		responseType: responseType,
	}
}

// Service provides all API methods for RailPredictions
type Service struct {
	client       *wmata.Client
	responseType wmata.ResponseType
}

// GetNextTrains retrieves realtime rail predictions for each station code passed.
// If no station codes passed, then all predictions will be retrieved
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/547636a6f9182302184cda78/operations/547636a6f918230da855363f
func (service *Service) GetNextTrains(stationCodes []string) (*GetNextTrainResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(railPredictionsServiceBaseURL)

	switch service.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/GetPrediction")
	case wmata.XML:
		requestUrl.WriteString("/GetPrediction")
	}

	if stationCodes == nil {
		requestUrl.WriteString("/All")

	} else {
		requestUrl.WriteString("/")
		for index, stationCode := range stationCodes {
			if index > 0 {
				requestUrl.WriteString(",")
			}

			requestUrl.WriteString(stationCode)
		}
	}

	nextTrain := GetNextTrainResponse{}

	return &nextTrain, service.client.BuildAndSendGetRequest(service.responseType, requestUrl.String(), nil, &nextTrain)
}
