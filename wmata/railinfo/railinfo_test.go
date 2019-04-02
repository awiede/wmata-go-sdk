package railinfo

import (
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestClient is a mock implementation of wmata.HTTPClient interface used for testing purposes
type TestClient struct{}

// ensure TestClient implements wmata.HTTPClient interface
var _ wmata.HTTPClient = (*TestClient)(nil)

// Do stubs out an httpClient.Do request
func (client *TestClient) Do(req *http.Request) (*http.Response, error) {
	testResponses, exist := testData[req.URL.Path]

	if !exist {
		return nil, errors.New("no test data found")
	}

	for _, response := range testResponses {
		if response.rawQuery == req.URL.RawQuery {
			rr := httptest.NewRecorder()
			rr.WriteHeader(http.StatusOK)
			_, writeErr := rr.Write([]byte(response.jsonResponse))

			return rr.Result(), writeErr
		}
	}

	return nil, errors.New("no data found")
}

type testResponseData struct {
	rawQuery             string
	stringParam1         string
	stringParam2         string
	requestType          interface{}
	requestParams        []interface{}
	jsonResponse         string
	unmarshalledResponse interface{}
}

var testData = map[string][]testResponseData{
	"/Rail.svc/json/jLines": {
		testResponseData{
			rawQuery:     "",
			jsonResponse: `{"Lines":[{"LineCode":"BL","DisplayName":"Blue","StartStationCode":"J03","EndStationCode":"G05","InternalDestination1":"","InternalDestination2":""},{"LineCode":"GR","DisplayName":"Green","StartStationCode":"F11","EndStationCode":"E10","InternalDestination1":"","InternalDestination2":""},{"LineCode":"OR","DisplayName":"Orange","StartStationCode":"K08","EndStationCode":"D13","InternalDestination1":"","InternalDestination2":""},{"LineCode":"RD","DisplayName":"Red","StartStationCode":"A15","EndStationCode":"B11","InternalDestination1":"A11","InternalDestination2":"B08"},{"LineCode":"SV","DisplayName":"Silver","StartStationCode":"N06","EndStationCode":"G05","InternalDestination1":"","InternalDestination2":""},{"LineCode":"YL","DisplayName":"Yellow","StartStationCode":"C15","EndStationCode":"E06","InternalDestination1":"E01","InternalDestination2":""}]}`,
			unmarshalledResponse: &GetLinesResponse{
				Lines: []LineResponse{
					{
						LineCode:             "BL",
						DisplayName:          "Blue",
						StartStationCode:     "J03",
						EndStationCode:       "G05",
						InternalDestination1: "",
						InternalDestination2: "",
					},
					{
						LineCode:             "GR",
						DisplayName:          "Green",
						StartStationCode:     "F11",
						EndStationCode:       "E10",
						InternalDestination1: "",
						InternalDestination2: "",
					},
					{
						LineCode:             "OR",
						DisplayName:          "Orange",
						StartStationCode:     "K08",
						EndStationCode:       "D13",
						InternalDestination1: "",
						InternalDestination2: "",
					},
					{
						LineCode:             "RD",
						DisplayName:          "Red",
						StartStationCode:     "A15",
						EndStationCode:       "B11",
						InternalDestination1: "A11",
						InternalDestination2: "B08",
					},
					{
						LineCode:             "SV",
						DisplayName:          "Silver",
						StartStationCode:     "N06",
						EndStationCode:       "G05",
						InternalDestination1: "",
						InternalDestination2: "",
					},
					{
						LineCode:             "YL",
						DisplayName:          "Yellow",
						StartStationCode:     "C15",
						EndStationCode:       "E06",
						InternalDestination1: "E01",
						InternalDestination2: "",
					},
				},
			},
		},
	},
	"/Rail.svc/json/jStationParking": {
		{
			rawQuery:     "StationCode=B08",
			stringParam1: "B08",
			jsonResponse: `{"StationsParking":[{"Code":"B08","Notes":"Parking is available at Montgomery County lots and garages.","AllDayParking":{"TotalCount":0,"RiderCost":null,"NonRiderCost":null,"SaturdayRiderCost":null,"SaturdayNonRiderCost":null},"ShortTermParking":{"TotalCount":0,"Notes":null}}]}`,
			unmarshalledResponse: &GetParkingInformationResponse{
				ParkingInformation: []StationParking{
					{
						StationCode: "B08",
						Notes:       "Parking is available at Montgomery County lots and garages.",
						AllDay: AllDayParking{
							TotalCount:           0,
							RiderCost:            0,
							NonRiderCost:         0,
							SaturdayRiderCost:    0,
							SaturdayNonRiderCost: 0,
						},
						ShortTerm: ShortTermParking{
							TotalCount: 0,
							Notes:      "",
						},
					},
				},
			},
		},
		{
			rawQuery: "StationCode=K08",
			stringParam1: "K08",
			jsonResponse: `{"StationsParking":[{"Code":"K08","Notes":"North Kiss & Ride - 45 short term metered spaces. South Kiss & Ride - 26 short term metered spaces.  101 spaces metered for 12-hr. max @ $1.00 per 60 mins. 17 spaces metered for 7-hr. max. @ $1.00 per 60 mins. Parking available from 8:30 AM to 2 AM.","AllDayParking":{"TotalCount":5169,"RiderCost":4.95,"NonRiderCost":4.95,"SaturdayRiderCost":0,"SaturdayNonRiderCost":0},"ShortTermParking":{"TotalCount":71,"Notes":"Parking available in section B between 8:30 AM - 3:30 PM and 7 PM - 2 AM, in section D between 10 AM - 2 PM."}}]}`,
			unmarshalledResponse: &GetParkingInformationResponse{
				ParkingInformation: []StationParking{
					{
						StationCode: "K08",
						Notes: "North Kiss & Ride - 45 short term metered spaces. South Kiss & Ride - 26 short term metered spaces.  101 spaces metered for 12-hr. max @ $1.00 per 60 mins. 17 spaces metered for 7-hr. max. @ $1.00 per 60 mins. Parking available from 8:30 AM to 2 AM.",
						AllDay: AllDayParking{
							TotalCount: 5169,
							RiderCost: 4.95,
							NonRiderCost: 4.95,
							SaturdayRiderCost: 0,
							SaturdayNonRiderCost: 0,
						},
						ShortTerm: ShortTermParking{
							TotalCount: 71,
							Notes: "Parking available in section B between 8:30 AM - 3:30 PM and 7 PM - 2 AM, in section D between 10 AM - 2 PM.",
						},
					},
				},
			},
		},
	},
}

// setupTestService creates a service struct with a mock http client
func setupTestService() *RailStationInfo {
	return &RailStationInfo{
		client: &wmata.Client{
			HTTPClient: &TestClient{},
		},
	}
}

func TestGetLines(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Rail.svc/json/jLines"]

	if !exist {
		t.Errorf("no data found for GetLines")
		return
	}

	for _, request := range testRequests {
		response, err := testService.GetLines()

		if err != nil {
			t.Errorf("error running GetLines %s", err.Error())
			return
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Errorf("unexpected response. Expected: %v but got: %v", response, request.unmarshalledResponse)
		}
	}

}

func TestGetParkingInformation(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Rail.svc/json/jStationParking"]

	if !exist {
		t.Errorf("no data found for GetLines")
		return
	}

	for _, request := range testRequests {
		response, err := testService.GetParkingInformation(request.stringParam1)

		if err != nil {
			t.Errorf("error running GetParkingInformation for station: %s Error: %s", request.stringParam1, err.Error())
			return
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Errorf("unexpected response. Expected: %v but got: %v", response, request.unmarshalledResponse)
		}
	}

}
