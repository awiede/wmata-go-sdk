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
						BlockNumber:     "WG-01",
						DateTime:        "2019-04-27T18:16:45",
						Deviation:       1,
						DirectionNumber: 1,
						DirectionText:   "WEST",
						Latitude:        38.908882,
						Longitude:       -77.06488,
						RouteID:         "G2",
						TripEndTime:     "2019-04-27T18:20:00",
						TripDestination: "GEORGETOWN UNIVERSITY",
						TripID:          "939481070",
						TripStartTime:   "2019-04-27T17:41:00",
						VehicleID:       "3072",
					},
					{
						BlockNumber:     "WG-03",
						DateTime:        "2019-04-27T18:16:39",
						Deviation:       1,
						DirectionNumber: 0,
						DirectionText:   "EAST",
						Latitude:        38.909653,
						Longitude:       -77.033432,
						RouteID:         "G2",
						TripEndTime:     "2019-04-27T18:35:00",
						TripDestination: "LEDROIT PARK - HOWARD UNIVERSITY",
						TripID:          "939526070",
						TripStartTime:   "2019-04-27T17:55:00",
						VehicleID:       "3076",
					},
					{
						BlockNumber:     "WG-02",
						DateTime:        "2019-04-27T18:16:21",
						Deviation:       5,
						DirectionNumber: 1,
						DirectionText:   "WEST",
						Latitude:        38.909637,
						Longitude:       -77.024284,
						RouteID:         "G2",
						TripEndTime:     "2019-04-27T18:50:00",
						TripDestination: "GEORGETOWN UNIVERSITY",
						TripID:          "939482070",
						TripStartTime:   "2019-04-27T18:11:00",
						VehicleID:       "3080",
					},
				},
			},
		},
	},
	"/Bus.svc/json/jRouteDetails": {
		{
			rawQuery:     "RouteID=G2",
			stringParam1: "G2",
			response:     `{"RouteID":"G2","Name":"G2 - HOWARD UNIV - GEORGETOWN UNIV","Direction0":{"TripHeadsign":"LEDROIT PARK - HOWARD UNIVERSITY","DirectionText":"EAST","DirectionNum":"0","Shape":[{"Lat":38.907367,"Lon":-77.071641,"SeqNum":1},{"Lat":38.90599,"Lon":-77.071515,"SeqNum":2},{"Lat":38.905889,"Lon":-77.071411,"SeqNum":3},{"Lat":38.905965,"Lon":-77.069018,"SeqNum":4},{"Lat":38.907698,"Lon":-77.069012,"SeqNum":5},{"Lat":38.907781,"Lon":-77.06886,"SeqNum":6},{"Lat":38.907803,"Lon":-77.068501,"SeqNum":7},{"Lat":38.907882,"Lon":-77.063647,"SeqNum":8},{"Lat":38.907674,"Lon":-77.0634,"SeqNum":9},{"Lat":38.907654,"Lon":-77.063209,"SeqNum":10},{"Lat":38.907527,"Lon":-77.063125,"SeqNum":11},{"Lat":38.907597,"Lon":-77.062759,"SeqNum":12},{"Lat":38.907547,"Lon":-77.059455,"SeqNum":13},{"Lat":38.907654,"Lon":-77.059196,"SeqNum":14},{"Lat":38.907674,"Lon":-77.057202,"SeqNum":15},{"Lat":38.907958,"Lon":-77.05699,"SeqNum":16},{"Lat":38.90926,"Lon":-77.056984,"SeqNum":17},{"Lat":38.909348,"Lon":-77.056327,"SeqNum":18},{"Lat":38.909367,"Lon":-77.053894,"SeqNum":19},{"Lat":38.909587,"Lon":-77.052527,"SeqNum":20},{"Lat":38.909537,"Lon":-77.049446,"SeqNum":21},{"Lat":38.909618,"Lon":-77.049248,"SeqNum":22},{"Lat":38.909633,"Lon":-77.048614,"SeqNum":23},{"Lat":38.909649,"Lon":-77.044439,"SeqNum":24},{"Lat":38.909559,"Lon":-77.04424,"SeqNum":25},{"Lat":38.90914,"Lon":-77.043909,"SeqNum":26},{"Lat":38.909019,"Lon":-77.04343,"SeqNum":27},{"Lat":38.90914,"Lon":-77.04298,"SeqNum":28},{"Lat":38.909419,"Lon":-77.04271,"SeqNum":29},{"Lat":38.909529,"Lon":-77.04267,"SeqNum":30},{"Lat":38.909619,"Lon":-77.042379,"SeqNum":31},{"Lat":38.909621,"Lon":-77.039001,"SeqNum":32},{"Lat":38.909573,"Lon":-77.038826,"SeqNum":33},{"Lat":38.90956,"Lon":-77.03682,"SeqNum":34},{"Lat":38.90961,"Lon":-77.036773,"SeqNum":35},{"Lat":38.909621,"Lon":-77.036124,"SeqNum":36},{"Lat":38.909618,"Lon":-77.032348,"SeqNum":37},{"Lat":38.909694,"Lon":-77.031632,"SeqNum":38},{"Lat":38.909652,"Lon":-77.030617,"SeqNum":39},{"Lat":38.909596,"Lon":-77.030411,"SeqNum":40},{"Lat":38.909229,"Lon":-77.030013,"SeqNum":41},{"Lat":38.909145,"Lon":-77.029456,"SeqNum":42},{"Lat":38.909289,"Lon":-77.029067,"SeqNum":43},{"Lat":38.909641,"Lon":-77.028724,"SeqNum":44},{"Lat":38.90963,"Lon":-77.016279,"SeqNum":45},{"Lat":38.910004,"Lon":-77.016434,"SeqNum":46},{"Lat":38.910038,"Lon":-77.016356,"SeqNum":47},{"Lat":38.910397,"Lon":-77.016575,"SeqNum":48},{"Lat":38.914099,"Lon":-77.017889,"SeqNum":49},{"Lat":38.914149,"Lon":-77.017319,"SeqNum":50},{"Lat":38.91424,"Lon":-77.017129,"SeqNum":51},{"Lat":38.91577,"Lon":-77.0173,"SeqNum":52},{"Lat":38.91848,"Lon":-77.017759,"SeqNum":53},{"Lat":38.91859,"Lon":-77.017691,"SeqNum":54},{"Lat":38.919297,"Lon":-77.017877,"SeqNum":55},{"Lat":38.919404,"Lon":-77.01779,"SeqNum":56},{"Lat":38.919674,"Lon":-77.015208,"SeqNum":57},{"Lat":38.91966,"Lon":-77.014799,"SeqNum":58},{"Lat":38.920501,"Lon":-77.014905,"SeqNum":59},{"Lat":38.920541,"Lon":-77.014838,"SeqNum":60},{"Lat":38.92067,"Lon":-77.01493,"SeqNum":61},{"Lat":38.920617,"Lon":-77.015862,"SeqNum":62}],"Stops":[{"StopID":"1001370","Name":"37TH ST NW + O ST NW","Lon":-77.071671,"Lat":38.907395,"Routes":["G2","G2v1"]},{"StopID":"1001345","Name":"PROSPECT ST NW + 36TH ST NW","Lon":-77.069828,"Lat":38.905833,"Routes":["G2","G2v1"]},{"StopID":"1001354","Name":"35TH ST NW + N ST NW","Lon":-77.06888,"Lat":38.906649,"Routes":["G2","G2v1"]},{"StopID":"1001385","Name":"O ST NW + 34TH ST NW","Lon":-77.067972,"Lat":38.9077,"Routes":["G2","G2v1"]},{"StopID":"1001379","Name":"DUMBARTON ST NW + WISCONSIN AVE NW","Lon":-77.063119,"Lat":38.90757,"Routes":["G2","G2v1"]},{"StopID":"1001381","Name":"DUMBARTON ST NW + 30TH ST NW","Lon":-77.059336,"Lat":38.907603,"Routes":["G2","G2v1"]},{"StopID":"1001389","Name":"28TH ST NW + DUMBARTON ST NW","Lon":-77.056995,"Lat":38.907989,"Routes":["G2","G2v1"]},{"StopID":"1001412","Name":"28TH ST NW + P ST NW","Lon":-77.057012,"Lat":38.909259,"Routes":["G2","G2v1"]},{"StopID":"1001418","Name":"P ST NW + 26TH ST NW","Lon":-77.054884,"Lat":38.909313,"Routes":["G2","G2v1"]},{"StopID":"1001435","Name":"P ST NW + ROCK CREEK & POTOMAC PKWY RAMP","Lon":-77.052,"Lat":38.909596,"Routes":["G2","G2v1"]},{"StopID":"1001433","Name":"P ST NW + 23RD ST NW","Lon":-77.049559,"Lat":38.909571,"Routes":["G2","G2v1"]},{"StopID":"1001440","Name":"P ST NW + 21ST ST NW","Lon":-77.046793,"Lat":38.909583,"Routes":["G2","G2v1"]},{"StopID":"1001439","Name":"P ST NW + 20TH ST NW","Lon":-77.045051,"Lat":38.909582,"Routes":["G2","G2v1"]},{"StopID":"1001436","Name":"P ST NW + 18TH ST NW","Lon":-77.041944,"Lat":38.909582,"Routes":["G2","G2v1"]},{"StopID":"1001434","Name":"P ST NW + 17TH ST NW","Lon":-77.038752,"Lat":38.909565,"Routes":["G2","G2v1"]},{"StopID":"1001438","Name":"P ST NW + 16TH ST NW","Lon":-77.036766,"Lat":38.909576,"Routes":["G2","G2v1"]},{"StopID":"1001437","Name":"P ST NW + 15TH ST NW","Lon":-77.034773,"Lat":38.909576,"Routes":["G2","G2v1"]},{"StopID":"1001452","Name":"P ST NW + 14TH ST NW","Lon":-77.032215,"Lat":38.909597,"Routes":["G2","G2v1"]},{"StopID":"1001456","Name":"P ST NW + LOGAN CIR NW","Lon":-77.03054,"Lat":38.909596,"Routes":["G2"]},{"StopID":"1001416","Name":"LOGAN CIR NW + P ST NW","Lon":-77.0289,"Lat":38.909349,"Routes":["63","G2"]},{"StopID":"1001455","Name":"P ST NW + 11TH ST NW","Lon":-77.02733,"Lat":38.909594,"Routes":["G2"]},{"StopID":"1001442","Name":"P ST NW + 9TH ST NW","Lon":-77.024196,"Lat":38.909599,"Routes":["G2"]},{"StopID":"1001445","Name":"P ST NW + 7TH ST NW","Lon":-77.021287,"Lat":38.909593,"Routes":["G2"]},{"StopID":"1001444","Name":"P ST NW + 5TH ST NW","Lon":-77.019082,"Lat":38.9096,"Routes":["G2"]},{"StopID":"1001472","Name":"NEW JERSEY AVE NW + P ST NW","Lon":-77.016348,"Lat":38.910066,"Routes":["96","96v1","96v2","96v3","G2"]},{"StopID":"1001544","Name":"NEW JERSEY AVE NW + R ST NW","Lon":-77.017218,"Lat":38.912503,"Routes":["96","96v1","96v2","96v3","G2"]},{"StopID":"1001583","Name":"NEW JERSEY AVE NW + RHODE ISLAND AVE NW","Lon":-77.017538,"Lat":38.913357,"Routes":["96","96v1","96v2","96v3","G2"]},{"StopID":"1003839","Name":"4TH ST NW + T ST NW","Lon":-77.017246,"Lat":38.915794,"Routes":["G2"]},{"StopID":"1001720","Name":"4TH ST NW + V ST NW","Lon":-77.017709,"Lat":38.918686,"Routes":["G2"]},{"StopID":"1001745","Name":"W ST NW + 4TH ST NW","Lon":-77.017086,"Lat":38.919395,"Routes":["G2"]},{"StopID":"1001748","Name":"W ST NW + 2ND ST NW","Lon":-77.014934,"Lat":38.919607,"Routes":["G2"]},{"StopID":"1001764","Name":"2ND ST NW + BRYANT ST NW","Lon":-77.014838,"Lat":38.920541,"Routes":["G2"]},{"StopID":"1001776","Name":"BRYANT ST NW + #301","Lon":-77.016037,"Lat":38.920612,"Routes":["G2"]}]},"Direction1":{"TripHeadsign":"GEORGETOWN UNIVERSITY","DirectionText":"WEST","DirectionNum":"1","Shape":[{"Lat":38.920625,"Lon":-77.016035,"SeqNum":1},{"Lat":38.920433,"Lon":-77.017861,"SeqNum":2},{"Lat":38.920332,"Lon":-77.018029,"SeqNum":3},{"Lat":38.918199,"Lon":-77.017782,"SeqNum":4},{"Lat":38.914495,"Lon":-77.017155,"SeqNum":5},{"Lat":38.914214,"Lon":-77.017227,"SeqNum":6},{"Lat":38.914107,"Lon":-77.017409,"SeqNum":7},{"Lat":38.914096,"Lon":-77.017837,"SeqNum":8},{"Lat":38.91396,"Lon":-77.01791,"SeqNum":9},{"Lat":38.913196,"Lon":-77.017679,"SeqNum":10},{"Lat":38.913333,"Lon":-77.017662,"SeqNum":11},{"Lat":38.912373,"Lon":-77.017364,"SeqNum":12},{"Lat":38.911257,"Lon":-77.016845,"SeqNum":13},{"Lat":38.909804,"Lon":-77.016418,"SeqNum":14},{"Lat":38.90963,"Lon":-77.016279,"SeqNum":15},{"Lat":38.909652,"Lon":-77.023735,"SeqNum":16},{"Lat":38.909599,"Lon":-77.026504,"SeqNum":17},{"Lat":38.909655,"Lon":-77.026779,"SeqNum":18},{"Lat":38.909672,"Lon":-77.028808,"SeqNum":19},{"Lat":38.910117,"Lon":-77.029335,"SeqNum":20},{"Lat":38.910187,"Lon":-77.029747,"SeqNum":21},{"Lat":38.91008,"Lon":-77.03012,"SeqNum":22},{"Lat":38.909927,"Lon":-77.030281,"SeqNum":23},{"Lat":38.909663,"Lon":-77.030456,"SeqNum":24},{"Lat":38.90964,"Lon":-77.042036,"SeqNum":25},{"Lat":38.909768,"Lon":-77.042291,"SeqNum":26},{"Lat":38.90968,"Lon":-77.042529,"SeqNum":27},{"Lat":38.910069,"Lon":-77.04286,"SeqNum":28},{"Lat":38.91026,"Lon":-77.043289,"SeqNum":29},{"Lat":38.910209,"Lon":-77.043769,"SeqNum":30},{"Lat":38.90997,"Lon":-77.04411,"SeqNum":31},{"Lat":38.9097,"Lon":-77.04425,"SeqNum":32},{"Lat":38.909649,"Lon":-77.044439,"SeqNum":33},{"Lat":38.909702,"Lon":-77.044663,"SeqNum":34},{"Lat":38.909644,"Lon":-77.048644,"SeqNum":35},{"Lat":38.90966,"Lon":-77.05225,"SeqNum":36},{"Lat":38.909725,"Lon":-77.052244,"SeqNum":37},{"Lat":38.90966,"Lon":-77.05225,"SeqNum":38},{"Lat":38.90966,"Lon":-77.0528,"SeqNum":39},{"Lat":38.909399,"Lon":-77.05395,"SeqNum":40},{"Lat":38.90937,"Lon":-77.05544,"SeqNum":41},{"Lat":38.909407,"Lon":-77.055442,"SeqNum":42},{"Lat":38.909328,"Lon":-77.057304,"SeqNum":43},{"Lat":38.909252,"Lon":-77.063934,"SeqNum":44},{"Lat":38.909303,"Lon":-77.064203,"SeqNum":45},{"Lat":38.908973,"Lon":-77.064211,"SeqNum":46},{"Lat":38.90889,"Lon":-77.064369,"SeqNum":47},{"Lat":38.908656,"Lon":-77.07161,"SeqNum":48},{"Lat":38.907421,"Lon":-77.071633,"SeqNum":49}],"Stops":[{"StopID":"1001776","Name":"BRYANT ST NW + #301","Lon":-77.016037,"Lat":38.920612,"Routes":["G2"]},{"StopID":"1001737","Name":"4TH ST NW + W ST NW","Lon":-77.017909,"Lat":38.91914,"Routes":["G2"]},{"StopID":"1003840","Name":"4TH ST NW + U ST NW","Lon":-77.01757,"Lat":38.916976,"Routes":["G2"]},{"StopID":"1001584","Name":"NEW JERSEY AVE NW + RHODE ISLAND AVE NW","Lon":-77.017662,"Lat":38.913167,"Routes":["96","96v1","96v4","96v5","G2"]},{"StopID":"1001541","Name":"NEW JERSEY AVE NW + R ST NW","Lon":-77.017378,"Lat":38.912372,"Routes":["96","96v1","96v4","96v5","G2"]},{"StopID":"1001465","Name":"NEW JERSEY AVE NW + P ST NW","Lon":-77.016463,"Lat":38.909815,"Routes":["96","96v1","96v4","96v5","G2"]},{"StopID":"1001441","Name":"P ST NW + 5TH ST NW","Lon":-77.018773,"Lat":38.909704,"Routes":["G2"]},{"StopID":"1001453","Name":"P ST NW + 7TH ST NW","Lon":-77.021658,"Lat":38.909709,"Routes":["G2"]},{"StopID":"1001454","Name":"P ST NW + 9TH ST NW","Lon":-77.023737,"Lat":38.909706,"Routes":["G2"]},{"StopID":"1001458","Name":"P ST NW + 11TH ST NW","Lon":-77.026814,"Lat":38.909717,"Routes":["G2"]},{"StopID":"1001451","Name":"P ST NW + LOGAN CIR NW","Lon":-77.028662,"Lat":38.909706,"Routes":["G2"]},{"StopID":"1001469","Name":"LOGAN CIR NW + 13TH ST NW","Lon":-77.030154,"Lat":38.910125,"Routes":["63","G2"]},{"StopID":"1001446","Name":"P ST NW + 14TH ST NW","Lon":-77.032411,"Lat":38.909726,"Routes":["G2"]},{"StopID":"1001457","Name":"P ST NW + 16TH ST NW","Lon":-77.03636,"Lat":38.909707,"Routes":["G2"]},{"StopID":"1001443","Name":"P ST NW + 17TH ST NW","Lon":-77.038234,"Lat":38.909723,"Routes":["G2"]},{"StopID":"1001459","Name":"P ST NW + DUPONT CIR NW","Lon":-77.042182,"Lat":38.909723,"Routes":["G2","N2","N4","N4v1","N6"]},{"StopID":"1001461","Name":"P ST NW + 20TH ST NW","Lon":-77.044753,"Lat":38.90972,"Routes":["D2","D2v1","G2"]},{"StopID":"1001449","Name":"P ST NW + 22ND ST NW","Lon":-77.048666,"Lat":38.909711,"Routes":["D1","D2","D2v1","D6","D6v3","G2"]},{"StopID":"1001463","Name":"P ST NW + ROCK CREEK & POTOMAC PKWY RAMP","Lon":-77.052244,"Lat":38.909725,"Routes":["G2"]},{"StopID":"1001424","Name":"P ST NW + 27TH ST NW","Lon":-77.055408,"Lat":38.909439,"Routes":["G2"]},{"StopID":"1001421","Name":"P ST NW + 28TH ST NW","Lon":-77.056951,"Lat":38.909424,"Routes":["G2"]},{"StopID":"1001420","Name":"P ST NW + 29TH ST NW","Lon":-77.057996,"Lat":38.909407,"Routes":["G2"]},{"StopID":"1001419","Name":"P ST NW + 30TH ST NW","Lon":-77.059097,"Lat":38.909379,"Routes":["G2"]},{"StopID":"1001417","Name":"P ST NW + 31ST ST NW","Lon":-77.061284,"Lat":38.909351,"Routes":["G2"]},{"StopID":"1001415","Name":"P ST NW + WISCONSIN AVE NW","Lon":-77.064005,"Lat":38.909307,"Routes":["G2"]},{"StopID":"1001401","Name":"P ST NW + 33RD ST NW","Lon":-77.065912,"Lat":38.908847,"Routes":["G2"]},{"StopID":"1001398","Name":"P ST NW + 35TH ST NW","Lon":-77.068904,"Lat":38.908794,"Routes":["G2"]},{"StopID":"1001370","Name":"37TH ST NW + O ST NW","Lon":-77.071671,"Lat":38.907395,"Routes":["G2","G2v1"]}]}}`,
			unmarshalledResponse: &GetRouteDetailsResponse{
				RouteID: "G2",
				Name:    "G2 - HOWARD UNIV - GEORGETOWN UNIV",
				Direction0: Direction{
					TripDestination: "LEDROIT PARK - HOWARD UNIVERSITY",
					DirectionText:   "EAST",
					DirectionNumber: "0",
					Shapes: []ShapePoint{
						{
							Latitude:       38.907367,
							Longitude:      -77.071641,
							SequenceNumber: 1,
						},
						{
							Latitude:       38.90599,
							Longitude:      -77.071515,
							SequenceNumber: 2,
						},
						{
							Latitude:       38.905889,
							Longitude:      -77.071411,
							SequenceNumber: 3,
						},
						{
							Latitude:       38.905965,
							Longitude:      -77.069018,
							SequenceNumber: 4,
						},
						{
							Latitude:       38.907698,
							Longitude:      -77.069012,
							SequenceNumber: 5,
						},
						{
							Latitude:       38.907781,
							Longitude:      -77.06886,
							SequenceNumber: 6,
						},
						{
							Latitude:       38.907803,
							Longitude:      -77.068501,
							SequenceNumber: 7,
						},
						{
							Latitude:       38.907882,
							Longitude:      -77.063647,
							SequenceNumber: 8,
						},
						{
							Latitude:       38.907674,
							Longitude:      -77.0634,
							SequenceNumber: 9,
						},
						{
							Latitude:       38.907654,
							Longitude:      -77.063209,
							SequenceNumber: 10,
						},
						{
							Latitude:       38.907527,
							Longitude:      -77.063125,
							SequenceNumber: 11,
						},
						{
							Latitude:       38.907597,
							Longitude:      -77.062759,
							SequenceNumber: 12,
						},
						{
							Latitude:       38.907547,
							Longitude:      -77.059455,
							SequenceNumber: 13,
						},
						{
							Latitude:       38.907654,
							Longitude:      -77.059196,
							SequenceNumber: 14,
						},
						{
							Latitude:       38.907674,
							Longitude:      -77.057202,
							SequenceNumber: 15,
						},
						{
							Latitude:       38.907958,
							Longitude:      -77.05699,
							SequenceNumber: 16,
						},
						{
							Latitude:       38.90926,
							Longitude:      -77.056984,
							SequenceNumber: 17,
						},
						{
							Latitude:       38.909348,
							Longitude:      -77.056327,
							SequenceNumber: 18,
						},
						{
							Latitude:       38.909367,
							Longitude:      -77.053894,
							SequenceNumber: 19,
						},
						{
							Latitude:       38.909587,
							Longitude:      -77.052527,
							SequenceNumber: 20,
						},
						{
							Latitude:       38.909537,
							Longitude:      -77.049446,
							SequenceNumber: 21,
						},
						{
							Latitude:       38.909618,
							Longitude:      -77.049248,
							SequenceNumber: 22,
						},
						{
							Latitude:       38.909633,
							Longitude:      -77.048614,
							SequenceNumber: 23,
						},
						{
							Latitude:       38.909649,
							Longitude:      -77.044439,
							SequenceNumber: 24,
						},
						{
							Latitude:       38.909559,
							Longitude:      -77.04424,
							SequenceNumber: 25,
						},
						{
							Latitude:       38.90914,
							Longitude:      -77.043909,
							SequenceNumber: 26,
						},
						{
							Latitude:       38.909019,
							Longitude:      -77.04343,
							SequenceNumber: 27,
						},
						{
							Latitude:       38.90914,
							Longitude:      -77.04298,
							SequenceNumber: 28,
						},
						{
							Latitude:       38.909419,
							Longitude:      -77.04271,
							SequenceNumber: 29,
						},
						{
							Latitude:       38.909529,
							Longitude:      -77.04267,
							SequenceNumber: 30,
						},
						{
							Latitude:       38.909619,
							Longitude:      -77.042379,
							SequenceNumber: 31,
						},
						{
							Latitude:       38.909621,
							Longitude:      -77.039001,
							SequenceNumber: 32,
						},
						{
							Latitude:       38.909573,
							Longitude:      -77.038826,
							SequenceNumber: 33,
						},
						{
							Latitude:       38.90956,
							Longitude:      -77.03682,
							SequenceNumber: 34,
						},
						{
							Latitude:       38.90961,
							Longitude:      -77.036773,
							SequenceNumber: 35,
						},
						{
							Latitude:       38.909621,
							Longitude:      -77.036124,
							SequenceNumber: 36,
						},
						{
							Latitude:       38.909618,
							Longitude:      -77.032348,
							SequenceNumber: 37,
						},
						{
							Latitude:       38.909694,
							Longitude:      -77.031632,
							SequenceNumber: 38,
						},
						{
							Latitude:       38.909652,
							Longitude:      -77.030617,
							SequenceNumber: 39,
						},
						{
							Latitude:       38.909596,
							Longitude:      -77.030411,
							SequenceNumber: 40,
						},
						{
							Latitude:       38.909229,
							Longitude:      -77.030013,
							SequenceNumber: 41,
						},
						{
							Latitude:       38.909145,
							Longitude:      -77.029456,
							SequenceNumber: 42,
						},
						{
							Latitude:       38.909289,
							Longitude:      -77.029067,
							SequenceNumber: 43,
						},
						{
							Latitude:       38.909641,
							Longitude:      -77.028724,
							SequenceNumber: 44,
						},
						{
							Latitude:       38.90963,
							Longitude:      -77.016279,
							SequenceNumber: 45,
						},
						{
							Latitude:       38.910004,
							Longitude:      -77.016434,
							SequenceNumber: 46,
						},
						{
							Latitude:       38.910038,
							Longitude:      -77.016356,
							SequenceNumber: 47,
						},
						{
							Latitude:       38.910397,
							Longitude:      -77.016575,
							SequenceNumber: 48,
						},
						{
							Latitude:       38.914099,
							Longitude:      -77.017889,
							SequenceNumber: 49,
						},
						{
							Latitude:       38.914149,
							Longitude:      -77.017319,
							SequenceNumber: 50,
						},
						{
							Latitude:       38.91424,
							Longitude:      -77.017129,
							SequenceNumber: 51,
						},
						{
							Latitude:       38.91577,
							Longitude:      -77.0173,
							SequenceNumber: 52,
						},
						{
							Latitude:       38.91848,
							Longitude:      -77.017759,
							SequenceNumber: 53,
						},
						{
							Latitude:       38.91859,
							Longitude:      -77.017691,
							SequenceNumber: 54,
						},
						{
							Latitude:       38.919297,
							Longitude:      -77.017877,
							SequenceNumber: 55,
						},
						{
							Latitude:       38.919404,
							Longitude:      -77.01779,
							SequenceNumber: 56,
						},
						{
							Latitude:       38.919674,
							Longitude:      -77.015208,
							SequenceNumber: 57,
						},
						{
							Latitude:       38.91966,
							Longitude:      -77.014799,
							SequenceNumber: 58,
						},
						{
							Latitude:       38.920501,
							Longitude:      -77.014905,
							SequenceNumber: 59,
						},
						{
							Latitude:       38.920541,
							Longitude:      -77.014838,
							SequenceNumber: 60,
						},
						{
							Latitude:       38.92067,
							Longitude:      -77.01493,
							SequenceNumber: 61,
						},
						{
							Latitude:       38.920617,
							Longitude:      -77.015862,
							SequenceNumber: 62,
						},
					},
					Stops: []Stop{
						{
							StopID:    "1001370",
							Name:      "37TH ST NW + O ST NW",
							Longitude: -77.071671,
							Latitude:  38.907395,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001345",
							Name:      "PROSPECT ST NW + 36TH ST NW",
							Longitude: -77.069828,
							Latitude:  38.905833,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001354",
							Name:      "35TH ST NW + N ST NW",
							Longitude: -77.06888,
							Latitude:  38.906649,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001385",
							Name:      "O ST NW + 34TH ST NW",
							Longitude: -77.067972,
							Latitude:  38.9077,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001379",
							Name:      "DUMBARTON ST NW + WISCONSIN AVE NW",
							Longitude: -77.063119,
							Latitude:  38.90757,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001381",
							Name:      "DUMBARTON ST NW + 30TH ST NW",
							Longitude: -77.059336,
							Latitude:  38.907603,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001389",
							Name:      "28TH ST NW + DUMBARTON ST NW",
							Longitude: -77.056995,
							Latitude:  38.907989,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001412",
							Name:      "28TH ST NW + P ST NW",
							Longitude: -77.057012,
							Latitude:  38.909259,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001418",
							Name:      "P ST NW + 26TH ST NW",
							Longitude: -77.054884,
							Latitude:  38.909313,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001435",
							Name:      "P ST NW + ROCK CREEK & POTOMAC PKWY RAMP",
							Longitude: -77.052,
							Latitude:  38.909596,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001433",
							Name:      "P ST NW + 23RD ST NW",
							Longitude: -77.049559,
							Latitude:  38.909571,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001440",
							Name:      "P ST NW + 21ST ST NW",
							Longitude: -77.046793,
							Latitude:  38.909583,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001439",
							Name:      "P ST NW + 20TH ST NW",
							Longitude: -77.045051,
							Latitude:  38.909582,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001436",
							Name:      "P ST NW + 18TH ST NW",
							Longitude: -77.041944,
							Latitude:  38.909582,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001434",
							Name:      "P ST NW + 17TH ST NW",
							Longitude: -77.038752,
							Latitude:  38.909565,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001438",
							Name:      "P ST NW + 16TH ST NW",
							Longitude: -77.036766,
							Latitude:  38.909576,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001437",
							Name:      "P ST NW + 15TH ST NW",
							Longitude: -77.034773,
							Latitude:  38.909576,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001452",
							Name:      "P ST NW + 14TH ST NW",
							Longitude: -77.032215,
							Latitude:  38.909597,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001456",
							Name:      "P ST NW + LOGAN CIR NW",
							Longitude: -77.03054,
							Latitude:  38.909596,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001416",
							Name:      "LOGAN CIR NW + P ST NW",
							Longitude: -77.0289,
							Latitude:  38.909349,
							Routes: []string{
								"63",
								"G2",
							},
						},
						{
							StopID:    "1001455",
							Name:      "P ST NW + 11TH ST NW",
							Longitude: -77.02733,
							Latitude:  38.909594,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001442",
							Name:      "P ST NW + 9TH ST NW",
							Longitude: -77.024196,
							Latitude:  38.909599,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001445",
							Name:      "P ST NW + 7TH ST NW",
							Longitude: -77.021287,
							Latitude:  38.909593,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001444",
							Name:      "P ST NW + 5TH ST NW",
							Longitude: -77.019082,
							Latitude:  38.9096,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001472",
							Name:      "NEW JERSEY AVE NW + P ST NW",
							Longitude: -77.016348,
							Latitude:  38.910066,
							Routes: []string{
								"96",
								"96v1",
								"96v2",
								"96v3",
								"G2",
							},
						},
						{
							StopID:    "1001544",
							Name:      "NEW JERSEY AVE NW + R ST NW",
							Longitude: -77.017218,
							Latitude:  38.912503,
							Routes: []string{
								"96",
								"96v1",
								"96v2",
								"96v3",
								"G2",
							},
						},
						{
							StopID:    "1001583",
							Name:      "NEW JERSEY AVE NW + RHODE ISLAND AVE NW",
							Longitude: -77.017538,
							Latitude:  38.913357,
							Routes: []string{
								"96",
								"96v1",
								"96v2",
								"96v3",
								"G2",
							},
						},
						{
							StopID:    "1003839",
							Name:      "4TH ST NW + T ST NW",
							Longitude: -77.017246,
							Latitude:  38.915794,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001720",
							Name:      "4TH ST NW + V ST NW",
							Longitude: -77.017709,
							Latitude:  38.918686,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001745",
							Name:      "W ST NW + 4TH ST NW",
							Longitude: -77.017086,
							Latitude:  38.919395,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001748",
							Name:      "W ST NW + 2ND ST NW",
							Longitude: -77.014934,
							Latitude:  38.919607,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001764",
							Name:      "2ND ST NW + BRYANT ST NW",
							Longitude: -77.014838,
							Latitude:  38.920541,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001776",
							Name:      "BRYANT ST NW + #301",
							Longitude: -77.016037,
							Latitude:  38.920612,
							Routes: []string{
								"G2",
							},
						},
					},
				},
				Direction1: Direction{
					TripDestination: "GEORGETOWN UNIVERSITY",
					DirectionText:   "WEST",
					DirectionNumber: "1",
					Shapes: []ShapePoint{
						{
							Latitude:       38.920625,
							Longitude:      -77.016035,
							SequenceNumber: 1,
						},
						{
							Latitude:       38.920433,
							Longitude:      -77.017861,
							SequenceNumber: 2,
						},
						{
							Latitude:       38.920332,
							Longitude:      -77.018029,
							SequenceNumber: 3,
						},
						{
							Latitude:       38.918199,
							Longitude:      -77.017782,
							SequenceNumber: 4,
						},
						{
							Latitude:       38.914495,
							Longitude:      -77.017155,
							SequenceNumber: 5,
						},
						{
							Latitude:       38.914214,
							Longitude:      -77.017227,
							SequenceNumber: 6,
						},
						{
							Latitude:       38.914107,
							Longitude:      -77.017409,
							SequenceNumber: 7,
						},
						{
							Latitude:       38.914096,
							Longitude:      -77.017837,
							SequenceNumber: 8,
						},
						{
							Latitude:       38.91396,
							Longitude:      -77.01791,
							SequenceNumber: 9,
						},
						{
							Latitude:       38.913196,
							Longitude:      -77.017679,
							SequenceNumber: 10,
						},
						{
							Latitude:       38.913333,
							Longitude:      -77.017662,
							SequenceNumber: 11,
						},
						{
							Latitude:       38.912373,
							Longitude:      -77.017364,
							SequenceNumber: 12,
						},
						{
							Latitude:       38.911257,
							Longitude:      -77.016845,
							SequenceNumber: 13,
						},
						{
							Latitude:       38.909804,
							Longitude:      -77.016418,
							SequenceNumber: 14,
						},
						{
							Latitude:       38.90963,
							Longitude:      -77.016279,
							SequenceNumber: 15,
						},
						{
							Latitude:       38.909652,
							Longitude:      -77.023735,
							SequenceNumber: 16,
						},
						{
							Latitude:       38.909599,
							Longitude:      -77.026504,
							SequenceNumber: 17,
						},
						{
							Latitude:       38.909655,
							Longitude:      -77.026779,
							SequenceNumber: 18,
						},
						{
							Latitude:       38.909672,
							Longitude:      -77.028808,
							SequenceNumber: 19,
						},
						{
							Latitude:       38.910117,
							Longitude:      -77.029335,
							SequenceNumber: 20,
						},
						{
							Latitude:       38.910187,
							Longitude:      -77.029747,
							SequenceNumber: 21,
						},
						{
							Latitude:       38.91008,
							Longitude:      -77.03012,
							SequenceNumber: 22,
						},
						{
							Latitude:       38.909927,
							Longitude:      -77.030281,
							SequenceNumber: 23,
						},
						{
							Latitude:       38.909663,
							Longitude:      -77.030456,
							SequenceNumber: 24,
						},
						{
							Latitude:       38.90964,
							Longitude:      -77.042036,
							SequenceNumber: 25,
						},
						{
							Latitude:       38.909768,
							Longitude:      -77.042291,
							SequenceNumber: 26,
						},
						{
							Latitude:       38.90968,
							Longitude:      -77.042529,
							SequenceNumber: 27,
						},
						{
							Latitude:       38.910069,
							Longitude:      -77.04286,
							SequenceNumber: 28,
						},
						{
							Latitude:       38.91026,
							Longitude:      -77.043289,
							SequenceNumber: 29,
						},
						{
							Latitude:       38.910209,
							Longitude:      -77.043769,
							SequenceNumber: 30,
						},
						{
							Latitude:       38.90997,
							Longitude:      -77.04411,
							SequenceNumber: 31,
						},
						{
							Latitude:       38.9097,
							Longitude:      -77.04425,
							SequenceNumber: 32,
						},
						{
							Latitude:       38.909649,
							Longitude:      -77.044439,
							SequenceNumber: 33,
						},
						{
							Latitude:       38.909702,
							Longitude:      -77.044663,
							SequenceNumber: 34,
						},
						{
							Latitude:       38.909644,
							Longitude:      -77.048644,
							SequenceNumber: 35,
						},
						{
							Latitude:       38.90966,
							Longitude:      -77.05225,
							SequenceNumber: 36,
						},
						{
							Latitude:       38.909725,
							Longitude:      -77.052244,
							SequenceNumber: 37,
						},
						{
							Latitude:       38.90966,
							Longitude:      -77.05225,
							SequenceNumber: 38,
						},
						{
							Latitude:       38.90966,
							Longitude:      -77.0528,
							SequenceNumber: 39,
						},
						{
							Latitude:       38.909399,
							Longitude:      -77.05395,
							SequenceNumber: 40,
						},
						{
							Latitude:       38.90937,
							Longitude:      -77.05544,
							SequenceNumber: 41,
						},
						{
							Latitude:       38.909407,
							Longitude:      -77.055442,
							SequenceNumber: 42,
						},
						{
							Latitude:       38.909328,
							Longitude:      -77.057304,
							SequenceNumber: 43,
						},
						{
							Latitude:       38.909252,
							Longitude:      -77.063934,
							SequenceNumber: 44,
						},
						{
							Latitude:       38.909303,
							Longitude:      -77.064203,
							SequenceNumber: 45,
						},
						{
							Latitude:       38.908973,
							Longitude:      -77.064211,
							SequenceNumber: 46,
						},
						{
							Latitude:       38.90889,
							Longitude:      -77.064369,
							SequenceNumber: 47,
						},
						{
							Latitude:       38.908656,
							Longitude:      -77.07161,
							SequenceNumber: 48,
						},
						{
							Latitude:       38.907421,
							Longitude:      -77.071633,
							SequenceNumber: 49,
						},
					},
					Stops: []Stop{
						{
							StopID:    "1001776",
							Name:      "BRYANT ST NW + #301",
							Longitude: -77.016037,
							Latitude:  38.920612,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001737",
							Name:      "4TH ST NW + W ST NW",
							Longitude: -77.017909,
							Latitude:  38.91914,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1003840",
							Name:      "4TH ST NW + U ST NW",
							Longitude: -77.01757,
							Latitude:  38.916976,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001584",
							Name:      "NEW JERSEY AVE NW + RHODE ISLAND AVE NW",
							Longitude: -77.017662,
							Latitude:  38.913167,
							Routes: []string{
								"96",
								"96v1",
								"96v4",
								"96v5",
								"G2",
							},
						},
						{
							StopID:    "1001541",
							Name:      "NEW JERSEY AVE NW + R ST NW",
							Longitude: -77.017378,
							Latitude:  38.912372,
							Routes: []string{
								"96",
								"96v1",
								"96v4",
								"96v5",
								"G2",
							},
						},
						{
							StopID:    "1001465",
							Name:      "NEW JERSEY AVE NW + P ST NW",
							Longitude: -77.016463,
							Latitude:  38.909815,
							Routes: []string{
								"96",
								"96v1",
								"96v4",
								"96v5",
								"G2",
							},
						},
						{
							StopID:    "1001441",
							Name:      "P ST NW + 5TH ST NW",
							Longitude: -77.018773,
							Latitude:  38.909704,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001453",
							Name:      "P ST NW + 7TH ST NW",
							Longitude: -77.021658,
							Latitude:  38.909709,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001454",
							Name:      "P ST NW + 9TH ST NW",
							Longitude: -77.023737,
							Latitude:  38.909706,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001458",
							Name:      "P ST NW + 11TH ST NW",
							Longitude: -77.026814,
							Latitude:  38.909717,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001451",
							Name:      "P ST NW + LOGAN CIR NW",
							Longitude: -77.028662,
							Latitude:  38.909706,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001469",
							Name:      "LOGAN CIR NW + 13TH ST NW",
							Longitude: -77.030154,
							Latitude:  38.910125,
							Routes: []string{
								"63",
								"G2",
							},
						},
						{
							StopID:    "1001446",
							Name:      "P ST NW + 14TH ST NW",
							Longitude: -77.032411,
							Latitude:  38.909726,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001457",
							Name:      "P ST NW + 16TH ST NW",
							Longitude: -77.03636,
							Latitude:  38.909707,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001443",
							Name:      "P ST NW + 17TH ST NW",
							Longitude: -77.038234,
							Latitude:  38.909723,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001459",
							Name:      "P ST NW + DUPONT CIR NW",
							Longitude: -77.042182,
							Latitude:  38.909723,
							Routes: []string{
								"G2",
								"N2",
								"N4",
								"N4v1",
								"N6",
							},
						},
						{
							StopID:    "1001461",
							Name:      "P ST NW + 20TH ST NW",
							Longitude: -77.044753,
							Latitude:  38.90972,
							Routes: []string{
								"D2",
								"D2v1",
								"G2",
							},
						},
						{
							StopID:    "1001449",
							Name:      "P ST NW + 22ND ST NW",
							Longitude: -77.048666,
							Latitude:  38.909711,
							Routes: []string{
								"D1",
								"D2",
								"D2v1",
								"D6",
								"D6v3",
								"G2",
							},
						},
						{
							StopID:    "1001463",
							Name:      "P ST NW + ROCK CREEK & POTOMAC PKWY RAMP",
							Longitude: -77.052244,
							Latitude:  38.909725,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001424",
							Name:      "P ST NW + 27TH ST NW",
							Longitude: -77.055408,
							Latitude:  38.909439,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001421",
							Name:      "P ST NW + 28TH ST NW",
							Longitude: -77.056951,
							Latitude:  38.909424,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001420",
							Name:      "P ST NW + 29TH ST NW",
							Longitude: -77.057996,
							Latitude:  38.909407,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001419",
							Name:      "P ST NW + 30TH ST NW",
							Longitude: -77.059097,
							Latitude:  38.909379,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001417",
							Name:      "P ST NW + 31ST ST NW",
							Longitude: -77.061284,
							Latitude:  38.909351,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001415",
							Name:      "P ST NW + WISCONSIN AVE NW",
							Longitude: -77.064005,
							Latitude:  38.909307,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001401",
							Name:      "P ST NW + 33RD ST NW",
							Longitude: -77.065912,
							Latitude:  38.908847,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001398",
							Name:      "P ST NW + 35TH ST NW",
							Longitude: -77.068904,
							Latitude:  38.908794,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001370",
							Name:      "37TH ST NW + O ST NW",
							Longitude: -77.071671,
							Latitude:  38.907395,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
					},
				},
			},
		},
	},
	"/Bus.svc/RouteDetails": {
		{
			rawQuery:     "RouteID=G2",
			stringParam1: "G2",
			response:     `<RouteDetailsInfo xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Direction0><DirectionNum>0</DirectionNum><DirectionText>EAST</DirectionText><Shape><ShapePoint><Lat>38.907367</Lat><Lon>-77.071641</Lon><SeqNum>1</SeqNum></ShapePoint><ShapePoint><Lat>38.905990</Lat><Lon>-77.071515</Lon><SeqNum>2</SeqNum></ShapePoint><ShapePoint><Lat>38.905889</Lat><Lon>-77.071411</Lon><SeqNum>3</SeqNum></ShapePoint><ShapePoint><Lat>38.905965</Lat><Lon>-77.069018</Lon><SeqNum>4</SeqNum></ShapePoint><ShapePoint><Lat>38.907698</Lat><Lon>-77.069012</Lon><SeqNum>5</SeqNum></ShapePoint><ShapePoint><Lat>38.907781</Lat><Lon>-77.068860</Lon><SeqNum>6</SeqNum></ShapePoint><ShapePoint><Lat>38.907803</Lat><Lon>-77.068501</Lon><SeqNum>7</SeqNum></ShapePoint><ShapePoint><Lat>38.907882</Lat><Lon>-77.063647</Lon><SeqNum>8</SeqNum></ShapePoint><ShapePoint><Lat>38.907674</Lat><Lon>-77.063400</Lon><SeqNum>9</SeqNum></ShapePoint><ShapePoint><Lat>38.907654</Lat><Lon>-77.063209</Lon><SeqNum>10</SeqNum></ShapePoint><ShapePoint><Lat>38.907527</Lat><Lon>-77.063125</Lon><SeqNum>11</SeqNum></ShapePoint><ShapePoint><Lat>38.907597</Lat><Lon>-77.062759</Lon><SeqNum>12</SeqNum></ShapePoint><ShapePoint><Lat>38.907547</Lat><Lon>-77.059455</Lon><SeqNum>13</SeqNum></ShapePoint><ShapePoint><Lat>38.907654</Lat><Lon>-77.059196</Lon><SeqNum>14</SeqNum></ShapePoint><ShapePoint><Lat>38.907674</Lat><Lon>-77.057202</Lon><SeqNum>15</SeqNum></ShapePoint><ShapePoint><Lat>38.907958</Lat><Lon>-77.056990</Lon><SeqNum>16</SeqNum></ShapePoint><ShapePoint><Lat>38.909260</Lat><Lon>-77.056984</Lon><SeqNum>17</SeqNum></ShapePoint><ShapePoint><Lat>38.909348</Lat><Lon>-77.056327</Lon><SeqNum>18</SeqNum></ShapePoint><ShapePoint><Lat>38.909367</Lat><Lon>-77.053894</Lon><SeqNum>19</SeqNum></ShapePoint><ShapePoint><Lat>38.909587</Lat><Lon>-77.052527</Lon><SeqNum>20</SeqNum></ShapePoint><ShapePoint><Lat>38.909537</Lat><Lon>-77.049446</Lon><SeqNum>21</SeqNum></ShapePoint><ShapePoint><Lat>38.909618</Lat><Lon>-77.049248</Lon><SeqNum>22</SeqNum></ShapePoint><ShapePoint><Lat>38.909633</Lat><Lon>-77.048614</Lon><SeqNum>23</SeqNum></ShapePoint><ShapePoint><Lat>38.909649</Lat><Lon>-77.044439</Lon><SeqNum>24</SeqNum></ShapePoint><ShapePoint><Lat>38.909559</Lat><Lon>-77.044240</Lon><SeqNum>25</SeqNum></ShapePoint><ShapePoint><Lat>38.909140</Lat><Lon>-77.043909</Lon><SeqNum>26</SeqNum></ShapePoint><ShapePoint><Lat>38.909019</Lat><Lon>-77.043430</Lon><SeqNum>27</SeqNum></ShapePoint><ShapePoint><Lat>38.909140</Lat><Lon>-77.042980</Lon><SeqNum>28</SeqNum></ShapePoint><ShapePoint><Lat>38.909419</Lat><Lon>-77.042710</Lon><SeqNum>29</SeqNum></ShapePoint><ShapePoint><Lat>38.909529</Lat><Lon>-77.042670</Lon><SeqNum>30</SeqNum></ShapePoint><ShapePoint><Lat>38.909619</Lat><Lon>-77.042379</Lon><SeqNum>31</SeqNum></ShapePoint><ShapePoint><Lat>38.909621</Lat><Lon>-77.039001</Lon><SeqNum>32</SeqNum></ShapePoint><ShapePoint><Lat>38.909573</Lat><Lon>-77.038826</Lon><SeqNum>33</SeqNum></ShapePoint><ShapePoint><Lat>38.909560</Lat><Lon>-77.036820</Lon><SeqNum>34</SeqNum></ShapePoint><ShapePoint><Lat>38.909610</Lat><Lon>-77.036773</Lon><SeqNum>35</SeqNum></ShapePoint><ShapePoint><Lat>38.909621</Lat><Lon>-77.036124</Lon><SeqNum>36</SeqNum></ShapePoint><ShapePoint><Lat>38.909618</Lat><Lon>-77.032348</Lon><SeqNum>37</SeqNum></ShapePoint><ShapePoint><Lat>38.909694</Lat><Lon>-77.031632</Lon><SeqNum>38</SeqNum></ShapePoint><ShapePoint><Lat>38.909652</Lat><Lon>-77.030617</Lon><SeqNum>39</SeqNum></ShapePoint><ShapePoint><Lat>38.909596</Lat><Lon>-77.030411</Lon><SeqNum>40</SeqNum></ShapePoint><ShapePoint><Lat>38.909229</Lat><Lon>-77.030013</Lon><SeqNum>41</SeqNum></ShapePoint><ShapePoint><Lat>38.909145</Lat><Lon>-77.029456</Lon><SeqNum>42</SeqNum></ShapePoint><ShapePoint><Lat>38.909289</Lat><Lon>-77.029067</Lon><SeqNum>43</SeqNum></ShapePoint><ShapePoint><Lat>38.909641</Lat><Lon>-77.028724</Lon><SeqNum>44</SeqNum></ShapePoint><ShapePoint><Lat>38.909630</Lat><Lon>-77.016279</Lon><SeqNum>45</SeqNum></ShapePoint><ShapePoint><Lat>38.910004</Lat><Lon>-77.016434</Lon><SeqNum>46</SeqNum></ShapePoint><ShapePoint><Lat>38.910038</Lat><Lon>-77.016356</Lon><SeqNum>47</SeqNum></ShapePoint><ShapePoint><Lat>38.910397</Lat><Lon>-77.016575</Lon><SeqNum>48</SeqNum></ShapePoint><ShapePoint><Lat>38.914099</Lat><Lon>-77.017889</Lon><SeqNum>49</SeqNum></ShapePoint><ShapePoint><Lat>38.914149</Lat><Lon>-77.017319</Lon><SeqNum>50</SeqNum></ShapePoint><ShapePoint><Lat>38.914240</Lat><Lon>-77.017129</Lon><SeqNum>51</SeqNum></ShapePoint><ShapePoint><Lat>38.915770</Lat><Lon>-77.017300</Lon><SeqNum>52</SeqNum></ShapePoint><ShapePoint><Lat>38.918480</Lat><Lon>-77.017759</Lon><SeqNum>53</SeqNum></ShapePoint><ShapePoint><Lat>38.918590</Lat><Lon>-77.017691</Lon><SeqNum>54</SeqNum></ShapePoint><ShapePoint><Lat>38.919297</Lat><Lon>-77.017877</Lon><SeqNum>55</SeqNum></ShapePoint><ShapePoint><Lat>38.919404</Lat><Lon>-77.017790</Lon><SeqNum>56</SeqNum></ShapePoint><ShapePoint><Lat>38.919674</Lat><Lon>-77.015208</Lon><SeqNum>57</SeqNum></ShapePoint><ShapePoint><Lat>38.919660</Lat><Lon>-77.014799</Lon><SeqNum>58</SeqNum></ShapePoint><ShapePoint><Lat>38.920501</Lat><Lon>-77.014905</Lon><SeqNum>59</SeqNum></ShapePoint><ShapePoint><Lat>38.920541</Lat><Lon>-77.014838</Lon><SeqNum>60</SeqNum></ShapePoint><ShapePoint><Lat>38.920670</Lat><Lon>-77.014930</Lon><SeqNum>61</SeqNum></ShapePoint><ShapePoint><Lat>38.920617</Lat><Lon>-77.015862</Lon><SeqNum>62</SeqNum></ShapePoint></Shape><Stops><Stop><Lat>38.907395</Lat><Lon>-77.071671</Lon><Name>37TH ST NW + O ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001370</StopID></Stop><Stop><Lat>38.905833</Lat><Lon>-77.069828</Lon><Name>PROSPECT ST NW + 36TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001345</StopID></Stop><Stop><Lat>38.906649</Lat><Lon>-77.068880</Lon><Name>35TH ST NW + N ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001354</StopID></Stop><Stop><Lat>38.907700</Lat><Lon>-77.067972</Lon><Name>O ST NW + 34TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001385</StopID></Stop><Stop><Lat>38.907570</Lat><Lon>-77.063119</Lon><Name>DUMBARTON ST NW + WISCONSIN AVE NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001379</StopID></Stop><Stop><Lat>38.907603</Lat><Lon>-77.059336</Lon><Name>DUMBARTON ST NW + 30TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001381</StopID></Stop><Stop><Lat>38.907989</Lat><Lon>-77.056995</Lon><Name>28TH ST NW + DUMBARTON ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001389</StopID></Stop><Stop><Lat>38.909259</Lat><Lon>-77.057012</Lon><Name>28TH ST NW + P ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001412</StopID></Stop><Stop><Lat>38.909313</Lat><Lon>-77.054884</Lon><Name>P ST NW + 26TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001418</StopID></Stop><Stop><Lat>38.909596</Lat><Lon>-77.052000</Lon><Name>P ST NW + ROCK CREEK &amp; POTOMAC PKWY RAMP</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001435</StopID></Stop><Stop><Lat>38.909571</Lat><Lon>-77.049559</Lon><Name>P ST NW + 23RD ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001433</StopID></Stop><Stop><Lat>38.909583</Lat><Lon>-77.046793</Lon><Name>P ST NW + 21ST ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001440</StopID></Stop><Stop><Lat>38.909582</Lat><Lon>-77.045051</Lon><Name>P ST NW + 20TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001439</StopID></Stop><Stop><Lat>38.909582</Lat><Lon>-77.041944</Lon><Name>P ST NW + 18TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001436</StopID></Stop><Stop><Lat>38.909565</Lat><Lon>-77.038752</Lon><Name>P ST NW + 17TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001434</StopID></Stop><Stop><Lat>38.909576</Lat><Lon>-77.036766</Lon><Name>P ST NW + 16TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001438</StopID></Stop><Stop><Lat>38.909576</Lat><Lon>-77.034773</Lon><Name>P ST NW + 15TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001437</StopID></Stop><Stop><Lat>38.909597</Lat><Lon>-77.032215</Lon><Name>P ST NW + 14TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001452</StopID></Stop><Stop><Lat>38.909596</Lat><Lon>-77.030540</Lon><Name>P ST NW + LOGAN CIR NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001456</StopID></Stop><Stop><Lat>38.909349</Lat><Lon>-77.028900</Lon><Name>LOGAN CIR NW + P ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>63</a:string><a:string>G2</a:string></Routes><StopID>1001416</StopID></Stop><Stop><Lat>38.909594</Lat><Lon>-77.027330</Lon><Name>P ST NW + 11TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001455</StopID></Stop><Stop><Lat>38.909599</Lat><Lon>-77.024196</Lon><Name>P ST NW + 9TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001442</StopID></Stop><Stop><Lat>38.909593</Lat><Lon>-77.021287</Lon><Name>P ST NW + 7TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001445</StopID></Stop><Stop><Lat>38.909600</Lat><Lon>-77.019082</Lon><Name>P ST NW + 5TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001444</StopID></Stop><Stop><Lat>38.910066</Lat><Lon>-77.016348</Lon><Name>NEW JERSEY AVE NW + P ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>96</a:string><a:string>96v1</a:string><a:string>96v2</a:string><a:string>96v3</a:string><a:string>G2</a:string></Routes><StopID>1001472</StopID></Stop><Stop><Lat>38.912503</Lat><Lon>-77.017218</Lon><Name>NEW JERSEY AVE NW + R ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>96</a:string><a:string>96v1</a:string><a:string>96v2</a:string><a:string>96v3</a:string><a:string>G2</a:string></Routes><StopID>1001544</StopID></Stop><Stop><Lat>38.913357</Lat><Lon>-77.017538</Lon><Name>NEW JERSEY AVE NW + RHODE ISLAND AVE NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>96</a:string><a:string>96v1</a:string><a:string>96v2</a:string><a:string>96v3</a:string><a:string>G2</a:string></Routes><StopID>1001583</StopID></Stop><Stop><Lat>38.915794</Lat><Lon>-77.017246</Lon><Name>4TH ST NW + T ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1003839</StopID></Stop><Stop><Lat>38.918686</Lat><Lon>-77.017709</Lon><Name>4TH ST NW + V ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001720</StopID></Stop><Stop><Lat>38.919395</Lat><Lon>-77.017086</Lon><Name>W ST NW + 4TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001745</StopID></Stop><Stop><Lat>38.919607</Lat><Lon>-77.014934</Lon><Name>W ST NW + 2ND ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001748</StopID></Stop><Stop><Lat>38.920541</Lat><Lon>-77.014838</Lon><Name>2ND ST NW + BRYANT ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001764</StopID></Stop><Stop><Lat>38.920612</Lat><Lon>-77.016037</Lon><Name>BRYANT ST NW + #301</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001776</StopID></Stop></Stops><TripHeadsign>LEDROIT PARK - HOWARD UNIVERSITY</TripHeadsign></Direction0><Direction1><DirectionNum>1</DirectionNum><DirectionText>WEST</DirectionText><Shape><ShapePoint><Lat>38.920625</Lat><Lon>-77.016035</Lon><SeqNum>1</SeqNum></ShapePoint><ShapePoint><Lat>38.920433</Lat><Lon>-77.017861</Lon><SeqNum>2</SeqNum></ShapePoint><ShapePoint><Lat>38.920332</Lat><Lon>-77.018029</Lon><SeqNum>3</SeqNum></ShapePoint><ShapePoint><Lat>38.918199</Lat><Lon>-77.017782</Lon><SeqNum>4</SeqNum></ShapePoint><ShapePoint><Lat>38.914495</Lat><Lon>-77.017155</Lon><SeqNum>5</SeqNum></ShapePoint><ShapePoint><Lat>38.914214</Lat><Lon>-77.017227</Lon><SeqNum>6</SeqNum></ShapePoint><ShapePoint><Lat>38.914107</Lat><Lon>-77.017409</Lon><SeqNum>7</SeqNum></ShapePoint><ShapePoint><Lat>38.914096</Lat><Lon>-77.017837</Lon><SeqNum>8</SeqNum></ShapePoint><ShapePoint><Lat>38.913960</Lat><Lon>-77.017910</Lon><SeqNum>9</SeqNum></ShapePoint><ShapePoint><Lat>38.913196</Lat><Lon>-77.017679</Lon><SeqNum>10</SeqNum></ShapePoint><ShapePoint><Lat>38.913333</Lat><Lon>-77.017662</Lon><SeqNum>11</SeqNum></ShapePoint><ShapePoint><Lat>38.912373</Lat><Lon>-77.017364</Lon><SeqNum>12</SeqNum></ShapePoint><ShapePoint><Lat>38.911257</Lat><Lon>-77.016845</Lon><SeqNum>13</SeqNum></ShapePoint><ShapePoint><Lat>38.909804</Lat><Lon>-77.016418</Lon><SeqNum>14</SeqNum></ShapePoint><ShapePoint><Lat>38.909630</Lat><Lon>-77.016279</Lon><SeqNum>15</SeqNum></ShapePoint><ShapePoint><Lat>38.909652</Lat><Lon>-77.023735</Lon><SeqNum>16</SeqNum></ShapePoint><ShapePoint><Lat>38.909599</Lat><Lon>-77.026504</Lon><SeqNum>17</SeqNum></ShapePoint><ShapePoint><Lat>38.909655</Lat><Lon>-77.026779</Lon><SeqNum>18</SeqNum></ShapePoint><ShapePoint><Lat>38.909672</Lat><Lon>-77.028808</Lon><SeqNum>19</SeqNum></ShapePoint><ShapePoint><Lat>38.910117</Lat><Lon>-77.029335</Lon><SeqNum>20</SeqNum></ShapePoint><ShapePoint><Lat>38.910187</Lat><Lon>-77.029747</Lon><SeqNum>21</SeqNum></ShapePoint><ShapePoint><Lat>38.910080</Lat><Lon>-77.030120</Lon><SeqNum>22</SeqNum></ShapePoint><ShapePoint><Lat>38.909927</Lat><Lon>-77.030281</Lon><SeqNum>23</SeqNum></ShapePoint><ShapePoint><Lat>38.909663</Lat><Lon>-77.030456</Lon><SeqNum>24</SeqNum></ShapePoint><ShapePoint><Lat>38.909640</Lat><Lon>-77.042036</Lon><SeqNum>25</SeqNum></ShapePoint><ShapePoint><Lat>38.909768</Lat><Lon>-77.042291</Lon><SeqNum>26</SeqNum></ShapePoint><ShapePoint><Lat>38.909680</Lat><Lon>-77.042529</Lon><SeqNum>27</SeqNum></ShapePoint><ShapePoint><Lat>38.910069</Lat><Lon>-77.042860</Lon><SeqNum>28</SeqNum></ShapePoint><ShapePoint><Lat>38.910260</Lat><Lon>-77.043289</Lon><SeqNum>29</SeqNum></ShapePoint><ShapePoint><Lat>38.910209</Lat><Lon>-77.043769</Lon><SeqNum>30</SeqNum></ShapePoint><ShapePoint><Lat>38.909970</Lat><Lon>-77.044110</Lon><SeqNum>31</SeqNum></ShapePoint><ShapePoint><Lat>38.909700</Lat><Lon>-77.044250</Lon><SeqNum>32</SeqNum></ShapePoint><ShapePoint><Lat>38.909649</Lat><Lon>-77.044439</Lon><SeqNum>33</SeqNum></ShapePoint><ShapePoint><Lat>38.909702</Lat><Lon>-77.044663</Lon><SeqNum>34</SeqNum></ShapePoint><ShapePoint><Lat>38.909644</Lat><Lon>-77.048644</Lon><SeqNum>35</SeqNum></ShapePoint><ShapePoint><Lat>38.909660</Lat><Lon>-77.052250</Lon><SeqNum>36</SeqNum></ShapePoint><ShapePoint><Lat>38.909725</Lat><Lon>-77.052244</Lon><SeqNum>37</SeqNum></ShapePoint><ShapePoint><Lat>38.909660</Lat><Lon>-77.052250</Lon><SeqNum>38</SeqNum></ShapePoint><ShapePoint><Lat>38.909660</Lat><Lon>-77.052800</Lon><SeqNum>39</SeqNum></ShapePoint><ShapePoint><Lat>38.909399</Lat><Lon>-77.053950</Lon><SeqNum>40</SeqNum></ShapePoint><ShapePoint><Lat>38.909370</Lat><Lon>-77.055440</Lon><SeqNum>41</SeqNum></ShapePoint><ShapePoint><Lat>38.909407</Lat><Lon>-77.055442</Lon><SeqNum>42</SeqNum></ShapePoint><ShapePoint><Lat>38.909328</Lat><Lon>-77.057304</Lon><SeqNum>43</SeqNum></ShapePoint><ShapePoint><Lat>38.909252</Lat><Lon>-77.063934</Lon><SeqNum>44</SeqNum></ShapePoint><ShapePoint><Lat>38.909303</Lat><Lon>-77.064203</Lon><SeqNum>45</SeqNum></ShapePoint><ShapePoint><Lat>38.908973</Lat><Lon>-77.064211</Lon><SeqNum>46</SeqNum></ShapePoint><ShapePoint><Lat>38.908890</Lat><Lon>-77.064369</Lon><SeqNum>47</SeqNum></ShapePoint><ShapePoint><Lat>38.908656</Lat><Lon>-77.071610</Lon><SeqNum>48</SeqNum></ShapePoint><ShapePoint><Lat>38.907421</Lat><Lon>-77.071633</Lon><SeqNum>49</SeqNum></ShapePoint></Shape><Stops><Stop><Lat>38.920612</Lat><Lon>-77.016037</Lon><Name>BRYANT ST NW + #301</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001776</StopID></Stop><Stop><Lat>38.919140</Lat><Lon>-77.017909</Lon><Name>4TH ST NW + W ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001737</StopID></Stop><Stop><Lat>38.916976</Lat><Lon>-77.017570</Lon><Name>4TH ST NW + U ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1003840</StopID></Stop><Stop><Lat>38.913167</Lat><Lon>-77.017662</Lon><Name>NEW JERSEY AVE NW + RHODE ISLAND AVE NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>96</a:string><a:string>96v1</a:string><a:string>96v4</a:string><a:string>96v5</a:string><a:string>G2</a:string></Routes><StopID>1001584</StopID></Stop><Stop><Lat>38.912372</Lat><Lon>-77.017378</Lon><Name>NEW JERSEY AVE NW + R ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>96</a:string><a:string>96v1</a:string><a:string>96v4</a:string><a:string>96v5</a:string><a:string>G2</a:string></Routes><StopID>1001541</StopID></Stop><Stop><Lat>38.909815</Lat><Lon>-77.016463</Lon><Name>NEW JERSEY AVE NW + P ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>96</a:string><a:string>96v1</a:string><a:string>96v4</a:string><a:string>96v5</a:string><a:string>G2</a:string></Routes><StopID>1001465</StopID></Stop><Stop><Lat>38.909704</Lat><Lon>-77.018773</Lon><Name>P ST NW + 5TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001441</StopID></Stop><Stop><Lat>38.909709</Lat><Lon>-77.021658</Lon><Name>P ST NW + 7TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001453</StopID></Stop><Stop><Lat>38.909706</Lat><Lon>-77.023737</Lon><Name>P ST NW + 9TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001454</StopID></Stop><Stop><Lat>38.909717</Lat><Lon>-77.026814</Lon><Name>P ST NW + 11TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001458</StopID></Stop><Stop><Lat>38.909706</Lat><Lon>-77.028662</Lon><Name>P ST NW + LOGAN CIR NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001451</StopID></Stop><Stop><Lat>38.910125</Lat><Lon>-77.030154</Lon><Name>LOGAN CIR NW + 13TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>63</a:string><a:string>G2</a:string></Routes><StopID>1001469</StopID></Stop><Stop><Lat>38.909726</Lat><Lon>-77.032411</Lon><Name>P ST NW + 14TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001446</StopID></Stop><Stop><Lat>38.909707</Lat><Lon>-77.036360</Lon><Name>P ST NW + 16TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001457</StopID></Stop><Stop><Lat>38.909723</Lat><Lon>-77.038234</Lon><Name>P ST NW + 17TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001443</StopID></Stop><Stop><Lat>38.909723</Lat><Lon>-77.042182</Lon><Name>P ST NW + DUPONT CIR NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>N2</a:string><a:string>N4</a:string><a:string>N4v1</a:string><a:string>N6</a:string></Routes><StopID>1001459</StopID></Stop><Stop><Lat>38.909720</Lat><Lon>-77.044753</Lon><Name>P ST NW + 20TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>D2</a:string><a:string>D2v1</a:string><a:string>G2</a:string></Routes><StopID>1001461</StopID></Stop><Stop><Lat>38.909711</Lat><Lon>-77.048666</Lon><Name>P ST NW + 22ND ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>D1</a:string><a:string>D2</a:string><a:string>D2v1</a:string><a:string>D6</a:string><a:string>D6v3</a:string><a:string>G2</a:string></Routes><StopID>1001449</StopID></Stop><Stop><Lat>38.909725</Lat><Lon>-77.052244</Lon><Name>P ST NW + ROCK CREEK &amp; POTOMAC PKWY RAMP</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001463</StopID></Stop><Stop><Lat>38.909439</Lat><Lon>-77.055408</Lon><Name>P ST NW + 27TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001424</StopID></Stop><Stop><Lat>38.909424</Lat><Lon>-77.056951</Lon><Name>P ST NW + 28TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001421</StopID></Stop><Stop><Lat>38.909407</Lat><Lon>-77.057996</Lon><Name>P ST NW + 29TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001420</StopID></Stop><Stop><Lat>38.909379</Lat><Lon>-77.059097</Lon><Name>P ST NW + 30TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001419</StopID></Stop><Stop><Lat>38.909351</Lat><Lon>-77.061284</Lon><Name>P ST NW + 31ST ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001417</StopID></Stop><Stop><Lat>38.909307</Lat><Lon>-77.064005</Lon><Name>P ST NW + WISCONSIN AVE NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001415</StopID></Stop><Stop><Lat>38.908847</Lat><Lon>-77.065912</Lon><Name>P ST NW + 33RD ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001401</StopID></Stop><Stop><Lat>38.908794</Lat><Lon>-77.068904</Lon><Name>P ST NW + 35TH ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string></Routes><StopID>1001398</StopID></Stop><Stop><Lat>38.907395</Lat><Lon>-77.071671</Lon><Name>37TH ST NW + O ST NW</Name><Routes xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"><a:string>G2</a:string><a:string>G2v1</a:string></Routes><StopID>1001370</StopID></Stop></Stops><TripHeadsign>GEORGETOWN UNIVERSITY</TripHeadsign></Direction1><Name>G2 - HOWARD UNIV - GEORGETOWN UNIV</Name><RouteID>G2</RouteID></RouteDetailsInfo>`,
			unmarshalledResponse: &GetRouteDetailsResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "RouteDetailsInfo",
				},
				RouteID: "G2",
				Name:    "G2 - HOWARD UNIV - GEORGETOWN UNIV",
				Direction0: Direction{
					TripDestination: "LEDROIT PARK - HOWARD UNIVERSITY",
					DirectionText:   "EAST",
					DirectionNumber: "0",
					Shapes: []ShapePoint{
						{
							Latitude:       38.907367,
							Longitude:      -77.071641,
							SequenceNumber: 1,
						},
						{
							Latitude:       38.90599,
							Longitude:      -77.071515,
							SequenceNumber: 2,
						},
						{
							Latitude:       38.905889,
							Longitude:      -77.071411,
							SequenceNumber: 3,
						},
						{
							Latitude:       38.905965,
							Longitude:      -77.069018,
							SequenceNumber: 4,
						},
						{
							Latitude:       38.907698,
							Longitude:      -77.069012,
							SequenceNumber: 5,
						},
						{
							Latitude:       38.907781,
							Longitude:      -77.06886,
							SequenceNumber: 6,
						},
						{
							Latitude:       38.907803,
							Longitude:      -77.068501,
							SequenceNumber: 7,
						},
						{
							Latitude:       38.907882,
							Longitude:      -77.063647,
							SequenceNumber: 8,
						},
						{
							Latitude:       38.907674,
							Longitude:      -77.0634,
							SequenceNumber: 9,
						},
						{
							Latitude:       38.907654,
							Longitude:      -77.063209,
							SequenceNumber: 10,
						},
						{
							Latitude:       38.907527,
							Longitude:      -77.063125,
							SequenceNumber: 11,
						},
						{
							Latitude:       38.907597,
							Longitude:      -77.062759,
							SequenceNumber: 12,
						},
						{
							Latitude:       38.907547,
							Longitude:      -77.059455,
							SequenceNumber: 13,
						},
						{
							Latitude:       38.907654,
							Longitude:      -77.059196,
							SequenceNumber: 14,
						},
						{
							Latitude:       38.907674,
							Longitude:      -77.057202,
							SequenceNumber: 15,
						},
						{
							Latitude:       38.907958,
							Longitude:      -77.05699,
							SequenceNumber: 16,
						},
						{
							Latitude:       38.90926,
							Longitude:      -77.056984,
							SequenceNumber: 17,
						},
						{
							Latitude:       38.909348,
							Longitude:      -77.056327,
							SequenceNumber: 18,
						},
						{
							Latitude:       38.909367,
							Longitude:      -77.053894,
							SequenceNumber: 19,
						},
						{
							Latitude:       38.909587,
							Longitude:      -77.052527,
							SequenceNumber: 20,
						},
						{
							Latitude:       38.909537,
							Longitude:      -77.049446,
							SequenceNumber: 21,
						},
						{
							Latitude:       38.909618,
							Longitude:      -77.049248,
							SequenceNumber: 22,
						},
						{
							Latitude:       38.909633,
							Longitude:      -77.048614,
							SequenceNumber: 23,
						},
						{
							Latitude:       38.909649,
							Longitude:      -77.044439,
							SequenceNumber: 24,
						},
						{
							Latitude:       38.909559,
							Longitude:      -77.04424,
							SequenceNumber: 25,
						},
						{
							Latitude:       38.90914,
							Longitude:      -77.043909,
							SequenceNumber: 26,
						},
						{
							Latitude:       38.909019,
							Longitude:      -77.04343,
							SequenceNumber: 27,
						},
						{
							Latitude:       38.90914,
							Longitude:      -77.04298,
							SequenceNumber: 28,
						},
						{
							Latitude:       38.909419,
							Longitude:      -77.04271,
							SequenceNumber: 29,
						},
						{
							Latitude:       38.909529,
							Longitude:      -77.04267,
							SequenceNumber: 30,
						},
						{
							Latitude:       38.909619,
							Longitude:      -77.042379,
							SequenceNumber: 31,
						},
						{
							Latitude:       38.909621,
							Longitude:      -77.039001,
							SequenceNumber: 32,
						},
						{
							Latitude:       38.909573,
							Longitude:      -77.038826,
							SequenceNumber: 33,
						},
						{
							Latitude:       38.90956,
							Longitude:      -77.03682,
							SequenceNumber: 34,
						},
						{
							Latitude:       38.90961,
							Longitude:      -77.036773,
							SequenceNumber: 35,
						},
						{
							Latitude:       38.909621,
							Longitude:      -77.036124,
							SequenceNumber: 36,
						},
						{
							Latitude:       38.909618,
							Longitude:      -77.032348,
							SequenceNumber: 37,
						},
						{
							Latitude:       38.909694,
							Longitude:      -77.031632,
							SequenceNumber: 38,
						},
						{
							Latitude:       38.909652,
							Longitude:      -77.030617,
							SequenceNumber: 39,
						},
						{
							Latitude:       38.909596,
							Longitude:      -77.030411,
							SequenceNumber: 40,
						},
						{
							Latitude:       38.909229,
							Longitude:      -77.030013,
							SequenceNumber: 41,
						},
						{
							Latitude:       38.909145,
							Longitude:      -77.029456,
							SequenceNumber: 42,
						},
						{
							Latitude:       38.909289,
							Longitude:      -77.029067,
							SequenceNumber: 43,
						},
						{
							Latitude:       38.909641,
							Longitude:      -77.028724,
							SequenceNumber: 44,
						},
						{
							Latitude:       38.90963,
							Longitude:      -77.016279,
							SequenceNumber: 45,
						},
						{
							Latitude:       38.910004,
							Longitude:      -77.016434,
							SequenceNumber: 46,
						},
						{
							Latitude:       38.910038,
							Longitude:      -77.016356,
							SequenceNumber: 47,
						},
						{
							Latitude:       38.910397,
							Longitude:      -77.016575,
							SequenceNumber: 48,
						},
						{
							Latitude:       38.914099,
							Longitude:      -77.017889,
							SequenceNumber: 49,
						},
						{
							Latitude:       38.914149,
							Longitude:      -77.017319,
							SequenceNumber: 50,
						},
						{
							Latitude:       38.91424,
							Longitude:      -77.017129,
							SequenceNumber: 51,
						},
						{
							Latitude:       38.91577,
							Longitude:      -77.0173,
							SequenceNumber: 52,
						},
						{
							Latitude:       38.91848,
							Longitude:      -77.017759,
							SequenceNumber: 53,
						},
						{
							Latitude:       38.91859,
							Longitude:      -77.017691,
							SequenceNumber: 54,
						},
						{
							Latitude:       38.919297,
							Longitude:      -77.017877,
							SequenceNumber: 55,
						},
						{
							Latitude:       38.919404,
							Longitude:      -77.01779,
							SequenceNumber: 56,
						},
						{
							Latitude:       38.919674,
							Longitude:      -77.015208,
							SequenceNumber: 57,
						},
						{
							Latitude:       38.91966,
							Longitude:      -77.014799,
							SequenceNumber: 58,
						},
						{
							Latitude:       38.920501,
							Longitude:      -77.014905,
							SequenceNumber: 59,
						},
						{
							Latitude:       38.920541,
							Longitude:      -77.014838,
							SequenceNumber: 60,
						},
						{
							Latitude:       38.92067,
							Longitude:      -77.01493,
							SequenceNumber: 61,
						},
						{
							Latitude:       38.920617,
							Longitude:      -77.015862,
							SequenceNumber: 62,
						},
					},
					Stops: []Stop{
						{
							StopID:    "1001370",
							Name:      "37TH ST NW + O ST NW",
							Longitude: -77.071671,
							Latitude:  38.907395,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001345",
							Name:      "PROSPECT ST NW + 36TH ST NW",
							Longitude: -77.069828,
							Latitude:  38.905833,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001354",
							Name:      "35TH ST NW + N ST NW",
							Longitude: -77.06888,
							Latitude:  38.906649,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001385",
							Name:      "O ST NW + 34TH ST NW",
							Longitude: -77.067972,
							Latitude:  38.9077,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001379",
							Name:      "DUMBARTON ST NW + WISCONSIN AVE NW",
							Longitude: -77.063119,
							Latitude:  38.90757,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001381",
							Name:      "DUMBARTON ST NW + 30TH ST NW",
							Longitude: -77.059336,
							Latitude:  38.907603,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001389",
							Name:      "28TH ST NW + DUMBARTON ST NW",
							Longitude: -77.056995,
							Latitude:  38.907989,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001412",
							Name:      "28TH ST NW + P ST NW",
							Longitude: -77.057012,
							Latitude:  38.909259,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001418",
							Name:      "P ST NW + 26TH ST NW",
							Longitude: -77.054884,
							Latitude:  38.909313,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001435",
							Name:      "P ST NW + ROCK CREEK & POTOMAC PKWY RAMP",
							Longitude: -77.052,
							Latitude:  38.909596,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001433",
							Name:      "P ST NW + 23RD ST NW",
							Longitude: -77.049559,
							Latitude:  38.909571,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001440",
							Name:      "P ST NW + 21ST ST NW",
							Longitude: -77.046793,
							Latitude:  38.909583,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001439",
							Name:      "P ST NW + 20TH ST NW",
							Longitude: -77.045051,
							Latitude:  38.909582,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001436",
							Name:      "P ST NW + 18TH ST NW",
							Longitude: -77.041944,
							Latitude:  38.909582,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001434",
							Name:      "P ST NW + 17TH ST NW",
							Longitude: -77.038752,
							Latitude:  38.909565,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001438",
							Name:      "P ST NW + 16TH ST NW",
							Longitude: -77.036766,
							Latitude:  38.909576,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001437",
							Name:      "P ST NW + 15TH ST NW",
							Longitude: -77.034773,
							Latitude:  38.909576,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001452",
							Name:      "P ST NW + 14TH ST NW",
							Longitude: -77.032215,
							Latitude:  38.909597,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
						{
							StopID:    "1001456",
							Name:      "P ST NW + LOGAN CIR NW",
							Longitude: -77.03054,
							Latitude:  38.909596,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001416",
							Name:      "LOGAN CIR NW + P ST NW",
							Longitude: -77.0289,
							Latitude:  38.909349,
							Routes: []string{
								"63",
								"G2",
							},
						},
						{
							StopID:    "1001455",
							Name:      "P ST NW + 11TH ST NW",
							Longitude: -77.02733,
							Latitude:  38.909594,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001442",
							Name:      "P ST NW + 9TH ST NW",
							Longitude: -77.024196,
							Latitude:  38.909599,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001445",
							Name:      "P ST NW + 7TH ST NW",
							Longitude: -77.021287,
							Latitude:  38.909593,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001444",
							Name:      "P ST NW + 5TH ST NW",
							Longitude: -77.019082,
							Latitude:  38.9096,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001472",
							Name:      "NEW JERSEY AVE NW + P ST NW",
							Longitude: -77.016348,
							Latitude:  38.910066,
							Routes: []string{
								"96",
								"96v1",
								"96v2",
								"96v3",
								"G2",
							},
						},
						{
							StopID:    "1001544",
							Name:      "NEW JERSEY AVE NW + R ST NW",
							Longitude: -77.017218,
							Latitude:  38.912503,
							Routes: []string{
								"96",
								"96v1",
								"96v2",
								"96v3",
								"G2",
							},
						},
						{
							StopID:    "1001583",
							Name:      "NEW JERSEY AVE NW + RHODE ISLAND AVE NW",
							Longitude: -77.017538,
							Latitude:  38.913357,
							Routes: []string{
								"96",
								"96v1",
								"96v2",
								"96v3",
								"G2",
							},
						},
						{
							StopID:    "1003839",
							Name:      "4TH ST NW + T ST NW",
							Longitude: -77.017246,
							Latitude:  38.915794,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001720",
							Name:      "4TH ST NW + V ST NW",
							Longitude: -77.017709,
							Latitude:  38.918686,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001745",
							Name:      "W ST NW + 4TH ST NW",
							Longitude: -77.017086,
							Latitude:  38.919395,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001748",
							Name:      "W ST NW + 2ND ST NW",
							Longitude: -77.014934,
							Latitude:  38.919607,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001764",
							Name:      "2ND ST NW + BRYANT ST NW",
							Longitude: -77.014838,
							Latitude:  38.920541,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001776",
							Name:      "BRYANT ST NW + #301",
							Longitude: -77.016037,
							Latitude:  38.920612,
							Routes: []string{
								"G2",
							},
						},
					},
				},
				Direction1: Direction{
					TripDestination: "GEORGETOWN UNIVERSITY",
					DirectionText:   "WEST",
					DirectionNumber: "1",
					Shapes: []ShapePoint{
						{
							Latitude:       38.920625,
							Longitude:      -77.016035,
							SequenceNumber: 1,
						},
						{
							Latitude:       38.920433,
							Longitude:      -77.017861,
							SequenceNumber: 2,
						},
						{
							Latitude:       38.920332,
							Longitude:      -77.018029,
							SequenceNumber: 3,
						},
						{
							Latitude:       38.918199,
							Longitude:      -77.017782,
							SequenceNumber: 4,
						},
						{
							Latitude:       38.914495,
							Longitude:      -77.017155,
							SequenceNumber: 5,
						},
						{
							Latitude:       38.914214,
							Longitude:      -77.017227,
							SequenceNumber: 6,
						},
						{
							Latitude:       38.914107,
							Longitude:      -77.017409,
							SequenceNumber: 7,
						},
						{
							Latitude:       38.914096,
							Longitude:      -77.017837,
							SequenceNumber: 8,
						},
						{
							Latitude:       38.91396,
							Longitude:      -77.01791,
							SequenceNumber: 9,
						},
						{
							Latitude:       38.913196,
							Longitude:      -77.017679,
							SequenceNumber: 10,
						},
						{
							Latitude:       38.913333,
							Longitude:      -77.017662,
							SequenceNumber: 11,
						},
						{
							Latitude:       38.912373,
							Longitude:      -77.017364,
							SequenceNumber: 12,
						},
						{
							Latitude:       38.911257,
							Longitude:      -77.016845,
							SequenceNumber: 13,
						},
						{
							Latitude:       38.909804,
							Longitude:      -77.016418,
							SequenceNumber: 14,
						},
						{
							Latitude:       38.90963,
							Longitude:      -77.016279,
							SequenceNumber: 15,
						},
						{
							Latitude:       38.909652,
							Longitude:      -77.023735,
							SequenceNumber: 16,
						},
						{
							Latitude:       38.909599,
							Longitude:      -77.026504,
							SequenceNumber: 17,
						},
						{
							Latitude:       38.909655,
							Longitude:      -77.026779,
							SequenceNumber: 18,
						},
						{
							Latitude:       38.909672,
							Longitude:      -77.028808,
							SequenceNumber: 19,
						},
						{
							Latitude:       38.910117,
							Longitude:      -77.029335,
							SequenceNumber: 20,
						},
						{
							Latitude:       38.910187,
							Longitude:      -77.029747,
							SequenceNumber: 21,
						},
						{
							Latitude:       38.91008,
							Longitude:      -77.03012,
							SequenceNumber: 22,
						},
						{
							Latitude:       38.909927,
							Longitude:      -77.030281,
							SequenceNumber: 23,
						},
						{
							Latitude:       38.909663,
							Longitude:      -77.030456,
							SequenceNumber: 24,
						},
						{
							Latitude:       38.90964,
							Longitude:      -77.042036,
							SequenceNumber: 25,
						},
						{
							Latitude:       38.909768,
							Longitude:      -77.042291,
							SequenceNumber: 26,
						},
						{
							Latitude:       38.90968,
							Longitude:      -77.042529,
							SequenceNumber: 27,
						},
						{
							Latitude:       38.910069,
							Longitude:      -77.04286,
							SequenceNumber: 28,
						},
						{
							Latitude:       38.91026,
							Longitude:      -77.043289,
							SequenceNumber: 29,
						},
						{
							Latitude:       38.910209,
							Longitude:      -77.043769,
							SequenceNumber: 30,
						},
						{
							Latitude:       38.90997,
							Longitude:      -77.04411,
							SequenceNumber: 31,
						},
						{
							Latitude:       38.9097,
							Longitude:      -77.04425,
							SequenceNumber: 32,
						},
						{
							Latitude:       38.909649,
							Longitude:      -77.044439,
							SequenceNumber: 33,
						},
						{
							Latitude:       38.909702,
							Longitude:      -77.044663,
							SequenceNumber: 34,
						},
						{
							Latitude:       38.909644,
							Longitude:      -77.048644,
							SequenceNumber: 35,
						},
						{
							Latitude:       38.90966,
							Longitude:      -77.05225,
							SequenceNumber: 36,
						},
						{
							Latitude:       38.909725,
							Longitude:      -77.052244,
							SequenceNumber: 37,
						},
						{
							Latitude:       38.90966,
							Longitude:      -77.05225,
							SequenceNumber: 38,
						},
						{
							Latitude:       38.90966,
							Longitude:      -77.0528,
							SequenceNumber: 39,
						},
						{
							Latitude:       38.909399,
							Longitude:      -77.05395,
							SequenceNumber: 40,
						},
						{
							Latitude:       38.90937,
							Longitude:      -77.05544,
							SequenceNumber: 41,
						},
						{
							Latitude:       38.909407,
							Longitude:      -77.055442,
							SequenceNumber: 42,
						},
						{
							Latitude:       38.909328,
							Longitude:      -77.057304,
							SequenceNumber: 43,
						},
						{
							Latitude:       38.909252,
							Longitude:      -77.063934,
							SequenceNumber: 44,
						},
						{
							Latitude:       38.909303,
							Longitude:      -77.064203,
							SequenceNumber: 45,
						},
						{
							Latitude:       38.908973,
							Longitude:      -77.064211,
							SequenceNumber: 46,
						},
						{
							Latitude:       38.90889,
							Longitude:      -77.064369,
							SequenceNumber: 47,
						},
						{
							Latitude:       38.908656,
							Longitude:      -77.07161,
							SequenceNumber: 48,
						},
						{
							Latitude:       38.907421,
							Longitude:      -77.071633,
							SequenceNumber: 49,
						},
					},
					Stops: []Stop{
						{
							StopID:    "1001776",
							Name:      "BRYANT ST NW + #301",
							Longitude: -77.016037,
							Latitude:  38.920612,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001737",
							Name:      "4TH ST NW + W ST NW",
							Longitude: -77.017909,
							Latitude:  38.91914,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1003840",
							Name:      "4TH ST NW + U ST NW",
							Longitude: -77.01757,
							Latitude:  38.916976,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001584",
							Name:      "NEW JERSEY AVE NW + RHODE ISLAND AVE NW",
							Longitude: -77.017662,
							Latitude:  38.913167,
							Routes: []string{
								"96",
								"96v1",
								"96v4",
								"96v5",
								"G2",
							},
						},
						{
							StopID:    "1001541",
							Name:      "NEW JERSEY AVE NW + R ST NW",
							Longitude: -77.017378,
							Latitude:  38.912372,
							Routes: []string{
								"96",
								"96v1",
								"96v4",
								"96v5",
								"G2",
							},
						},
						{
							StopID:    "1001465",
							Name:      "NEW JERSEY AVE NW + P ST NW",
							Longitude: -77.016463,
							Latitude:  38.909815,
							Routes: []string{
								"96",
								"96v1",
								"96v4",
								"96v5",
								"G2",
							},
						},
						{
							StopID:    "1001441",
							Name:      "P ST NW + 5TH ST NW",
							Longitude: -77.018773,
							Latitude:  38.909704,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001453",
							Name:      "P ST NW + 7TH ST NW",
							Longitude: -77.021658,
							Latitude:  38.909709,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001454",
							Name:      "P ST NW + 9TH ST NW",
							Longitude: -77.023737,
							Latitude:  38.909706,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001458",
							Name:      "P ST NW + 11TH ST NW",
							Longitude: -77.026814,
							Latitude:  38.909717,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001451",
							Name:      "P ST NW + LOGAN CIR NW",
							Longitude: -77.028662,
							Latitude:  38.909706,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001469",
							Name:      "LOGAN CIR NW + 13TH ST NW",
							Longitude: -77.030154,
							Latitude:  38.910125,
							Routes: []string{
								"63",
								"G2",
							},
						},
						{
							StopID:    "1001446",
							Name:      "P ST NW + 14TH ST NW",
							Longitude: -77.032411,
							Latitude:  38.909726,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001457",
							Name:      "P ST NW + 16TH ST NW",
							Longitude: -77.03636,
							Latitude:  38.909707,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001443",
							Name:      "P ST NW + 17TH ST NW",
							Longitude: -77.038234,
							Latitude:  38.909723,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001459",
							Name:      "P ST NW + DUPONT CIR NW",
							Longitude: -77.042182,
							Latitude:  38.909723,
							Routes: []string{
								"G2",
								"N2",
								"N4",
								"N4v1",
								"N6",
							},
						},
						{
							StopID:    "1001461",
							Name:      "P ST NW + 20TH ST NW",
							Longitude: -77.044753,
							Latitude:  38.90972,
							Routes: []string{
								"D2",
								"D2v1",
								"G2",
							},
						},
						{
							StopID:    "1001449",
							Name:      "P ST NW + 22ND ST NW",
							Longitude: -77.048666,
							Latitude:  38.909711,
							Routes: []string{
								"D1",
								"D2",
								"D2v1",
								"D6",
								"D6v3",
								"G2",
							},
						},
						{
							StopID:    "1001463",
							Name:      "P ST NW + ROCK CREEK & POTOMAC PKWY RAMP",
							Longitude: -77.052244,
							Latitude:  38.909725,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001424",
							Name:      "P ST NW + 27TH ST NW",
							Longitude: -77.055408,
							Latitude:  38.909439,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001421",
							Name:      "P ST NW + 28TH ST NW",
							Longitude: -77.056951,
							Latitude:  38.909424,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001420",
							Name:      "P ST NW + 29TH ST NW",
							Longitude: -77.057996,
							Latitude:  38.909407,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001419",
							Name:      "P ST NW + 30TH ST NW",
							Longitude: -77.059097,
							Latitude:  38.909379,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001417",
							Name:      "P ST NW + 31ST ST NW",
							Longitude: -77.061284,
							Latitude:  38.909351,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001415",
							Name:      "P ST NW + WISCONSIN AVE NW",
							Longitude: -77.064005,
							Latitude:  38.909307,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001401",
							Name:      "P ST NW + 33RD ST NW",
							Longitude: -77.065912,
							Latitude:  38.908847,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001398",
							Name:      "P ST NW + 35TH ST NW",
							Longitude: -77.068904,
							Latitude:  38.908794,
							Routes: []string{
								"G2",
							},
						},
						{
							StopID:    "1001370",
							Name:      "37TH ST NW + O ST NW",
							Longitude: -77.071671,
							Latitude:  38.907395,
							Routes: []string{
								"G2",
								"G2v1",
							},
						},
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

func TestGetRouteDetails(t *testing.T) {
	jsonAndXmlPaths := []string{"/Bus.svc/json/jRouteDetails", "/Bus.svc/RouteDetails"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetRouteDetails")
			return
		}

		for _, request := range testRequests {
			response, err := testService.GetRouteDetails(request.stringParam1, request.stringParam2)

			if err != nil {
				t.Errorf("error calling GetPostions, routeID: %s date: %s error: %s", request.stringParam1, request.stringParam2, err.Error())
				continue
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}