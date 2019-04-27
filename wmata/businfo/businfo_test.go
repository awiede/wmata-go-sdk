package businfo

import (
	"encoding/xml"
	"errors"
	"fmt"
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

	return nil, fmt.Errorf("no data found for query: %s", req.URL.RawQuery)
}

type testResponseData struct {
	rawQuery             string
	stringParam1         string
	stringParam2         string
	requestType          interface{}
	response             string
	unmarshalledResponse interface{}
}

var testData = map[string][]testResponseData{
	"/Bus.svc/json/jBusPositions": {
		{
			rawQuery: "RouteID=S2",
			requestType: GetPositionsRequest{
				RouteID: "S2",
			},
			response: `{"BusPositions":[{"VehicleID":"7205","Lat":38.894608,"Lon":-77.026634,"Deviation":9,"DateTime":"2019-04-27T18:05:38","TripID":"897053070","RouteID":"S2","DirectionNum":1,"DirectionText":"SOUTH","TripHeadsign":"FEDERAL TRIANGLE","TripStartTime":"2019-04-27T17:07:00","TripEndTime":"2019-04-27T17:58:00","BlockNumber":"NS-11"},{"VehicleID":"7342","Lat":38.952606,"Lon":-77.036331,"Deviation":3,"DateTime":"2019-04-27T18:05:51","TripID":"896961070","RouteID":"S2","DirectionNum":0,"DirectionText":"NORTH","TripHeadsign":"SILVER SPRING STATION","TripStartTime":"2019-04-27T17:24:00","TripEndTime":"2019-04-27T18:19:00","BlockNumber":"NS-04"},{"VehicleID":"7357","Lat":38.921047,"Lon":-77.036545,"Deviation":-1,"DateTime":"2019-04-27T18:05:54","TripID":"897054070","RouteID":"S2","DirectionNum":1,"DirectionText":"SOUTH","TripHeadsign":"FEDERAL TRIANGLE","TripStartTime":"2019-04-27T17:37:00","TripEndTime":"2019-04-27T18:28:00","BlockNumber":"NS-12"},{"VehicleID":"7204","Lat":38.90134,"Lon":-77.0317,"Deviation":1,"DateTime":"2019-04-27T18:05:42","TripID":"896962070","RouteID":"S2","DirectionNum":0,"DirectionText":"NORTH","TripHeadsign":"SILVER SPRING STATION","TripStartTime":"2019-04-27T17:54:00","TripEndTime":"2019-04-27T18:49:00","BlockNumber":"NS-09"}]}`,
			unmarshalledResponse: &GetPositionsResponse{
				BusPositions: []BusPosition{
					{
						VehicleID:       "7205",
						Latitude:        38.894608,
						Longitude:       -77.026634,
						Deviation:       9,
						DateTime:        "2019-04-27T18:05:38",
						TripID:          "897053070",
						RouteID:         "S2",
						DirectionNumber: 1,
						DirectionText:   "SOUTH",
						TripDestination: "FEDERAL TRIANGLE",
						TripStartTime:   "2019-04-27T17:07:00",
						TripEndTime:     "2019-04-27T17:58:00",
						BlockNumber:     "NS-11",
					},
					{
						VehicleID:       "7342",
						Latitude:        38.952606,
						Longitude:       -77.036331,
						Deviation:       3,
						DateTime:        "2019-04-27T18:05:51",
						TripID:          "896961070",
						RouteID:         "S2",
						DirectionNumber: 0,
						DirectionText:   "NORTH",
						TripDestination: "SILVER SPRING STATION",
						TripStartTime:   "2019-04-27T17:24:00",
						TripEndTime:     "2019-04-27T18:19:00",
						BlockNumber:     "NS-04",
					},
					{
						VehicleID:       "7357",
						Latitude:        38.921047,
						Longitude:       -77.036545,
						Deviation:       -1,
						DateTime:        "2019-04-27T18:05:54",
						TripID:          "897054070",
						RouteID:         "S2",
						DirectionNumber: 1,
						DirectionText:   "SOUTH",
						TripDestination: "FEDERAL TRIANGLE",
						TripStartTime:   "2019-04-27T17:37:00",
						TripEndTime:     "2019-04-27T18:28:00",
						BlockNumber:     "NS-12",
					},
					{
						VehicleID:       "7204",
						Latitude:        38.90134,
						Longitude:       -77.0317,
						Deviation:       1,
						DateTime:        "2019-04-27T18:05:42",
						TripID:          "896962070",
						RouteID:         "S2",
						DirectionNumber: 0,
						DirectionText:   "NORTH",
						TripDestination: "SILVER SPRING STATION",
						TripStartTime:   "2019-04-27T17:54:00",
						TripEndTime:     "2019-04-27T18:49:00",
						BlockNumber:     "NS-09",
					},
				},
			},
		},
	},
	"/Bus.svc/BusPositions": {
		{
			rawQuery: "RouteID=G2",
			requestType: GetPositionsRequest{
				RouteID: "G2",
			},
			response: `<BusPositionsResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><BusPositions><BusPosition><BlockNumber>WG-01</BlockNumber><DateTime>2019-04-27T18:16:45</DateTime><Deviation>1</Deviation><DirectionNum>1</DirectionNum><DirectionText>WEST</DirectionText><Lat>38.908882</Lat><Lon>-77.06488</Lon><RouteID>G2</RouteID><TripEndTime>2019-04-27T18:20:00</TripEndTime><TripHeadsign>GEORGETOWN UNIVERSITY</TripHeadsign><TripID>939481070</TripID><TripStartTime>2019-04-27T17:41:00</TripStartTime><VehicleID>3072</VehicleID></BusPosition><BusPosition><BlockNumber>WG-03</BlockNumber><DateTime>2019-04-27T18:16:39</DateTime><Deviation>1</Deviation><DirectionNum>0</DirectionNum><DirectionText>EAST</DirectionText><Lat>38.909653</Lat><Lon>-77.033432</Lon><RouteID>G2</RouteID><TripEndTime>2019-04-27T18:35:00</TripEndTime><TripHeadsign>LEDROIT PARK - HOWARD UNIVERSITY</TripHeadsign><TripID>939526070</TripID><TripStartTime>2019-04-27T17:55:00</TripStartTime><VehicleID>3076</VehicleID></BusPosition><BusPosition><BlockNumber>WG-02</BlockNumber><DateTime>2019-04-27T18:16:21</DateTime><Deviation>5</Deviation><DirectionNum>1</DirectionNum><DirectionText>WEST</DirectionText><Lat>38.909637</Lat><Lon>-77.024284</Lon><RouteID>G2</RouteID><TripEndTime>2019-04-27T18:50:00</TripEndTime><TripHeadsign>GEORGETOWN UNIVERSITY</TripHeadsign><TripID>939482070</TripID><TripStartTime>2019-04-27T18:11:00</TripStartTime><VehicleID>3080</VehicleID></BusPosition></BusPositions></BusPositionsResp>`,
			unmarshalledResponse: &GetPositionsResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "BusPositionsResp",
				},
				BusPositions: []BusPosition{
					{
						BlockNumber: "WG-01",
						DateTime: "2019-04-27T18:16:45",
						Deviation: 1,
						DirectionNumber: 1,
						DirectionText: "WEST",
						Latitude: 38.908882,
						Longitude: -77.06488,
						RouteID: "G2",
						TripEndTime: "2019-04-27T18:20:00",
						TripDestination: "GEORGETOWN UNIVERSITY",
						TripID: "939481070",
						TripStartTime: "2019-04-27T17:41:00",
						VehicleID: "3072",
					},
					{
						BlockNumber: "WG-03",
						DateTime: "2019-04-27T18:16:39",
						Deviation: 1,
						DirectionNumber: 0,
						DirectionText: "EAST",
						Latitude: 38.909653,
						Longitude: -77.033432,
						RouteID: "G2",
						TripEndTime: "2019-04-27T18:35:00",
						TripDestination: "LEDROIT PARK - HOWARD UNIVERSITY",
						TripID: "939526070",
						TripStartTime: "2019-04-27T17:55:00",
						VehicleID: "3076",
					},
					{
						BlockNumber: "WG-02",
						DateTime: "2019-04-27T18:16:21",
						Deviation: 5,
						DirectionNumber: 1,
						DirectionText: "WEST",
						Latitude: 38.909637,
						Longitude: -77.024284,
						RouteID: "G2",
						TripEndTime: "2019-04-27T18:50:00",
						TripDestination: "GEORGETOWN UNIVERSITY",
						TripID: "939482070",
						TripStartTime: "2019-04-27T18:11:00",
						VehicleID: "3080",
					},
				},
			},
		},
	},
}

// setupTestService creates a service struct with a mock http client
func setupTestService(responseType wmata.ResponseType) *Service {
	return &Service{
		client: &wmata.Client{
			HTTPClient: &testClient{},
		},
		responseType: responseType,
	}
}

func TestGetPositions(t *testing.T) {
	jsonAndXmlPaths := []string{"/Bus.svc/json/jBusPositions", "/Bus.svc/BusPositions"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetPositions")
			return
		}

		for _, request := range testRequests {
			getPositionsRequest := request.requestType.(GetPositionsRequest)
			response, err := testService.GetPositions(&getPositionsRequest)

			if err != nil {
				t.Errorf("error calling GetPostions, request: %v error: %s", getPositionsRequest, err.Error())
				continue
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}
