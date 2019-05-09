package buspredictions

import (
	"encoding/xml"
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"strings"
)

const busPredictionsServiceBaseURL = "https://api.wmata.com/NextBusService.svc"

type GetNextBusResponse struct {
	XMLName            xml.Name            `json:"-" xml:"http://www.wmata.com NextBusResponse"`
	NextBusPredictions []NextBusPrediction `json:"Predictions" xml:"Predictions>NextBusPrediction"`
	StopName           string              `json:"StopName" xml:"StopName"`
}

type NextBusPrediction struct {
	DirectionNumber string `json:"DirectionNum" xml:"DirectionNum"`
	DirectionText   string `json:"DirectionText" xml:"DirectionText"`
	Minutes         int    `json:"Minutes" xml:"Minutes"`
	RouteID         string `json:"RouteID" xml:"RouteID"`
	TripID          string `json:"TripID" xml:"TripID"`
	VehicleID       string `json:"VehicleID" xml:"VehicleID"`
}

// BusPredictions defines the method available in the WMATA "Real-Time Bus Predictions" API
type BusPredictions interface {
	GetNextBuses(stopID string) (*GetNextBusResponse, error)
}

var _ BusPredictions = (*Service)(nil)

// NewService returns a new Incidents service with a reference to an existing wmata.Client
func NewService(client *wmata.Client, responseType wmata.ResponseType) *Service {
	return &Service{
		client:       client,
		responseType: responseType,
	}
}

// Service provides all API methods for BusPredictions
type Service struct {
	client       *wmata.Client
	responseType wmata.ResponseType
}

// GetNexBuses retrieves next bus arrival times by stopID
// Documentation on service structure can be found here: https://developer.wmata.com/docs/services/5476365e031f590f38092508/operations/5476365e031f5909e4fe331d
func (service *Service) GetNextBuses(stopID string) (*GetNextBusResponse, error) {
	if stopID == "" {
		return nil, errors.New("stopID is required")
	}

	var requestUrl strings.Builder
	requestUrl.WriteString(busPredictionsServiceBaseURL)

	switch service.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jPredictions")
	case wmata.XML:
		requestUrl.WriteString("/Predictions")
	}

	nextBus := GetNextBusResponse{}

	return &nextBus, service.client.BuildAndSendGetRequest(service.responseType, requestUrl.String(), map[string]string{"StopID": stopID}, &nextBus)

}
