package buspredictions

import (
	"encoding/xml"
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"github.com/kr/pretty"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// testClient is a mock implementation of wmata.HTTPClient interface used for testing purposes
type testClient struct{}

// ensure testClient implements wmata.HTTPClient interface
var _ wmata.HTTPClient = (*testClient)(nil)

// Do stubs out an httpClient.Do request
func (client *testClient) Do(req *http.Request) (*http.Response, error) {
	testResponses, exist := testData[req.URL.Path]

	if !exist {
		return nil, errors.New("no test data found")
	}

	for _, response := range testResponses {
		if response.rawQuery == req.URL.RawQuery {
			rr := httptest.NewRecorder()
			rr.WriteHeader(http.StatusOK)
			_, writeErr := rr.Write([]byte(response.response))

			return rr.Result(), writeErr
		}
	}

	return nil, errors.New("no data found")
}

type testResponseData struct {
	rawQuery             string
	param                string
	response             string
	expectedError        error
	unmarshalledResponse interface{}
}

var testData = map[string][]testResponseData{
	"/NextBusService.svc/json/jPredictions": {
		{
			rawQuery: "StopID=1001370",
			param:    "1001370",
			response: `{"StopName":"37th St Nw + O St Nw","Predictions":[{"RouteID":"G2","DirectionText":"East to Ledroit Park - Howard University","DirectionNum":"0","Minutes":2,"VehicleID":"3072","TripID":"939584010"},{"RouteID":"G2","DirectionText":"East to Ledroit Park - Howard University","DirectionNum":"0","Minutes":38,"VehicleID":"3081","TripID":"939585010"}]}`,
			unmarshalledResponse: &GetNextBusResponse{
				StopName: "37th St Nw + O St Nw",
				NextBusPredictions: []NextBusPrediction{
					{
						RouteID:         "G2",
						DirectionText:   "East to Ledroit Park - Howard University",
						DirectionNumber: "0",
						Minutes:         2,
						VehicleID:       "3072",
						TripID:          "939584010",
					},
					{
						RouteID:         "G2",
						DirectionText:   "East to Ledroit Park - Howard University",
						DirectionNumber: "0",
						Minutes:         38,
						VehicleID:       "3081",
						TripID:          "939585010",
					},
				},
			},
		},
		{
			param:         "",
			expectedError: errors.New("stopID is required"),
		},
	},
	"/NextBusService.svc/Predictions": {
		{
			rawQuery: "StopID=1001370",
			param:    "1001370",
			response: `<NextBusResponse xmlns:i="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.wmata.com"><Predictions><NextBusPrediction><DirectionNum>0</DirectionNum><DirectionText>East to Ledroit Park - Howard University</DirectionText><Minutes>2</Minutes><RouteID>G2</RouteID><TripID>939584010</TripID><VehicleID>3072</VehicleID></NextBusPrediction><NextBusPrediction><DirectionNum>0</DirectionNum><DirectionText>East to Ledroit Park - Howard University</DirectionText><Minutes>38</Minutes><RouteID>G2</RouteID><TripID>939585010</TripID><VehicleID>3081</VehicleID></NextBusPrediction></Predictions><StopName>37th St Nw + O St Nw</StopName></NextBusResponse>`,
			unmarshalledResponse: &GetNextBusResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "NextBusResponse",
				},
				StopName: "37th St Nw + O St Nw",
				NextBusPredictions: []NextBusPrediction{
					{
						RouteID:         "G2",
						DirectionText:   "East to Ledroit Park - Howard University",
						DirectionNumber: "0",
						Minutes:         2,
						VehicleID:       "3072",
						TripID:          "939584010",
					},
					{
						RouteID:         "G2",
						DirectionText:   "East to Ledroit Park - Howard University",
						DirectionNumber: "0",
						Minutes:         38,
						VehicleID:       "3081",
						TripID:          "939585010",
					},
				},
			},
		},
		{
			param:         "",
			expectedError: errors.New("stopID is required"),
		},
	},
}

// setupTestService creates a service struct with a mock http client
func setupTestService(responseType wmata.ResponseType) *Service {
	wmataClient := wmata.Client{
		HTTPClient: &testClient{},
	}
	return NewService(&wmataClient, responseType)
}

func TestGetNextBuses(t *testing.T) {
	jsonAndXmlPaths := []string{"/NextBusService.svc/json/jPredictions", "/NextBusService.svc/Predictions"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetNextBuses")
			return
		}

		for _, request := range testRequests {
			response, responseErr := testService.GetNextBuses(request.param)

			if responseErr != nil {
				if request.expectedError == nil || responseErr.Error() != request.expectedError.Error() {
					t.Errorf("unexpected error: %s", responseErr)
				}
				continue
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}
