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
	"/Bus.svc/json/jRoutes": {
		{
			rawQuery: "",
			response: `{"Routes":[{"RouteID":"10A","Name":"10A - HUNTINGTON STA - PENTAGON","LineDescription":"Alexandria-Pentagon Line"},{"RouteID":"10B","Name":"10B - HUNTING POINT - BALLSTON STA","LineDescription":"Hunting Point-Ballston Line"},{"RouteID":"10E","Name":"10E - HUNTING POINT - PENTAGON","LineDescription":"Alexandria-Pentagon Line"},{"RouteID":"10N","Name":"10N - NATIONAL AIRPORT - PENTAGON","LineDescription":"Alexandria-Pentagon Line"},{"RouteID":"11Y","Name":"11Y - MOUNT VERNON - POTOMAC PARK","LineDescription":"Mt Vernon Express Line"},{"RouteID":"11Yv1","Name":"11Y - POTOMAC PARK - HUNTING POINT","LineDescription":"Mt Vernon Express Line"},{"RouteID":"11Yv2","Name":"11Y - HUNTING POINT - POTOMAC PARK","LineDescription":"Mt Vernon Express Line"},{"RouteID":"15K","Name":"15K - E FALLS CH STA-MCLEAN LANGLEY","LineDescription":"Chain Bridge Road Line"},{"RouteID":"15Kv1","Name":"15K - CIA - EAST FALLS CHURCH STA","LineDescription":"Chain Bridge Road Line"},{"RouteID":"16A","Name":"16A - PATRIOT+AMERICANA - PENTAGON","LineDescription":"Columbia Pike Line"},{"RouteID":"16C","Name":"16C - CULMORE - PENTAGON","LineDescription":"Columbia Pike Line"},{"RouteID":"16Cv1","Name":"16C - CULMORE - FEDERAL TRIANGLE","LineDescription":"Columbia Pike Line"},{"RouteID":"16E","Name":"16E - GLN CARLYN+VISTA - FRANKLIN SQ","LineDescription":"Columbia Pike Line"},{"RouteID":"16G","Name":"16G - DINWIDDIE+COLUMBIA - PENTAGON CITY","LineDescription":"Columbia Pike-Pentagon City Line"},{"RouteID":"16Gv1","Name":"16G - COL PIKE+CARLIN SPR - DINWD+COL PK","LineDescription":"Columbia Pike-Pentagon City Line"},{"RouteID":"16H","Name":"16H - SKYLINE CITY - PENTAGON CITY STA","LineDescription":"Columbia Pike-Pentagon City Line"},{"RouteID":"16L","Name":"16L - ANNANDALE - PENTAGON HOV","LineDescription":"Annandale-Skyline City-Pentagon Line"},{"RouteID":"16Y","Name":"16Y - FOUR MILE+COLUMBIA - MCPHERSON SQ","LineDescription":"Columbia Pike-Farragut Square Line"},{"RouteID":"16Yv1","Name":"16Y - COLUMBIA PIKE - MCPHERSON SQ","LineDescription":"Columbia Pike-Farragut Square Line"},{"RouteID":"17B","Name":"17B - BURKE CENTRE - PENTAGON HOV","LineDescription":"Kings Park-North Springfield Line"},{"RouteID":"17G","Name":"17G - G MASON UNIV - PENTAGON","LineDescription":"Kings Park Express Line"},{"RouteID":"17H","Name":"17H - TWNBRK RD+TWNBRK RN- PENTAGON","LineDescription":"Kings Park Express Line"},{"RouteID":"17K","Name":"17K - TWNBRK RD+TWNBRK RN- PENTAGON","LineDescription":"Kings Park Express Line"},{"RouteID":"17L","Name":"17L - TWNBRK RD+TWNBRK RUN-PENTAGON","LineDescription":"Kings Park Express Line"},{"RouteID":"17M","Name":"17M - EDSALL+CANARD - PENTAGON","LineDescription":"Kings Park-North Springfield Line"},{"RouteID":"18G","Name":"18G - ROLLING VALLEY - PENTAGON","LineDescription":"Orange Hunt Line"},{"RouteID":"18H","Name":"18H - HUNTSMAN+CORK CTY - PENTAGON","LineDescription":"Orange Hunt Line"},{"RouteID":"18J","Name":"18J - ROLLING VALLEY - PENTAGON","LineDescription":"Orange Hunt Line"},{"RouteID":"18P","Name":"18P - BURKE CENTRE - PENTAGON","LineDescription":"Burke Centre Line"},{"RouteID":"18Pv1","Name":"18P - PENTAGON - ROLLING VALLEY MALL","LineDescription":"Burke Centre Line"},{"RouteID":"18Pv2","Name":"18P - ROLLING VALLEY P+R - PENTAGON","LineDescription":"Burke Centre Line"},{"RouteID":"1A","Name":"1A - VIENNA-BALLSTON (7 CORNERS)","LineDescription":"Wilson Blvd-Vienna Line"},{"RouteID":"1B","Name":"1B - DUNN LORING  - BALLSTON","LineDescription":"Wilson Blvd-Vienna Line"},{"RouteID":"1C","Name":"1C - WEST OX DIV-DUNN LORING (VIA MALL)","LineDescription":"Fair Oaks-Fairfax Blvd Line"},{"RouteID":"1Cv1","Name":"1C - FAIRFAX CO GOV CTR - DUNN LORING","LineDescription":"Fair Oaks-Fairfax Blvd Line"},{"RouteID":"1Cv2","Name":"1C - WEST OX DIV - DUNN LORING (NO MALL)","LineDescription":"Fair Oaks-Fairfax Blvd Line"},{"RouteID":"1Cv3","Name":"1C - FAIR OAKS MALL - DUNN LORING","LineDescription":"Fair Oaks-Fairfax Blvd Line"},{"RouteID":"1Cv4","Name":"1C - DUNN LORING - FAIR OAKS MALL","LineDescription":"Fair Oaks-Fairfax Blvd Line"},{"RouteID":"21A","Name":"21A - S REYNOLDS+EOS 21 CONDOS - PENTAGON","LineDescription":"Landmark-Bren Mar Park-Pentagon Line"},{"RouteID":"21D","Name":"21D - LANDMARK MEWS -PENTAGON","LineDescription":"Landmark-Bren Mar Park-Pentagon Line"},{"RouteID":"22A","Name":"22A - BALLSTON STA - PENTAGON","LineDescription":"Barcroft-South Fairlington Line"},{"RouteID":"22Av1","Name":"22A - SHIRLINGTON - BALLSTON  STA","LineDescription":"Barcroft-South Fairlington Line"},{"RouteID":"22C","Name":"22C - BALLSTON STA - PENTAGON","LineDescription":"Barcroft-South Fairlington Line"},{"RouteID":"22F","Name":"22F - NVCC - PENTAGON VIA HOV","LineDescription":"Barcroft-South Fairlington Line"},{"RouteID":"23A","Name":"23A - TYSONS CORNER CTR - CRYSTAL CTY","LineDescription":"McLean-Crystal City Line"},{"RouteID":"23B","Name":"23B - BALLSTON STA - CRYSTAL CTY","LineDescription":"McLean-Crystal City Line"},{"RouteID":"23Bv1","Name":"23B - LINDEN RESOURCES - BALLSTON STATION","LineDescription":"McLean-Crystal City Line"},{"RouteID":"23T","Name":"23T - TYSONS CORNER CTR - SHIRLINGTON","LineDescription":"McLean-Crystal City Line"},{"RouteID":"25B","Name":"25B - VAN DORN - BALLSTON","LineDescription":"Landmark-Ballston Line"},{"RouteID":"25Bv1","Name":"25B - SOUTHERN TOWERS - BALLSTON","LineDescription":"Landmark-Ballston Line"},{"RouteID":"25Bv2","Name":"25B - VAN DORN - BALLSTON/NO LDMRK CTR","LineDescription":"Landmark-Ballston Line"},{"RouteID":"25Bv3","Name":"25B - BALLSTON - SOUTHERN TOWERS","LineDescription":"Landmark-Ballston Line"},{"RouteID":"26A","Name":"26A - NVCC ANNANDALE - E FALLS CHURCH STA","LineDescription":"Annandale-East Falls Church Line"},{"RouteID":"28A","Name":"28A - TYSONS CORNER STA-KING ST STA","LineDescription":"Leesburg Pike Line"},{"RouteID":"28Av1","Name":"28A - SOUTHERN TOWERS-TYSONS CORNER STA","LineDescription":"Leesburg Pike Line"},{"RouteID":"28F","Name":"28F - BLDG 5113 G MASON DR - PENTAGON","LineDescription":"Skyline City Line"},{"RouteID":"28G","Name":"28G - BLDG 5113 G MASON DR - PENTAGON","LineDescription":"Skyline City Line"},{"RouteID":"29C","Name":"29C - NVCC ANNANDALE - PENTAGON","LineDescription":"Annandale Line"},{"RouteID":"29G","Name":"29G - AMERICANA+HERITAGE - PENTAGON","LineDescription":"Annandale Line"},{"RouteID":"29K","Name":"29K - GMU - KING ST STA","LineDescription":"Alexandria-Fairfax Line"},{"RouteID":"29Kv1","Name":"29K - GMU - KING ST/NO LDMRK","LineDescription":"Alexandria-Fairfax Line"},{"RouteID":"29N","Name":"29N - VIENNA STA - KING ST (VIA MALL)","LineDescription":"Alexandria-Fairfax Line"},{"RouteID":"29Nv1","Name":"29N - VIENNA STA-KING ST (NO MALL)","LineDescription":"Alexandria-Fairfax Line"},{"RouteID":"29W","Name":"29W - NVCC ANNANDALE - PENTAGON","LineDescription":"Braeburn Drive-Pentagon Express Line"},{"RouteID":"2A","Name":"2A - DUNN LORING STA - BALLSTON STA","LineDescription":"Washington Blvd.-Dunn Loring Line"},{"RouteID":"2B","Name":"2B - WEST OX RD DIV-DUNN LORING STATION","LineDescription":"Fair Oaks-Jermantown Road Line"},{"RouteID":"2Bv1","Name":"2B - WEST OX RD-DUNN LORING STA(NO MALL)","LineDescription":"Fair Oaks-Jermantown Road Line"},{"RouteID":"2Bv2","Name":"2B - FAIR OAKS MALL-DUNN LORING STATION","LineDescription":"Fair Oaks-Jermantown Road Line"},{"RouteID":"2Bv3","Name":"2B - DUNN LORING STA - FAIR OAKS MALL","LineDescription":"Fair Oaks-Jermantown Road Line"},{"RouteID":"30N","Name":"30N - FRIENDSHIP HGTS- NAYLOR RD STA","LineDescription":"Friendship Heights-Southeast Line"},{"RouteID":"30S","Name":"30S - FRIENDSHIP HGTS- SOUTHRN AVE STA","LineDescription":"Friendship Heights-Southeast Line"},{"RouteID":"31","Name":"31 - POTOMAC PARK-FRIENDSHIP HGTS","LineDescription":"Wisconsin Avenue Line"},{"RouteID":"32","Name":"32 - VIRGINIA AVE+E ST- SOUTHRN AVE","LineDescription":"Pennsylvania Avenue Line"},{"RouteID":"32v1","Name":"32 - PENN AVE + 8TH ST - SOUTHRN AVE","LineDescription":"Pennsylvania Avenue Line"},{"RouteID":"33","Name":"33 - 10TH ST+PA AV NW - FRIENDSHIP HGTS","LineDescription":"Wisconsin Avenue Line"},{"RouteID":"34","Name":"34 - 10TH ST + PA AVE- NAYLOR RD STA","LineDescription":"Pennsylvania Avenue Line"},{"RouteID":"36","Name":"36 - VIRGINIA AVE+E ST - NAYLOR RD STA","LineDescription":"Pennsylvania Avenue Line"},{"RouteID":"37","Name":"37 - 10TH ST+PA AV NW - FRIENDSHIP HGTS","LineDescription":"Wisconsin Avenue Limited Line"},{"RouteID":"38B","Name":"38B - BALLSTON - FARRAGUT","LineDescription":"Ballston-Farragut Square Line"},{"RouteID":"38Bv1","Name":"38B - WASH & QUINCY - FARRAGUT","LineDescription":"Ballston-Farragut Square Line"},{"RouteID":"38Bv2","Name":"38B - WASHINGTON-LEE HS - FARRAGUT SQ","LineDescription":"Ballston-Farragut Square Line"},{"RouteID":"39","Name":"39 - VIRGINIA AVE+21ST NW- NAYLOR RD STA","LineDescription":"Pennsylvania Avenue Limited Line"},{"RouteID":"3A","Name":"3A - ANNANDALE - E FALLS CHURCH","LineDescription":"Annandale Road Line"},{"RouteID":"3Av1","Name":"3A - ANNANDALE - 7 CORNERS","LineDescription":"Annandale Road Line"},{"RouteID":"3T","Name":"3T - MCLEAN STATION - E FALLS CH STA","LineDescription":"Pimmit Hills Line"},{"RouteID":"3Tv1","Name":"3T - MCLEAN STATION - W FALLS CHURCH","LineDescription":"Pimmit Hills Line"},{"RouteID":"3Y","Name":"3Y - E FALLS CHURCH - MCPHERSON SQ","LineDescription":"Lee Highway-Farragut Square Line"},{"RouteID":"42","Name":"42 - 9TH + F ST  - MT PLEASANT","LineDescription":"Mount Pleasant Line"},{"RouteID":"43","Name":"43 - I + 13TH NW - MT PLEASANT","LineDescription":"Mount Pleasant Line"},{"RouteID":"4A","Name":"4A - SEVEN CORNERS - ROSSLYN","LineDescription":"Pershing Dr-Arlington Blvd Line"},{"RouteID":"4B","Name":"4B - SEVEN CORNERS - ROSSLYN","LineDescription":"Pershing Dr-Arlington Blvd Line"},{"RouteID":"52","Name":"52 - L ENFNT PLAZA - TAKOMA STATION","LineDescription":"14th Street Line"},{"RouteID":"52v1","Name":"52 - L ENFNT PLAZA - 14TH+COLORADO","LineDescription":"14th Street Line"},{"RouteID":"52v2","Name":"52 - 14TH+COLORADO - L ENFANT PLAZA","LineDescription":"14th Street Line"},{"RouteID":"52v3","Name":"52 - 14TH & U - TAKOMA STATION","LineDescription":"14th Street Line"},{"RouteID":"54","Name":"54 - METRO CENTER - TAKOMA STA","LineDescription":"14th Street Line"},{"RouteID":"54v1","Name":"54 - 14TH+COLORADO - METRO CENTER","LineDescription":"14th Street Line"},{"RouteID":"59","Name":"59 - FEDERAL TRIANGLE - TAKOMA STATION","LineDescription":"14th Street Limited Line"},{"RouteID":"5A","Name":"5A - DULLES AIRPORT - LENFANT PLAZA","LineDescription":"DC-Dulles Line"},{"RouteID":"60","Name":"60 - GEORGIA + PETWORTH - FT TOTTEN","LineDescription":"Fort Totten-Petworth Line"},{"RouteID":"62","Name":"62 - GEORGIA+PETWORTH - TAKOMA STATION","LineDescription":"Takoma-Petworth Line"},{"RouteID":"62v1","Name":"62 - COOLIDGE HS - GEORGIA + PETWORTH","LineDescription":"Takoma-Petworth Line"},{"RouteID":"63","Name":"63 - FED TRIANGLE - TAKOMA STA","LineDescription":"Takoma-Petworth Line"},{"RouteID":"64","Name":"64 - FEDERAL TRIANGLE -FORT TOTTEN","LineDescription":"Fort Totten-Petworth Line"},{"RouteID":"64v1","Name":"64 - GEORGIA + PETWOTH - FT TOTTEN","LineDescription":"Fort Totten-Petworth Line"},{"RouteID":"70","Name":"70 - ARCHIVES - SILVER SPRING","LineDescription":"Georgia Avenue-7th Street Line"},{"RouteID":"70v1","Name":"70 - GEORGIA & EUCLID  - ARCHIVES","LineDescription":"Georgia Avenue-7th Street Line"},{"RouteID":"74","Name":"74 - NATIONALS PARK - CONVENTION CTR","LineDescription":"Convention Center-Southwest Waterfront Line"},{"RouteID":"79","Name":"79 - ARCHIVES - SILVER SPRING STA","LineDescription":"Georgia Avenue MetroExtra Line"},{"RouteID":"7A","Name":"7A - LINCOLNIA+QUANTRELL - PENTAGON","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"7Av1","Name":"7A - PENTAGON - SOUTHERN TWRS","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"7Av2","Name":"7A - SOUTHERN TWRS - PENTAGON","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"7Av3","Name":"7A - LINCOLNIA/QUANTRELL - PENT VIA PENT","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"7C","Name":"7C - PARK CENTER - PENTAGON","LineDescription":"Park Center-Pentagon Line"},{"RouteID":"7F","Name":"7F - LINCOLNIA+QUANTRELL - PENTAGON","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"7Fv1","Name":"7F - LINC + QUANT - PENT CITY - PENTAGON","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"7M","Name":"7M - MARK CENTER - PENTAGON (NON-HOV)","LineDescription":"Mark Center-Pentagon Line"},{"RouteID":"7Mv1","Name":"7M - MARK CENTER - PENTAGON (HOV)","LineDescription":"Mark Center-Pentagon Line"},{"RouteID":"7P","Name":"7P - PARK CTR - PENTAGON","LineDescription":"Park Center-Pentagon Line"},{"RouteID":"7W","Name":"7W - LNCLNA+QUANTRLL- PENTAGON","LineDescription":"Lincolnia-Pentagon Line"},{"RouteID":"7Y","Name":"7Y - SOUTHERN TWRS - H+17TH ST","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"7Yv1","Name":"7Y - PENTAGON - H+17TH ST","LineDescription":"Lincolnia-North Fairlington Line"},{"RouteID":"80","Name":"80 - KENNEDY CTR   - FORT TOTTEN STA","LineDescription":"North Capitol Street Line"},{"RouteID":"80v1","Name":"80 - MCPHERSON SQ  - BROOKLAND","LineDescription":"North Capitol Street Line"},{"RouteID":"80v2","Name":"80 - MCPHERSON SQ  - FORT TOTTEN STA","LineDescription":"North Capitol Street Line"},{"RouteID":"80v3","Name":"80 - KENNEDY CTR   - BROOKLAND STA","LineDescription":"North Capitol Street Line"},{"RouteID":"83","Name":"83 - RHODE ISLAND AVE STA-CHERRY HILL","LineDescription":"College Park Line"},{"RouteID":"83v1","Name":"83 - MT RAINIER - RHODE ISLAND","LineDescription":"College Park Line"},{"RouteID":"83v2","Name":"83 - RHODE ISLAND - MT RAINIER","LineDescription":"College Park Line"},{"RouteID":"83v3","Name":"83 - RHODE ISLAND AVE STA-COLLEGE PARK","LineDescription":"College Park Line"},{"RouteID":"83v4","Name":"83 - COLLEGE PARK-RHODE ISLAND AVE STA","LineDescription":"College Park Line"},{"RouteID":"86","Name":"86 - RHODE ISLAND AVE STA- CALVERTON","LineDescription":"College Park Line"},{"RouteID":"86v1","Name":"86 - RHODE ISLAND AVE STA- COLLEGE PARK","LineDescription":"College Park Line"},{"RouteID":"86v2","Name":"86 - MT RAINIER   - CALVERTON","LineDescription":"College Park Line"},{"RouteID":"87","Name":"87 - NEW CARROLTON -CYPRESS+LAURL LAKES","LineDescription":"Laurel Express Line"},{"RouteID":"87v1","Name":"87 - GRNBELT STA -CYPRESS+LAURL LAKES","LineDescription":"Laurel Express Line"},{"RouteID":"87v2","Name":"87 - GRNBELT-CYP+LRL LAKES (NO P&R)","LineDescription":"Laurel Express Line"},{"RouteID":"87v3","Name":"87 - GRNBELT STA - BALTIMORE+MAIN ST","LineDescription":"Laurel Express Line"},{"RouteID":"87v4","Name":"87 - BALTIMORE+MAIN ST - GRNBELT STA","LineDescription":"Laurel Express Line"},{"RouteID":"87v5","Name":"87 - CYPRESS+LAURL LAKES -GRNBELT STA","LineDescription":"Laurel Express Line"},{"RouteID":"89","Name":"89 - GREENBELT STA - 4TH & GREEN HILL","LineDescription":"Laurel Line"},{"RouteID":"89v1","Name":"89 - GREENBELT STA - CHERRY LA+4TH ST","LineDescription":"Laurel Line"},{"RouteID":"89M","Name":"89M - GREENBELT STA - S LAUREL P+R","LineDescription":"Laurel Line"},{"RouteID":"8S","Name":"8S - RADFORD+QUAKER - PENTAGON","LineDescription":"Foxchase-Seminary Valley Line"},{"RouteID":"8W","Name":"8W - MARK CENTER - PENTAGON V FOXCHASE","LineDescription":"Foxchase-Seminary Valley Line"},{"RouteID":"8Z","Name":"8Z - QUAKER+OSAGE - PENTAGON","LineDescription":"Foxchase-Seminary Valley Line"},{"RouteID":"90","Name":"90 - ANACOSTIA - DK ELLNGTN BRDG","LineDescription":"U Street-Garfield Line"},{"RouteID":"90v1","Name":"90 - 8TH ST + L ST  - DK ELLNGTN BRDG","LineDescription":"U Street-Garfield Line"},{"RouteID":"90v2","Name":"90 - KIPP DC PREP- ANACOSTIA","LineDescription":"U Street-Garfield Line"},{"RouteID":"92","Name":"92 - CONGRESS HTS STA - REEVES CTR","LineDescription":"U Street-Garfield Line"},{"RouteID":"92v1","Name":"92 - EASTERN MARKET - CONGRESS HGTS STA","LineDescription":"U Street-Garfield Line"},{"RouteID":"92v2","Name":"92 - CONGRESS HTS STA- EASTERN MARKET","LineDescription":"U Street-Garfield Line"},{"RouteID":"96","Name":"96 - TENLEYTOWN STA - CAPITOL HTS STA","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"96v1","Name":"96 - ELLINGTON BR - CAPITOL HTS","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"96v2","Name":"96 - CAPITOL HTS  - REEVES CTR","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"96v3","Name":"96 - CAPITOL HTS - ELLINGTON BRDG","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"96v4","Name":"96 - REEVES CTR  - CAPITOL HTS","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"96v5","Name":"96 - TENLEYTOWN STA - STADIUM ARMORY STA","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"97","Name":"97 - UNION STATION - CAPITOL HTS","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"97v1","Name":"97 - EASTERN HS - CAPITOL HTS","LineDescription":"East Capitol Street-Cardoza Line"},{"RouteID":"A12","Name":"A12 - ADDISON RD STA - CAPITAL PLAZA","LineDescription":"Martin Luther King Jr Highway Line"},{"RouteID":"A12v1","Name":"A12 - BARLWE+MATHEW HENSON - ADDISON RD","LineDescription":"Martin Luther King Jr Highway Line"},{"RouteID":"A12v2","Name":"A12 - CAPITOL HTS STA - CAPITAL PLAZA","LineDescription":"Martin Luther King Jr Highway Line"},{"RouteID":"A12v3","Name":"A12 - CAPITAL PLAZA - CAPITOL HTS","LineDescription":"Martin Luther King Jr Highway Line"},{"RouteID":"A2","Name":"A2 - SOUTHERN AVE - ANACOSTIA (VIA HOSP)","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A2v1","Name":"A2 - ANACOSTIA - MISS+ATLANTIC","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A2v2","Name":"A2 - SOUTHERN AVE STA - ANACOSTIA","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A2v3","Name":"A2 - MISS+ATLANTIC - ANACOSTIA","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A31","Name":"A31 - ANACOSTIA HIGH - MINNESOTA STA","LineDescription":"Minnesota Ave-Anacostia Line"},{"RouteID":"A32","Name":"A32 - ANACOSTIA HIGH - SOUTHRN AVE STA","LineDescription":"Minnesota Ave-Anacostia Line"},{"RouteID":"A33","Name":"A33 - ANACOSTIA HIGH - ANACOSTIA STA","LineDescription":"Minnesota Ave-Anacostia Line"},{"RouteID":"A4","Name":"A4 - DC VILLAGE - ANACOSTIA","LineDescription":"Anacostia-Fort Drum Line"},{"RouteID":"A4v1","Name":"A4 - USCG-FT DRUM (VIA ANAC)","LineDescription":"Anacostia-Fort Drum Line"},{"RouteID":"A4v2","Name":"A4 - FT DRUM - ANACOSTIA","LineDescription":"Anacostia-Fort Drum Line"},{"RouteID":"A4v3","Name":"A4 - ANACOSTIA - FT DRUM","LineDescription":"Anacostia-Fort Drum Line"},{"RouteID":"A4v4","Name":"A4 - FT DRUM-USCG (VIA ANAC)","LineDescription":"Anacostia-Fort Drum Line"},{"RouteID":"A4v5","Name":"A4 - DC VILL-USCG (VIA ANAC)","LineDescription":"Anacostia-Fort Drum Line"},{"RouteID":"A6","Name":"A6 - 4501 3RD ST - ANACOSTIA","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A6v1","Name":"A6 - SOUTHERN AVE+S CAPITOL - ANACOSTIA","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A7","Name":"A7 - SOUTHRN+S CAPITOL - ANACOSTIA","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A8","Name":"A8 - 4501 3RD ST - ANACOSTIA","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A8v1","Name":"A8 - SOUTHERN + S CAPITOL - ANACOSTIA","LineDescription":"Anacostia-Congress Heights Line"},{"RouteID":"A9","Name":"A9 - LIVINGSTON - MCPHERSON SQUARE","LineDescription":"Martin Luther King Jr Ave Limited Line"},{"RouteID":"B2","Name":"B2 - ANACOSTIA STA - MOUNT RAINIER","LineDescription":"Bladensburg Road-Anacostia Line"},{"RouteID":"B2v1","Name":"B2 - ANACOSTIA STA - BLDNSGRG+VST NE","LineDescription":"Bladensburg Road-Anacostia Line"},{"RouteID":"B2v2","Name":"B2 - POTOMAC AVE - MOUNT RAINIER","LineDescription":"Bladensburg Road-Anacostia Line"},{"RouteID":"B2v3","Name":"B2 - BLDNSBRG+26TH - ANACOSTIA STA","LineDescription":"Bladensburg Road-Anacostia Line"},{"RouteID":"B2v4","Name":"B2 - EASTERN HS - ANACOSTIA STA","LineDescription":"Bladensburg Road-Anacostia Line"},{"RouteID":"B21","Name":"B21 - NEW CARROLTON STA - BOWIE STATE","LineDescription":"Bowie State University Line"},{"RouteID":"B22","Name":"B22 - NEW CARROLTON STA - BOWIE STATE","LineDescription":"Bowie State University Line"},{"RouteID":"B22v1","Name":"B22 - OLD CHAPEL & 197 - NEW CARRLTN STA","LineDescription":"Bowie State University Line"},{"RouteID":"B24","Name":"B24 - NEW CARLTN STA-BOWIE P+R(VIA BHC)","LineDescription":"Bowie-Belair Line"},{"RouteID":"B24v1","Name":"B24 - NEW CARROLTON STA - BOWIE P+R","LineDescription":"Bowie-Belair Line"},{"RouteID":"B27","Name":"B27 - NEW CARROLTON STA - BOWIE STATE","LineDescription":"Bowie-New Carrollton Line"},{"RouteID":"B29","Name":"B29 - NEW CARROLTON - CROFTON CC (PM)(PR)","LineDescription":"Crofton-New Carrollton Line"},{"RouteID":"B29v1","Name":"B29 - NEW CARROLLTON STA - GATEWAY CTR","LineDescription":"Crofton-New Carrollton Line"},{"RouteID":"B29v2","Name":"B29 - NEW CARROLTON - CROFTON CC (NO PR)","LineDescription":"Crofton-New Carrollton Line"},{"RouteID":"B30","Name":"B30 - GREENBELT STA - BWI LT RAIL STA","LineDescription":"Greenbelt-BWI Thurgood Marshall Airport Express Line"},{"RouteID":"B8","Name":"B8 - RHODE ISLAND AV -PETERSBRG APTS","LineDescription":"Fort Lincoln Shuttle Line"},{"RouteID":"B8v1","Name":"B8 - BLDNSBRG+S DKTA -PETERSBRG APTS","LineDescription":"Fort Lincoln Shuttle Line"},{"RouteID":"B8v2","Name":"B8 - PETERSBRG APTS  -BLDNSBRG+S DKTA","LineDescription":"Fort Lincoln Shuttle Line"},{"RouteID":"B9","Name":"B9 - RHODE ISLAND AVE - COLMAR MANOR","LineDescription":"Fort Lincoln Shuttle Line"},{"RouteID":"C11","Name":"C11 - CLINTON P+R - BRANCH AVE STA","LineDescription":"Clinton Line"},{"RouteID":"C12","Name":"C12 - NAYLOR RD STA  - BRANCH AVE STA","LineDescription":"Hillcrest Heights Line"},{"RouteID":"C13","Name":"C13 - CLINTON P+R - BRANCH AVE STA","LineDescription":"Clinton Line"},{"RouteID":"C14","Name":"C14 - NAYLOR RD STA  - BRANCH AVE STA","LineDescription":"Hillcrest Heights Line"},{"RouteID":"C2","Name":"C2 - WHEATN STA - GRNBELT STA UMD ALT","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C2v1","Name":"C2 - TAKOMA LANGLEY XROADS - GRNBELT STA","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C2v2","Name":"C2 - GREENBELT STA - RANDOLPH + PARKLAWN","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C2v3","Name":"C2 - TAKOMA LANGLEY XROADS - WHEATON","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C21","Name":"C21 - ADDISON RD STA - COLLINGTON","LineDescription":"Central Avenue Line"},{"RouteID":"C21v1","Name":"C21 - ADDISON RD STA - POINTER RIDGE","LineDescription":"Central Avenue Line"},{"RouteID":"C21v2","Name":"C21 - ADDISON RD STA - CAMPUS WAY S","LineDescription":"Central Avenue Line"},{"RouteID":"C22","Name":"C22 - ADDISON RD STA - COLLINGTON","LineDescription":"Central Avenue Line"},{"RouteID":"C22v1","Name":"C22 - ADDISON RD STA- POINTER RIDGE","LineDescription":"Central Avenue Line"},{"RouteID":"C26","Name":"C26 - LARGO TOWN CTR - WATKNS+CHESTERTON","LineDescription":"Central Avenue Line"},{"RouteID":"C26v1","Name":"C26 - LARGO TOWN CTR-WATKINS+CAMBLETON","LineDescription":"Central Avenue Line"},{"RouteID":"C28","Name":"C28 - PT RIDGE -NEW CAROLTN VIA HLTH CTR","LineDescription":"Pointer Ridge Line"},{"RouteID":"C28v1","Name":"C28 - PT RIDGE - NEW CARROLLTON STA","LineDescription":"Pointer Ridge Line"},{"RouteID":"C29*1","Name":"C29 - POINTER RIDGE - ADDISON RD STA","LineDescription":"Central Avenue Line"},{"RouteID":"C29*2","Name":"C29 - WATKNS+CAMBLETON - ADDISN RD STA","LineDescription":"Central Avenue Line"},{"RouteID":"C29*4","Name":"C29 - ADDISON RD STA - BOWIE STATE","LineDescription":"Central Avenue Line"},{"RouteID":"C29/","Name":"C29 - ADDISON RD STA - POINTER RIDGE","LineDescription":"Central Avenue Line"},{"RouteID":"C4","Name":"C4 - TWINBROOK STA - PG PLAZA STA","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C4v1","Name":"C4 - TLTC-TWINBROOK","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C4v2","Name":"C4 - TWINBROOK STA - WHEATON STA","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C4v3","Name":"C4 - PG PLAZA STA - WHEATON STA","LineDescription":"Greenbelt-Twinbrook Line"},{"RouteID":"C8","Name":"C8 - WHITE FLINT - COLLEGE PARK","LineDescription":"College Park-White Flint Line"},{"RouteID":"C8v1","Name":"C8 - WHITE FLNT-COLLEGE PK (NO FDA/ARCH)","LineDescription":"College Park-White Flint Line"},{"RouteID":"C8v2","Name":"C8 - WHITE FLINT-COLLEGE PARK (NO FDA)","LineDescription":"College Park-White Flint Line"},{"RouteID":"C8v3","Name":"C8 - GLENMONT-COLLEGE PK (NO FDA/ARCH)","LineDescription":"College Park-White Flint Line"},{"RouteID":"D1","Name":"D1 - GLOVER PARK - FRANKLIN SQUARE","LineDescription":"Glover Park-Franklin Square Line"},{"RouteID":"D12","Name":"D12 - SOUTHERN AVE STA - SUITLAND STA","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D12v1","Name":"D12 - SUITLAND STA - SOUTHERN AVE STA","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D12v2","Name":"D12 - ST BARNABAS RD   - SUITLAND STA","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D13","Name":"D13 - SOUTHERN AVE STA - SUITLAND STA","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D13v1","Name":"D13 - SOUTHRN STA - ALLENTWN+OLD BRNCH","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D14","Name":"D14 - SOUTHERN AVE STA - SUITLAND STA","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D14v1","Name":"D14 - ALLENTWN+OLD BRNCH - SUITLND STA","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D14v2","Name":"D14 - SUITLND STA - ALLENTWN+OLD BRNCH","LineDescription":"Oxon Hill-Suitland Line"},{"RouteID":"D2","Name":"D2 - GLOVER PARK   - CONNETICUT +Q ST","LineDescription":"Glover Park-Dupont Circle Line"},{"RouteID":"D2v1","Name":"D2 - DUPONT CIR - GLOVER PARK","LineDescription":"Glover Park-Dupont Circle Line"},{"RouteID":"D31","Name":"D31 - TENLEYTOWN STA - 16TH + EASTERN","LineDescription":"16th Street-Tenleytown Line"},{"RouteID":"D32","Name":"D32 - TENLEYTOWN STA - COLUMBIA HTS STA","LineDescription":"16th Street-Tenleytown Line"},{"RouteID":"D33","Name":"D33 - TENLEYTOWN STA - 16TH + SHEPHERD","LineDescription":"16th Street-Tenleytown Line"},{"RouteID":"D34","Name":"D34 - TENLEYTOWN STA - 14TH + COLORADO","LineDescription":"16th Street-Tenleytown Line"},{"RouteID":"D4","Name":"D4 - DUPONT CIRCLE - IVY CITY","LineDescription":"Ivy City-Franklin Square Line"},{"RouteID":"D4v1","Name":"D4 - FRANKLIN SQUARE - IVY CITY","LineDescription":"Ivy City-Franklin Square Line"},{"RouteID":"D4v2","Name":"D4 - IVY CITY - FRANKLIN SQUARE","LineDescription":"Ivy City-Franklin Square Line"},{"RouteID":"D5","Name":"D5 - MASS LTTL FLWR CHRCH- FARRGT SQR","LineDescription":"MacArthur Blvd-Georgetown Line"},{"RouteID":"D6","Name":"D6 - SIBLEY HOSP - STADIUM ARMRY","LineDescription":"Sibley Hospital–Stadium-Armory Line"},{"RouteID":"D6v1","Name":"D6 - STADIUM ARMRY STA - FARRAGUT SQUARE","LineDescription":"Sibley Hospital–Stadium-Armory Line"},{"RouteID":"D6v2","Name":"D6 - SIBLEY HOSP - FARRAGUT SQ","LineDescription":"Sibley Hospital–Stadium-Armory Line"},{"RouteID":"D6v3","Name":"D6 - FARRAGUT SQUARE - STADIUM ARMRY","LineDescription":"Sibley Hospital–Stadium-Armory Line"},{"RouteID":"D8","Name":"D8 - UNION STATION - VA MED CTR","LineDescription":"Hospital Center Line"},{"RouteID":"D8v1","Name":"D8 - RHODE ISLAND STA - UNION STA","LineDescription":"Hospital Center Line"},{"RouteID":"E2","Name":"E2 - IVY CITY - FT TOTTEN","LineDescription":"Ivy City-Fort Totten Line"},{"RouteID":"E4","Name":"E4 - FRIENDSHP HTS - RIGGS PK","LineDescription":"Military Road-Crosstown Line"},{"RouteID":"E4v1","Name":"E4 - FRIENDSHP HTS - FT TOTTEN","LineDescription":"Military Road-Crosstown Line"},{"RouteID":"E4v2","Name":"E4 - FT TOTTEN - FRIENDSHP HTS","LineDescription":"Military Road-Crosstown Line"},{"RouteID":"E6","Name":"E6 - FRIENDSHP HTS  -KNOLLWOOD","LineDescription":"Chevy Chase Line"},{"RouteID":"F1","Name":"F1 - CHEVERLY STA - TAKOMA STA","LineDescription":"Chillum Road Line"},{"RouteID":"F12","Name":"F12 - CHEVERLY STA - NEW CARROLTON STA","LineDescription":"Ardwick Industrial Park Shuttle Line"},{"RouteID":"F12v1","Name":"F12 - CHEVERLY STA - LANDOVER STA","LineDescription":"Ardwick Industrial Park Shuttle Line"},{"RouteID":"F13","Name":"F13 - CHEVERLY STA-WASHINGTON BUS PARK","LineDescription":"Cheverly-Washington Business Park Line"},{"RouteID":"F13v1","Name":"F13 - CHEVERLY STA - WASHINGTON BUS PARK","LineDescription":"Cheverly-Washington Business Park Line"},{"RouteID":"F13v2","Name":"F13 - NEW CARROLTON - WASHINGTON BUS PARK","LineDescription":"Cheverly-Washington Business Park Line"},{"RouteID":"F13v3","Name":"F13 - WHITFIELD + VOLTA - CHEVERLY STA","LineDescription":"Cheverly-Washington Business Park Line"},{"RouteID":"F14","Name":"F14 - NAYLOR RD STA -NEW CARROLTON STA","LineDescription":"Sheriff Road-Capitol Heights Line"},{"RouteID":"F14v1","Name":"F14 - BRADBURY HGTS -NEW CARROLTON STA","LineDescription":"Sheriff Road-Capitol Heights Line"},{"RouteID":"F2","Name":"F2 - CHEVERLY STA - TAKOMA STA","LineDescription":"Chillum Road Line"},{"RouteID":"F2v1","Name":"F2 - CHEVRLY STA-QUEENS CHAPL+CARSN CIR","LineDescription":"Chillum Road Line"},{"RouteID":"F2v2","Name":"F2 - 34TH + OTIS - TAKOMA STA","LineDescription":"Chillum Road Line"},{"RouteID":"F4","Name":"F4 - SILVR SPRING STA - NEW CARRLLTON","LineDescription":"New Carrollton-Silver Spring Line"},{"RouteID":"F4v1","Name":"F4 - PG PLAZA STA -NEW CARROLTON STA","LineDescription":"New Carrollton-Silver Spring Line"},{"RouteID":"F4v2","Name":"F4 - NEW CARRLLTON - PG PLAZA STA","LineDescription":"New Carrollton-Silver Spring Line"},{"RouteID":"F6","Name":"F6 - FT TOTTEN - NEW CARROLLTN","LineDescription":"New Carrollton-Fort Totten Line"},{"RouteID":"F6v1","Name":"F6 - PG PLAZA - NEW CARROLLTON STA","LineDescription":"New Carrollton-Fort Totten Line"},{"RouteID":"F6v2","Name":"F6 - NEW CARROLLTON - PG PLAZA STA","LineDescription":"New Carrollton-Fort Totten Line"},{"RouteID":"F8","Name":"F8 - CHEVERLY STA - TLTC","LineDescription":"Langley Park-Cheverly Line"},{"RouteID":"G12","Name":"G12 - GREENBELT STA - NEW CARROLLTON STA","LineDescription":"Greenbelt-New Carrollton Line"},{"RouteID":"G12v1","Name":"G12 - GREENBELT STA - ROOSEVELT CENTER","LineDescription":"Greenbelt-New Carrollton Line"},{"RouteID":"G12v2","Name":"G12 - ROOSEVELT CTR - NEW CARROLLTON STA","LineDescription":"Greenbelt-New Carrollton Line"},{"RouteID":"G14","Name":"G14 - GREENBELT STA -NEW CARROLTON STA","LineDescription":"Greenbelt-New Carrollton Line"},{"RouteID":"G14v1","Name":"G14 - ROOSEVELT CTR -NEW CARROLTON STA","LineDescription":"Greenbelt-New Carrollton Line"},{"RouteID":"G14v2","Name":"G14 - GREENBELT STA - NEW CARROLTON STA","LineDescription":"Greenbelt-New Carrollton Line"},{"RouteID":"G2","Name":"G2 - GEORGETOWN UNIV - HOWARD UNIV","LineDescription":"P Street-LeDroit Park Line"},{"RouteID":"G2v1","Name":"G2 - GEORGETOWN UNIV - P+14TH ST","LineDescription":"P Street-LeDroit Park Line"},{"RouteID":"G8","Name":"G8 - FARRAGUT SQUARE - AVONDALE","LineDescription":"Rhode Island Avenue Line"},{"RouteID":"G8v1","Name":"G8 - RHODE ISLAND AVE STA-FARRAGUT SQ","LineDescription":"Rhode Island Avenue Line"},{"RouteID":"G8v2","Name":"G8 - FARRAGUT SQ - RHODE ISLAND STA","LineDescription":"Rhode Island Avenue Line"},{"RouteID":"G8v3","Name":"G8 - BROOKLAND STA - AVONDALE","LineDescription":"Rhode Island Avenue Line"},{"RouteID":"G9","Name":"G9 - FRANKLIN SQ - RHODE ISLD + EAST PM","LineDescription":"Rhode Island Avenue Limited Line"},{"RouteID":"G9v1","Name":"G9 - FRANKLIN SQ - RHODE ISLD + EASTERN","LineDescription":"Rhode Island Avenue Limited Line"},{"RouteID":"H1","Name":"H1 - C + 17TH ST - BROOKLAND CUA STA","LineDescription":"Brookland-Potomac Park Line"},{"RouteID":"H11","Name":"H11 - HEATHER HILL - NAYLOR RD STA","LineDescription":"Marlow Heights-Temple Hills Line"},{"RouteID":"H12","Name":"H12 - HEATHER HILL-NAYLOR RD (MACYS)","LineDescription":"Marlow Heights-Temple Hills Line"},{"RouteID":"H12v1","Name":"H12 - HEATHER HILL-NAYLOR RD STA (P&R)","LineDescription":"Marlow Heights-Temple Hills Line"},{"RouteID":"H13","Name":"H13 - HEATHER HILL - NAYLOR RD STA","LineDescription":"Marlow Heights-Temple Hills Line"},{"RouteID":"H2","Name":"H2 - TENLEYTOWN AU-ST - BROOKLND CUA STA","LineDescription":"Crosstown Line"},{"RouteID":"H3","Name":"H3 - TENLEYTOWN STA -BROOKLND CUA STA","LineDescription":"Crosstown Line"},{"RouteID":"H4","Name":"H4 - TENLEYTWN AU-STA -BROOKLAND CUA STA","LineDescription":"Crosstown Line"},{"RouteID":"H4v1","Name":"H4 - COLUMBIA RD+14TH -TENLEYTOWN-AU STA","LineDescription":"Crosstown Line"},{"RouteID":"H6","Name":"H6 - BROOKLAND CUA STA - FORT LINCOLN","LineDescription":"Brookland-Fort Lincoln Line"},{"RouteID":"H6v1","Name":"H6 - FORT LINCOLN - FORT LINCOLN","LineDescription":"Brookland-Fort Lincoln Line"},{"RouteID":"H8","Name":"H8 - MT PLEASANT+17TH - RHODE ISLAND","LineDescription":"Park Road-Brookland Line"},{"RouteID":"H9","Name":"H9 - RHODE ISLAND - FORT DR + 1ST","LineDescription":"Park Road-Brookland Line"},{"RouteID":"J1","Name":"J1 - MEDCAL CTR STA - SILVR SPRNG STA","LineDescription":"Bethesda-Silver Spring Line"},{"RouteID":"J1v1","Name":"J1 - SILVR SPRNG STA - MONT MALL/BAT LN","LineDescription":"Bethesda-Silver Spring Line"},{"RouteID":"J12","Name":"J12 - ADDISON RD STA - FORESTVILLE","LineDescription":"Marlboro Pike Line"},{"RouteID":"J12v1","Name":"J12 - ADDISON RD-FORESTVILLE VIA PRES PKY","LineDescription":"Marlboro Pike Line"},{"RouteID":"J2","Name":"J2 - MONTGOMRY MALL - SILVR SPRNG STA","LineDescription":"Bethesda-Silver Spring Line"},{"RouteID":"J2v1","Name":"J2 - MEDICAL CTR STA - SILVR SPRNG STA","LineDescription":"Bethesda-Silver Spring Line"},{"RouteID":"J2v2","Name":"J2 - SILVER SPRNG - MONT MALL/BATTRY LA","LineDescription":"Bethesda-Silver Spring Line"},{"RouteID":"J4","Name":"J4 - BETHESDA STA - COLLEGE PARK STA","LineDescription":"College Park-Bethesda Limited"},{"RouteID":"K12","Name":"K12 - BRANCH AVE ST-SUITLAND ST","LineDescription":"Forestville Line"},{"RouteID":"K12v1","Name":"K12 - PENN MAR - SUITLAND ST","LineDescription":"Forestville Line"},{"RouteID":"K12v2","Name":"K12 - BRANCH AVE - PENN MAR","LineDescription":"Forestville Line"},{"RouteID":"K2","Name":"K2 - FT TOTTEN STA - TAKOMA STA","LineDescription":"Takoma-Fort Totten Line"},{"RouteID":"K6","Name":"K6 - FORT TOTTEN STA - WHITE OAK","LineDescription":"New Hampshire Ave-Maryland Line"},{"RouteID":"K6v1","Name":"K6 - TLTC - FORT TOTTEN STA","LineDescription":"New Hampshire Ave-Maryland Line"},{"RouteID":"K9","Name":"K9 - FORT TOTTEN STA - FDA/FRC (PM)","LineDescription":"New Hampshire Ave-Maryland Limited Line"},{"RouteID":"K9v1","Name":"K9 - FORT TOTTEN STA - FDA/FRC (AM)","LineDescription":"New Hampshire Ave-Maryland Limited Line"},{"RouteID":"L1","Name":"L1 - POTOMAC PK  - CHEVY CHASE","LineDescription":"Connecticut Ave Line"},{"RouteID":"L2","Name":"L2 - FARRAGUT SQ - CHEVY CHASE","LineDescription":"Connecticut Ave Line"},{"RouteID":"L2v1","Name":"L2 - VAN NESS-UDC STA - CHEVY CHASE","LineDescription":"Connecticut Ave Line"},{"RouteID":"L2v2","Name":"L2 - FARRAGUT SQ - BETHESDA","LineDescription":"Connecticut Ave Line"},{"RouteID":"L8","Name":"L8 - FRIENDSHIP HTS STA - ASPEN HILL","LineDescription":"Connecticut Ave-Maryland Line"},{"RouteID":"M4","Name":"M4 - SIBLEY HOSPITAL - TENLEYTOWN/AU STA","LineDescription":"Nebraska Ave Line"},{"RouteID":"M4v1","Name":"M4 - TENLEYTOWN   - PINEHRST CIR","LineDescription":"Nebraska Ave Line"},{"RouteID":"M4v2","Name":"M4 - PINEHRST CIR - TENLEYTOWN","LineDescription":"Nebraska Ave Line"},{"RouteID":"M6","Name":"M6 - POTOMAC AVE - ALABAMA + PENN","LineDescription":"Fairfax Village Line"},{"RouteID":"M6v1","Name":"M6 - ALABAMA + PENN - FAIRFAX VILLAGE","LineDescription":"Fairfax Village Line"},{"RouteID":"MW1","Name":"MW1 - BRADDOCK RD - PENTAGON CITY","LineDescription":"Metroway-Potomac Yard Line"},{"RouteID":"MW1v1","Name":"MW1 - POTOMAC YARD - CRYSTAL CITY","LineDescription":"Metroway-Potomac Yard Line"},{"RouteID":"MW1v2","Name":"MW1 - CRYSTAL CITY - BRADDOCK RD","LineDescription":"Metroway-Potomac Yard Line"},{"RouteID":"MW1v3","Name":"MW1 - BRADDOCK RD - CRYSTAL CITY","LineDescription":"Metroway-Potomac Yard Line"},{"RouteID":"N2","Name":"N2 - FRIENDSHIP HTS - FARRAGUT SQ","LineDescription":"Massachusetts Ave Line"},{"RouteID":"N4","Name":"N4 - FRIENDSHP HTS - POTOMAC PARK","LineDescription":"Massachusetts Ave Line"},{"RouteID":"N4v1","Name":"N4 - FRIENDSHP HTS - FARRAGUT  SQ","LineDescription":"Massachusetts Ave Line"},{"RouteID":"N6","Name":"N6 - FRNDSHIP HTS - FARRAGUT SQ","LineDescription":"Massachusetts Ave Line"},{"RouteID":"NH1","Name":"NH1 - NATIONAL HARBOR-SOUTHERN AVE STA","LineDescription":"National Harbor-Southern Avenue Line"},{"RouteID":"NH2","Name":"NH2 - HUNTINGTON STA-NATIONAL HARBOR","LineDescription":"National Harbor-Alexandria Line"},{"RouteID":"P12","Name":"P12 - EASTOVER - ADDISON RD STA (NO HOSP)","LineDescription":"Eastover-Addison Road Line"},{"RouteID":"P12v1","Name":"P12 - IVERSON MALL - ADDISON RD STA","LineDescription":"Eastover-Addison Road Line"},{"RouteID":"P12v2","Name":"P12 - SUITLAND STA - ADDISON RD STA","LineDescription":"Eastover-Addison Road Line"},{"RouteID":"P18","Name":"P18 - FT WASH P+R LOT - SOUTHERN AVE","LineDescription":"Oxon Hill-Fort Washington Line"},{"RouteID":"P19","Name":"P19 - FT WASH P+R LOT-SOUTHERN AVE STA","LineDescription":"Oxon Hill-Fort Washington Line"},{"RouteID":"P6","Name":"P6 - ANACOSTIA STA - RHODE ISLAND STA","LineDescription":"Anacostia-Eckington Line"},{"RouteID":"P6v1","Name":"P6 - ECKINGTON - RHODE ISLAND AVE","LineDescription":"Anacostia-Eckington Line"},{"RouteID":"P6v2","Name":"P6 - RHODE ISLAND AVE - ECKINGTON","LineDescription":"Anacostia-Eckington Line"},{"RouteID":"P6v3","Name":"P6 - ANACOSTIA STA - ARCHIVES","LineDescription":"Anacostia-Eckington Line"},{"RouteID":"P6v4","Name":"P6 - ARCHIVES - ANACOSTIA","LineDescription":"Anacostia-Eckington Line"},{"RouteID":"Q1","Name":"Q1 - SILVR SPRNG STA - SHADY GRVE STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q2","Name":"Q2 - SILVR SPRNG STA - SHADY GRVE STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q2v1","Name":"Q2 - MONT COLLEGE - SILVR SPRNG STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q2v2","Name":"Q2 - SILVR SPRNG STA - MONT COLLEGE","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q4","Name":"Q4 - SILVER SPRNG STA - ROCKVILLE STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q4v1","Name":"Q4 - WHEATON STA - ROCKVILLE STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q5","Name":"Q5 - WHEATON STA - SHADY GRVE STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q6","Name":"Q6 - WHEATON STA - SHADY GRVE STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"Q6v1","Name":"Q6 - ROCKVILLE STA - WHEATON STA","LineDescription":"Viers Mill Road Line"},{"RouteID":"R1","Name":"R1 - FORT TOTTEN STA - ADELPHI","LineDescription":"Riggs Road Line"},{"RouteID":"R12","Name":"R12 - DEANWOOD STA - GREENBELT STA","LineDescription":"Kenilworth Avenue Line"},{"RouteID":"R12v1","Name":"R12 - DEANWOOD STA - GREENBELT STA","LineDescription":"Kenilworth Avenue Line"},{"RouteID":"R2","Name":"R2 - FORT TOTTEN - CALVERTON","LineDescription":"Riggs Road Line"},{"RouteID":"R2v1","Name":"R2 - HIGH POINT HS - FORT TOTTEN","LineDescription":"Riggs Road Line"},{"RouteID":"R2v2","Name":"R2 - POWDER MILL+CHERRY HILL - CALVERTON","LineDescription":"Riggs Road Line"},{"RouteID":"R4","Name":"R4 - BROOKLAND STA- HIGHVIEW","LineDescription":"Queens Chapel Road Line"},{"RouteID":"REX","Name":"REX - FT BELVOIR POST - KING ST STA","LineDescription":"Richmond Highway Express"},{"RouteID":"REXv1","Name":"REX - FT BELVOIR COMM HOSP - KING ST STA","LineDescription":"Richmond Highway Express"},{"RouteID":"REXv2","Name":"REX - KING ST STA - FT BELVOIR COMM HOSP","LineDescription":"Richmond Highway Express"},{"RouteID":"REXv3","Name":"REX - KING ST STA - WOODLAWN","LineDescription":"Richmond Highway Express"},{"RouteID":"REXv4","Name":"REX - WOODLAWN - KING ST STA","LineDescription":"Richmond Highway Express"},{"RouteID":"S1","Name":"S1 - NORTHERN DIVISION - POTOMAC PK","LineDescription":"16th Street-Potomac Park Line"},{"RouteID":"S1v1","Name":"S1 - VIRGINIA+E - COLORDO+16TH","LineDescription":"16th Street-Potomac Park Line"},{"RouteID":"S2","Name":"S2 - FED TRIANGLE  - SILVER SPRNG","LineDescription":"16th Street Line"},{"RouteID":"S2v1","Name":"S2 - 16TH & HARVARD - MCPHERSON SQ","LineDescription":"16th Street Line"},{"RouteID":"S35","Name":"S35 - BRANCH + RANDLE CIR - FT DUPONT","LineDescription":"Fort Dupont Shuttle Line"},{"RouteID":"S4","Name":"S4 - FED TRIANGLE - SILVER SPRNG","LineDescription":"16th Street Line"},{"RouteID":"S4v1","Name":"S4 - SILVER SPRING STA - FRANKLIN SQ","LineDescription":"16th Street Line"},{"RouteID":"S41","Name":"S41 - CARVER TERRACE - RHODE ISLAND AVE","LineDescription":"Rhode Island Ave-Carver Terrace Line"},{"RouteID":"S80","Name":"S80 - FRANCONIA-SPRNGFLD - METRO PARK","LineDescription":"Springfield Circulator-Metro Park Shuttle (TAGS)"},{"RouteID":"S80v1","Name":"S80 - FRANCONIA-SPRINGFLD - HILTON","LineDescription":"Springfield Circulator-Metro Park Shuttle (TAGS)"},{"RouteID":"S80v2","Name":"S80 - HILTON - FRANCONIA-SPRNGFLD","LineDescription":"Springfield Circulator-Metro Park Shuttle (TAGS)"},{"RouteID":"S9","Name":"S9 - FRANKLIN SQ - SILVER SPRING STA","LineDescription":"16th Street Limited Line"},{"RouteID":"S9v1","Name":"S9 - FRANKLIN SQ - COLORADO + 16TH","LineDescription":"16th Street Limited Line"},{"RouteID":"S91","Name":"S91 - FRANCONIA SPRINGFLD STA SHUTTLE","LineDescription":"Springfield Circulator-Metro Park Shuttle (TAGS)"},{"RouteID":"S91v1","Name":"S91 - FRANCONIA SPRINGFLD STA SHUTTLE","LineDescription":"Springfield Circulator-Metro Park Shuttle (TAGS)"},{"RouteID":"T14","Name":"T14 - RHD ISLND AVE STA-NEW CARRLTN STA","LineDescription":"Rhode Island Ave-New Carrollton Line"},{"RouteID":"T14v1","Name":"T14 - MT RAINIER - NEW CARRLTN STA","LineDescription":"Rhode Island Ave-New Carrollton Line"},{"RouteID":"T18","Name":"T18 - R I AVE STA - NEW CARROLLTON STA","LineDescription":"Annapolis Road Line"},{"RouteID":"T18v1","Name":"T18 - BLADENSBURG HS - NEW CARROLLTON STA","LineDescription":"Annapolis Road Line"},{"RouteID":"T2","Name":"T2 - FRIENDSHIP HTS - ROCKVILLE STA","LineDescription":"River Road Line"},{"RouteID":"U4","Name":"U4 - MINNESOTA AVE - SHERIFF RD","LineDescription":"Sheriff Road-River Terrace Line"},{"RouteID":"U4v1","Name":"U4 - SHERIFF RD - MINNESOTA STA","LineDescription":"Sheriff Road-River Terrace Line"},{"RouteID":"U4v2","Name":"U4 - RIVER TERRACE - MINNESOTA AVE","LineDescription":"Sheriff Road-River Terrace Line"},{"RouteID":"U5","Name":"U5 - MINNESOTA AVE-MARSHALL HTS","LineDescription":"Marshall Heights Line"},{"RouteID":"U6","Name":"U6 - MINNESOTA - LINCOLN HEIGHTS","LineDescription":"Marshall Heights Line"},{"RouteID":"U6v1","Name":"U6 - 37TH + RIDGE - PLUMMER ES","LineDescription":"Marshall Heights Line"},{"RouteID":"U6v2","Name":"U6 - LINCOLN HTS - E CAP + 47TH ST NE","LineDescription":"Marshall Heights Line"},{"RouteID":"U7","Name":"U7 - RIDGE + ANACOSTIA - DEANWOOD","LineDescription":"Deanwood-Minnesota Ave Line"},{"RouteID":"U7v1","Name":"U7 - MINNESOTA STA - KENILWORTH HAYES","LineDescription":"Deanwood-Minnesota Ave Line"},{"RouteID":"U7v2","Name":"U7 - KENILWORTH HAYES - MINNESOTA STA","LineDescription":"Deanwood-Minnesota Ave Line"},{"RouteID":"U7v3","Name":"U7 - MINNESOTA STA - DEANWOOD","LineDescription":"Deanwood-Minnesota Ave Line"},{"RouteID":"U7v4","Name":"U7 - DEANWOOD - MINNESOTA AVE","LineDescription":"Deanwood-Minnesota Ave Line"},{"RouteID":"V1","Name":"V1 - BUR OF ENGRVNG - BENNING HTS","LineDescription":"Benning Heights-M Street Line"},{"RouteID":"V12","Name":"V12 - SUITLAND STA - ADDISON RD STA","LineDescription":"District Heights-Suitland Line"},{"RouteID":"V14","Name":"V14 - PENN MAR - DEANWOOD STA","LineDescription":"District Heights - Seat Pleasant Line"},{"RouteID":"V14v1","Name":"V14 - PENN MAR  - ADDISON RD STA","LineDescription":"District Heights - Seat Pleasant Line"},{"RouteID":"V2","Name":"V2 - ANACOSTIA - CAPITOL HGTS","LineDescription":"Capitol Heights-Minnesota Avenue Line"},{"RouteID":"V2v1","Name":"V2 - MINNESOTA AVE - ANACOSTIA","LineDescription":"Capitol Heights-Minnesota Avenue Line"},{"RouteID":"V4","Name":"V4 - 1ST + K ST SE - CAPITOL HGTS","LineDescription":"Capitol Heights-Minnesota Avenue Line"},{"RouteID":"V4v1","Name":"V4 - MINN AVE STA-CAPITOL HGTS","LineDescription":"Capitol Heights-Minnesota Avenue Line"},{"RouteID":"V7","Name":"V7 - CONGRESS HGTS - MINN STA","LineDescription":"Benning Heights-Alabama Ave Line"},{"RouteID":"V8","Name":"V8 - BENNING HGTS - MINN STA","LineDescription":"Benning Heights-Alabama Ave Line"},{"RouteID":"W1","Name":"W1 - FORT DRUM  - SOUTHERN AVE STA","LineDescription":"Shipley Terrace - Fort Drum Line"},{"RouteID":"W14","Name":"W14 - FT WASHINGTON-SOUTHERN AVE STA","LineDescription":"Bock Road Line"},{"RouteID":"W14v1","Name":"W14 - SOUTHERN AVE - FRIENDLY","LineDescription":"Bock Road Line"},{"RouteID":"W14v2","Name":"W14 - FRIENDLY - SOUTHERN AVE","LineDescription":"Bock Road Line"},{"RouteID":"W2","Name":"W2 - MALCOM X+OAKWD - UNITED MEDICAL CTR","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W2v1","Name":"W2 - NYLDR+GOOD HOP- HOWARD+ANACOSTIA","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W2v2","Name":"W2 - NAYLOR+GOODHOPE - MALCM X+OAKWOOD","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W2v3","Name":"W2 - HOWRD+ANACOSTIA-NAYLOR+GOODHOPE","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W2v4","Name":"W2 - ANACOSTIA STA - UNITED MEDICAL CTR","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W2v5","Name":"W2 - ANACOSTIA STA - SOUTHERN AVE STA","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W2v6","Name":"W2 - MELLN+M L KNG - UNITED MEDICAL CTR","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W2v7","Name":"W2 - SOUTHERN AVE-WASHINGTON OVERLOOK","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W3","Name":"W3 - MALCOM X+OAKWD - UNITED MEDICAL CTR","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W3v1","Name":"W3 - MELLN+M L KNG - UNITED MEDICAL CTR","LineDescription":"United Medical Center-Anacostia Line"},{"RouteID":"W4","Name":"W4 - ANACOSTIA STA - DEANWOOD STA","LineDescription":"Deanwood-Alabama Ave Line"},{"RouteID":"W4v1","Name":"W4 - MALCOLM X & PORTLAND - DEANWOOD","LineDescription":"Deanwood-Alabama Ave Line"},{"RouteID":"W4v2","Name":"W4 - DEANWOOD STA - MALCLM X+ ML KING","LineDescription":"Deanwood-Alabama Ave Line"},{"RouteID":"W45","Name":"W45 - TENLEYTOWN STA - 16TH + SHEPHERD","LineDescription":"Mt Pleasant-Tenleytown Line"},{"RouteID":"W47","Name":"W47 - TENLEYTOWN STA - COLUMBIA HTS STA","LineDescription":"Mt Pleasant-Tenleytown Line"},{"RouteID":"W5","Name":"W5 - DC VILLAGE-USCG (VIA ANAC)","LineDescription":"Anacostia-Fort Drum Line"},{"RouteID":"W6","Name":"W6 - ANACOSTIA - ANACOSTIA","LineDescription":"Garfield-Anacostia Loop Line"},{"RouteID":"W6v1","Name":"W6 - NAYLOR+GOODHOPE - ANACOSTIA","LineDescription":"Garfield-Anacostia Loop Line"},{"RouteID":"W8","Name":"W8 - ANACOSTIA - ANACOSTIA","LineDescription":"Garfield-Anacostia Loop Line"},{"RouteID":"W8v1","Name":"W8 - ANACOSTIA-NAYLOR+GOODHOPE","LineDescription":"Garfield-Anacostia Loop Line"},{"RouteID":"W8v2","Name":"W8 - NAYLOR+GOODHOPE - ANACOSTIA","LineDescription":"Garfield-Anacostia Loop Line"},{"RouteID":"X1","Name":"X1 - FOGGY BOTTOM+GWU- MINNESOTA STA","LineDescription":"Benning Road Line"},{"RouteID":"X2","Name":"X2 - LAFAYETTE SQ - MINNESOTA STA","LineDescription":"Benning Road-H Street Line"},{"RouteID":"X2v1","Name":"X2 - PHELPS HS - MINNESOTA STA","LineDescription":"Benning Road-H Street Line"},{"RouteID":"X2v2","Name":"X2 - FRIENDSHIP EDISON PCS-MINNESOTA STA","LineDescription":"Benning Road-H Street Line"},{"RouteID":"X2v3","Name":"X2 - PHELPS HS - LAFAYTTE SQ","LineDescription":"Benning Road-H Street Line"},{"RouteID":"X3","Name":"X3 - DUKE ELLINGTON BR - MINNESOTA STN","LineDescription":"Benning Road Line"},{"RouteID":"X3v1","Name":"X3 - KIPP DC - MINNESOTA AVE STN","LineDescription":"Benning Road Line"},{"RouteID":"X8","Name":"X8 - UNION STATION - CARVER TERR","LineDescription":"Maryland Ave Line"},{"RouteID":"X9","Name":"X9 - NY AVE & 12TH NW - CAPITOL HTS STA","LineDescription":"Benning Road-H Street Limited Line"},{"RouteID":"X9v1","Name":"X9 - NY AVE & 12TH ST NW - MINNESOTA AVE","LineDescription":"Benning Road-H Street Limited Line"},{"RouteID":"X9v2","Name":"X9 - MINNESOTA AVE ST - NY AVE & 12TH ST","LineDescription":"Benning Road-H Street Limited Line"},{"RouteID":"Y2","Name":"Y2 - SILVER SPRING STA - MONTG MED CTR","LineDescription":"Georgia Ave-Maryland Line"},{"RouteID":"Y7","Name":"Y7 - SILVER SPRING STA - ICC P&R","LineDescription":"Georgia Ave-Maryland Line"},{"RouteID":"Y8","Name":"Y8 - SILVER SPR STA - MONTGOMERY MED CTR","LineDescription":"Georgia Ave-Maryland Line"},{"RouteID":"Z11","Name":"Z11 - SILVR SPRING - BURTONSVILLE P&R","LineDescription":"Greencastle-Briggs Chaney Express Line"},{"RouteID":"Z11v1","Name":"Z11 - GREENCASTLE - SILVR SPRING","LineDescription":"Greencastle-Briggs Chaney Express Line"},{"RouteID":"Z2","Name":"Z2 - SILVR SPRING - OLNEY (NO BLAKE HS)","LineDescription":"Colesville - Ashton Line"},{"RouteID":"Z2v1","Name":"Z2 - SILVER SPRING - BONIFANT & NH","LineDescription":"Colesville - Ashton Line"},{"RouteID":"Z2v2","Name":"Z2 - SILVER SPRING STA-OLNEY (BLAKE HS)","LineDescription":"Colesville - Ashton Line"},{"RouteID":"Z2v3","Name":"Z2 - COLESVILLE-SILVER SPRING","LineDescription":"Colesville - Ashton Line"},{"RouteID":"Z6","Name":"Z6 - SILVR SPRNG STA - BURTONSVILLE","LineDescription":"Calverton-Westfarm Line"},{"RouteID":"Z6v1","Name":"Z6 - CASTLE BLVD - SILVER SPRING STA","LineDescription":"Calverton-Westfarm Line"},{"RouteID":"Z6v2","Name":"Z6 - SILVER SPRING -CASTLE BLVD","LineDescription":"Calverton-Westfarm Line"},{"RouteID":"Z7","Name":"Z7 - SLVER SPRNG STA-S LAUREL P&R (4COR)","LineDescription":"Laurel-Burtonsville Express Line"},{"RouteID":"Z7v1","Name":"Z7 - S LAUREL P&R-SILVR SPR (NO 4COR)","LineDescription":"Laurel-Burtonsville Express Line"},{"RouteID":"Z8","Name":"Z8 - SILVER SRING STA - BRIGGS CHANEY","LineDescription":"Fairland Line"},{"RouteID":"Z8v1","Name":"Z8 - WHITE OAK - SILVER SPRING STA","LineDescription":"Fairland Line"},{"RouteID":"Z8v2","Name":"Z8 - SILVR SPRNG - CSTLE BLVD VERIZON","LineDescription":"Fairland Line"},{"RouteID":"Z8v3","Name":"Z8 - SILVER SPRING STA - CASTLE BLVD","LineDescription":"Fairland Line"},{"RouteID":"Z8v4","Name":"Z8 - SILVER SPRING - GRNCSTLE (VERIZON)","LineDescription":"Fairland Line"},{"RouteID":"Z8v5","Name":"Z8 - SILVER SPRING STA - GREENCASTLE","LineDescription":"Fairland Line"},{"RouteID":"Z8v6","Name":"Z8 - SILVER SPRING STA - WHITE OAK","LineDescription":"Fairland Line"}]}`,
			unmarshalledResponse: &GetRoutesResponse{
				Routes: []Route{
					{
						RouteID:         "10A",
						Name:            "10A - HUNTINGTON STA - PENTAGON",
						LineDescription: "Alexandria-Pentagon Line",
					},
					{
						RouteID:         "10B",
						Name:            "10B - HUNTING POINT - BALLSTON STA",
						LineDescription: "Hunting Point-Ballston Line",
					},
					{
						RouteID:         "10E",
						Name:            "10E - HUNTING POINT - PENTAGON",
						LineDescription: "Alexandria-Pentagon Line",
					},
					{
						RouteID:         "10N",
						Name:            "10N - NATIONAL AIRPORT - PENTAGON",
						LineDescription: "Alexandria-Pentagon Line",
					},
					{
						RouteID:         "11Y",
						Name:            "11Y - MOUNT VERNON - POTOMAC PARK",
						LineDescription: "Mt Vernon Express Line",
					},
					{
						RouteID:         "11Yv1",
						Name:            "11Y - POTOMAC PARK - HUNTING POINT",
						LineDescription: "Mt Vernon Express Line",
					},
					{
						RouteID:         "11Yv2",
						Name:            "11Y - HUNTING POINT - POTOMAC PARK",
						LineDescription: "Mt Vernon Express Line",
					},
					{
						RouteID:         "15K",
						Name:            "15K - E FALLS CH STA-MCLEAN LANGLEY",
						LineDescription: "Chain Bridge Road Line",
					},
					{
						RouteID:         "15Kv1",
						Name:            "15K - CIA - EAST FALLS CHURCH STA",
						LineDescription: "Chain Bridge Road Line",
					},
					{
						RouteID:         "16A",
						Name:            "16A - PATRIOT+AMERICANA - PENTAGON",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16C",
						Name:            "16C - CULMORE - PENTAGON",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16Cv1",
						Name:            "16C - CULMORE - FEDERAL TRIANGLE",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16E",
						Name:            "16E - GLN CARLYN+VISTA - FRANKLIN SQ",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16G",
						Name:            "16G - DINWIDDIE+COLUMBIA - PENTAGON CITY",
						LineDescription: "Columbia Pike-Pentagon City Line",
					},
					{
						RouteID:         "16Gv1",
						Name:            "16G - COL PIKE+CARLIN SPR - DINWD+COL PK",
						LineDescription: "Columbia Pike-Pentagon City Line",
					},
					{
						RouteID:         "16H",
						Name:            "16H - SKYLINE CITY - PENTAGON CITY STA",
						LineDescription: "Columbia Pike-Pentagon City Line",
					},
					{
						RouteID:         "16L",
						Name:            "16L - ANNANDALE - PENTAGON HOV",
						LineDescription: "Annandale-Skyline City-Pentagon Line",
					},
					{
						RouteID:         "16Y",
						Name:            "16Y - FOUR MILE+COLUMBIA - MCPHERSON SQ",
						LineDescription: "Columbia Pike-Farragut Square Line",
					},
					{
						RouteID:         "16Yv1",
						Name:            "16Y - COLUMBIA PIKE - MCPHERSON SQ",
						LineDescription: "Columbia Pike-Farragut Square Line",
					},
					{
						RouteID:         "17B",
						Name:            "17B - BURKE CENTRE - PENTAGON HOV",
						LineDescription: "Kings Park-North Springfield Line",
					},
					{
						RouteID:         "17G",
						Name:            "17G - G MASON UNIV - PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17H",
						Name:            "17H - TWNBRK RD+TWNBRK RN- PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17K",
						Name:            "17K - TWNBRK RD+TWNBRK RN- PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17L",
						Name:            "17L - TWNBRK RD+TWNBRK RUN-PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17M",
						Name:            "17M - EDSALL+CANARD - PENTAGON",
						LineDescription: "Kings Park-North Springfield Line",
					},
					{
						RouteID:         "18G",
						Name:            "18G - ROLLING VALLEY - PENTAGON",
						LineDescription: "Orange Hunt Line",
					},
					{
						RouteID:         "18H",
						Name:            "18H - HUNTSMAN+CORK CTY - PENTAGON",
						LineDescription: "Orange Hunt Line",
					},
					{
						RouteID:         "18J",
						Name:            "18J - ROLLING VALLEY - PENTAGON",
						LineDescription: "Orange Hunt Line",
					},
					{
						RouteID:         "18P",
						Name:            "18P - BURKE CENTRE - PENTAGON",
						LineDescription: "Burke Centre Line",
					},
					{
						RouteID:         "18Pv1",
						Name:            "18P - PENTAGON - ROLLING VALLEY MALL",
						LineDescription: "Burke Centre Line",
					},
					{
						RouteID:         "18Pv2",
						Name:            "18P - ROLLING VALLEY P+R - PENTAGON",
						LineDescription: "Burke Centre Line",
					},
					{
						RouteID:         "1A",
						Name:            "1A - VIENNA-BALLSTON (7 CORNERS)",
						LineDescription: "Wilson Blvd-Vienna Line",
					},
					{
						RouteID:         "1B",
						Name:            "1B - DUNN LORING  - BALLSTON",
						LineDescription: "Wilson Blvd-Vienna Line",
					},
					{
						RouteID:         "1C",
						Name:            "1C - WEST OX DIV-DUNN LORING (VIA MALL)",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv1",
						Name:            "1C - FAIRFAX CO GOV CTR - DUNN LORING",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv2",
						Name:            "1C - WEST OX DIV - DUNN LORING (NO MALL)",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv3",
						Name:            "1C - FAIR OAKS MALL - DUNN LORING",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv4",
						Name:            "1C - DUNN LORING - FAIR OAKS MALL",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "21A",
						Name:            "21A - S REYNOLDS+EOS 21 CONDOS - PENTAGON",
						LineDescription: "Landmark-Bren Mar Park-Pentagon Line",
					},
					{
						RouteID:         "21D",
						Name:            "21D - LANDMARK MEWS -PENTAGON",
						LineDescription: "Landmark-Bren Mar Park-Pentagon Line",
					},
					{
						RouteID:         "22A",
						Name:            "22A - BALLSTON STA - PENTAGON",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "22Av1",
						Name:            "22A - SHIRLINGTON - BALLSTON  STA",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "22C",
						Name:            "22C - BALLSTON STA - PENTAGON",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "22F",
						Name:            "22F - NVCC - PENTAGON VIA HOV",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "23A",
						Name:            "23A - TYSONS CORNER CTR - CRYSTAL CTY",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "23B",
						Name:            "23B - BALLSTON STA - CRYSTAL CTY",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "23Bv1",
						Name:            "23B - LINDEN RESOURCES - BALLSTON STATION",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "23T",
						Name:            "23T - TYSONS CORNER CTR - SHIRLINGTON",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "25B",
						Name:            "25B - VAN DORN - BALLSTON",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "25Bv1",
						Name:            "25B - SOUTHERN TOWERS - BALLSTON",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "25Bv2",
						Name:            "25B - VAN DORN - BALLSTON/NO LDMRK CTR",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "25Bv3",
						Name:            "25B - BALLSTON - SOUTHERN TOWERS",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "26A",
						Name:            "26A - NVCC ANNANDALE - E FALLS CHURCH STA",
						LineDescription: "Annandale-East Falls Church Line",
					},
					{
						RouteID:         "28A",
						Name:            "28A - TYSONS CORNER STA-KING ST STA",
						LineDescription: "Leesburg Pike Line",
					},
					{
						RouteID:         "28Av1",
						Name:            "28A - SOUTHERN TOWERS-TYSONS CORNER STA",
						LineDescription: "Leesburg Pike Line",
					},
					{
						RouteID:         "28F",
						Name:            "28F - BLDG 5113 G MASON DR - PENTAGON",
						LineDescription: "Skyline City Line",
					},
					{
						RouteID:         "28G",
						Name:            "28G - BLDG 5113 G MASON DR - PENTAGON",
						LineDescription: "Skyline City Line",
					},
					{
						RouteID:         "29C",
						Name:            "29C - NVCC ANNANDALE - PENTAGON",
						LineDescription: "Annandale Line",
					},
					{
						RouteID:         "29G",
						Name:            "29G - AMERICANA+HERITAGE - PENTAGON",
						LineDescription: "Annandale Line",
					},
					{
						RouteID:         "29K",
						Name:            "29K - GMU - KING ST STA",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29Kv1",
						Name:            "29K - GMU - KING ST/NO LDMRK",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29N",
						Name:            "29N - VIENNA STA - KING ST (VIA MALL)",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29Nv1",
						Name:            "29N - VIENNA STA-KING ST (NO MALL)",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29W",
						Name:            "29W - NVCC ANNANDALE - PENTAGON",
						LineDescription: "Braeburn Drive-Pentagon Express Line",
					},
					{
						RouteID:         "2A",
						Name:            "2A - DUNN LORING STA - BALLSTON STA",
						LineDescription: "Washington Blvd.-Dunn Loring Line",
					},
					{
						RouteID:         "2B",
						Name:            "2B - WEST OX RD DIV-DUNN LORING STATION",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "2Bv1",
						Name:            "2B - WEST OX RD-DUNN LORING STA(NO MALL)",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "2Bv2",
						Name:            "2B - FAIR OAKS MALL-DUNN LORING STATION",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "2Bv3",
						Name:            "2B - DUNN LORING STA - FAIR OAKS MALL",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "30N",
						Name:            "30N - FRIENDSHIP HGTS- NAYLOR RD STA",
						LineDescription: "Friendship Heights-Southeast Line",
					},
					{
						RouteID:         "30S",
						Name:            "30S - FRIENDSHIP HGTS- SOUTHRN AVE STA",
						LineDescription: "Friendship Heights-Southeast Line",
					},
					{
						RouteID:         "31",
						Name:            "31 - POTOMAC PARK-FRIENDSHIP HGTS",
						LineDescription: "Wisconsin Avenue Line",
					},
					{
						RouteID:         "32",
						Name:            "32 - VIRGINIA AVE+E ST- SOUTHRN AVE",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "32v1",
						Name:            "32 - PENN AVE + 8TH ST - SOUTHRN AVE",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "33",
						Name:            "33 - 10TH ST+PA AV NW - FRIENDSHIP HGTS",
						LineDescription: "Wisconsin Avenue Line",
					},
					{
						RouteID:         "34",
						Name:            "34 - 10TH ST + PA AVE- NAYLOR RD STA",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "36",
						Name:            "36 - VIRGINIA AVE+E ST - NAYLOR RD STA",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "37",
						Name:            "37 - 10TH ST+PA AV NW - FRIENDSHIP HGTS",
						LineDescription: "Wisconsin Avenue Limited Line",
					},
					{
						RouteID:         "38B",
						Name:            "38B - BALLSTON - FARRAGUT",
						LineDescription: "Ballston-Farragut Square Line",
					},
					{
						RouteID:         "38Bv1",
						Name:            "38B - WASH & QUINCY - FARRAGUT",
						LineDescription: "Ballston-Farragut Square Line",
					},
					{
						RouteID:         "38Bv2",
						Name:            "38B - WASHINGTON-LEE HS - FARRAGUT SQ",
						LineDescription: "Ballston-Farragut Square Line",
					},
					{
						RouteID:         "39",
						Name:            "39 - VIRGINIA AVE+21ST NW- NAYLOR RD STA",
						LineDescription: "Pennsylvania Avenue Limited Line",
					},
					{
						RouteID:         "3A",
						Name:            "3A - ANNANDALE - E FALLS CHURCH",
						LineDescription: "Annandale Road Line",
					},
					{
						RouteID:         "3Av1",
						Name:            "3A - ANNANDALE - 7 CORNERS",
						LineDescription: "Annandale Road Line",
					},
					{
						RouteID:         "3T",
						Name:            "3T - MCLEAN STATION - E FALLS CH STA",
						LineDescription: "Pimmit Hills Line",
					},
					{
						RouteID:         "3Tv1",
						Name:            "3T - MCLEAN STATION - W FALLS CHURCH",
						LineDescription: "Pimmit Hills Line",
					},
					{
						RouteID:         "3Y",
						Name:            "3Y - E FALLS CHURCH - MCPHERSON SQ",
						LineDescription: "Lee Highway-Farragut Square Line",
					},
					{
						RouteID:         "42",
						Name:            "42 - 9TH + F ST  - MT PLEASANT",
						LineDescription: "Mount Pleasant Line",
					},
					{
						RouteID:         "43",
						Name:            "43 - I + 13TH NW - MT PLEASANT",
						LineDescription: "Mount Pleasant Line",
					},
					{
						RouteID:         "4A",
						Name:            "4A - SEVEN CORNERS - ROSSLYN",
						LineDescription: "Pershing Dr-Arlington Blvd Line",
					},
					{
						RouteID:         "4B",
						Name:            "4B - SEVEN CORNERS - ROSSLYN",
						LineDescription: "Pershing Dr-Arlington Blvd Line",
					},
					{
						RouteID:         "52",
						Name:            "52 - L ENFNT PLAZA - TAKOMA STATION",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "52v1",
						Name:            "52 - L ENFNT PLAZA - 14TH+COLORADO",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "52v2",
						Name:            "52 - 14TH+COLORADO - L ENFANT PLAZA",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "52v3",
						Name:            "52 - 14TH & U - TAKOMA STATION",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "54",
						Name:            "54 - METRO CENTER - TAKOMA STA",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "54v1",
						Name:            "54 - 14TH+COLORADO - METRO CENTER",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "59",
						Name:            "59 - FEDERAL TRIANGLE - TAKOMA STATION",
						LineDescription: "14th Street Limited Line",
					},
					{
						RouteID:         "5A",
						Name:            "5A - DULLES AIRPORT - LENFANT PLAZA",
						LineDescription: "DC-Dulles Line",
					},
					{
						RouteID:         "60",
						Name:            "60 - GEORGIA + PETWORTH - FT TOTTEN",
						LineDescription: "Fort Totten-Petworth Line",
					},
					{
						RouteID:         "62",
						Name:            "62 - GEORGIA+PETWORTH - TAKOMA STATION",
						LineDescription: "Takoma-Petworth Line",
					},
					{
						RouteID:         "62v1",
						Name:            "62 - COOLIDGE HS - GEORGIA + PETWORTH",
						LineDescription: "Takoma-Petworth Line",
					},
					{
						RouteID:         "63",
						Name:            "63 - FED TRIANGLE - TAKOMA STA",
						LineDescription: "Takoma-Petworth Line",
					},
					{
						RouteID:         "64",
						Name:            "64 - FEDERAL TRIANGLE -FORT TOTTEN",
						LineDescription: "Fort Totten-Petworth Line",
					},
					{
						RouteID:         "64v1",
						Name:            "64 - GEORGIA + PETWOTH - FT TOTTEN",
						LineDescription: "Fort Totten-Petworth Line",
					},
					{
						RouteID:         "70",
						Name:            "70 - ARCHIVES - SILVER SPRING",
						LineDescription: "Georgia Avenue-7th Street Line",
					},
					{
						RouteID:         "70v1",
						Name:            "70 - GEORGIA & EUCLID  - ARCHIVES",
						LineDescription: "Georgia Avenue-7th Street Line",
					},
					{
						RouteID:         "74",
						Name:            "74 - NATIONALS PARK - CONVENTION CTR",
						LineDescription: "Convention Center-Southwest Waterfront Line",
					},
					{
						RouteID:         "79",
						Name:            "79 - ARCHIVES - SILVER SPRING STA",
						LineDescription: "Georgia Avenue MetroExtra Line",
					},
					{
						RouteID:         "7A",
						Name:            "7A - LINCOLNIA+QUANTRELL - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Av1",
						Name:            "7A - PENTAGON - SOUTHERN TWRS",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Av2",
						Name:            "7A - SOUTHERN TWRS - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Av3",
						Name:            "7A - LINCOLNIA/QUANTRELL - PENT VIA PENT",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7C",
						Name:            "7C - PARK CENTER - PENTAGON",
						LineDescription: "Park Center-Pentagon Line",
					},
					{
						RouteID:         "7F",
						Name:            "7F - LINCOLNIA+QUANTRELL - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Fv1",
						Name:            "7F - LINC + QUANT - PENT CITY - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7M",
						Name:            "7M - MARK CENTER - PENTAGON (NON-HOV)",
						LineDescription: "Mark Center-Pentagon Line",
					},
					{
						RouteID:         "7Mv1",
						Name:            "7M - MARK CENTER - PENTAGON (HOV)",
						LineDescription: "Mark Center-Pentagon Line",
					},
					{
						RouteID:         "7P",
						Name:            "7P - PARK CTR - PENTAGON",
						LineDescription: "Park Center-Pentagon Line",
					},
					{
						RouteID:         "7W",
						Name:            "7W - LNCLNA+QUANTRLL- PENTAGON",
						LineDescription: "Lincolnia-Pentagon Line",
					},
					{
						RouteID:         "7Y",
						Name:            "7Y - SOUTHERN TWRS - H+17TH ST",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Yv1",
						Name:            "7Y - PENTAGON - H+17TH ST",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "80",
						Name:            "80 - KENNEDY CTR   - FORT TOTTEN STA",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "80v1",
						Name:            "80 - MCPHERSON SQ  - BROOKLAND",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "80v2",
						Name:            "80 - MCPHERSON SQ  - FORT TOTTEN STA",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "80v3",
						Name:            "80 - KENNEDY CTR   - BROOKLAND STA",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "83",
						Name:            "83 - RHODE ISLAND AVE STA-CHERRY HILL",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v1",
						Name:            "83 - MT RAINIER - RHODE ISLAND",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v2",
						Name:            "83 - RHODE ISLAND - MT RAINIER",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v3",
						Name:            "83 - RHODE ISLAND AVE STA-COLLEGE PARK",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v4",
						Name:            "83 - COLLEGE PARK-RHODE ISLAND AVE STA",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "86",
						Name:            "86 - RHODE ISLAND AVE STA- CALVERTON",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "86v1",
						Name:            "86 - RHODE ISLAND AVE STA- COLLEGE PARK",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "86v2",
						Name:            "86 - MT RAINIER   - CALVERTON",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "87",
						Name:            "87 - NEW CARROLTON -CYPRESS+LAURL LAKES",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v1",
						Name:            "87 - GRNBELT STA -CYPRESS+LAURL LAKES",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v2",
						Name:            "87 - GRNBELT-CYP+LRL LAKES (NO P&R)",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v3",
						Name:            "87 - GRNBELT STA - BALTIMORE+MAIN ST",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v4",
						Name:            "87 - BALTIMORE+MAIN ST - GRNBELT STA",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v5",
						Name:            "87 - CYPRESS+LAURL LAKES -GRNBELT STA",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "89",
						Name:            "89 - GREENBELT STA - 4TH & GREEN HILL",
						LineDescription: "Laurel Line",
					},
					{
						RouteID:         "89v1",
						Name:            "89 - GREENBELT STA - CHERRY LA+4TH ST",
						LineDescription: "Laurel Line",
					},
					{
						RouteID:         "89M",
						Name:            "89M - GREENBELT STA - S LAUREL P+R",
						LineDescription: "Laurel Line",
					},
					{
						RouteID:         "8S",
						Name:            "8S - RADFORD+QUAKER - PENTAGON",
						LineDescription: "Foxchase-Seminary Valley Line",
					},
					{
						RouteID:         "8W",
						Name:            "8W - MARK CENTER - PENTAGON V FOXCHASE",
						LineDescription: "Foxchase-Seminary Valley Line",
					},
					{
						RouteID:         "8Z",
						Name:            "8Z - QUAKER+OSAGE - PENTAGON",
						LineDescription: "Foxchase-Seminary Valley Line",
					},
					{
						RouteID:         "90",
						Name:            "90 - ANACOSTIA - DK ELLNGTN BRDG",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "90v1",
						Name:            "90 - 8TH ST + L ST  - DK ELLNGTN BRDG",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "90v2",
						Name:            "90 - KIPP DC PREP- ANACOSTIA",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "92",
						Name:            "92 - CONGRESS HTS STA - REEVES CTR",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "92v1",
						Name:            "92 - EASTERN MARKET - CONGRESS HGTS STA",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "92v2",
						Name:            "92 - CONGRESS HTS STA- EASTERN MARKET",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "96",
						Name:            "96 - TENLEYTOWN STA - CAPITOL HTS STA",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v1",
						Name:            "96 - ELLINGTON BR - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v2",
						Name:            "96 - CAPITOL HTS  - REEVES CTR",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v3",
						Name:            "96 - CAPITOL HTS - ELLINGTON BRDG",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v4",
						Name:            "96 - REEVES CTR  - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v5",
						Name:            "96 - TENLEYTOWN STA - STADIUM ARMORY STA",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "97",
						Name:            "97 - UNION STATION - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "97v1",
						Name:            "97 - EASTERN HS - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "A12",
						Name:            "A12 - ADDISON RD STA - CAPITAL PLAZA",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A12v1",
						Name:            "A12 - BARLWE+MATHEW HENSON - ADDISON RD",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A12v2",
						Name:            "A12 - CAPITOL HTS STA - CAPITAL PLAZA",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A12v3",
						Name:            "A12 - CAPITAL PLAZA - CAPITOL HTS",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A2",
						Name:            "A2 - SOUTHERN AVE - ANACOSTIA (VIA HOSP)",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A2v1",
						Name:            "A2 - ANACOSTIA - MISS+ATLANTIC",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A2v2",
						Name:            "A2 - SOUTHERN AVE STA - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A2v3",
						Name:            "A2 - MISS+ATLANTIC - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A31",
						Name:            "A31 - ANACOSTIA HIGH - MINNESOTA STA",
						LineDescription: "Minnesota Ave-Anacostia Line",
					},
					{
						RouteID:         "A32",
						Name:            "A32 - ANACOSTIA HIGH - SOUTHRN AVE STA",
						LineDescription: "Minnesota Ave-Anacostia Line",
					},
					{
						RouteID:         "A33",
						Name:            "A33 - ANACOSTIA HIGH - ANACOSTIA STA",
						LineDescription: "Minnesota Ave-Anacostia Line",
					},
					{
						RouteID:         "A4",
						Name:            "A4 - DC VILLAGE - ANACOSTIA",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v1",
						Name:            "A4 - USCG-FT DRUM (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v2",
						Name:            "A4 - FT DRUM - ANACOSTIA",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v3",
						Name:            "A4 - ANACOSTIA - FT DRUM",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v4",
						Name:            "A4 - FT DRUM-USCG (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v5",
						Name:            "A4 - DC VILL-USCG (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A6",
						Name:            "A6 - 4501 3RD ST - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A6v1",
						Name:            "A6 - SOUTHERN AVE+S CAPITOL - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A7",
						Name:            "A7 - SOUTHRN+S CAPITOL - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A8",
						Name:            "A8 - 4501 3RD ST - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A8v1",
						Name:            "A8 - SOUTHERN + S CAPITOL - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A9",
						Name:            "A9 - LIVINGSTON - MCPHERSON SQUARE",
						LineDescription: "Martin Luther King Jr Ave Limited Line",
					},
					{
						RouteID:         "B2",
						Name:            "B2 - ANACOSTIA STA - MOUNT RAINIER",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v1",
						Name:            "B2 - ANACOSTIA STA - BLDNSGRG+VST NE",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v2",
						Name:            "B2 - POTOMAC AVE - MOUNT RAINIER",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v3",
						Name:            "B2 - BLDNSBRG+26TH - ANACOSTIA STA",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v4",
						Name:            "B2 - EASTERN HS - ANACOSTIA STA",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B21",
						Name:            "B21 - NEW CARROLTON STA - BOWIE STATE",
						LineDescription: "Bowie State University Line",
					},
					{
						RouteID:         "B22",
						Name:            "B22 - NEW CARROLTON STA - BOWIE STATE",
						LineDescription: "Bowie State University Line",
					},
					{
						RouteID:         "B22v1",
						Name:            "B22 - OLD CHAPEL & 197 - NEW CARRLTN STA",
						LineDescription: "Bowie State University Line",
					},
					{
						RouteID:         "B24",
						Name:            "B24 - NEW CARLTN STA-BOWIE P+R(VIA BHC)",
						LineDescription: "Bowie-Belair Line",
					},
					{
						RouteID:         "B24v1",
						Name:            "B24 - NEW CARROLTON STA - BOWIE P+R",
						LineDescription: "Bowie-Belair Line",
					},
					{
						RouteID:         "B27",
						Name:            "B27 - NEW CARROLTON STA - BOWIE STATE",
						LineDescription: "Bowie-New Carrollton Line",
					},
					{
						RouteID:         "B29",
						Name:            "B29 - NEW CARROLTON - CROFTON CC (PM)(PR)",
						LineDescription: "Crofton-New Carrollton Line",
					},
					{
						RouteID:         "B29v1",
						Name:            "B29 - NEW CARROLLTON STA - GATEWAY CTR",
						LineDescription: "Crofton-New Carrollton Line",
					},
					{
						RouteID:         "B29v2",
						Name:            "B29 - NEW CARROLTON - CROFTON CC (NO PR)",
						LineDescription: "Crofton-New Carrollton Line",
					},
					{
						RouteID:         "B30",
						Name:            "B30 - GREENBELT STA - BWI LT RAIL STA",
						LineDescription: "Greenbelt-BWI Thurgood Marshall Airport Express Line",
					},
					{
						RouteID:         "B8",
						Name:            "B8 - RHODE ISLAND AV -PETERSBRG APTS",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "B8v1",
						Name:            "B8 - BLDNSBRG+S DKTA -PETERSBRG APTS",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "B8v2",
						Name:            "B8 - PETERSBRG APTS  -BLDNSBRG+S DKTA",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "B9",
						Name:            "B9 - RHODE ISLAND AVE - COLMAR MANOR",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "C11",
						Name:            "C11 - CLINTON P+R - BRANCH AVE STA",
						LineDescription: "Clinton Line",
					},
					{
						RouteID:         "C12",
						Name:            "C12 - NAYLOR RD STA  - BRANCH AVE STA",
						LineDescription: "Hillcrest Heights Line",
					},
					{
						RouteID:         "C13",
						Name:            "C13 - CLINTON P+R - BRANCH AVE STA",
						LineDescription: "Clinton Line",
					},
					{
						RouteID:         "C14",
						Name:            "C14 - NAYLOR RD STA  - BRANCH AVE STA",
						LineDescription: "Hillcrest Heights Line",
					},
					{
						RouteID:         "C2",
						Name:            "C2 - WHEATN STA - GRNBELT STA UMD ALT",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C2v1",
						Name:            "C2 - TAKOMA LANGLEY XROADS - GRNBELT STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C2v2",
						Name:            "C2 - GREENBELT STA - RANDOLPH + PARKLAWN",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C2v3",
						Name:            "C2 - TAKOMA LANGLEY XROADS - WHEATON",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C21",
						Name:            "C21 - ADDISON RD STA - COLLINGTON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C21v1",
						Name:            "C21 - ADDISON RD STA - POINTER RIDGE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C21v2",
						Name:            "C21 - ADDISON RD STA - CAMPUS WAY S",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C22",
						Name:            "C22 - ADDISON RD STA - COLLINGTON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C22v1",
						Name:            "C22 - ADDISON RD STA- POINTER RIDGE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C26",
						Name:            "C26 - LARGO TOWN CTR - WATKNS+CHESTERTON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C26v1",
						Name:            "C26 - LARGO TOWN CTR-WATKINS+CAMBLETON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C28",
						Name:            "C28 - PT RIDGE -NEW CAROLTN VIA HLTH CTR",
						LineDescription: "Pointer Ridge Line",
					},
					{
						RouteID:         "C28v1",
						Name:            "C28 - PT RIDGE - NEW CARROLLTON STA",
						LineDescription: "Pointer Ridge Line",
					},
					{
						RouteID:         "C29*1",
						Name:            "C29 - POINTER RIDGE - ADDISON RD STA",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C29*2",
						Name:            "C29 - WATKNS+CAMBLETON - ADDISN RD STA",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C29*4",
						Name:            "C29 - ADDISON RD STA - BOWIE STATE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C29/",
						Name:            "C29 - ADDISON RD STA - POINTER RIDGE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C4",
						Name:            "C4 - TWINBROOK STA - PG PLAZA STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C4v1",
						Name:            "C4 - TLTC-TWINBROOK",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C4v2",
						Name:            "C4 - TWINBROOK STA - WHEATON STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C4v3",
						Name:            "C4 - PG PLAZA STA - WHEATON STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C8",
						Name:            "C8 - WHITE FLINT - COLLEGE PARK",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "C8v1",
						Name:            "C8 - WHITE FLNT-COLLEGE PK (NO FDA/ARCH)",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "C8v2",
						Name:            "C8 - WHITE FLINT-COLLEGE PARK (NO FDA)",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "C8v3",
						Name:            "C8 - GLENMONT-COLLEGE PK (NO FDA/ARCH)",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "D1",
						Name:            "D1 - GLOVER PARK - FRANKLIN SQUARE",
						LineDescription: "Glover Park-Franklin Square Line",
					},
					{
						RouteID:         "D12",
						Name:            "D12 - SOUTHERN AVE STA - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D12v1",
						Name:            "D12 - SUITLAND STA - SOUTHERN AVE STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D12v2",
						Name:            "D12 - ST BARNABAS RD   - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D13",
						Name:            "D13 - SOUTHERN AVE STA - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D13v1",
						Name:            "D13 - SOUTHRN STA - ALLENTWN+OLD BRNCH",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D14",
						Name:            "D14 - SOUTHERN AVE STA - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D14v1",
						Name:            "D14 - ALLENTWN+OLD BRNCH - SUITLND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D14v2",
						Name:            "D14 - SUITLND STA - ALLENTWN+OLD BRNCH",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D2",
						Name:            "D2 - GLOVER PARK   - CONNETICUT +Q ST",
						LineDescription: "Glover Park-Dupont Circle Line",
					},
					{
						RouteID:         "D2v1",
						Name:            "D2 - DUPONT CIR - GLOVER PARK",
						LineDescription: "Glover Park-Dupont Circle Line",
					},
					{
						RouteID:         "D31",
						Name:            "D31 - TENLEYTOWN STA - 16TH + EASTERN",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D32",
						Name:            "D32 - TENLEYTOWN STA - COLUMBIA HTS STA",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D33",
						Name:            "D33 - TENLEYTOWN STA - 16TH + SHEPHERD",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D34",
						Name:            "D34 - TENLEYTOWN STA - 14TH + COLORADO",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D4",
						Name:            "D4 - DUPONT CIRCLE - IVY CITY",
						LineDescription: "Ivy City-Franklin Square Line",
					},
					{
						RouteID:         "D4v1",
						Name:            "D4 - FRANKLIN SQUARE - IVY CITY",
						LineDescription: "Ivy City-Franklin Square Line",
					},
					{
						RouteID:         "D4v2",
						Name:            "D4 - IVY CITY - FRANKLIN SQUARE",
						LineDescription: "Ivy City-Franklin Square Line",
					},
					{
						RouteID:         "D5",
						Name:            "D5 - MASS LTTL FLWR CHRCH- FARRGT SQR",
						LineDescription: "MacArthur Blvd-Georgetown Line",
					},
					{
						RouteID:         "D6",
						Name:            "D6 - SIBLEY HOSP - STADIUM ARMRY",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D6v1",
						Name:            "D6 - STADIUM ARMRY STA - FARRAGUT SQUARE",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D6v2",
						Name:            "D6 - SIBLEY HOSP - FARRAGUT SQ",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D6v3",
						Name:            "D6 - FARRAGUT SQUARE - STADIUM ARMRY",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D8",
						Name:            "D8 - UNION STATION - VA MED CTR",
						LineDescription: "Hospital Center Line",
					},
					{
						RouteID:         "D8v1",
						Name:            "D8 - RHODE ISLAND STA - UNION STA",
						LineDescription: "Hospital Center Line",
					},
					{
						RouteID:         "E2",
						Name:            "E2 - IVY CITY - FT TOTTEN",
						LineDescription: "Ivy City-Fort Totten Line",
					},
					{
						RouteID:         "E4",
						Name:            "E4 - FRIENDSHP HTS - RIGGS PK",
						LineDescription: "Military Road-Crosstown Line",
					},
					{
						RouteID:         "E4v1",
						Name:            "E4 - FRIENDSHP HTS - FT TOTTEN",
						LineDescription: "Military Road-Crosstown Line",
					},
					{
						RouteID:         "E4v2",
						Name:            "E4 - FT TOTTEN - FRIENDSHP HTS",
						LineDescription: "Military Road-Crosstown Line",
					},
					{
						RouteID:         "E6",
						Name:            "E6 - FRIENDSHP HTS  -KNOLLWOOD",
						LineDescription: "Chevy Chase Line",
					},
					{
						RouteID:         "F1",
						Name:            "F1 - CHEVERLY STA - TAKOMA STA",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F12",
						Name:            "F12 - CHEVERLY STA - NEW CARROLTON STA",
						LineDescription: "Ardwick Industrial Park Shuttle Line",
					},
					{
						RouteID:         "F12v1",
						Name:            "F12 - CHEVERLY STA - LANDOVER STA",
						LineDescription: "Ardwick Industrial Park Shuttle Line",
					},
					{
						RouteID:         "F13",
						Name:            "F13 - CHEVERLY STA-WASHINGTON BUS PARK",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F13v1",
						Name:            "F13 - CHEVERLY STA - WASHINGTON BUS PARK",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F13v2",
						Name:            "F13 - NEW CARROLTON - WASHINGTON BUS PARK",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F13v3",
						Name:            "F13 - WHITFIELD + VOLTA - CHEVERLY STA",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F14",
						Name:            "F14 - NAYLOR RD STA -NEW CARROLTON STA",
						LineDescription: "Sheriff Road-Capitol Heights Line",
					},
					{
						RouteID:         "F14v1",
						Name:            "F14 - BRADBURY HGTS -NEW CARROLTON STA",
						LineDescription: "Sheriff Road-Capitol Heights Line",
					},
					{
						RouteID:         "F2",
						Name:            "F2 - CHEVERLY STA - TAKOMA STA",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F2v1",
						Name:            "F2 - CHEVRLY STA-QUEENS CHAPL+CARSN CIR",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F2v2",
						Name:            "F2 - 34TH + OTIS - TAKOMA STA",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F4",
						Name:            "F4 - SILVR SPRING STA - NEW CARRLLTON",
						LineDescription: "New Carrollton-Silver Spring Line",
					},
					{
						RouteID:         "F4v1",
						Name:            "F4 - PG PLAZA STA -NEW CARROLTON STA",
						LineDescription: "New Carrollton-Silver Spring Line",
					},
					{
						RouteID:         "F4v2",
						Name:            "F4 - NEW CARRLLTON - PG PLAZA STA",
						LineDescription: "New Carrollton-Silver Spring Line",
					},
					{
						RouteID:         "F6",
						Name:            "F6 - FT TOTTEN - NEW CARROLLTN",
						LineDescription: "New Carrollton-Fort Totten Line",
					},
					{
						RouteID:         "F6v1",
						Name:            "F6 - PG PLAZA - NEW CARROLLTON STA",
						LineDescription: "New Carrollton-Fort Totten Line",
					},
					{
						RouteID:         "F6v2",
						Name:            "F6 - NEW CARROLLTON - PG PLAZA STA",
						LineDescription: "New Carrollton-Fort Totten Line",
					},
					{
						RouteID:         "F8",
						Name:            "F8 - CHEVERLY STA - TLTC",
						LineDescription: "Langley Park-Cheverly Line",
					},
					{
						RouteID:         "G12",
						Name:            "G12 - GREENBELT STA - NEW CARROLLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G12v1",
						Name:            "G12 - GREENBELT STA - ROOSEVELT CENTER",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G12v2",
						Name:            "G12 - ROOSEVELT CTR - NEW CARROLLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G14",
						Name:            "G14 - GREENBELT STA -NEW CARROLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G14v1",
						Name:            "G14 - ROOSEVELT CTR -NEW CARROLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G14v2",
						Name:            "G14 - GREENBELT STA - NEW CARROLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G2",
						Name:            "G2 - GEORGETOWN UNIV - HOWARD UNIV",
						LineDescription: "P Street-LeDroit Park Line",
					},
					{
						RouteID:         "G2v1",
						Name:            "G2 - GEORGETOWN UNIV - P+14TH ST",
						LineDescription: "P Street-LeDroit Park Line",
					},
					{
						RouteID:         "G8",
						Name:            "G8 - FARRAGUT SQUARE - AVONDALE",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G8v1",
						Name:            "G8 - RHODE ISLAND AVE STA-FARRAGUT SQ",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G8v2",
						Name:            "G8 - FARRAGUT SQ - RHODE ISLAND STA",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G8v3",
						Name:            "G8 - BROOKLAND STA - AVONDALE",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G9",
						Name:            "G9 - FRANKLIN SQ - RHODE ISLD + EAST PM",
						LineDescription: "Rhode Island Avenue Limited Line",
					},
					{
						RouteID:         "G9v1",
						Name:            "G9 - FRANKLIN SQ - RHODE ISLD + EASTERN",
						LineDescription: "Rhode Island Avenue Limited Line",
					},
					{
						RouteID:         "H1",
						Name:            "H1 - C + 17TH ST - BROOKLAND CUA STA",
						LineDescription: "Brookland-Potomac Park Line",
					},
					{
						RouteID:         "H11",
						Name:            "H11 - HEATHER HILL - NAYLOR RD STA",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H12",
						Name:            "H12 - HEATHER HILL-NAYLOR RD (MACYS)",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H12v1",
						Name:            "H12 - HEATHER HILL-NAYLOR RD STA (P&R)",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H13",
						Name:            "H13 - HEATHER HILL - NAYLOR RD STA",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H2",
						Name:            "H2 - TENLEYTOWN AU-ST - BROOKLND CUA STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H3",
						Name:            "H3 - TENLEYTOWN STA -BROOKLND CUA STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H4",
						Name:            "H4 - TENLEYTWN AU-STA -BROOKLAND CUA STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H4v1",
						Name:            "H4 - COLUMBIA RD+14TH -TENLEYTOWN-AU STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H6",
						Name:            "H6 - BROOKLAND CUA STA - FORT LINCOLN",
						LineDescription: "Brookland-Fort Lincoln Line",
					},
					{
						RouteID:         "H6v1",
						Name:            "H6 - FORT LINCOLN - FORT LINCOLN",
						LineDescription: "Brookland-Fort Lincoln Line",
					},
					{
						RouteID:         "H8",
						Name:            "H8 - MT PLEASANT+17TH - RHODE ISLAND",
						LineDescription: "Park Road-Brookland Line",
					},
					{
						RouteID:         "H9",
						Name:            "H9 - RHODE ISLAND - FORT DR + 1ST",
						LineDescription: "Park Road-Brookland Line",
					},
					{
						RouteID:         "J1",
						Name:            "J1 - MEDCAL CTR STA - SILVR SPRNG STA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J1v1",
						Name:            "J1 - SILVR SPRNG STA - MONT MALL/BAT LN",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J12",
						Name:            "J12 - ADDISON RD STA - FORESTVILLE",
						LineDescription: "Marlboro Pike Line",
					},
					{
						RouteID:         "J12v1",
						Name:            "J12 - ADDISON RD-FORESTVILLE VIA PRES PKY",
						LineDescription: "Marlboro Pike Line",
					},
					{
						RouteID:         "J2",
						Name:            "J2 - MONTGOMRY MALL - SILVR SPRNG STA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J2v1",
						Name:            "J2 - MEDICAL CTR STA - SILVR SPRNG STA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J2v2",
						Name:            "J2 - SILVER SPRNG - MONT MALL/BATTRY LA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J4",
						Name:            "J4 - BETHESDA STA - COLLEGE PARK STA",
						LineDescription: "College Park-Bethesda Limited",
					},
					{
						RouteID:         "K12",
						Name:            "K12 - BRANCH AVE ST-SUITLAND ST",
						LineDescription: "Forestville Line",
					},
					{
						RouteID:         "K12v1",
						Name:            "K12 - PENN MAR - SUITLAND ST",
						LineDescription: "Forestville Line",
					},
					{
						RouteID:         "K12v2",
						Name:            "K12 - BRANCH AVE - PENN MAR",
						LineDescription: "Forestville Line",
					},
					{
						RouteID:         "K2",
						Name:            "K2 - FT TOTTEN STA - TAKOMA STA",
						LineDescription: "Takoma-Fort Totten Line",
					},
					{
						RouteID:         "K6",
						Name:            "K6 - FORT TOTTEN STA - WHITE OAK",
						LineDescription: "New Hampshire Ave-Maryland Line",
					},
					{
						RouteID:         "K6v1",
						Name:            "K6 - TLTC - FORT TOTTEN STA",
						LineDescription: "New Hampshire Ave-Maryland Line",
					},
					{
						RouteID:         "K9",
						Name:            "K9 - FORT TOTTEN STA - FDA/FRC (PM)",
						LineDescription: "New Hampshire Ave-Maryland Limited Line",
					},
					{
						RouteID:         "K9v1",
						Name:            "K9 - FORT TOTTEN STA - FDA/FRC (AM)",
						LineDescription: "New Hampshire Ave-Maryland Limited Line",
					},
					{
						RouteID:         "L1",
						Name:            "L1 - POTOMAC PK  - CHEVY CHASE",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L2",
						Name:            "L2 - FARRAGUT SQ - CHEVY CHASE",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L2v1",
						Name:            "L2 - VAN NESS-UDC STA - CHEVY CHASE",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L2v2",
						Name:            "L2 - FARRAGUT SQ - BETHESDA",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L8",
						Name:            "L8 - FRIENDSHIP HTS STA - ASPEN HILL",
						LineDescription: "Connecticut Ave-Maryland Line",
					},
					{
						RouteID:         "M4",
						Name:            "M4 - SIBLEY HOSPITAL - TENLEYTOWN/AU STA",
						LineDescription: "Nebraska Ave Line",
					},
					{
						RouteID:         "M4v1",
						Name:            "M4 - TENLEYTOWN   - PINEHRST CIR",
						LineDescription: "Nebraska Ave Line",
					},
					{
						RouteID:         "M4v2",
						Name:            "M4 - PINEHRST CIR - TENLEYTOWN",
						LineDescription: "Nebraska Ave Line",
					},
					{
						RouteID:         "M6",
						Name:            "M6 - POTOMAC AVE - ALABAMA + PENN",
						LineDescription: "Fairfax Village Line",
					},
					{
						RouteID:         "M6v1",
						Name:            "M6 - ALABAMA + PENN - FAIRFAX VILLAGE",
						LineDescription: "Fairfax Village Line",
					},
					{
						RouteID:         "MW1",
						Name:            "MW1 - BRADDOCK RD - PENTAGON CITY",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "MW1v1",
						Name:            "MW1 - POTOMAC YARD - CRYSTAL CITY",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "MW1v2",
						Name:            "MW1 - CRYSTAL CITY - BRADDOCK RD",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "MW1v3",
						Name:            "MW1 - BRADDOCK RD - CRYSTAL CITY",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "N2",
						Name:            "N2 - FRIENDSHIP HTS - FARRAGUT SQ",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "N4",
						Name:            "N4 - FRIENDSHP HTS - POTOMAC PARK",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "N4v1",
						Name:            "N4 - FRIENDSHP HTS - FARRAGUT  SQ",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "N6",
						Name:            "N6 - FRNDSHIP HTS - FARRAGUT SQ",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "NH1",
						Name:            "NH1 - NATIONAL HARBOR-SOUTHERN AVE STA",
						LineDescription: "National Harbor-Southern Avenue Line",
					},
					{
						RouteID:         "NH2",
						Name:            "NH2 - HUNTINGTON STA-NATIONAL HARBOR",
						LineDescription: "National Harbor-Alexandria Line",
					},
					{
						RouteID:         "P12",
						Name:            "P12 - EASTOVER - ADDISON RD STA (NO HOSP)",
						LineDescription: "Eastover-Addison Road Line",
					},
					{
						RouteID:         "P12v1",
						Name:            "P12 - IVERSON MALL - ADDISON RD STA",
						LineDescription: "Eastover-Addison Road Line",
					},
					{
						RouteID:         "P12v2",
						Name:            "P12 - SUITLAND STA - ADDISON RD STA",
						LineDescription: "Eastover-Addison Road Line",
					},
					{
						RouteID:         "P18",
						Name:            "P18 - FT WASH P+R LOT - SOUTHERN AVE",
						LineDescription: "Oxon Hill-Fort Washington Line",
					},
					{
						RouteID:         "P19",
						Name:            "P19 - FT WASH P+R LOT-SOUTHERN AVE STA",
						LineDescription: "Oxon Hill-Fort Washington Line",
					},
					{
						RouteID:         "P6",
						Name:            "P6 - ANACOSTIA STA - RHODE ISLAND STA",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v1",
						Name:            "P6 - ECKINGTON - RHODE ISLAND AVE",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v2",
						Name:            "P6 - RHODE ISLAND AVE - ECKINGTON",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v3",
						Name:            "P6 - ANACOSTIA STA - ARCHIVES",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v4",
						Name:            "P6 - ARCHIVES - ANACOSTIA",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "Q1",
						Name:            "Q1 - SILVR SPRNG STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q2",
						Name:            "Q2 - SILVR SPRNG STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q2v1",
						Name:            "Q2 - MONT COLLEGE - SILVR SPRNG STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q2v2",
						Name:            "Q2 - SILVR SPRNG STA - MONT COLLEGE",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q4",
						Name:            "Q4 - SILVER SPRNG STA - ROCKVILLE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q4v1",
						Name:            "Q4 - WHEATON STA - ROCKVILLE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q5",
						Name:            "Q5 - WHEATON STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q6",
						Name:            "Q6 - WHEATON STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q6v1",
						Name:            "Q6 - ROCKVILLE STA - WHEATON STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "R1",
						Name:            "R1 - FORT TOTTEN STA - ADELPHI",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R12",
						Name:            "R12 - DEANWOOD STA - GREENBELT STA",
						LineDescription: "Kenilworth Avenue Line",
					},
					{
						RouteID:         "R12v1",
						Name:            "R12 - DEANWOOD STA - GREENBELT STA",
						LineDescription: "Kenilworth Avenue Line",
					},
					{
						RouteID:         "R2",
						Name:            "R2 - FORT TOTTEN - CALVERTON",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R2v1",
						Name:            "R2 - HIGH POINT HS - FORT TOTTEN",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R2v2",
						Name:            "R2 - POWDER MILL+CHERRY HILL - CALVERTON",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R4",
						Name:            "R4 - BROOKLAND STA- HIGHVIEW",
						LineDescription: "Queens Chapel Road Line",
					},
					{
						RouteID:         "REX",
						Name:            "REX - FT BELVOIR POST - KING ST STA",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv1",
						Name:            "REX - FT BELVOIR COMM HOSP - KING ST STA",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv2",
						Name:            "REX - KING ST STA - FT BELVOIR COMM HOSP",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv3",
						Name:            "REX - KING ST STA - WOODLAWN",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv4",
						Name:            "REX - WOODLAWN - KING ST STA",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "S1",
						Name:            "S1 - NORTHERN DIVISION - POTOMAC PK",
						LineDescription: "16th Street-Potomac Park Line",
					},
					{
						RouteID:         "S1v1",
						Name:            "S1 - VIRGINIA+E - COLORDO+16TH",
						LineDescription: "16th Street-Potomac Park Line",
					},
					{
						RouteID:         "S2",
						Name:            "S2 - FED TRIANGLE  - SILVER SPRNG",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S2v1",
						Name:            "S2 - 16TH & HARVARD - MCPHERSON SQ",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S35",
						Name:            "S35 - BRANCH + RANDLE CIR - FT DUPONT",
						LineDescription: "Fort Dupont Shuttle Line",
					},
					{
						RouteID:         "S4",
						Name:            "S4 - FED TRIANGLE - SILVER SPRNG",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S4v1",
						Name:            "S4 - SILVER SPRING STA - FRANKLIN SQ",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S41",
						Name:            "S41 - CARVER TERRACE - RHODE ISLAND AVE",
						LineDescription: "Rhode Island Ave-Carver Terrace Line",
					},
					{
						RouteID:         "S80",
						Name:            "S80 - FRANCONIA-SPRNGFLD - METRO PARK",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S80v1",
						Name:            "S80 - FRANCONIA-SPRINGFLD - HILTON",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S80v2",
						Name:            "S80 - HILTON - FRANCONIA-SPRNGFLD",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S9",
						Name:            "S9 - FRANKLIN SQ - SILVER SPRING STA",
						LineDescription: "16th Street Limited Line",
					},
					{
						RouteID:         "S9v1",
						Name:            "S9 - FRANKLIN SQ - COLORADO + 16TH",
						LineDescription: "16th Street Limited Line",
					},
					{
						RouteID:         "S91",
						Name:            "S91 - FRANCONIA SPRINGFLD STA SHUTTLE",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S91v1",
						Name:            "S91 - FRANCONIA SPRINGFLD STA SHUTTLE",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "T14",
						Name:            "T14 - RHD ISLND AVE STA-NEW CARRLTN STA",
						LineDescription: "Rhode Island Ave-New Carrollton Line",
					},
					{
						RouteID:         "T14v1",
						Name:            "T14 - MT RAINIER - NEW CARRLTN STA",
						LineDescription: "Rhode Island Ave-New Carrollton Line",
					},
					{
						RouteID:         "T18",
						Name:            "T18 - R I AVE STA - NEW CARROLLTON STA",
						LineDescription: "Annapolis Road Line",
					},
					{
						RouteID:         "T18v1",
						Name:            "T18 - BLADENSBURG HS - NEW CARROLLTON STA",
						LineDescription: "Annapolis Road Line",
					},
					{
						RouteID:         "T2",
						Name:            "T2 - FRIENDSHIP HTS - ROCKVILLE STA",
						LineDescription: "River Road Line",
					},
					{
						RouteID:         "U4",
						Name:            "U4 - MINNESOTA AVE - SHERIFF RD",
						LineDescription: "Sheriff Road-River Terrace Line",
					},
					{
						RouteID:         "U4v1",
						Name:            "U4 - SHERIFF RD - MINNESOTA STA",
						LineDescription: "Sheriff Road-River Terrace Line",
					},
					{
						RouteID:         "U4v2",
						Name:            "U4 - RIVER TERRACE - MINNESOTA AVE",
						LineDescription: "Sheriff Road-River Terrace Line",
					},
					{
						RouteID:         "U5",
						Name:            "U5 - MINNESOTA AVE-MARSHALL HTS",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U6",
						Name:            "U6 - MINNESOTA - LINCOLN HEIGHTS",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U6v1",
						Name:            "U6 - 37TH + RIDGE - PLUMMER ES",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U6v2",
						Name:            "U6 - LINCOLN HTS - E CAP + 47TH ST NE",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U7",
						Name:            "U7 - RIDGE + ANACOSTIA - DEANWOOD",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v1",
						Name:            "U7 - MINNESOTA STA - KENILWORTH HAYES",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v2",
						Name:            "U7 - KENILWORTH HAYES - MINNESOTA STA",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v3",
						Name:            "U7 - MINNESOTA STA - DEANWOOD",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v4",
						Name:            "U7 - DEANWOOD - MINNESOTA AVE",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "V1",
						Name:            "V1 - BUR OF ENGRVNG - BENNING HTS",
						LineDescription: "Benning Heights-M Street Line",
					},
					{
						RouteID:         "V12",
						Name:            "V12 - SUITLAND STA - ADDISON RD STA",
						LineDescription: "District Heights-Suitland Line",
					},
					{
						RouteID:         "V14",
						Name:            "V14 - PENN MAR - DEANWOOD STA",
						LineDescription: "District Heights - Seat Pleasant Line",
					},
					{
						RouteID:         "V14v1",
						Name:            "V14 - PENN MAR  - ADDISON RD STA",
						LineDescription: "District Heights - Seat Pleasant Line",
					},
					{
						RouteID:         "V2",
						Name:            "V2 - ANACOSTIA - CAPITOL HGTS",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V2v1",
						Name:            "V2 - MINNESOTA AVE - ANACOSTIA",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V4",
						Name:            "V4 - 1ST + K ST SE - CAPITOL HGTS",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V4v1",
						Name:            "V4 - MINN AVE STA-CAPITOL HGTS",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V7",
						Name:            "V7 - CONGRESS HGTS - MINN STA",
						LineDescription: "Benning Heights-Alabama Ave Line",
					},
					{
						RouteID:         "V8",
						Name:            "V8 - BENNING HGTS - MINN STA",
						LineDescription: "Benning Heights-Alabama Ave Line",
					},
					{
						RouteID:         "W1",
						Name:            "W1 - FORT DRUM  - SOUTHERN AVE STA",
						LineDescription: "Shipley Terrace - Fort Drum Line",
					},
					{
						RouteID:         "W14",
						Name:            "W14 - FT WASHINGTON-SOUTHERN AVE STA",
						LineDescription: "Bock Road Line",
					},
					{
						RouteID:         "W14v1",
						Name:            "W14 - SOUTHERN AVE - FRIENDLY",
						LineDescription: "Bock Road Line",
					},
					{
						RouteID:         "W14v2",
						Name:            "W14 - FRIENDLY - SOUTHERN AVE",
						LineDescription: "Bock Road Line",
					},
					{
						RouteID:         "W2",
						Name:            "W2 - MALCOM X+OAKWD - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v1",
						Name:            "W2 - NYLDR+GOOD HOP- HOWARD+ANACOSTIA",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v2",
						Name:            "W2 - NAYLOR+GOODHOPE - MALCM X+OAKWOOD",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v3",
						Name:            "W2 - HOWRD+ANACOSTIA-NAYLOR+GOODHOPE",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v4",
						Name:            "W2 - ANACOSTIA STA - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v5",
						Name:            "W2 - ANACOSTIA STA - SOUTHERN AVE STA",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v6",
						Name:            "W2 - MELLN+M L KNG - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v7",
						Name:            "W2 - SOUTHERN AVE-WASHINGTON OVERLOOK",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W3",
						Name:            "W3 - MALCOM X+OAKWD - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W3v1",
						Name:            "W3 - MELLN+M L KNG - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W4",
						Name:            "W4 - ANACOSTIA STA - DEANWOOD STA",
						LineDescription: "Deanwood-Alabama Ave Line",
					},
					{
						RouteID:         "W4v1",
						Name:            "W4 - MALCOLM X & PORTLAND - DEANWOOD",
						LineDescription: "Deanwood-Alabama Ave Line",
					},
					{
						RouteID:         "W4v2",
						Name:            "W4 - DEANWOOD STA - MALCLM X+ ML KING",
						LineDescription: "Deanwood-Alabama Ave Line",
					},
					{
						RouteID:         "W45",
						Name:            "W45 - TENLEYTOWN STA - 16TH + SHEPHERD",
						LineDescription: "Mt Pleasant-Tenleytown Line",
					},
					{
						RouteID:         "W47",
						Name:            "W47 - TENLEYTOWN STA - COLUMBIA HTS STA",
						LineDescription: "Mt Pleasant-Tenleytown Line",
					},
					{
						RouteID:         "W5",
						Name:            "W5 - DC VILLAGE-USCG (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "W6",
						Name:            "W6 - ANACOSTIA - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W6v1",
						Name:            "W6 - NAYLOR+GOODHOPE - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W8",
						Name:            "W8 - ANACOSTIA - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W8v1",
						Name:            "W8 - ANACOSTIA-NAYLOR+GOODHOPE",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W8v2",
						Name:            "W8 - NAYLOR+GOODHOPE - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "X1",
						Name:            "X1 - FOGGY BOTTOM+GWU- MINNESOTA STA",
						LineDescription: "Benning Road Line",
					},
					{
						RouteID:         "X2",
						Name:            "X2 - LAFAYETTE SQ - MINNESOTA STA",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X2v1",
						Name:            "X2 - PHELPS HS - MINNESOTA STA",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X2v2",
						Name:            "X2 - FRIENDSHIP EDISON PCS-MINNESOTA STA",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X2v3",
						Name:            "X2 - PHELPS HS - LAFAYTTE SQ",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X3",
						Name:            "X3 - DUKE ELLINGTON BR - MINNESOTA STN",
						LineDescription: "Benning Road Line",
					},
					{
						RouteID:         "X3v1",
						Name:            "X3 - KIPP DC - MINNESOTA AVE STN",
						LineDescription: "Benning Road Line",
					},
					{
						RouteID:         "X8",
						Name:            "X8 - UNION STATION - CARVER TERR",
						LineDescription: "Maryland Ave Line",
					},
					{
						RouteID:         "X9",
						Name:            "X9 - NY AVE & 12TH NW - CAPITOL HTS STA",
						LineDescription: "Benning Road-H Street Limited Line",
					},
					{
						RouteID:         "X9v1",
						Name:            "X9 - NY AVE & 12TH ST NW - MINNESOTA AVE",
						LineDescription: "Benning Road-H Street Limited Line",
					},
					{
						RouteID:         "X9v2",
						Name:            "X9 - MINNESOTA AVE ST - NY AVE & 12TH ST",
						LineDescription: "Benning Road-H Street Limited Line",
					},
					{
						RouteID:         "Y2",
						Name:            "Y2 - SILVER SPRING STA - MONTG MED CTR",
						LineDescription: "Georgia Ave-Maryland Line",
					},
					{
						RouteID:         "Y7",
						Name:            "Y7 - SILVER SPRING STA - ICC P&R",
						LineDescription: "Georgia Ave-Maryland Line",
					},
					{
						RouteID:         "Y8",
						Name:            "Y8 - SILVER SPR STA - MONTGOMERY MED CTR",
						LineDescription: "Georgia Ave-Maryland Line",
					},
					{
						RouteID:         "Z11",
						Name:            "Z11 - SILVR SPRING - BURTONSVILLE P&R",
						LineDescription: "Greencastle-Briggs Chaney Express Line",
					},
					{
						RouteID:         "Z11v1",
						Name:            "Z11 - GREENCASTLE - SILVR SPRING",
						LineDescription: "Greencastle-Briggs Chaney Express Line",
					},
					{
						RouteID:         "Z2",
						Name:            "Z2 - SILVR SPRING - OLNEY (NO BLAKE HS)",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z2v1",
						Name:            "Z2 - SILVER SPRING - BONIFANT & NH",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z2v2",
						Name:            "Z2 - SILVER SPRING STA-OLNEY (BLAKE HS)",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z2v3",
						Name:            "Z2 - COLESVILLE-SILVER SPRING",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z6",
						Name:            "Z6 - SILVR SPRNG STA - BURTONSVILLE",
						LineDescription: "Calverton-Westfarm Line",
					},
					{
						RouteID:         "Z6v1",
						Name:            "Z6 - CASTLE BLVD - SILVER SPRING STA",
						LineDescription: "Calverton-Westfarm Line",
					},
					{
						RouteID:         "Z6v2",
						Name:            "Z6 - SILVER SPRING -CASTLE BLVD",
						LineDescription: "Calverton-Westfarm Line",
					},
					{
						RouteID:         "Z7",
						Name:            "Z7 - SLVER SPRNG STA-S LAUREL P&R (4COR)",
						LineDescription: "Laurel-Burtonsville Express Line",
					},
					{
						RouteID:         "Z7v1",
						Name:            "Z7 - S LAUREL P&R-SILVR SPR (NO 4COR)",
						LineDescription: "Laurel-Burtonsville Express Line",
					},
					{
						RouteID:         "Z8",
						Name:            "Z8 - SILVER SRING STA - BRIGGS CHANEY",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v1",
						Name:            "Z8 - WHITE OAK - SILVER SPRING STA",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v2",
						Name:            "Z8 - SILVR SPRNG - CSTLE BLVD VERIZON",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v3",
						Name:            "Z8 - SILVER SPRING STA - CASTLE BLVD",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v4",
						Name:            "Z8 - SILVER SPRING - GRNCSTLE (VERIZON)",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v5",
						Name:            "Z8 - SILVER SPRING STA - GREENCASTLE",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v6",
						Name:            "Z8 - SILVER SPRING STA - WHITE OAK",
						LineDescription: "Fairland Line",
					},
				},
			},
		},
	},
	"/Bus.svc/Routes": {
		{
			rawQuery: "",
			response: `<RoutesResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Routes><Route><LineDescription>Alexandria-Pentagon Line</LineDescription><Name>10A - HUNTINGTON STA - PENTAGON</Name><RouteID>10A</RouteID></Route><Route><LineDescription>Hunting Point-Ballston Line</LineDescription><Name>10B - HUNTING POINT - BALLSTON STA</Name><RouteID>10B</RouteID></Route><Route><LineDescription>Alexandria-Pentagon Line</LineDescription><Name>10E - HUNTING POINT - PENTAGON</Name><RouteID>10E</RouteID></Route><Route><LineDescription>Alexandria-Pentagon Line</LineDescription><Name>10N - NATIONAL AIRPORT - PENTAGON</Name><RouteID>10N</RouteID></Route><Route><LineDescription>Mt Vernon Express Line</LineDescription><Name>11Y - MOUNT VERNON - POTOMAC PARK</Name><RouteID>11Y</RouteID></Route><Route><LineDescription>Mt Vernon Express Line</LineDescription><Name>11Y - POTOMAC PARK - HUNTING POINT</Name><RouteID>11Yv1</RouteID></Route><Route><LineDescription>Mt Vernon Express Line</LineDescription><Name>11Y - HUNTING POINT - POTOMAC PARK</Name><RouteID>11Yv2</RouteID></Route><Route><LineDescription>Chain Bridge Road Line</LineDescription><Name>15K - E FALLS CH STA-MCLEAN LANGLEY</Name><RouteID>15K</RouteID></Route><Route><LineDescription>Chain Bridge Road Line</LineDescription><Name>15K - CIA - EAST FALLS CHURCH STA</Name><RouteID>15Kv1</RouteID></Route><Route><LineDescription>Columbia Pike Line</LineDescription><Name>16A - PATRIOT+AMERICANA - PENTAGON</Name><RouteID>16A</RouteID></Route><Route><LineDescription>Columbia Pike Line</LineDescription><Name>16C - CULMORE - PENTAGON</Name><RouteID>16C</RouteID></Route><Route><LineDescription>Columbia Pike Line</LineDescription><Name>16C - CULMORE - FEDERAL TRIANGLE</Name><RouteID>16Cv1</RouteID></Route><Route><LineDescription>Columbia Pike Line</LineDescription><Name>16E - GLN CARLYN+VISTA - FRANKLIN SQ</Name><RouteID>16E</RouteID></Route><Route><LineDescription>Columbia Pike-Pentagon City Line</LineDescription><Name>16G - DINWIDDIE+COLUMBIA - PENTAGON CITY</Name><RouteID>16G</RouteID></Route><Route><LineDescription>Columbia Pike-Pentagon City Line</LineDescription><Name>16G - COL PIKE+CARLIN SPR - DINWD+COL PK</Name><RouteID>16Gv1</RouteID></Route><Route><LineDescription>Columbia Pike-Pentagon City Line</LineDescription><Name>16H - SKYLINE CITY - PENTAGON CITY STA</Name><RouteID>16H</RouteID></Route><Route><LineDescription>Annandale-Skyline City-Pentagon Line</LineDescription><Name>16L - ANNANDALE - PENTAGON HOV</Name><RouteID>16L</RouteID></Route><Route><LineDescription>Columbia Pike-Farragut Square Line</LineDescription><Name>16Y - FOUR MILE+COLUMBIA - MCPHERSON SQ</Name><RouteID>16Y</RouteID></Route><Route><LineDescription>Columbia Pike-Farragut Square Line</LineDescription><Name>16Y - COLUMBIA PIKE - MCPHERSON SQ</Name><RouteID>16Yv1</RouteID></Route><Route><LineDescription>Kings Park-North Springfield Line</LineDescription><Name>17B - BURKE CENTRE - PENTAGON HOV</Name><RouteID>17B</RouteID></Route><Route><LineDescription>Kings Park Express Line</LineDescription><Name>17G - G MASON UNIV - PENTAGON</Name><RouteID>17G</RouteID></Route><Route><LineDescription>Kings Park Express Line</LineDescription><Name>17H - TWNBRK RD+TWNBRK RN- PENTAGON</Name><RouteID>17H</RouteID></Route><Route><LineDescription>Kings Park Express Line</LineDescription><Name>17K - TWNBRK RD+TWNBRK RN- PENTAGON</Name><RouteID>17K</RouteID></Route><Route><LineDescription>Kings Park Express Line</LineDescription><Name>17L - TWNBRK RD+TWNBRK RUN-PENTAGON</Name><RouteID>17L</RouteID></Route><Route><LineDescription>Kings Park-North Springfield Line</LineDescription><Name>17M - EDSALL+CANARD - PENTAGON</Name><RouteID>17M</RouteID></Route><Route><LineDescription>Orange Hunt Line</LineDescription><Name>18G - ROLLING VALLEY - PENTAGON</Name><RouteID>18G</RouteID></Route><Route><LineDescription>Orange Hunt Line</LineDescription><Name>18H - HUNTSMAN+CORK CTY - PENTAGON</Name><RouteID>18H</RouteID></Route><Route><LineDescription>Orange Hunt Line</LineDescription><Name>18J - ROLLING VALLEY - PENTAGON</Name><RouteID>18J</RouteID></Route><Route><LineDescription>Burke Centre Line</LineDescription><Name>18P - BURKE CENTRE - PENTAGON</Name><RouteID>18P</RouteID></Route><Route><LineDescription>Burke Centre Line</LineDescription><Name>18P - PENTAGON - ROLLING VALLEY MALL</Name><RouteID>18Pv1</RouteID></Route><Route><LineDescription>Burke Centre Line</LineDescription><Name>18P - ROLLING VALLEY P+R - PENTAGON</Name><RouteID>18Pv2</RouteID></Route><Route><LineDescription>Wilson Blvd-Vienna Line</LineDescription><Name>1A - VIENNA-BALLSTON (7 CORNERS)</Name><RouteID>1A</RouteID></Route><Route><LineDescription>Wilson Blvd-Vienna Line</LineDescription><Name>1B - DUNN LORING  - BALLSTON</Name><RouteID>1B</RouteID></Route><Route><LineDescription>Fair Oaks-Fairfax Blvd Line</LineDescription><Name>1C - WEST OX DIV-DUNN LORING (VIA MALL)</Name><RouteID>1C</RouteID></Route><Route><LineDescription>Fair Oaks-Fairfax Blvd Line</LineDescription><Name>1C - FAIRFAX CO GOV CTR - DUNN LORING</Name><RouteID>1Cv1</RouteID></Route><Route><LineDescription>Fair Oaks-Fairfax Blvd Line</LineDescription><Name>1C - WEST OX DIV - DUNN LORING (NO MALL)</Name><RouteID>1Cv2</RouteID></Route><Route><LineDescription>Fair Oaks-Fairfax Blvd Line</LineDescription><Name>1C - FAIR OAKS MALL - DUNN LORING</Name><RouteID>1Cv3</RouteID></Route><Route><LineDescription>Fair Oaks-Fairfax Blvd Line</LineDescription><Name>1C - DUNN LORING - FAIR OAKS MALL</Name><RouteID>1Cv4</RouteID></Route><Route><LineDescription>Landmark-Bren Mar Park-Pentagon Line</LineDescription><Name>21A - S REYNOLDS+EOS 21 CONDOS - PENTAGON</Name><RouteID>21A</RouteID></Route><Route><LineDescription>Landmark-Bren Mar Park-Pentagon Line</LineDescription><Name>21D - LANDMARK MEWS -PENTAGON</Name><RouteID>21D</RouteID></Route><Route><LineDescription>Barcroft-South Fairlington Line</LineDescription><Name>22A - BALLSTON STA - PENTAGON</Name><RouteID>22A</RouteID></Route><Route><LineDescription>Barcroft-South Fairlington Line</LineDescription><Name>22A - SHIRLINGTON - BALLSTON  STA</Name><RouteID>22Av1</RouteID></Route><Route><LineDescription>Barcroft-South Fairlington Line</LineDescription><Name>22C - BALLSTON STA - PENTAGON</Name><RouteID>22C</RouteID></Route><Route><LineDescription>Barcroft-South Fairlington Line</LineDescription><Name>22F - NVCC - PENTAGON VIA HOV</Name><RouteID>22F</RouteID></Route><Route><LineDescription>McLean-Crystal City Line</LineDescription><Name>23A - TYSONS CORNER CTR - CRYSTAL CTY</Name><RouteID>23A</RouteID></Route><Route><LineDescription>McLean-Crystal City Line</LineDescription><Name>23B - BALLSTON STA - CRYSTAL CTY</Name><RouteID>23B</RouteID></Route><Route><LineDescription>McLean-Crystal City Line</LineDescription><Name>23B - LINDEN RESOURCES - BALLSTON STATION</Name><RouteID>23Bv1</RouteID></Route><Route><LineDescription>McLean-Crystal City Line</LineDescription><Name>23T - TYSONS CORNER CTR - SHIRLINGTON</Name><RouteID>23T</RouteID></Route><Route><LineDescription>Landmark-Ballston Line</LineDescription><Name>25B - VAN DORN - BALLSTON</Name><RouteID>25B</RouteID></Route><Route><LineDescription>Landmark-Ballston Line</LineDescription><Name>25B - SOUTHERN TOWERS - BALLSTON</Name><RouteID>25Bv1</RouteID></Route><Route><LineDescription>Landmark-Ballston Line</LineDescription><Name>25B - VAN DORN - BALLSTON/NO LDMRK CTR</Name><RouteID>25Bv2</RouteID></Route><Route><LineDescription>Landmark-Ballston Line</LineDescription><Name>25B - BALLSTON - SOUTHERN TOWERS</Name><RouteID>25Bv3</RouteID></Route><Route><LineDescription>Annandale-East Falls Church Line</LineDescription><Name>26A - NVCC ANNANDALE - E FALLS CHURCH STA</Name><RouteID>26A</RouteID></Route><Route><LineDescription>Leesburg Pike Line</LineDescription><Name>28A - TYSONS CORNER STA-KING ST STA</Name><RouteID>28A</RouteID></Route><Route><LineDescription>Leesburg Pike Line</LineDescription><Name>28A - SOUTHERN TOWERS-TYSONS CORNER STA</Name><RouteID>28Av1</RouteID></Route><Route><LineDescription>Skyline City Line</LineDescription><Name>28F - BLDG 5113 G MASON DR - PENTAGON</Name><RouteID>28F</RouteID></Route><Route><LineDescription>Skyline City Line</LineDescription><Name>28G - BLDG 5113 G MASON DR - PENTAGON</Name><RouteID>28G</RouteID></Route><Route><LineDescription>Annandale Line</LineDescription><Name>29C - NVCC ANNANDALE - PENTAGON</Name><RouteID>29C</RouteID></Route><Route><LineDescription>Annandale Line</LineDescription><Name>29G - AMERICANA+HERITAGE - PENTAGON</Name><RouteID>29G</RouteID></Route><Route><LineDescription>Alexandria-Fairfax Line</LineDescription><Name>29K - GMU - KING ST STA</Name><RouteID>29K</RouteID></Route><Route><LineDescription>Alexandria-Fairfax Line</LineDescription><Name>29K - GMU - KING ST/NO LDMRK</Name><RouteID>29Kv1</RouteID></Route><Route><LineDescription>Alexandria-Fairfax Line</LineDescription><Name>29N - VIENNA STA - KING ST (VIA MALL)</Name><RouteID>29N</RouteID></Route><Route><LineDescription>Alexandria-Fairfax Line</LineDescription><Name>29N - VIENNA STA-KING ST (NO MALL)</Name><RouteID>29Nv1</RouteID></Route><Route><LineDescription>Braeburn Drive-Pentagon Express Line</LineDescription><Name>29W - NVCC ANNANDALE - PENTAGON</Name><RouteID>29W</RouteID></Route><Route><LineDescription>Washington Blvd.-Dunn Loring Line</LineDescription><Name>2A - DUNN LORING STA - BALLSTON STA</Name><RouteID>2A</RouteID></Route><Route><LineDescription>Fair Oaks-Jermantown Road Line</LineDescription><Name>2B - WEST OX RD DIV-DUNN LORING STATION</Name><RouteID>2B</RouteID></Route><Route><LineDescription>Fair Oaks-Jermantown Road Line</LineDescription><Name>2B - WEST OX RD-DUNN LORING STA(NO MALL)</Name><RouteID>2Bv1</RouteID></Route><Route><LineDescription>Fair Oaks-Jermantown Road Line</LineDescription><Name>2B - FAIR OAKS MALL-DUNN LORING STATION</Name><RouteID>2Bv2</RouteID></Route><Route><LineDescription>Fair Oaks-Jermantown Road Line</LineDescription><Name>2B - DUNN LORING STA - FAIR OAKS MALL</Name><RouteID>2Bv3</RouteID></Route><Route><LineDescription>Friendship Heights-Southeast Line</LineDescription><Name>30N - FRIENDSHIP HGTS- NAYLOR RD STA</Name><RouteID>30N</RouteID></Route><Route><LineDescription>Friendship Heights-Southeast Line</LineDescription><Name>30S - FRIENDSHIP HGTS- SOUTHRN AVE STA</Name><RouteID>30S</RouteID></Route><Route><LineDescription>Wisconsin Avenue Line</LineDescription><Name>31 - POTOMAC PARK-FRIENDSHIP HGTS</Name><RouteID>31</RouteID></Route><Route><LineDescription>Pennsylvania Avenue Line</LineDescription><Name>32 - VIRGINIA AVE+E ST- SOUTHRN AVE</Name><RouteID>32</RouteID></Route><Route><LineDescription>Pennsylvania Avenue Line</LineDescription><Name>32 - PENN AVE + 8TH ST - SOUTHRN AVE</Name><RouteID>32v1</RouteID></Route><Route><LineDescription>Wisconsin Avenue Line</LineDescription><Name>33 - 10TH ST+PA AV NW - FRIENDSHIP HGTS</Name><RouteID>33</RouteID></Route><Route><LineDescription>Pennsylvania Avenue Line</LineDescription><Name>34 - 10TH ST + PA AVE- NAYLOR RD STA</Name><RouteID>34</RouteID></Route><Route><LineDescription>Pennsylvania Avenue Line</LineDescription><Name>36 - VIRGINIA AVE+E ST - NAYLOR RD STA</Name><RouteID>36</RouteID></Route><Route><LineDescription>Wisconsin Avenue Limited Line</LineDescription><Name>37 - 10TH ST+PA AV NW - FRIENDSHIP HGTS</Name><RouteID>37</RouteID></Route><Route><LineDescription>Ballston-Farragut Square Line</LineDescription><Name>38B - BALLSTON - FARRAGUT</Name><RouteID>38B</RouteID></Route><Route><LineDescription>Ballston-Farragut Square Line</LineDescription><Name>38B - WASH &amp; QUINCY - FARRAGUT</Name><RouteID>38Bv1</RouteID></Route><Route><LineDescription>Ballston-Farragut Square Line</LineDescription><Name>38B - WASHINGTON-LEE HS - FARRAGUT SQ</Name><RouteID>38Bv2</RouteID></Route><Route><LineDescription>Pennsylvania Avenue Limited Line</LineDescription><Name>39 - VIRGINIA AVE+21ST NW- NAYLOR RD STA</Name><RouteID>39</RouteID></Route><Route><LineDescription>Annandale Road Line</LineDescription><Name>3A - ANNANDALE - E FALLS CHURCH</Name><RouteID>3A</RouteID></Route><Route><LineDescription>Annandale Road Line</LineDescription><Name>3A - ANNANDALE - 7 CORNERS</Name><RouteID>3Av1</RouteID></Route><Route><LineDescription>Pimmit Hills Line</LineDescription><Name>3T - MCLEAN STATION - E FALLS CH STA</Name><RouteID>3T</RouteID></Route><Route><LineDescription>Pimmit Hills Line</LineDescription><Name>3T - MCLEAN STATION - W FALLS CHURCH</Name><RouteID>3Tv1</RouteID></Route><Route><LineDescription>Lee Highway-Farragut Square Line</LineDescription><Name>3Y - E FALLS CHURCH - MCPHERSON SQ</Name><RouteID>3Y</RouteID></Route><Route><LineDescription>Mount Pleasant Line</LineDescription><Name>42 - 9TH + F ST  - MT PLEASANT</Name><RouteID>42</RouteID></Route><Route><LineDescription>Mount Pleasant Line</LineDescription><Name>43 - I + 13TH NW - MT PLEASANT</Name><RouteID>43</RouteID></Route><Route><LineDescription>Pershing Dr-Arlington Blvd Line</LineDescription><Name>4A - SEVEN CORNERS - ROSSLYN</Name><RouteID>4A</RouteID></Route><Route><LineDescription>Pershing Dr-Arlington Blvd Line</LineDescription><Name>4B - SEVEN CORNERS - ROSSLYN</Name><RouteID>4B</RouteID></Route><Route><LineDescription>14th Street Line</LineDescription><Name>52 - L ENFNT PLAZA - TAKOMA STATION</Name><RouteID>52</RouteID></Route><Route><LineDescription>14th Street Line</LineDescription><Name>52 - L ENFNT PLAZA - 14TH+COLORADO</Name><RouteID>52v1</RouteID></Route><Route><LineDescription>14th Street Line</LineDescription><Name>52 - 14TH+COLORADO - L ENFANT PLAZA</Name><RouteID>52v2</RouteID></Route><Route><LineDescription>14th Street Line</LineDescription><Name>52 - 14TH &amp; U - TAKOMA STATION</Name><RouteID>52v3</RouteID></Route><Route><LineDescription>14th Street Line</LineDescription><Name>54 - METRO CENTER - TAKOMA STA</Name><RouteID>54</RouteID></Route><Route><LineDescription>14th Street Line</LineDescription><Name>54 - 14TH+COLORADO - METRO CENTER</Name><RouteID>54v1</RouteID></Route><Route><LineDescription>14th Street Limited Line</LineDescription><Name>59 - FEDERAL TRIANGLE - TAKOMA STATION</Name><RouteID>59</RouteID></Route><Route><LineDescription>DC-Dulles Line</LineDescription><Name>5A - DULLES AIRPORT - LENFANT PLAZA</Name><RouteID>5A</RouteID></Route><Route><LineDescription>Fort Totten-Petworth Line</LineDescription><Name>60 - GEORGIA + PETWORTH - FT TOTTEN</Name><RouteID>60</RouteID></Route><Route><LineDescription>Takoma-Petworth Line</LineDescription><Name>62 - GEORGIA+PETWORTH - TAKOMA STATION</Name><RouteID>62</RouteID></Route><Route><LineDescription>Takoma-Petworth Line</LineDescription><Name>62 - COOLIDGE HS - GEORGIA + PETWORTH</Name><RouteID>62v1</RouteID></Route><Route><LineDescription>Takoma-Petworth Line</LineDescription><Name>63 - FED TRIANGLE - TAKOMA STA</Name><RouteID>63</RouteID></Route><Route><LineDescription>Fort Totten-Petworth Line</LineDescription><Name>64 - FEDERAL TRIANGLE -FORT TOTTEN</Name><RouteID>64</RouteID></Route><Route><LineDescription>Fort Totten-Petworth Line</LineDescription><Name>64 - GEORGIA + PETWOTH - FT TOTTEN</Name><RouteID>64v1</RouteID></Route><Route><LineDescription>Georgia Avenue-7th Street Line</LineDescription><Name>70 - ARCHIVES - SILVER SPRING</Name><RouteID>70</RouteID></Route><Route><LineDescription>Georgia Avenue-7th Street Line</LineDescription><Name>70 - GEORGIA &amp; EUCLID  - ARCHIVES</Name><RouteID>70v1</RouteID></Route><Route><LineDescription>Convention Center-Southwest Waterfront Line</LineDescription><Name>74 - NATIONALS PARK - CONVENTION CTR</Name><RouteID>74</RouteID></Route><Route><LineDescription>Georgia Avenue MetroExtra Line</LineDescription><Name>79 - ARCHIVES - SILVER SPRING STA</Name><RouteID>79</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7A - LINCOLNIA+QUANTRELL - PENTAGON</Name><RouteID>7A</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7A - PENTAGON - SOUTHERN TWRS</Name><RouteID>7Av1</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7A - SOUTHERN TWRS - PENTAGON</Name><RouteID>7Av2</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7A - LINCOLNIA/QUANTRELL - PENT VIA PENT</Name><RouteID>7Av3</RouteID></Route><Route><LineDescription>Park Center-Pentagon Line</LineDescription><Name>7C - PARK CENTER - PENTAGON</Name><RouteID>7C</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7F - LINCOLNIA+QUANTRELL - PENTAGON</Name><RouteID>7F</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7F - LINC + QUANT - PENT CITY - PENTAGON</Name><RouteID>7Fv1</RouteID></Route><Route><LineDescription>Mark Center-Pentagon Line</LineDescription><Name>7M - MARK CENTER - PENTAGON (NON-HOV)</Name><RouteID>7M</RouteID></Route><Route><LineDescription>Mark Center-Pentagon Line</LineDescription><Name>7M - MARK CENTER - PENTAGON (HOV)</Name><RouteID>7Mv1</RouteID></Route><Route><LineDescription>Park Center-Pentagon Line</LineDescription><Name>7P - PARK CTR - PENTAGON</Name><RouteID>7P</RouteID></Route><Route><LineDescription>Lincolnia-Pentagon Line</LineDescription><Name>7W - LNCLNA+QUANTRLL- PENTAGON</Name><RouteID>7W</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7Y - SOUTHERN TWRS - H+17TH ST</Name><RouteID>7Y</RouteID></Route><Route><LineDescription>Lincolnia-North Fairlington Line</LineDescription><Name>7Y - PENTAGON - H+17TH ST</Name><RouteID>7Yv1</RouteID></Route><Route><LineDescription>North Capitol Street Line</LineDescription><Name>80 - KENNEDY CTR   - FORT TOTTEN STA</Name><RouteID>80</RouteID></Route><Route><LineDescription>North Capitol Street Line</LineDescription><Name>80 - MCPHERSON SQ  - BROOKLAND</Name><RouteID>80v1</RouteID></Route><Route><LineDescription>North Capitol Street Line</LineDescription><Name>80 - MCPHERSON SQ  - FORT TOTTEN STA</Name><RouteID>80v2</RouteID></Route><Route><LineDescription>North Capitol Street Line</LineDescription><Name>80 - KENNEDY CTR   - BROOKLAND STA</Name><RouteID>80v3</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>83 - RHODE ISLAND AVE STA-CHERRY HILL</Name><RouteID>83</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>83 - MT RAINIER - RHODE ISLAND</Name><RouteID>83v1</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>83 - RHODE ISLAND - MT RAINIER</Name><RouteID>83v2</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>83 - RHODE ISLAND AVE STA-COLLEGE PARK</Name><RouteID>83v3</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>83 - COLLEGE PARK-RHODE ISLAND AVE STA</Name><RouteID>83v4</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>86 - RHODE ISLAND AVE STA- CALVERTON</Name><RouteID>86</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>86 - RHODE ISLAND AVE STA- COLLEGE PARK</Name><RouteID>86v1</RouteID></Route><Route><LineDescription>College Park Line</LineDescription><Name>86 - MT RAINIER   - CALVERTON</Name><RouteID>86v2</RouteID></Route><Route><LineDescription>Laurel Express Line</LineDescription><Name>87 - NEW CARROLTON -CYPRESS+LAURL LAKES</Name><RouteID>87</RouteID></Route><Route><LineDescription>Laurel Express Line</LineDescription><Name>87 - GRNBELT STA -CYPRESS+LAURL LAKES</Name><RouteID>87v1</RouteID></Route><Route><LineDescription>Laurel Express Line</LineDescription><Name>87 - GRNBELT-CYP+LRL LAKES (NO P&amp;R)</Name><RouteID>87v2</RouteID></Route><Route><LineDescription>Laurel Express Line</LineDescription><Name>87 - GRNBELT STA - BALTIMORE+MAIN ST</Name><RouteID>87v3</RouteID></Route><Route><LineDescription>Laurel Express Line</LineDescription><Name>87 - BALTIMORE+MAIN ST - GRNBELT STA</Name><RouteID>87v4</RouteID></Route><Route><LineDescription>Laurel Express Line</LineDescription><Name>87 - CYPRESS+LAURL LAKES -GRNBELT STA</Name><RouteID>87v5</RouteID></Route><Route><LineDescription>Laurel Line</LineDescription><Name>89 - GREENBELT STA - 4TH &amp; GREEN HILL</Name><RouteID>89</RouteID></Route><Route><LineDescription>Laurel Line</LineDescription><Name>89 - GREENBELT STA - CHERRY LA+4TH ST</Name><RouteID>89v1</RouteID></Route><Route><LineDescription>Laurel Line</LineDescription><Name>89M - GREENBELT STA - S LAUREL P+R</Name><RouteID>89M</RouteID></Route><Route><LineDescription>Foxchase-Seminary Valley Line</LineDescription><Name>8S - RADFORD+QUAKER - PENTAGON</Name><RouteID>8S</RouteID></Route><Route><LineDescription>Foxchase-Seminary Valley Line</LineDescription><Name>8W - MARK CENTER - PENTAGON V FOXCHASE</Name><RouteID>8W</RouteID></Route><Route><LineDescription>Foxchase-Seminary Valley Line</LineDescription><Name>8Z - QUAKER+OSAGE - PENTAGON</Name><RouteID>8Z</RouteID></Route><Route><LineDescription>U Street-Garfield Line</LineDescription><Name>90 - ANACOSTIA - DK ELLNGTN BRDG</Name><RouteID>90</RouteID></Route><Route><LineDescription>U Street-Garfield Line</LineDescription><Name>90 - 8TH ST + L ST  - DK ELLNGTN BRDG</Name><RouteID>90v1</RouteID></Route><Route><LineDescription>U Street-Garfield Line</LineDescription><Name>90 - KIPP DC PREP- ANACOSTIA</Name><RouteID>90v2</RouteID></Route><Route><LineDescription>U Street-Garfield Line</LineDescription><Name>92 - CONGRESS HTS STA - REEVES CTR</Name><RouteID>92</RouteID></Route><Route><LineDescription>U Street-Garfield Line</LineDescription><Name>92 - EASTERN MARKET - CONGRESS HGTS STA</Name><RouteID>92v1</RouteID></Route><Route><LineDescription>U Street-Garfield Line</LineDescription><Name>92 - CONGRESS HTS STA- EASTERN MARKET</Name><RouteID>92v2</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>96 - TENLEYTOWN STA - CAPITOL HTS STA</Name><RouteID>96</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>96 - ELLINGTON BR - CAPITOL HTS</Name><RouteID>96v1</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>96 - CAPITOL HTS  - REEVES CTR</Name><RouteID>96v2</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>96 - CAPITOL HTS - ELLINGTON BRDG</Name><RouteID>96v3</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>96 - REEVES CTR  - CAPITOL HTS</Name><RouteID>96v4</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>96 - TENLEYTOWN STA - STADIUM ARMORY STA</Name><RouteID>96v5</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>97 - UNION STATION - CAPITOL HTS</Name><RouteID>97</RouteID></Route><Route><LineDescription>East Capitol Street-Cardoza Line</LineDescription><Name>97 - EASTERN HS - CAPITOL HTS</Name><RouteID>97v1</RouteID></Route><Route><LineDescription>Martin Luther King Jr Highway Line</LineDescription><Name>A12 - ADDISON RD STA - CAPITAL PLAZA</Name><RouteID>A12</RouteID></Route><Route><LineDescription>Martin Luther King Jr Highway Line</LineDescription><Name>A12 - BARLWE+MATHEW HENSON - ADDISON RD</Name><RouteID>A12v1</RouteID></Route><Route><LineDescription>Martin Luther King Jr Highway Line</LineDescription><Name>A12 - CAPITOL HTS STA - CAPITAL PLAZA</Name><RouteID>A12v2</RouteID></Route><Route><LineDescription>Martin Luther King Jr Highway Line</LineDescription><Name>A12 - CAPITAL PLAZA - CAPITOL HTS</Name><RouteID>A12v3</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A2 - SOUTHERN AVE - ANACOSTIA (VIA HOSP)</Name><RouteID>A2</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A2 - ANACOSTIA - MISS+ATLANTIC</Name><RouteID>A2v1</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A2 - SOUTHERN AVE STA - ANACOSTIA</Name><RouteID>A2v2</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A2 - MISS+ATLANTIC - ANACOSTIA</Name><RouteID>A2v3</RouteID></Route><Route><LineDescription>Minnesota Ave-Anacostia Line</LineDescription><Name>A31 - ANACOSTIA HIGH - MINNESOTA STA</Name><RouteID>A31</RouteID></Route><Route><LineDescription>Minnesota Ave-Anacostia Line</LineDescription><Name>A32 - ANACOSTIA HIGH - SOUTHRN AVE STA</Name><RouteID>A32</RouteID></Route><Route><LineDescription>Minnesota Ave-Anacostia Line</LineDescription><Name>A33 - ANACOSTIA HIGH - ANACOSTIA STA</Name><RouteID>A33</RouteID></Route><Route><LineDescription>Anacostia-Fort Drum Line</LineDescription><Name>A4 - DC VILLAGE - ANACOSTIA</Name><RouteID>A4</RouteID></Route><Route><LineDescription>Anacostia-Fort Drum Line</LineDescription><Name>A4 - USCG-FT DRUM (VIA ANAC)</Name><RouteID>A4v1</RouteID></Route><Route><LineDescription>Anacostia-Fort Drum Line</LineDescription><Name>A4 - FT DRUM - ANACOSTIA</Name><RouteID>A4v2</RouteID></Route><Route><LineDescription>Anacostia-Fort Drum Line</LineDescription><Name>A4 - ANACOSTIA - FT DRUM</Name><RouteID>A4v3</RouteID></Route><Route><LineDescription>Anacostia-Fort Drum Line</LineDescription><Name>A4 - FT DRUM-USCG (VIA ANAC)</Name><RouteID>A4v4</RouteID></Route><Route><LineDescription>Anacostia-Fort Drum Line</LineDescription><Name>A4 - DC VILL-USCG (VIA ANAC)</Name><RouteID>A4v5</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A6 - 4501 3RD ST - ANACOSTIA</Name><RouteID>A6</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A6 - SOUTHERN AVE+S CAPITOL - ANACOSTIA</Name><RouteID>A6v1</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A7 - SOUTHRN+S CAPITOL - ANACOSTIA</Name><RouteID>A7</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A8 - 4501 3RD ST - ANACOSTIA</Name><RouteID>A8</RouteID></Route><Route><LineDescription>Anacostia-Congress Heights Line</LineDescription><Name>A8 - SOUTHERN + S CAPITOL - ANACOSTIA</Name><RouteID>A8v1</RouteID></Route><Route><LineDescription>Martin Luther King Jr Ave Limited Line</LineDescription><Name>A9 - LIVINGSTON - MCPHERSON SQUARE</Name><RouteID>A9</RouteID></Route><Route><LineDescription>Bladensburg Road-Anacostia Line</LineDescription><Name>B2 - ANACOSTIA STA - MOUNT RAINIER</Name><RouteID>B2</RouteID></Route><Route><LineDescription>Bladensburg Road-Anacostia Line</LineDescription><Name>B2 - ANACOSTIA STA - BLDNSGRG+VST NE</Name><RouteID>B2v1</RouteID></Route><Route><LineDescription>Bladensburg Road-Anacostia Line</LineDescription><Name>B2 - POTOMAC AVE - MOUNT RAINIER</Name><RouteID>B2v2</RouteID></Route><Route><LineDescription>Bladensburg Road-Anacostia Line</LineDescription><Name>B2 - BLDNSBRG+26TH - ANACOSTIA STA</Name><RouteID>B2v3</RouteID></Route><Route><LineDescription>Bladensburg Road-Anacostia Line</LineDescription><Name>B2 - EASTERN HS - ANACOSTIA STA</Name><RouteID>B2v4</RouteID></Route><Route><LineDescription>Bowie State University Line</LineDescription><Name>B21 - NEW CARROLTON STA - BOWIE STATE</Name><RouteID>B21</RouteID></Route><Route><LineDescription>Bowie State University Line</LineDescription><Name>B22 - NEW CARROLTON STA - BOWIE STATE</Name><RouteID>B22</RouteID></Route><Route><LineDescription>Bowie State University Line</LineDescription><Name>B22 - OLD CHAPEL &amp; 197 - NEW CARRLTN STA</Name><RouteID>B22v1</RouteID></Route><Route><LineDescription>Bowie-Belair Line</LineDescription><Name>B24 - NEW CARLTN STA-BOWIE P+R(VIA BHC)</Name><RouteID>B24</RouteID></Route><Route><LineDescription>Bowie-Belair Line</LineDescription><Name>B24 - NEW CARROLTON STA - BOWIE P+R</Name><RouteID>B24v1</RouteID></Route><Route><LineDescription>Bowie-New Carrollton Line</LineDescription><Name>B27 - NEW CARROLTON STA - BOWIE STATE</Name><RouteID>B27</RouteID></Route><Route><LineDescription>Crofton-New Carrollton Line</LineDescription><Name>B29 - NEW CARROLTON - CROFTON CC (PM)(PR)</Name><RouteID>B29</RouteID></Route><Route><LineDescription>Crofton-New Carrollton Line</LineDescription><Name>B29 - NEW CARROLLTON STA - GATEWAY CTR</Name><RouteID>B29v1</RouteID></Route><Route><LineDescription>Crofton-New Carrollton Line</LineDescription><Name>B29 - NEW CARROLTON - CROFTON CC (NO PR)</Name><RouteID>B29v2</RouteID></Route><Route><LineDescription>Greenbelt-BWI Thurgood Marshall Airport Express Line</LineDescription><Name>B30 - GREENBELT STA - BWI LT RAIL STA</Name><RouteID>B30</RouteID></Route><Route><LineDescription>Fort Lincoln Shuttle Line</LineDescription><Name>B8 - RHODE ISLAND AV -PETERSBRG APTS</Name><RouteID>B8</RouteID></Route><Route><LineDescription>Fort Lincoln Shuttle Line</LineDescription><Name>B8 - BLDNSBRG+S DKTA -PETERSBRG APTS</Name><RouteID>B8v1</RouteID></Route><Route><LineDescription>Fort Lincoln Shuttle Line</LineDescription><Name>B8 - PETERSBRG APTS  -BLDNSBRG+S DKTA</Name><RouteID>B8v2</RouteID></Route><Route><LineDescription>Fort Lincoln Shuttle Line</LineDescription><Name>B9 - RHODE ISLAND AVE - COLMAR MANOR</Name><RouteID>B9</RouteID></Route><Route><LineDescription>Clinton Line</LineDescription><Name>C11 - CLINTON P+R - BRANCH AVE STA</Name><RouteID>C11</RouteID></Route><Route><LineDescription>Hillcrest Heights Line</LineDescription><Name>C12 - NAYLOR RD STA  - BRANCH AVE STA</Name><RouteID>C12</RouteID></Route><Route><LineDescription>Clinton Line</LineDescription><Name>C13 - CLINTON P+R - BRANCH AVE STA</Name><RouteID>C13</RouteID></Route><Route><LineDescription>Hillcrest Heights Line</LineDescription><Name>C14 - NAYLOR RD STA  - BRANCH AVE STA</Name><RouteID>C14</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C2 - WHEATN STA - GRNBELT STA UMD ALT</Name><RouteID>C2</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C2 - TAKOMA LANGLEY XROADS - GRNBELT STA</Name><RouteID>C2v1</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C2 - GREENBELT STA - RANDOLPH + PARKLAWN</Name><RouteID>C2v2</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C2 - TAKOMA LANGLEY XROADS - WHEATON</Name><RouteID>C2v3</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C21 - ADDISON RD STA - COLLINGTON</Name><RouteID>C21</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C21 - ADDISON RD STA - POINTER RIDGE</Name><RouteID>C21v1</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C21 - ADDISON RD STA - CAMPUS WAY S</Name><RouteID>C21v2</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C22 - ADDISON RD STA - COLLINGTON</Name><RouteID>C22</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C22 - ADDISON RD STA- POINTER RIDGE</Name><RouteID>C22v1</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C26 - LARGO TOWN CTR - WATKNS+CHESTERTON</Name><RouteID>C26</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C26 - LARGO TOWN CTR-WATKINS+CAMBLETON</Name><RouteID>C26v1</RouteID></Route><Route><LineDescription>Pointer Ridge Line</LineDescription><Name>C28 - PT RIDGE -NEW CAROLTN VIA HLTH CTR</Name><RouteID>C28</RouteID></Route><Route><LineDescription>Pointer Ridge Line</LineDescription><Name>C28 - PT RIDGE - NEW CARROLLTON STA</Name><RouteID>C28v1</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C29 - POINTER RIDGE - ADDISON RD STA</Name><RouteID>C29*1</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C29 - WATKNS+CAMBLETON - ADDISN RD STA</Name><RouteID>C29*2</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C29 - ADDISON RD STA - BOWIE STATE</Name><RouteID>C29*4</RouteID></Route><Route><LineDescription>Central Avenue Line</LineDescription><Name>C29 - ADDISON RD STA - POINTER RIDGE</Name><RouteID>C29/</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C4 - TWINBROOK STA - PG PLAZA STA</Name><RouteID>C4</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C4 - TLTC-TWINBROOK</Name><RouteID>C4v1</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C4 - TWINBROOK STA - WHEATON STA</Name><RouteID>C4v2</RouteID></Route><Route><LineDescription>Greenbelt-Twinbrook Line</LineDescription><Name>C4 - PG PLAZA STA - WHEATON STA</Name><RouteID>C4v3</RouteID></Route><Route><LineDescription>College Park-White Flint Line</LineDescription><Name>C8 - WHITE FLINT - COLLEGE PARK</Name><RouteID>C8</RouteID></Route><Route><LineDescription>College Park-White Flint Line</LineDescription><Name>C8 - WHITE FLNT-COLLEGE PK (NO FDA/ARCH)</Name><RouteID>C8v1</RouteID></Route><Route><LineDescription>College Park-White Flint Line</LineDescription><Name>C8 - WHITE FLINT-COLLEGE PARK (NO FDA)</Name><RouteID>C8v2</RouteID></Route><Route><LineDescription>College Park-White Flint Line</LineDescription><Name>C8 - GLENMONT-COLLEGE PK (NO FDA/ARCH)</Name><RouteID>C8v3</RouteID></Route><Route><LineDescription>Glover Park-Franklin Square Line</LineDescription><Name>D1 - GLOVER PARK - FRANKLIN SQUARE</Name><RouteID>D1</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D12 - SOUTHERN AVE STA - SUITLAND STA</Name><RouteID>D12</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D12 - SUITLAND STA - SOUTHERN AVE STA</Name><RouteID>D12v1</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D12 - ST BARNABAS RD   - SUITLAND STA</Name><RouteID>D12v2</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D13 - SOUTHERN AVE STA - SUITLAND STA</Name><RouteID>D13</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D13 - SOUTHRN STA - ALLENTWN+OLD BRNCH</Name><RouteID>D13v1</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D14 - SOUTHERN AVE STA - SUITLAND STA</Name><RouteID>D14</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D14 - ALLENTWN+OLD BRNCH - SUITLND STA</Name><RouteID>D14v1</RouteID></Route><Route><LineDescription>Oxon Hill-Suitland Line</LineDescription><Name>D14 - SUITLND STA - ALLENTWN+OLD BRNCH</Name><RouteID>D14v2</RouteID></Route><Route><LineDescription>Glover Park-Dupont Circle Line</LineDescription><Name>D2 - GLOVER PARK   - CONNETICUT +Q ST</Name><RouteID>D2</RouteID></Route><Route><LineDescription>Glover Park-Dupont Circle Line</LineDescription><Name>D2 - DUPONT CIR - GLOVER PARK</Name><RouteID>D2v1</RouteID></Route><Route><LineDescription>16th Street-Tenleytown Line</LineDescription><Name>D31 - TENLEYTOWN STA - 16TH + EASTERN</Name><RouteID>D31</RouteID></Route><Route><LineDescription>16th Street-Tenleytown Line</LineDescription><Name>D32 - TENLEYTOWN STA - COLUMBIA HTS STA</Name><RouteID>D32</RouteID></Route><Route><LineDescription>16th Street-Tenleytown Line</LineDescription><Name>D33 - TENLEYTOWN STA - 16TH + SHEPHERD</Name><RouteID>D33</RouteID></Route><Route><LineDescription>16th Street-Tenleytown Line</LineDescription><Name>D34 - TENLEYTOWN STA - 14TH + COLORADO</Name><RouteID>D34</RouteID></Route><Route><LineDescription>Ivy City-Franklin Square Line</LineDescription><Name>D4 - DUPONT CIRCLE - IVY CITY</Name><RouteID>D4</RouteID></Route><Route><LineDescription>Ivy City-Franklin Square Line</LineDescription><Name>D4 - FRANKLIN SQUARE - IVY CITY</Name><RouteID>D4v1</RouteID></Route><Route><LineDescription>Ivy City-Franklin Square Line</LineDescription><Name>D4 - IVY CITY - FRANKLIN SQUARE</Name><RouteID>D4v2</RouteID></Route><Route><LineDescription>MacArthur Blvd-Georgetown Line</LineDescription><Name>D5 - MASS LTTL FLWR CHRCH- FARRGT SQR</Name><RouteID>D5</RouteID></Route><Route><LineDescription>Sibley Hospital–Stadium-Armory Line</LineDescription><Name>D6 - SIBLEY HOSP - STADIUM ARMRY</Name><RouteID>D6</RouteID></Route><Route><LineDescription>Sibley Hospital–Stadium-Armory Line</LineDescription><Name>D6 - STADIUM ARMRY STA - FARRAGUT SQUARE</Name><RouteID>D6v1</RouteID></Route><Route><LineDescription>Sibley Hospital–Stadium-Armory Line</LineDescription><Name>D6 - SIBLEY HOSP - FARRAGUT SQ</Name><RouteID>D6v2</RouteID></Route><Route><LineDescription>Sibley Hospital–Stadium-Armory Line</LineDescription><Name>D6 - FARRAGUT SQUARE - STADIUM ARMRY</Name><RouteID>D6v3</RouteID></Route><Route><LineDescription>Hospital Center Line</LineDescription><Name>D8 - UNION STATION - VA MED CTR</Name><RouteID>D8</RouteID></Route><Route><LineDescription>Hospital Center Line</LineDescription><Name>D8 - RHODE ISLAND STA - UNION STA</Name><RouteID>D8v1</RouteID></Route><Route><LineDescription>Ivy City-Fort Totten Line</LineDescription><Name>E2 - IVY CITY - FT TOTTEN</Name><RouteID>E2</RouteID></Route><Route><LineDescription>Military Road-Crosstown Line</LineDescription><Name>E4 - FRIENDSHP HTS - RIGGS PK</Name><RouteID>E4</RouteID></Route><Route><LineDescription>Military Road-Crosstown Line</LineDescription><Name>E4 - FRIENDSHP HTS - FT TOTTEN</Name><RouteID>E4v1</RouteID></Route><Route><LineDescription>Military Road-Crosstown Line</LineDescription><Name>E4 - FT TOTTEN - FRIENDSHP HTS</Name><RouteID>E4v2</RouteID></Route><Route><LineDescription>Chevy Chase Line</LineDescription><Name>E6 - FRIENDSHP HTS  -KNOLLWOOD</Name><RouteID>E6</RouteID></Route><Route><LineDescription>Chillum Road Line</LineDescription><Name>F1 - CHEVERLY STA - TAKOMA STA</Name><RouteID>F1</RouteID></Route><Route><LineDescription>Ardwick Industrial Park Shuttle Line</LineDescription><Name>F12 - CHEVERLY STA - NEW CARROLTON STA</Name><RouteID>F12</RouteID></Route><Route><LineDescription>Ardwick Industrial Park Shuttle Line</LineDescription><Name>F12 - CHEVERLY STA - LANDOVER STA</Name><RouteID>F12v1</RouteID></Route><Route><LineDescription>Cheverly-Washington Business Park Line</LineDescription><Name>F13 - CHEVERLY STA-WASHINGTON BUS PARK</Name><RouteID>F13</RouteID></Route><Route><LineDescription>Cheverly-Washington Business Park Line</LineDescription><Name>F13 - CHEVERLY STA - WASHINGTON BUS PARK</Name><RouteID>F13v1</RouteID></Route><Route><LineDescription>Cheverly-Washington Business Park Line</LineDescription><Name>F13 - NEW CARROLTON - WASHINGTON BUS PARK</Name><RouteID>F13v2</RouteID></Route><Route><LineDescription>Cheverly-Washington Business Park Line</LineDescription><Name>F13 - WHITFIELD + VOLTA - CHEVERLY STA</Name><RouteID>F13v3</RouteID></Route><Route><LineDescription>Sheriff Road-Capitol Heights Line</LineDescription><Name>F14 - NAYLOR RD STA -NEW CARROLTON STA</Name><RouteID>F14</RouteID></Route><Route><LineDescription>Sheriff Road-Capitol Heights Line</LineDescription><Name>F14 - BRADBURY HGTS -NEW CARROLTON STA</Name><RouteID>F14v1</RouteID></Route><Route><LineDescription>Chillum Road Line</LineDescription><Name>F2 - CHEVERLY STA - TAKOMA STA</Name><RouteID>F2</RouteID></Route><Route><LineDescription>Chillum Road Line</LineDescription><Name>F2 - CHEVRLY STA-QUEENS CHAPL+CARSN CIR</Name><RouteID>F2v1</RouteID></Route><Route><LineDescription>Chillum Road Line</LineDescription><Name>F2 - 34TH + OTIS - TAKOMA STA</Name><RouteID>F2v2</RouteID></Route><Route><LineDescription>New Carrollton-Silver Spring Line</LineDescription><Name>F4 - SILVR SPRING STA - NEW CARRLLTON</Name><RouteID>F4</RouteID></Route><Route><LineDescription>New Carrollton-Silver Spring Line</LineDescription><Name>F4 - PG PLAZA STA -NEW CARROLTON STA</Name><RouteID>F4v1</RouteID></Route><Route><LineDescription>New Carrollton-Silver Spring Line</LineDescription><Name>F4 - NEW CARRLLTON - PG PLAZA STA</Name><RouteID>F4v2</RouteID></Route><Route><LineDescription>New Carrollton-Fort Totten Line</LineDescription><Name>F6 - FT TOTTEN - NEW CARROLLTN</Name><RouteID>F6</RouteID></Route><Route><LineDescription>New Carrollton-Fort Totten Line</LineDescription><Name>F6 - PG PLAZA - NEW CARROLLTON STA</Name><RouteID>F6v1</RouteID></Route><Route><LineDescription>New Carrollton-Fort Totten Line</LineDescription><Name>F6 - NEW CARROLLTON - PG PLAZA STA</Name><RouteID>F6v2</RouteID></Route><Route><LineDescription>Langley Park-Cheverly Line</LineDescription><Name>F8 - CHEVERLY STA - TLTC</Name><RouteID>F8</RouteID></Route><Route><LineDescription>Greenbelt-New Carrollton Line</LineDescription><Name>G12 - GREENBELT STA - NEW CARROLLTON STA</Name><RouteID>G12</RouteID></Route><Route><LineDescription>Greenbelt-New Carrollton Line</LineDescription><Name>G12 - GREENBELT STA - ROOSEVELT CENTER</Name><RouteID>G12v1</RouteID></Route><Route><LineDescription>Greenbelt-New Carrollton Line</LineDescription><Name>G12 - ROOSEVELT CTR - NEW CARROLLTON STA</Name><RouteID>G12v2</RouteID></Route><Route><LineDescription>Greenbelt-New Carrollton Line</LineDescription><Name>G14 - GREENBELT STA -NEW CARROLTON STA</Name><RouteID>G14</RouteID></Route><Route><LineDescription>Greenbelt-New Carrollton Line</LineDescription><Name>G14 - ROOSEVELT CTR -NEW CARROLTON STA</Name><RouteID>G14v1</RouteID></Route><Route><LineDescription>Greenbelt-New Carrollton Line</LineDescription><Name>G14 - GREENBELT STA - NEW CARROLTON STA</Name><RouteID>G14v2</RouteID></Route><Route><LineDescription>P Street-LeDroit Park Line</LineDescription><Name>G2 - GEORGETOWN UNIV - HOWARD UNIV</Name><RouteID>G2</RouteID></Route><Route><LineDescription>P Street-LeDroit Park Line</LineDescription><Name>G2 - GEORGETOWN UNIV - P+14TH ST</Name><RouteID>G2v1</RouteID></Route><Route><LineDescription>Rhode Island Avenue Line</LineDescription><Name>G8 - FARRAGUT SQUARE - AVONDALE</Name><RouteID>G8</RouteID></Route><Route><LineDescription>Rhode Island Avenue Line</LineDescription><Name>G8 - RHODE ISLAND AVE STA-FARRAGUT SQ</Name><RouteID>G8v1</RouteID></Route><Route><LineDescription>Rhode Island Avenue Line</LineDescription><Name>G8 - FARRAGUT SQ - RHODE ISLAND STA</Name><RouteID>G8v2</RouteID></Route><Route><LineDescription>Rhode Island Avenue Line</LineDescription><Name>G8 - BROOKLAND STA - AVONDALE</Name><RouteID>G8v3</RouteID></Route><Route><LineDescription>Rhode Island Avenue Limited Line</LineDescription><Name>G9 - FRANKLIN SQ - RHODE ISLD + EAST PM</Name><RouteID>G9</RouteID></Route><Route><LineDescription>Rhode Island Avenue Limited Line</LineDescription><Name>G9 - FRANKLIN SQ - RHODE ISLD + EASTERN</Name><RouteID>G9v1</RouteID></Route><Route><LineDescription>Brookland-Potomac Park Line</LineDescription><Name>H1 - C + 17TH ST - BROOKLAND CUA STA</Name><RouteID>H1</RouteID></Route><Route><LineDescription>Marlow Heights-Temple Hills Line</LineDescription><Name>H11 - HEATHER HILL - NAYLOR RD STA</Name><RouteID>H11</RouteID></Route><Route><LineDescription>Marlow Heights-Temple Hills Line</LineDescription><Name>H12 - HEATHER HILL-NAYLOR RD (MACYS)</Name><RouteID>H12</RouteID></Route><Route><LineDescription>Marlow Heights-Temple Hills Line</LineDescription><Name>H12 - HEATHER HILL-NAYLOR RD STA (P&amp;R)</Name><RouteID>H12v1</RouteID></Route><Route><LineDescription>Marlow Heights-Temple Hills Line</LineDescription><Name>H13 - HEATHER HILL - NAYLOR RD STA</Name><RouteID>H13</RouteID></Route><Route><LineDescription>Crosstown Line</LineDescription><Name>H2 - TENLEYTOWN AU-ST - BROOKLND CUA STA</Name><RouteID>H2</RouteID></Route><Route><LineDescription>Crosstown Line</LineDescription><Name>H3 - TENLEYTOWN STA -BROOKLND CUA STA</Name><RouteID>H3</RouteID></Route><Route><LineDescription>Crosstown Line</LineDescription><Name>H4 - TENLEYTWN AU-STA -BROOKLAND CUA STA</Name><RouteID>H4</RouteID></Route><Route><LineDescription>Crosstown Line</LineDescription><Name>H4 - COLUMBIA RD+14TH -TENLEYTOWN-AU STA</Name><RouteID>H4v1</RouteID></Route><Route><LineDescription>Brookland-Fort Lincoln Line</LineDescription><Name>H6 - BROOKLAND CUA STA - FORT LINCOLN</Name><RouteID>H6</RouteID></Route><Route><LineDescription>Brookland-Fort Lincoln Line</LineDescription><Name>H6 - FORT LINCOLN - FORT LINCOLN</Name><RouteID>H6v1</RouteID></Route><Route><LineDescription>Park Road-Brookland Line</LineDescription><Name>H8 - MT PLEASANT+17TH - RHODE ISLAND</Name><RouteID>H8</RouteID></Route><Route><LineDescription>Park Road-Brookland Line</LineDescription><Name>H9 - RHODE ISLAND - FORT DR + 1ST</Name><RouteID>H9</RouteID></Route><Route><LineDescription>Bethesda-Silver Spring Line</LineDescription><Name>J1 - MEDCAL CTR STA - SILVR SPRNG STA</Name><RouteID>J1</RouteID></Route><Route><LineDescription>Bethesda-Silver Spring Line</LineDescription><Name>J1 - SILVR SPRNG STA - MONT MALL/BAT LN</Name><RouteID>J1v1</RouteID></Route><Route><LineDescription>Marlboro Pike Line</LineDescription><Name>J12 - ADDISON RD STA - FORESTVILLE</Name><RouteID>J12</RouteID></Route><Route><LineDescription>Marlboro Pike Line</LineDescription><Name>J12 - ADDISON RD-FORESTVILLE VIA PRES PKY</Name><RouteID>J12v1</RouteID></Route><Route><LineDescription>Bethesda-Silver Spring Line</LineDescription><Name>J2 - MONTGOMRY MALL - SILVR SPRNG STA</Name><RouteID>J2</RouteID></Route><Route><LineDescription>Bethesda-Silver Spring Line</LineDescription><Name>J2 - MEDICAL CTR STA - SILVR SPRNG STA</Name><RouteID>J2v1</RouteID></Route><Route><LineDescription>Bethesda-Silver Spring Line</LineDescription><Name>J2 - SILVER SPRNG - MONT MALL/BATTRY LA</Name><RouteID>J2v2</RouteID></Route><Route><LineDescription>College Park-Bethesda Limited</LineDescription><Name>J4 - BETHESDA STA - COLLEGE PARK STA</Name><RouteID>J4</RouteID></Route><Route><LineDescription>Forestville Line</LineDescription><Name>K12 - BRANCH AVE ST-SUITLAND ST</Name><RouteID>K12</RouteID></Route><Route><LineDescription>Forestville Line</LineDescription><Name>K12 - PENN MAR - SUITLAND ST</Name><RouteID>K12v1</RouteID></Route><Route><LineDescription>Forestville Line</LineDescription><Name>K12 - BRANCH AVE - PENN MAR</Name><RouteID>K12v2</RouteID></Route><Route><LineDescription>Takoma-Fort Totten Line</LineDescription><Name>K2 - FT TOTTEN STA - TAKOMA STA</Name><RouteID>K2</RouteID></Route><Route><LineDescription>New Hampshire Ave-Maryland Line</LineDescription><Name>K6 - FORT TOTTEN STA - WHITE OAK</Name><RouteID>K6</RouteID></Route><Route><LineDescription>New Hampshire Ave-Maryland Line</LineDescription><Name>K6 - TLTC - FORT TOTTEN STA</Name><RouteID>K6v1</RouteID></Route><Route><LineDescription>New Hampshire Ave-Maryland Limited Line</LineDescription><Name>K9 - FORT TOTTEN STA - FDA/FRC (PM)</Name><RouteID>K9</RouteID></Route><Route><LineDescription>New Hampshire Ave-Maryland Limited Line</LineDescription><Name>K9 - FORT TOTTEN STA - FDA/FRC (AM)</Name><RouteID>K9v1</RouteID></Route><Route><LineDescription>Connecticut Ave Line</LineDescription><Name>L1 - POTOMAC PK  - CHEVY CHASE</Name><RouteID>L1</RouteID></Route><Route><LineDescription>Connecticut Ave Line</LineDescription><Name>L2 - FARRAGUT SQ - CHEVY CHASE</Name><RouteID>L2</RouteID></Route><Route><LineDescription>Connecticut Ave Line</LineDescription><Name>L2 - VAN NESS-UDC STA - CHEVY CHASE</Name><RouteID>L2v1</RouteID></Route><Route><LineDescription>Connecticut Ave Line</LineDescription><Name>L2 - FARRAGUT SQ - BETHESDA</Name><RouteID>L2v2</RouteID></Route><Route><LineDescription>Connecticut Ave-Maryland Line</LineDescription><Name>L8 - FRIENDSHIP HTS STA - ASPEN HILL</Name><RouteID>L8</RouteID></Route><Route><LineDescription>Nebraska Ave Line</LineDescription><Name>M4 - SIBLEY HOSPITAL - TENLEYTOWN/AU STA</Name><RouteID>M4</RouteID></Route><Route><LineDescription>Nebraska Ave Line</LineDescription><Name>M4 - TENLEYTOWN   - PINEHRST CIR</Name><RouteID>M4v1</RouteID></Route><Route><LineDescription>Nebraska Ave Line</LineDescription><Name>M4 - PINEHRST CIR - TENLEYTOWN</Name><RouteID>M4v2</RouteID></Route><Route><LineDescription>Fairfax Village Line</LineDescription><Name>M6 - POTOMAC AVE - ALABAMA + PENN</Name><RouteID>M6</RouteID></Route><Route><LineDescription>Fairfax Village Line</LineDescription><Name>M6 - ALABAMA + PENN - FAIRFAX VILLAGE</Name><RouteID>M6v1</RouteID></Route><Route><LineDescription>Metroway-Potomac Yard Line</LineDescription><Name>MW1 - BRADDOCK RD - PENTAGON CITY</Name><RouteID>MW1</RouteID></Route><Route><LineDescription>Metroway-Potomac Yard Line</LineDescription><Name>MW1 - POTOMAC YARD - CRYSTAL CITY</Name><RouteID>MW1v1</RouteID></Route><Route><LineDescription>Metroway-Potomac Yard Line</LineDescription><Name>MW1 - CRYSTAL CITY - BRADDOCK RD</Name><RouteID>MW1v2</RouteID></Route><Route><LineDescription>Metroway-Potomac Yard Line</LineDescription><Name>MW1 - BRADDOCK RD - CRYSTAL CITY</Name><RouteID>MW1v3</RouteID></Route><Route><LineDescription>Massachusetts Ave Line</LineDescription><Name>N2 - FRIENDSHIP HTS - FARRAGUT SQ</Name><RouteID>N2</RouteID></Route><Route><LineDescription>Massachusetts Ave Line</LineDescription><Name>N4 - FRIENDSHP HTS - POTOMAC PARK</Name><RouteID>N4</RouteID></Route><Route><LineDescription>Massachusetts Ave Line</LineDescription><Name>N4 - FRIENDSHP HTS - FARRAGUT  SQ</Name><RouteID>N4v1</RouteID></Route><Route><LineDescription>Massachusetts Ave Line</LineDescription><Name>N6 - FRNDSHIP HTS - FARRAGUT SQ</Name><RouteID>N6</RouteID></Route><Route><LineDescription>National Harbor-Southern Avenue Line</LineDescription><Name>NH1 - NATIONAL HARBOR-SOUTHERN AVE STA</Name><RouteID>NH1</RouteID></Route><Route><LineDescription>National Harbor-Alexandria Line</LineDescription><Name>NH2 - HUNTINGTON STA-NATIONAL HARBOR</Name><RouteID>NH2</RouteID></Route><Route><LineDescription>Eastover-Addison Road Line</LineDescription><Name>P12 - EASTOVER - ADDISON RD STA (NO HOSP)</Name><RouteID>P12</RouteID></Route><Route><LineDescription>Eastover-Addison Road Line</LineDescription><Name>P12 - IVERSON MALL - ADDISON RD STA</Name><RouteID>P12v1</RouteID></Route><Route><LineDescription>Eastover-Addison Road Line</LineDescription><Name>P12 - SUITLAND STA - ADDISON RD STA</Name><RouteID>P12v2</RouteID></Route><Route><LineDescription>Oxon Hill-Fort Washington Line</LineDescription><Name>P18 - FT WASH P+R LOT - SOUTHERN AVE</Name><RouteID>P18</RouteID></Route><Route><LineDescription>Oxon Hill-Fort Washington Line</LineDescription><Name>P19 - FT WASH P+R LOT-SOUTHERN AVE STA</Name><RouteID>P19</RouteID></Route><Route><LineDescription>Anacostia-Eckington Line</LineDescription><Name>P6 - ANACOSTIA STA - RHODE ISLAND STA</Name><RouteID>P6</RouteID></Route><Route><LineDescription>Anacostia-Eckington Line</LineDescription><Name>P6 - ECKINGTON - RHODE ISLAND AVE</Name><RouteID>P6v1</RouteID></Route><Route><LineDescription>Anacostia-Eckington Line</LineDescription><Name>P6 - RHODE ISLAND AVE - ECKINGTON</Name><RouteID>P6v2</RouteID></Route><Route><LineDescription>Anacostia-Eckington Line</LineDescription><Name>P6 - ANACOSTIA STA - ARCHIVES</Name><RouteID>P6v3</RouteID></Route><Route><LineDescription>Anacostia-Eckington Line</LineDescription><Name>P6 - ARCHIVES - ANACOSTIA</Name><RouteID>P6v4</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q1 - SILVR SPRNG STA - SHADY GRVE STA</Name><RouteID>Q1</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q2 - SILVR SPRNG STA - SHADY GRVE STA</Name><RouteID>Q2</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q2 - MONT COLLEGE - SILVR SPRNG STA</Name><RouteID>Q2v1</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q2 - SILVR SPRNG STA - MONT COLLEGE</Name><RouteID>Q2v2</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q4 - SILVER SPRNG STA - ROCKVILLE STA</Name><RouteID>Q4</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q4 - WHEATON STA - ROCKVILLE STA</Name><RouteID>Q4v1</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q5 - WHEATON STA - SHADY GRVE STA</Name><RouteID>Q5</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q6 - WHEATON STA - SHADY GRVE STA</Name><RouteID>Q6</RouteID></Route><Route><LineDescription>Viers Mill Road Line</LineDescription><Name>Q6 - ROCKVILLE STA - WHEATON STA</Name><RouteID>Q6v1</RouteID></Route><Route><LineDescription>Riggs Road Line</LineDescription><Name>R1 - FORT TOTTEN STA - ADELPHI</Name><RouteID>R1</RouteID></Route><Route><LineDescription>Kenilworth Avenue Line</LineDescription><Name>R12 - DEANWOOD STA - GREENBELT STA</Name><RouteID>R12</RouteID></Route><Route><LineDescription>Kenilworth Avenue Line</LineDescription><Name>R12 - DEANWOOD STA - GREENBELT STA</Name><RouteID>R12v1</RouteID></Route><Route><LineDescription>Riggs Road Line</LineDescription><Name>R2 - FORT TOTTEN - CALVERTON</Name><RouteID>R2</RouteID></Route><Route><LineDescription>Riggs Road Line</LineDescription><Name>R2 - HIGH POINT HS - FORT TOTTEN</Name><RouteID>R2v1</RouteID></Route><Route><LineDescription>Riggs Road Line</LineDescription><Name>R2 - POWDER MILL+CHERRY HILL - CALVERTON</Name><RouteID>R2v2</RouteID></Route><Route><LineDescription>Queens Chapel Road Line</LineDescription><Name>R4 - BROOKLAND STA- HIGHVIEW</Name><RouteID>R4</RouteID></Route><Route><LineDescription>Richmond Highway Express</LineDescription><Name>REX - FT BELVOIR POST - KING ST STA</Name><RouteID>REX</RouteID></Route><Route><LineDescription>Richmond Highway Express</LineDescription><Name>REX - FT BELVOIR COMM HOSP - KING ST STA</Name><RouteID>REXv1</RouteID></Route><Route><LineDescription>Richmond Highway Express</LineDescription><Name>REX - KING ST STA - FT BELVOIR COMM HOSP</Name><RouteID>REXv2</RouteID></Route><Route><LineDescription>Richmond Highway Express</LineDescription><Name>REX - KING ST STA - WOODLAWN</Name><RouteID>REXv3</RouteID></Route><Route><LineDescription>Richmond Highway Express</LineDescription><Name>REX - WOODLAWN - KING ST STA</Name><RouteID>REXv4</RouteID></Route><Route><LineDescription>16th Street-Potomac Park Line</LineDescription><Name>S1 - NORTHERN DIVISION - POTOMAC PK</Name><RouteID>S1</RouteID></Route><Route><LineDescription>16th Street-Potomac Park Line</LineDescription><Name>S1 - VIRGINIA+E - COLORDO+16TH</Name><RouteID>S1v1</RouteID></Route><Route><LineDescription>16th Street Line</LineDescription><Name>S2 - FED TRIANGLE  - SILVER SPRNG</Name><RouteID>S2</RouteID></Route><Route><LineDescription>16th Street Line</LineDescription><Name>S2 - 16TH &amp; HARVARD - MCPHERSON SQ</Name><RouteID>S2v1</RouteID></Route><Route><LineDescription>Fort Dupont Shuttle Line</LineDescription><Name>S35 - BRANCH + RANDLE CIR - FT DUPONT</Name><RouteID>S35</RouteID></Route><Route><LineDescription>16th Street Line</LineDescription><Name>S4 - FED TRIANGLE - SILVER SPRNG</Name><RouteID>S4</RouteID></Route><Route><LineDescription>16th Street Line</LineDescription><Name>S4 - SILVER SPRING STA - FRANKLIN SQ</Name><RouteID>S4v1</RouteID></Route><Route><LineDescription>Rhode Island Ave-Carver Terrace Line</LineDescription><Name>S41 - CARVER TERRACE - RHODE ISLAND AVE</Name><RouteID>S41</RouteID></Route><Route><LineDescription>Springfield Circulator-Metro Park Shuttle (TAGS)</LineDescription><Name>S80 - FRANCONIA-SPRNGFLD - METRO PARK</Name><RouteID>S80</RouteID></Route><Route><LineDescription>Springfield Circulator-Metro Park Shuttle (TAGS)</LineDescription><Name>S80 - FRANCONIA-SPRINGFLD - HILTON</Name><RouteID>S80v1</RouteID></Route><Route><LineDescription>Springfield Circulator-Metro Park Shuttle (TAGS)</LineDescription><Name>S80 - HILTON - FRANCONIA-SPRNGFLD</Name><RouteID>S80v2</RouteID></Route><Route><LineDescription>16th Street Limited Line</LineDescription><Name>S9 - FRANKLIN SQ - SILVER SPRING STA</Name><RouteID>S9</RouteID></Route><Route><LineDescription>16th Street Limited Line</LineDescription><Name>S9 - FRANKLIN SQ - COLORADO + 16TH</Name><RouteID>S9v1</RouteID></Route><Route><LineDescription>Springfield Circulator-Metro Park Shuttle (TAGS)</LineDescription><Name>S91 - FRANCONIA SPRINGFLD STA SHUTTLE</Name><RouteID>S91</RouteID></Route><Route><LineDescription>Springfield Circulator-Metro Park Shuttle (TAGS)</LineDescription><Name>S91 - FRANCONIA SPRINGFLD STA SHUTTLE</Name><RouteID>S91v1</RouteID></Route><Route><LineDescription>Rhode Island Ave-New Carrollton Line</LineDescription><Name>T14 - RHD ISLND AVE STA-NEW CARRLTN STA</Name><RouteID>T14</RouteID></Route><Route><LineDescription>Rhode Island Ave-New Carrollton Line</LineDescription><Name>T14 - MT RAINIER - NEW CARRLTN STA</Name><RouteID>T14v1</RouteID></Route><Route><LineDescription>Annapolis Road Line</LineDescription><Name>T18 - R I AVE STA - NEW CARROLLTON STA</Name><RouteID>T18</RouteID></Route><Route><LineDescription>Annapolis Road Line</LineDescription><Name>T18 - BLADENSBURG HS - NEW CARROLLTON STA</Name><RouteID>T18v1</RouteID></Route><Route><LineDescription>River Road Line</LineDescription><Name>T2 - FRIENDSHIP HTS - ROCKVILLE STA</Name><RouteID>T2</RouteID></Route><Route><LineDescription>Sheriff Road-River Terrace Line</LineDescription><Name>U4 - MINNESOTA AVE - SHERIFF RD</Name><RouteID>U4</RouteID></Route><Route><LineDescription>Sheriff Road-River Terrace Line</LineDescription><Name>U4 - SHERIFF RD - MINNESOTA STA</Name><RouteID>U4v1</RouteID></Route><Route><LineDescription>Sheriff Road-River Terrace Line</LineDescription><Name>U4 - RIVER TERRACE - MINNESOTA AVE</Name><RouteID>U4v2</RouteID></Route><Route><LineDescription>Marshall Heights Line</LineDescription><Name>U5 - MINNESOTA AVE-MARSHALL HTS</Name><RouteID>U5</RouteID></Route><Route><LineDescription>Marshall Heights Line</LineDescription><Name>U6 - MINNESOTA - LINCOLN HEIGHTS</Name><RouteID>U6</RouteID></Route><Route><LineDescription>Marshall Heights Line</LineDescription><Name>U6 - 37TH + RIDGE - PLUMMER ES</Name><RouteID>U6v1</RouteID></Route><Route><LineDescription>Marshall Heights Line</LineDescription><Name>U6 - LINCOLN HTS - E CAP + 47TH ST NE</Name><RouteID>U6v2</RouteID></Route><Route><LineDescription>Deanwood-Minnesota Ave Line</LineDescription><Name>U7 - RIDGE + ANACOSTIA - DEANWOOD</Name><RouteID>U7</RouteID></Route><Route><LineDescription>Deanwood-Minnesota Ave Line</LineDescription><Name>U7 - MINNESOTA STA - KENILWORTH HAYES</Name><RouteID>U7v1</RouteID></Route><Route><LineDescription>Deanwood-Minnesota Ave Line</LineDescription><Name>U7 - KENILWORTH HAYES - MINNESOTA STA</Name><RouteID>U7v2</RouteID></Route><Route><LineDescription>Deanwood-Minnesota Ave Line</LineDescription><Name>U7 - MINNESOTA STA - DEANWOOD</Name><RouteID>U7v3</RouteID></Route><Route><LineDescription>Deanwood-Minnesota Ave Line</LineDescription><Name>U7 - DEANWOOD - MINNESOTA AVE</Name><RouteID>U7v4</RouteID></Route><Route><LineDescription>Benning Heights-M Street Line</LineDescription><Name>V1 - BUR OF ENGRVNG - BENNING HTS</Name><RouteID>V1</RouteID></Route><Route><LineDescription>District Heights-Suitland Line</LineDescription><Name>V12 - SUITLAND STA - ADDISON RD STA</Name><RouteID>V12</RouteID></Route><Route><LineDescription>District Heights - Seat Pleasant Line</LineDescription><Name>V14 - PENN MAR - DEANWOOD STA</Name><RouteID>V14</RouteID></Route><Route><LineDescription>District Heights - Seat Pleasant Line</LineDescription><Name>V14 - PENN MAR  - ADDISON RD STA</Name><RouteID>V14v1</RouteID></Route><Route><LineDescription>Capitol Heights-Minnesota Avenue Line</LineDescription><Name>V2 - ANACOSTIA - CAPITOL HGTS</Name><RouteID>V2</RouteID></Route><Route><LineDescription>Capitol Heights-Minnesota Avenue Line</LineDescription><Name>V2 - MINNESOTA AVE - ANACOSTIA</Name><RouteID>V2v1</RouteID></Route><Route><LineDescription>Capitol Heights-Minnesota Avenue Line</LineDescription><Name>V4 - 1ST + K ST SE - CAPITOL HGTS</Name><RouteID>V4</RouteID></Route><Route><LineDescription>Capitol Heights-Minnesota Avenue Line</LineDescription><Name>V4 - MINN AVE STA-CAPITOL HGTS</Name><RouteID>V4v1</RouteID></Route><Route><LineDescription>Benning Heights-Alabama Ave Line</LineDescription><Name>V7 - CONGRESS HGTS - MINN STA</Name><RouteID>V7</RouteID></Route><Route><LineDescription>Benning Heights-Alabama Ave Line</LineDescription><Name>V8 - BENNING HGTS - MINN STA</Name><RouteID>V8</RouteID></Route><Route><LineDescription>Shipley Terrace - Fort Drum Line</LineDescription><Name>W1 - FORT DRUM  - SOUTHERN AVE STA</Name><RouteID>W1</RouteID></Route><Route><LineDescription>Bock Road Line</LineDescription><Name>W14 - FT WASHINGTON-SOUTHERN AVE STA</Name><RouteID>W14</RouteID></Route><Route><LineDescription>Bock Road Line</LineDescription><Name>W14 - SOUTHERN AVE - FRIENDLY</Name><RouteID>W14v1</RouteID></Route><Route><LineDescription>Bock Road Line</LineDescription><Name>W14 - FRIENDLY - SOUTHERN AVE</Name><RouteID>W14v2</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - MALCOM X+OAKWD - UNITED MEDICAL CTR</Name><RouteID>W2</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - NYLDR+GOOD HOP- HOWARD+ANACOSTIA</Name><RouteID>W2v1</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - NAYLOR+GOODHOPE - MALCM X+OAKWOOD</Name><RouteID>W2v2</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - HOWRD+ANACOSTIA-NAYLOR+GOODHOPE</Name><RouteID>W2v3</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - ANACOSTIA STA - UNITED MEDICAL CTR</Name><RouteID>W2v4</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - ANACOSTIA STA - SOUTHERN AVE STA</Name><RouteID>W2v5</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - MELLN+M L KNG - UNITED MEDICAL CTR</Name><RouteID>W2v6</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W2 - SOUTHERN AVE-WASHINGTON OVERLOOK</Name><RouteID>W2v7</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W3 - MALCOM X+OAKWD - UNITED MEDICAL CTR</Name><RouteID>W3</RouteID></Route><Route><LineDescription>United Medical Center-Anacostia Line</LineDescription><Name>W3 - MELLN+M L KNG - UNITED MEDICAL CTR</Name><RouteID>W3v1</RouteID></Route><Route><LineDescription>Deanwood-Alabama Ave Line</LineDescription><Name>W4 - ANACOSTIA STA - DEANWOOD STA</Name><RouteID>W4</RouteID></Route><Route><LineDescription>Deanwood-Alabama Ave Line</LineDescription><Name>W4 - MALCOLM X &amp; PORTLAND - DEANWOOD</Name><RouteID>W4v1</RouteID></Route><Route><LineDescription>Deanwood-Alabama Ave Line</LineDescription><Name>W4 - DEANWOOD STA - MALCLM X+ ML KING</Name><RouteID>W4v2</RouteID></Route><Route><LineDescription>Mt Pleasant-Tenleytown Line</LineDescription><Name>W45 - TENLEYTOWN STA - 16TH + SHEPHERD</Name><RouteID>W45</RouteID></Route><Route><LineDescription>Mt Pleasant-Tenleytown Line</LineDescription><Name>W47 - TENLEYTOWN STA - COLUMBIA HTS STA</Name><RouteID>W47</RouteID></Route><Route><LineDescription>Anacostia-Fort Drum Line</LineDescription><Name>W5 - DC VILLAGE-USCG (VIA ANAC)</Name><RouteID>W5</RouteID></Route><Route><LineDescription>Garfield-Anacostia Loop Line</LineDescription><Name>W6 - ANACOSTIA - ANACOSTIA</Name><RouteID>W6</RouteID></Route><Route><LineDescription>Garfield-Anacostia Loop Line</LineDescription><Name>W6 - NAYLOR+GOODHOPE - ANACOSTIA</Name><RouteID>W6v1</RouteID></Route><Route><LineDescription>Garfield-Anacostia Loop Line</LineDescription><Name>W8 - ANACOSTIA - ANACOSTIA</Name><RouteID>W8</RouteID></Route><Route><LineDescription>Garfield-Anacostia Loop Line</LineDescription><Name>W8 - ANACOSTIA-NAYLOR+GOODHOPE</Name><RouteID>W8v1</RouteID></Route><Route><LineDescription>Garfield-Anacostia Loop Line</LineDescription><Name>W8 - NAYLOR+GOODHOPE - ANACOSTIA</Name><RouteID>W8v2</RouteID></Route><Route><LineDescription>Benning Road Line</LineDescription><Name>X1 - FOGGY BOTTOM+GWU- MINNESOTA STA</Name><RouteID>X1</RouteID></Route><Route><LineDescription>Benning Road-H Street Line</LineDescription><Name>X2 - LAFAYETTE SQ - MINNESOTA STA</Name><RouteID>X2</RouteID></Route><Route><LineDescription>Benning Road-H Street Line</LineDescription><Name>X2 - PHELPS HS - MINNESOTA STA</Name><RouteID>X2v1</RouteID></Route><Route><LineDescription>Benning Road-H Street Line</LineDescription><Name>X2 - FRIENDSHIP EDISON PCS-MINNESOTA STA</Name><RouteID>X2v2</RouteID></Route><Route><LineDescription>Benning Road-H Street Line</LineDescription><Name>X2 - PHELPS HS - LAFAYTTE SQ</Name><RouteID>X2v3</RouteID></Route><Route><LineDescription>Benning Road Line</LineDescription><Name>X3 - DUKE ELLINGTON BR - MINNESOTA STN</Name><RouteID>X3</RouteID></Route><Route><LineDescription>Benning Road Line</LineDescription><Name>X3 - KIPP DC - MINNESOTA AVE STN</Name><RouteID>X3v1</RouteID></Route><Route><LineDescription>Maryland Ave Line</LineDescription><Name>X8 - UNION STATION - CARVER TERR</Name><RouteID>X8</RouteID></Route><Route><LineDescription>Benning Road-H Street Limited Line</LineDescription><Name>X9 - NY AVE &amp; 12TH NW - CAPITOL HTS STA</Name><RouteID>X9</RouteID></Route><Route><LineDescription>Benning Road-H Street Limited Line</LineDescription><Name>X9 - NY AVE &amp; 12TH ST NW - MINNESOTA AVE</Name><RouteID>X9v1</RouteID></Route><Route><LineDescription>Benning Road-H Street Limited Line</LineDescription><Name>X9 - MINNESOTA AVE ST - NY AVE &amp; 12TH ST</Name><RouteID>X9v2</RouteID></Route><Route><LineDescription>Georgia Ave-Maryland Line</LineDescription><Name>Y2 - SILVER SPRING STA - MONTG MED CTR</Name><RouteID>Y2</RouteID></Route><Route><LineDescription>Georgia Ave-Maryland Line</LineDescription><Name>Y7 - SILVER SPRING STA - ICC P&amp;R</Name><RouteID>Y7</RouteID></Route><Route><LineDescription>Georgia Ave-Maryland Line</LineDescription><Name>Y8 - SILVER SPR STA - MONTGOMERY MED CTR</Name><RouteID>Y8</RouteID></Route><Route><LineDescription>Greencastle-Briggs Chaney Express Line</LineDescription><Name>Z11 - SILVR SPRING - BURTONSVILLE P&amp;R</Name><RouteID>Z11</RouteID></Route><Route><LineDescription>Greencastle-Briggs Chaney Express Line</LineDescription><Name>Z11 - GREENCASTLE - SILVR SPRING</Name><RouteID>Z11v1</RouteID></Route><Route><LineDescription>Colesville - Ashton Line</LineDescription><Name>Z2 - SILVR SPRING - OLNEY (NO BLAKE HS)</Name><RouteID>Z2</RouteID></Route><Route><LineDescription>Colesville - Ashton Line</LineDescription><Name>Z2 - SILVER SPRING - BONIFANT &amp; NH</Name><RouteID>Z2v1</RouteID></Route><Route><LineDescription>Colesville - Ashton Line</LineDescription><Name>Z2 - SILVER SPRING STA-OLNEY (BLAKE HS)</Name><RouteID>Z2v2</RouteID></Route><Route><LineDescription>Colesville - Ashton Line</LineDescription><Name>Z2 - COLESVILLE-SILVER SPRING</Name><RouteID>Z2v3</RouteID></Route><Route><LineDescription>Calverton-Westfarm Line</LineDescription><Name>Z6 - SILVR SPRNG STA - BURTONSVILLE</Name><RouteID>Z6</RouteID></Route><Route><LineDescription>Calverton-Westfarm Line</LineDescription><Name>Z6 - CASTLE BLVD - SILVER SPRING STA</Name><RouteID>Z6v1</RouteID></Route><Route><LineDescription>Calverton-Westfarm Line</LineDescription><Name>Z6 - SILVER SPRING -CASTLE BLVD</Name><RouteID>Z6v2</RouteID></Route><Route><LineDescription>Laurel-Burtonsville Express Line</LineDescription><Name>Z7 - SLVER SPRNG STA-S LAUREL P&amp;R (4COR)</Name><RouteID>Z7</RouteID></Route><Route><LineDescription>Laurel-Burtonsville Express Line</LineDescription><Name>Z7 - S LAUREL P&amp;R-SILVR SPR (NO 4COR)</Name><RouteID>Z7v1</RouteID></Route><Route><LineDescription>Fairland Line</LineDescription><Name>Z8 - SILVER SRING STA - BRIGGS CHANEY</Name><RouteID>Z8</RouteID></Route><Route><LineDescription>Fairland Line</LineDescription><Name>Z8 - WHITE OAK - SILVER SPRING STA</Name><RouteID>Z8v1</RouteID></Route><Route><LineDescription>Fairland Line</LineDescription><Name>Z8 - SILVR SPRNG - CSTLE BLVD VERIZON</Name><RouteID>Z8v2</RouteID></Route><Route><LineDescription>Fairland Line</LineDescription><Name>Z8 - SILVER SPRING STA - CASTLE BLVD</Name><RouteID>Z8v3</RouteID></Route><Route><LineDescription>Fairland Line</LineDescription><Name>Z8 - SILVER SPRING - GRNCSTLE (VERIZON)</Name><RouteID>Z8v4</RouteID></Route><Route><LineDescription>Fairland Line</LineDescription><Name>Z8 - SILVER SPRING STA - GREENCASTLE</Name><RouteID>Z8v5</RouteID></Route><Route><LineDescription>Fairland Line</LineDescription><Name>Z8 - SILVER SPRING STA - WHITE OAK</Name><RouteID>Z8v6</RouteID></Route></Routes></RoutesResp>`,
			unmarshalledResponse: &GetRoutesResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "RoutesResp",
				},
				Routes: []Route{
					{
						RouteID:         "10A",
						Name:            "10A - HUNTINGTON STA - PENTAGON",
						LineDescription: "Alexandria-Pentagon Line",
					},
					{
						RouteID:         "10B",
						Name:            "10B - HUNTING POINT - BALLSTON STA",
						LineDescription: "Hunting Point-Ballston Line",
					},
					{
						RouteID:         "10E",
						Name:            "10E - HUNTING POINT - PENTAGON",
						LineDescription: "Alexandria-Pentagon Line",
					},
					{
						RouteID:         "10N",
						Name:            "10N - NATIONAL AIRPORT - PENTAGON",
						LineDescription: "Alexandria-Pentagon Line",
					},
					{
						RouteID:         "11Y",
						Name:            "11Y - MOUNT VERNON - POTOMAC PARK",
						LineDescription: "Mt Vernon Express Line",
					},
					{
						RouteID:         "11Yv1",
						Name:            "11Y - POTOMAC PARK - HUNTING POINT",
						LineDescription: "Mt Vernon Express Line",
					},
					{
						RouteID:         "11Yv2",
						Name:            "11Y - HUNTING POINT - POTOMAC PARK",
						LineDescription: "Mt Vernon Express Line",
					},
					{
						RouteID:         "15K",
						Name:            "15K - E FALLS CH STA-MCLEAN LANGLEY",
						LineDescription: "Chain Bridge Road Line",
					},
					{
						RouteID:         "15Kv1",
						Name:            "15K - CIA - EAST FALLS CHURCH STA",
						LineDescription: "Chain Bridge Road Line",
					},
					{
						RouteID:         "16A",
						Name:            "16A - PATRIOT+AMERICANA - PENTAGON",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16C",
						Name:            "16C - CULMORE - PENTAGON",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16Cv1",
						Name:            "16C - CULMORE - FEDERAL TRIANGLE",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16E",
						Name:            "16E - GLN CARLYN+VISTA - FRANKLIN SQ",
						LineDescription: "Columbia Pike Line",
					},
					{
						RouteID:         "16G",
						Name:            "16G - DINWIDDIE+COLUMBIA - PENTAGON CITY",
						LineDescription: "Columbia Pike-Pentagon City Line",
					},
					{
						RouteID:         "16Gv1",
						Name:            "16G - COL PIKE+CARLIN SPR - DINWD+COL PK",
						LineDescription: "Columbia Pike-Pentagon City Line",
					},
					{
						RouteID:         "16H",
						Name:            "16H - SKYLINE CITY - PENTAGON CITY STA",
						LineDescription: "Columbia Pike-Pentagon City Line",
					},
					{
						RouteID:         "16L",
						Name:            "16L - ANNANDALE - PENTAGON HOV",
						LineDescription: "Annandale-Skyline City-Pentagon Line",
					},
					{
						RouteID:         "16Y",
						Name:            "16Y - FOUR MILE+COLUMBIA - MCPHERSON SQ",
						LineDescription: "Columbia Pike-Farragut Square Line",
					},
					{
						RouteID:         "16Yv1",
						Name:            "16Y - COLUMBIA PIKE - MCPHERSON SQ",
						LineDescription: "Columbia Pike-Farragut Square Line",
					},
					{
						RouteID:         "17B",
						Name:            "17B - BURKE CENTRE - PENTAGON HOV",
						LineDescription: "Kings Park-North Springfield Line",
					},
					{
						RouteID:         "17G",
						Name:            "17G - G MASON UNIV - PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17H",
						Name:            "17H - TWNBRK RD+TWNBRK RN- PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17K",
						Name:            "17K - TWNBRK RD+TWNBRK RN- PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17L",
						Name:            "17L - TWNBRK RD+TWNBRK RUN-PENTAGON",
						LineDescription: "Kings Park Express Line",
					},
					{
						RouteID:         "17M",
						Name:            "17M - EDSALL+CANARD - PENTAGON",
						LineDescription: "Kings Park-North Springfield Line",
					},
					{
						RouteID:         "18G",
						Name:            "18G - ROLLING VALLEY - PENTAGON",
						LineDescription: "Orange Hunt Line",
					},
					{
						RouteID:         "18H",
						Name:            "18H - HUNTSMAN+CORK CTY - PENTAGON",
						LineDescription: "Orange Hunt Line",
					},
					{
						RouteID:         "18J",
						Name:            "18J - ROLLING VALLEY - PENTAGON",
						LineDescription: "Orange Hunt Line",
					},
					{
						RouteID:         "18P",
						Name:            "18P - BURKE CENTRE - PENTAGON",
						LineDescription: "Burke Centre Line",
					},
					{
						RouteID:         "18Pv1",
						Name:            "18P - PENTAGON - ROLLING VALLEY MALL",
						LineDescription: "Burke Centre Line",
					},
					{
						RouteID:         "18Pv2",
						Name:            "18P - ROLLING VALLEY P+R - PENTAGON",
						LineDescription: "Burke Centre Line",
					},
					{
						RouteID:         "1A",
						Name:            "1A - VIENNA-BALLSTON (7 CORNERS)",
						LineDescription: "Wilson Blvd-Vienna Line",
					},
					{
						RouteID:         "1B",
						Name:            "1B - DUNN LORING  - BALLSTON",
						LineDescription: "Wilson Blvd-Vienna Line",
					},
					{
						RouteID:         "1C",
						Name:            "1C - WEST OX DIV-DUNN LORING (VIA MALL)",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv1",
						Name:            "1C - FAIRFAX CO GOV CTR - DUNN LORING",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv2",
						Name:            "1C - WEST OX DIV - DUNN LORING (NO MALL)",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv3",
						Name:            "1C - FAIR OAKS MALL - DUNN LORING",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "1Cv4",
						Name:            "1C - DUNN LORING - FAIR OAKS MALL",
						LineDescription: "Fair Oaks-Fairfax Blvd Line",
					},
					{
						RouteID:         "21A",
						Name:            "21A - S REYNOLDS+EOS 21 CONDOS - PENTAGON",
						LineDescription: "Landmark-Bren Mar Park-Pentagon Line",
					},
					{
						RouteID:         "21D",
						Name:            "21D - LANDMARK MEWS -PENTAGON",
						LineDescription: "Landmark-Bren Mar Park-Pentagon Line",
					},
					{
						RouteID:         "22A",
						Name:            "22A - BALLSTON STA - PENTAGON",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "22Av1",
						Name:            "22A - SHIRLINGTON - BALLSTON  STA",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "22C",
						Name:            "22C - BALLSTON STA - PENTAGON",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "22F",
						Name:            "22F - NVCC - PENTAGON VIA HOV",
						LineDescription: "Barcroft-South Fairlington Line",
					},
					{
						RouteID:         "23A",
						Name:            "23A - TYSONS CORNER CTR - CRYSTAL CTY",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "23B",
						Name:            "23B - BALLSTON STA - CRYSTAL CTY",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "23Bv1",
						Name:            "23B - LINDEN RESOURCES - BALLSTON STATION",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "23T",
						Name:            "23T - TYSONS CORNER CTR - SHIRLINGTON",
						LineDescription: "McLean-Crystal City Line",
					},
					{
						RouteID:         "25B",
						Name:            "25B - VAN DORN - BALLSTON",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "25Bv1",
						Name:            "25B - SOUTHERN TOWERS - BALLSTON",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "25Bv2",
						Name:            "25B - VAN DORN - BALLSTON/NO LDMRK CTR",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "25Bv3",
						Name:            "25B - BALLSTON - SOUTHERN TOWERS",
						LineDescription: "Landmark-Ballston Line",
					},
					{
						RouteID:         "26A",
						Name:            "26A - NVCC ANNANDALE - E FALLS CHURCH STA",
						LineDescription: "Annandale-East Falls Church Line",
					},
					{
						RouteID:         "28A",
						Name:            "28A - TYSONS CORNER STA-KING ST STA",
						LineDescription: "Leesburg Pike Line",
					},
					{
						RouteID:         "28Av1",
						Name:            "28A - SOUTHERN TOWERS-TYSONS CORNER STA",
						LineDescription: "Leesburg Pike Line",
					},
					{
						RouteID:         "28F",
						Name:            "28F - BLDG 5113 G MASON DR - PENTAGON",
						LineDescription: "Skyline City Line",
					},
					{
						RouteID:         "28G",
						Name:            "28G - BLDG 5113 G MASON DR - PENTAGON",
						LineDescription: "Skyline City Line",
					},
					{
						RouteID:         "29C",
						Name:            "29C - NVCC ANNANDALE - PENTAGON",
						LineDescription: "Annandale Line",
					},
					{
						RouteID:         "29G",
						Name:            "29G - AMERICANA+HERITAGE - PENTAGON",
						LineDescription: "Annandale Line",
					},
					{
						RouteID:         "29K",
						Name:            "29K - GMU - KING ST STA",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29Kv1",
						Name:            "29K - GMU - KING ST/NO LDMRK",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29N",
						Name:            "29N - VIENNA STA - KING ST (VIA MALL)",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29Nv1",
						Name:            "29N - VIENNA STA-KING ST (NO MALL)",
						LineDescription: "Alexandria-Fairfax Line",
					},
					{
						RouteID:         "29W",
						Name:            "29W - NVCC ANNANDALE - PENTAGON",
						LineDescription: "Braeburn Drive-Pentagon Express Line",
					},
					{
						RouteID:         "2A",
						Name:            "2A - DUNN LORING STA - BALLSTON STA",
						LineDescription: "Washington Blvd.-Dunn Loring Line",
					},
					{
						RouteID:         "2B",
						Name:            "2B - WEST OX RD DIV-DUNN LORING STATION",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "2Bv1",
						Name:            "2B - WEST OX RD-DUNN LORING STA(NO MALL)",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "2Bv2",
						Name:            "2B - FAIR OAKS MALL-DUNN LORING STATION",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "2Bv3",
						Name:            "2B - DUNN LORING STA - FAIR OAKS MALL",
						LineDescription: "Fair Oaks-Jermantown Road Line",
					},
					{
						RouteID:         "30N",
						Name:            "30N - FRIENDSHIP HGTS- NAYLOR RD STA",
						LineDescription: "Friendship Heights-Southeast Line",
					},
					{
						RouteID:         "30S",
						Name:            "30S - FRIENDSHIP HGTS- SOUTHRN AVE STA",
						LineDescription: "Friendship Heights-Southeast Line",
					},
					{
						RouteID:         "31",
						Name:            "31 - POTOMAC PARK-FRIENDSHIP HGTS",
						LineDescription: "Wisconsin Avenue Line",
					},
					{
						RouteID:         "32",
						Name:            "32 - VIRGINIA AVE+E ST- SOUTHRN AVE",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "32v1",
						Name:            "32 - PENN AVE + 8TH ST - SOUTHRN AVE",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "33",
						Name:            "33 - 10TH ST+PA AV NW - FRIENDSHIP HGTS",
						LineDescription: "Wisconsin Avenue Line",
					},
					{
						RouteID:         "34",
						Name:            "34 - 10TH ST + PA AVE- NAYLOR RD STA",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "36",
						Name:            "36 - VIRGINIA AVE+E ST - NAYLOR RD STA",
						LineDescription: "Pennsylvania Avenue Line",
					},
					{
						RouteID:         "37",
						Name:            "37 - 10TH ST+PA AV NW - FRIENDSHIP HGTS",
						LineDescription: "Wisconsin Avenue Limited Line",
					},
					{
						RouteID:         "38B",
						Name:            "38B - BALLSTON - FARRAGUT",
						LineDescription: "Ballston-Farragut Square Line",
					},
					{
						RouteID:         "38Bv1",
						Name:            "38B - WASH & QUINCY - FARRAGUT",
						LineDescription: "Ballston-Farragut Square Line",
					},
					{
						RouteID:         "38Bv2",
						Name:            "38B - WASHINGTON-LEE HS - FARRAGUT SQ",
						LineDescription: "Ballston-Farragut Square Line",
					},
					{
						RouteID:         "39",
						Name:            "39 - VIRGINIA AVE+21ST NW- NAYLOR RD STA",
						LineDescription: "Pennsylvania Avenue Limited Line",
					},
					{
						RouteID:         "3A",
						Name:            "3A - ANNANDALE - E FALLS CHURCH",
						LineDescription: "Annandale Road Line",
					},
					{
						RouteID:         "3Av1",
						Name:            "3A - ANNANDALE - 7 CORNERS",
						LineDescription: "Annandale Road Line",
					},
					{
						RouteID:         "3T",
						Name:            "3T - MCLEAN STATION - E FALLS CH STA",
						LineDescription: "Pimmit Hills Line",
					},
					{
						RouteID:         "3Tv1",
						Name:            "3T - MCLEAN STATION - W FALLS CHURCH",
						LineDescription: "Pimmit Hills Line",
					},
					{
						RouteID:         "3Y",
						Name:            "3Y - E FALLS CHURCH - MCPHERSON SQ",
						LineDescription: "Lee Highway-Farragut Square Line",
					},
					{
						RouteID:         "42",
						Name:            "42 - 9TH + F ST  - MT PLEASANT",
						LineDescription: "Mount Pleasant Line",
					},
					{
						RouteID:         "43",
						Name:            "43 - I + 13TH NW - MT PLEASANT",
						LineDescription: "Mount Pleasant Line",
					},
					{
						RouteID:         "4A",
						Name:            "4A - SEVEN CORNERS - ROSSLYN",
						LineDescription: "Pershing Dr-Arlington Blvd Line",
					},
					{
						RouteID:         "4B",
						Name:            "4B - SEVEN CORNERS - ROSSLYN",
						LineDescription: "Pershing Dr-Arlington Blvd Line",
					},
					{
						RouteID:         "52",
						Name:            "52 - L ENFNT PLAZA - TAKOMA STATION",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "52v1",
						Name:            "52 - L ENFNT PLAZA - 14TH+COLORADO",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "52v2",
						Name:            "52 - 14TH+COLORADO - L ENFANT PLAZA",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "52v3",
						Name:            "52 - 14TH & U - TAKOMA STATION",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "54",
						Name:            "54 - METRO CENTER - TAKOMA STA",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "54v1",
						Name:            "54 - 14TH+COLORADO - METRO CENTER",
						LineDescription: "14th Street Line",
					},
					{
						RouteID:         "59",
						Name:            "59 - FEDERAL TRIANGLE - TAKOMA STATION",
						LineDescription: "14th Street Limited Line",
					},
					{
						RouteID:         "5A",
						Name:            "5A - DULLES AIRPORT - LENFANT PLAZA",
						LineDescription: "DC-Dulles Line",
					},
					{
						RouteID:         "60",
						Name:            "60 - GEORGIA + PETWORTH - FT TOTTEN",
						LineDescription: "Fort Totten-Petworth Line",
					},
					{
						RouteID:         "62",
						Name:            "62 - GEORGIA+PETWORTH - TAKOMA STATION",
						LineDescription: "Takoma-Petworth Line",
					},
					{
						RouteID:         "62v1",
						Name:            "62 - COOLIDGE HS - GEORGIA + PETWORTH",
						LineDescription: "Takoma-Petworth Line",
					},
					{
						RouteID:         "63",
						Name:            "63 - FED TRIANGLE - TAKOMA STA",
						LineDescription: "Takoma-Petworth Line",
					},
					{
						RouteID:         "64",
						Name:            "64 - FEDERAL TRIANGLE -FORT TOTTEN",
						LineDescription: "Fort Totten-Petworth Line",
					},
					{
						RouteID:         "64v1",
						Name:            "64 - GEORGIA + PETWOTH - FT TOTTEN",
						LineDescription: "Fort Totten-Petworth Line",
					},
					{
						RouteID:         "70",
						Name:            "70 - ARCHIVES - SILVER SPRING",
						LineDescription: "Georgia Avenue-7th Street Line",
					},
					{
						RouteID:         "70v1",
						Name:            "70 - GEORGIA & EUCLID  - ARCHIVES",
						LineDescription: "Georgia Avenue-7th Street Line",
					},
					{
						RouteID:         "74",
						Name:            "74 - NATIONALS PARK - CONVENTION CTR",
						LineDescription: "Convention Center-Southwest Waterfront Line",
					},
					{
						RouteID:         "79",
						Name:            "79 - ARCHIVES - SILVER SPRING STA",
						LineDescription: "Georgia Avenue MetroExtra Line",
					},
					{
						RouteID:         "7A",
						Name:            "7A - LINCOLNIA+QUANTRELL - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Av1",
						Name:            "7A - PENTAGON - SOUTHERN TWRS",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Av2",
						Name:            "7A - SOUTHERN TWRS - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Av3",
						Name:            "7A - LINCOLNIA/QUANTRELL - PENT VIA PENT",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7C",
						Name:            "7C - PARK CENTER - PENTAGON",
						LineDescription: "Park Center-Pentagon Line",
					},
					{
						RouteID:         "7F",
						Name:            "7F - LINCOLNIA+QUANTRELL - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Fv1",
						Name:            "7F - LINC + QUANT - PENT CITY - PENTAGON",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7M",
						Name:            "7M - MARK CENTER - PENTAGON (NON-HOV)",
						LineDescription: "Mark Center-Pentagon Line",
					},
					{
						RouteID:         "7Mv1",
						Name:            "7M - MARK CENTER - PENTAGON (HOV)",
						LineDescription: "Mark Center-Pentagon Line",
					},
					{
						RouteID:         "7P",
						Name:            "7P - PARK CTR - PENTAGON",
						LineDescription: "Park Center-Pentagon Line",
					},
					{
						RouteID:         "7W",
						Name:            "7W - LNCLNA+QUANTRLL- PENTAGON",
						LineDescription: "Lincolnia-Pentagon Line",
					},
					{
						RouteID:         "7Y",
						Name:            "7Y - SOUTHERN TWRS - H+17TH ST",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "7Yv1",
						Name:            "7Y - PENTAGON - H+17TH ST",
						LineDescription: "Lincolnia-North Fairlington Line",
					},
					{
						RouteID:         "80",
						Name:            "80 - KENNEDY CTR   - FORT TOTTEN STA",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "80v1",
						Name:            "80 - MCPHERSON SQ  - BROOKLAND",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "80v2",
						Name:            "80 - MCPHERSON SQ  - FORT TOTTEN STA",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "80v3",
						Name:            "80 - KENNEDY CTR   - BROOKLAND STA",
						LineDescription: "North Capitol Street Line",
					},
					{
						RouteID:         "83",
						Name:            "83 - RHODE ISLAND AVE STA-CHERRY HILL",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v1",
						Name:            "83 - MT RAINIER - RHODE ISLAND",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v2",
						Name:            "83 - RHODE ISLAND - MT RAINIER",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v3",
						Name:            "83 - RHODE ISLAND AVE STA-COLLEGE PARK",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "83v4",
						Name:            "83 - COLLEGE PARK-RHODE ISLAND AVE STA",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "86",
						Name:            "86 - RHODE ISLAND AVE STA- CALVERTON",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "86v1",
						Name:            "86 - RHODE ISLAND AVE STA- COLLEGE PARK",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "86v2",
						Name:            "86 - MT RAINIER   - CALVERTON",
						LineDescription: "College Park Line",
					},
					{
						RouteID:         "87",
						Name:            "87 - NEW CARROLTON -CYPRESS+LAURL LAKES",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v1",
						Name:            "87 - GRNBELT STA -CYPRESS+LAURL LAKES",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v2",
						Name:            "87 - GRNBELT-CYP+LRL LAKES (NO P&R)",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v3",
						Name:            "87 - GRNBELT STA - BALTIMORE+MAIN ST",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v4",
						Name:            "87 - BALTIMORE+MAIN ST - GRNBELT STA",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "87v5",
						Name:            "87 - CYPRESS+LAURL LAKES -GRNBELT STA",
						LineDescription: "Laurel Express Line",
					},
					{
						RouteID:         "89",
						Name:            "89 - GREENBELT STA - 4TH & GREEN HILL",
						LineDescription: "Laurel Line",
					},
					{
						RouteID:         "89v1",
						Name:            "89 - GREENBELT STA - CHERRY LA+4TH ST",
						LineDescription: "Laurel Line",
					},
					{
						RouteID:         "89M",
						Name:            "89M - GREENBELT STA - S LAUREL P+R",
						LineDescription: "Laurel Line",
					},
					{
						RouteID:         "8S",
						Name:            "8S - RADFORD+QUAKER - PENTAGON",
						LineDescription: "Foxchase-Seminary Valley Line",
					},
					{
						RouteID:         "8W",
						Name:            "8W - MARK CENTER - PENTAGON V FOXCHASE",
						LineDescription: "Foxchase-Seminary Valley Line",
					},
					{
						RouteID:         "8Z",
						Name:            "8Z - QUAKER+OSAGE - PENTAGON",
						LineDescription: "Foxchase-Seminary Valley Line",
					},
					{
						RouteID:         "90",
						Name:            "90 - ANACOSTIA - DK ELLNGTN BRDG",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "90v1",
						Name:            "90 - 8TH ST + L ST  - DK ELLNGTN BRDG",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "90v2",
						Name:            "90 - KIPP DC PREP- ANACOSTIA",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "92",
						Name:            "92 - CONGRESS HTS STA - REEVES CTR",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "92v1",
						Name:            "92 - EASTERN MARKET - CONGRESS HGTS STA",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "92v2",
						Name:            "92 - CONGRESS HTS STA- EASTERN MARKET",
						LineDescription: "U Street-Garfield Line",
					},
					{
						RouteID:         "96",
						Name:            "96 - TENLEYTOWN STA - CAPITOL HTS STA",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v1",
						Name:            "96 - ELLINGTON BR - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v2",
						Name:            "96 - CAPITOL HTS  - REEVES CTR",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v3",
						Name:            "96 - CAPITOL HTS - ELLINGTON BRDG",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v4",
						Name:            "96 - REEVES CTR  - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "96v5",
						Name:            "96 - TENLEYTOWN STA - STADIUM ARMORY STA",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "97",
						Name:            "97 - UNION STATION - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "97v1",
						Name:            "97 - EASTERN HS - CAPITOL HTS",
						LineDescription: "East Capitol Street-Cardoza Line",
					},
					{
						RouteID:         "A12",
						Name:            "A12 - ADDISON RD STA - CAPITAL PLAZA",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A12v1",
						Name:            "A12 - BARLWE+MATHEW HENSON - ADDISON RD",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A12v2",
						Name:            "A12 - CAPITOL HTS STA - CAPITAL PLAZA",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A12v3",
						Name:            "A12 - CAPITAL PLAZA - CAPITOL HTS",
						LineDescription: "Martin Luther King Jr Highway Line",
					},
					{
						RouteID:         "A2",
						Name:            "A2 - SOUTHERN AVE - ANACOSTIA (VIA HOSP)",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A2v1",
						Name:            "A2 - ANACOSTIA - MISS+ATLANTIC",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A2v2",
						Name:            "A2 - SOUTHERN AVE STA - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A2v3",
						Name:            "A2 - MISS+ATLANTIC - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A31",
						Name:            "A31 - ANACOSTIA HIGH - MINNESOTA STA",
						LineDescription: "Minnesota Ave-Anacostia Line",
					},
					{
						RouteID:         "A32",
						Name:            "A32 - ANACOSTIA HIGH - SOUTHRN AVE STA",
						LineDescription: "Minnesota Ave-Anacostia Line",
					},
					{
						RouteID:         "A33",
						Name:            "A33 - ANACOSTIA HIGH - ANACOSTIA STA",
						LineDescription: "Minnesota Ave-Anacostia Line",
					},
					{
						RouteID:         "A4",
						Name:            "A4 - DC VILLAGE - ANACOSTIA",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v1",
						Name:            "A4 - USCG-FT DRUM (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v2",
						Name:            "A4 - FT DRUM - ANACOSTIA",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v3",
						Name:            "A4 - ANACOSTIA - FT DRUM",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v4",
						Name:            "A4 - FT DRUM-USCG (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A4v5",
						Name:            "A4 - DC VILL-USCG (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "A6",
						Name:            "A6 - 4501 3RD ST - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A6v1",
						Name:            "A6 - SOUTHERN AVE+S CAPITOL - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A7",
						Name:            "A7 - SOUTHRN+S CAPITOL - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A8",
						Name:            "A8 - 4501 3RD ST - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A8v1",
						Name:            "A8 - SOUTHERN + S CAPITOL - ANACOSTIA",
						LineDescription: "Anacostia-Congress Heights Line",
					},
					{
						RouteID:         "A9",
						Name:            "A9 - LIVINGSTON - MCPHERSON SQUARE",
						LineDescription: "Martin Luther King Jr Ave Limited Line",
					},
					{
						RouteID:         "B2",
						Name:            "B2 - ANACOSTIA STA - MOUNT RAINIER",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v1",
						Name:            "B2 - ANACOSTIA STA - BLDNSGRG+VST NE",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v2",
						Name:            "B2 - POTOMAC AVE - MOUNT RAINIER",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v3",
						Name:            "B2 - BLDNSBRG+26TH - ANACOSTIA STA",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B2v4",
						Name:            "B2 - EASTERN HS - ANACOSTIA STA",
						LineDescription: "Bladensburg Road-Anacostia Line",
					},
					{
						RouteID:         "B21",
						Name:            "B21 - NEW CARROLTON STA - BOWIE STATE",
						LineDescription: "Bowie State University Line",
					},
					{
						RouteID:         "B22",
						Name:            "B22 - NEW CARROLTON STA - BOWIE STATE",
						LineDescription: "Bowie State University Line",
					},
					{
						RouteID:         "B22v1",
						Name:            "B22 - OLD CHAPEL & 197 - NEW CARRLTN STA",
						LineDescription: "Bowie State University Line",
					},
					{
						RouteID:         "B24",
						Name:            "B24 - NEW CARLTN STA-BOWIE P+R(VIA BHC)",
						LineDescription: "Bowie-Belair Line",
					},
					{
						RouteID:         "B24v1",
						Name:            "B24 - NEW CARROLTON STA - BOWIE P+R",
						LineDescription: "Bowie-Belair Line",
					},
					{
						RouteID:         "B27",
						Name:            "B27 - NEW CARROLTON STA - BOWIE STATE",
						LineDescription: "Bowie-New Carrollton Line",
					},
					{
						RouteID:         "B29",
						Name:            "B29 - NEW CARROLTON - CROFTON CC (PM)(PR)",
						LineDescription: "Crofton-New Carrollton Line",
					},
					{
						RouteID:         "B29v1",
						Name:            "B29 - NEW CARROLLTON STA - GATEWAY CTR",
						LineDescription: "Crofton-New Carrollton Line",
					},
					{
						RouteID:         "B29v2",
						Name:            "B29 - NEW CARROLTON - CROFTON CC (NO PR)",
						LineDescription: "Crofton-New Carrollton Line",
					},
					{
						RouteID:         "B30",
						Name:            "B30 - GREENBELT STA - BWI LT RAIL STA",
						LineDescription: "Greenbelt-BWI Thurgood Marshall Airport Express Line",
					},
					{
						RouteID:         "B8",
						Name:            "B8 - RHODE ISLAND AV -PETERSBRG APTS",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "B8v1",
						Name:            "B8 - BLDNSBRG+S DKTA -PETERSBRG APTS",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "B8v2",
						Name:            "B8 - PETERSBRG APTS  -BLDNSBRG+S DKTA",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "B9",
						Name:            "B9 - RHODE ISLAND AVE - COLMAR MANOR",
						LineDescription: "Fort Lincoln Shuttle Line",
					},
					{
						RouteID:         "C11",
						Name:            "C11 - CLINTON P+R - BRANCH AVE STA",
						LineDescription: "Clinton Line",
					},
					{
						RouteID:         "C12",
						Name:            "C12 - NAYLOR RD STA  - BRANCH AVE STA",
						LineDescription: "Hillcrest Heights Line",
					},
					{
						RouteID:         "C13",
						Name:            "C13 - CLINTON P+R - BRANCH AVE STA",
						LineDescription: "Clinton Line",
					},
					{
						RouteID:         "C14",
						Name:            "C14 - NAYLOR RD STA  - BRANCH AVE STA",
						LineDescription: "Hillcrest Heights Line",
					},
					{
						RouteID:         "C2",
						Name:            "C2 - WHEATN STA - GRNBELT STA UMD ALT",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C2v1",
						Name:            "C2 - TAKOMA LANGLEY XROADS - GRNBELT STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C2v2",
						Name:            "C2 - GREENBELT STA - RANDOLPH + PARKLAWN",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C2v3",
						Name:            "C2 - TAKOMA LANGLEY XROADS - WHEATON",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C21",
						Name:            "C21 - ADDISON RD STA - COLLINGTON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C21v1",
						Name:            "C21 - ADDISON RD STA - POINTER RIDGE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C21v2",
						Name:            "C21 - ADDISON RD STA - CAMPUS WAY S",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C22",
						Name:            "C22 - ADDISON RD STA - COLLINGTON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C22v1",
						Name:            "C22 - ADDISON RD STA- POINTER RIDGE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C26",
						Name:            "C26 - LARGO TOWN CTR - WATKNS+CHESTERTON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C26v1",
						Name:            "C26 - LARGO TOWN CTR-WATKINS+CAMBLETON",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C28",
						Name:            "C28 - PT RIDGE -NEW CAROLTN VIA HLTH CTR",
						LineDescription: "Pointer Ridge Line",
					},
					{
						RouteID:         "C28v1",
						Name:            "C28 - PT RIDGE - NEW CARROLLTON STA",
						LineDescription: "Pointer Ridge Line",
					},
					{
						RouteID:         "C29*1",
						Name:            "C29 - POINTER RIDGE - ADDISON RD STA",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C29*2",
						Name:            "C29 - WATKNS+CAMBLETON - ADDISN RD STA",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C29*4",
						Name:            "C29 - ADDISON RD STA - BOWIE STATE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C29/",
						Name:            "C29 - ADDISON RD STA - POINTER RIDGE",
						LineDescription: "Central Avenue Line",
					},
					{
						RouteID:         "C4",
						Name:            "C4 - TWINBROOK STA - PG PLAZA STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C4v1",
						Name:            "C4 - TLTC-TWINBROOK",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C4v2",
						Name:            "C4 - TWINBROOK STA - WHEATON STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C4v3",
						Name:            "C4 - PG PLAZA STA - WHEATON STA",
						LineDescription: "Greenbelt-Twinbrook Line",
					},
					{
						RouteID:         "C8",
						Name:            "C8 - WHITE FLINT - COLLEGE PARK",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "C8v1",
						Name:            "C8 - WHITE FLNT-COLLEGE PK (NO FDA/ARCH)",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "C8v2",
						Name:            "C8 - WHITE FLINT-COLLEGE PARK (NO FDA)",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "C8v3",
						Name:            "C8 - GLENMONT-COLLEGE PK (NO FDA/ARCH)",
						LineDescription: "College Park-White Flint Line",
					},
					{
						RouteID:         "D1",
						Name:            "D1 - GLOVER PARK - FRANKLIN SQUARE",
						LineDescription: "Glover Park-Franklin Square Line",
					},
					{
						RouteID:         "D12",
						Name:            "D12 - SOUTHERN AVE STA - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D12v1",
						Name:            "D12 - SUITLAND STA - SOUTHERN AVE STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D12v2",
						Name:            "D12 - ST BARNABAS RD   - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D13",
						Name:            "D13 - SOUTHERN AVE STA - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D13v1",
						Name:            "D13 - SOUTHRN STA - ALLENTWN+OLD BRNCH",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D14",
						Name:            "D14 - SOUTHERN AVE STA - SUITLAND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D14v1",
						Name:            "D14 - ALLENTWN+OLD BRNCH - SUITLND STA",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D14v2",
						Name:            "D14 - SUITLND STA - ALLENTWN+OLD BRNCH",
						LineDescription: "Oxon Hill-Suitland Line",
					},
					{
						RouteID:         "D2",
						Name:            "D2 - GLOVER PARK   - CONNETICUT +Q ST",
						LineDescription: "Glover Park-Dupont Circle Line",
					},
					{
						RouteID:         "D2v1",
						Name:            "D2 - DUPONT CIR - GLOVER PARK",
						LineDescription: "Glover Park-Dupont Circle Line",
					},
					{
						RouteID:         "D31",
						Name:            "D31 - TENLEYTOWN STA - 16TH + EASTERN",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D32",
						Name:            "D32 - TENLEYTOWN STA - COLUMBIA HTS STA",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D33",
						Name:            "D33 - TENLEYTOWN STA - 16TH + SHEPHERD",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D34",
						Name:            "D34 - TENLEYTOWN STA - 14TH + COLORADO",
						LineDescription: "16th Street-Tenleytown Line",
					},
					{
						RouteID:         "D4",
						Name:            "D4 - DUPONT CIRCLE - IVY CITY",
						LineDescription: "Ivy City-Franklin Square Line",
					},
					{
						RouteID:         "D4v1",
						Name:            "D4 - FRANKLIN SQUARE - IVY CITY",
						LineDescription: "Ivy City-Franklin Square Line",
					},
					{
						RouteID:         "D4v2",
						Name:            "D4 - IVY CITY - FRANKLIN SQUARE",
						LineDescription: "Ivy City-Franklin Square Line",
					},
					{
						RouteID:         "D5",
						Name:            "D5 - MASS LTTL FLWR CHRCH- FARRGT SQR",
						LineDescription: "MacArthur Blvd-Georgetown Line",
					},
					{
						RouteID:         "D6",
						Name:            "D6 - SIBLEY HOSP - STADIUM ARMRY",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D6v1",
						Name:            "D6 - STADIUM ARMRY STA - FARRAGUT SQUARE",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D6v2",
						Name:            "D6 - SIBLEY HOSP - FARRAGUT SQ",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D6v3",
						Name:            "D6 - FARRAGUT SQUARE - STADIUM ARMRY",
						LineDescription: "Sibley Hospital–Stadium-Armory Line",
					},
					{
						RouteID:         "D8",
						Name:            "D8 - UNION STATION - VA MED CTR",
						LineDescription: "Hospital Center Line",
					},
					{
						RouteID:         "D8v1",
						Name:            "D8 - RHODE ISLAND STA - UNION STA",
						LineDescription: "Hospital Center Line",
					},
					{
						RouteID:         "E2",
						Name:            "E2 - IVY CITY - FT TOTTEN",
						LineDescription: "Ivy City-Fort Totten Line",
					},
					{
						RouteID:         "E4",
						Name:            "E4 - FRIENDSHP HTS - RIGGS PK",
						LineDescription: "Military Road-Crosstown Line",
					},
					{
						RouteID:         "E4v1",
						Name:            "E4 - FRIENDSHP HTS - FT TOTTEN",
						LineDescription: "Military Road-Crosstown Line",
					},
					{
						RouteID:         "E4v2",
						Name:            "E4 - FT TOTTEN - FRIENDSHP HTS",
						LineDescription: "Military Road-Crosstown Line",
					},
					{
						RouteID:         "E6",
						Name:            "E6 - FRIENDSHP HTS  -KNOLLWOOD",
						LineDescription: "Chevy Chase Line",
					},
					{
						RouteID:         "F1",
						Name:            "F1 - CHEVERLY STA - TAKOMA STA",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F12",
						Name:            "F12 - CHEVERLY STA - NEW CARROLTON STA",
						LineDescription: "Ardwick Industrial Park Shuttle Line",
					},
					{
						RouteID:         "F12v1",
						Name:            "F12 - CHEVERLY STA - LANDOVER STA",
						LineDescription: "Ardwick Industrial Park Shuttle Line",
					},
					{
						RouteID:         "F13",
						Name:            "F13 - CHEVERLY STA-WASHINGTON BUS PARK",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F13v1",
						Name:            "F13 - CHEVERLY STA - WASHINGTON BUS PARK",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F13v2",
						Name:            "F13 - NEW CARROLTON - WASHINGTON BUS PARK",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F13v3",
						Name:            "F13 - WHITFIELD + VOLTA - CHEVERLY STA",
						LineDescription: "Cheverly-Washington Business Park Line",
					},
					{
						RouteID:         "F14",
						Name:            "F14 - NAYLOR RD STA -NEW CARROLTON STA",
						LineDescription: "Sheriff Road-Capitol Heights Line",
					},
					{
						RouteID:         "F14v1",
						Name:            "F14 - BRADBURY HGTS -NEW CARROLTON STA",
						LineDescription: "Sheriff Road-Capitol Heights Line",
					},
					{
						RouteID:         "F2",
						Name:            "F2 - CHEVERLY STA - TAKOMA STA",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F2v1",
						Name:            "F2 - CHEVRLY STA-QUEENS CHAPL+CARSN CIR",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F2v2",
						Name:            "F2 - 34TH + OTIS - TAKOMA STA",
						LineDescription: "Chillum Road Line",
					},
					{
						RouteID:         "F4",
						Name:            "F4 - SILVR SPRING STA - NEW CARRLLTON",
						LineDescription: "New Carrollton-Silver Spring Line",
					},
					{
						RouteID:         "F4v1",
						Name:            "F4 - PG PLAZA STA -NEW CARROLTON STA",
						LineDescription: "New Carrollton-Silver Spring Line",
					},
					{
						RouteID:         "F4v2",
						Name:            "F4 - NEW CARRLLTON - PG PLAZA STA",
						LineDescription: "New Carrollton-Silver Spring Line",
					},
					{
						RouteID:         "F6",
						Name:            "F6 - FT TOTTEN - NEW CARROLLTN",
						LineDescription: "New Carrollton-Fort Totten Line",
					},
					{
						RouteID:         "F6v1",
						Name:            "F6 - PG PLAZA - NEW CARROLLTON STA",
						LineDescription: "New Carrollton-Fort Totten Line",
					},
					{
						RouteID:         "F6v2",
						Name:            "F6 - NEW CARROLLTON - PG PLAZA STA",
						LineDescription: "New Carrollton-Fort Totten Line",
					},
					{
						RouteID:         "F8",
						Name:            "F8 - CHEVERLY STA - TLTC",
						LineDescription: "Langley Park-Cheverly Line",
					},
					{
						RouteID:         "G12",
						Name:            "G12 - GREENBELT STA - NEW CARROLLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G12v1",
						Name:            "G12 - GREENBELT STA - ROOSEVELT CENTER",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G12v2",
						Name:            "G12 - ROOSEVELT CTR - NEW CARROLLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G14",
						Name:            "G14 - GREENBELT STA -NEW CARROLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G14v1",
						Name:            "G14 - ROOSEVELT CTR -NEW CARROLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G14v2",
						Name:            "G14 - GREENBELT STA - NEW CARROLTON STA",
						LineDescription: "Greenbelt-New Carrollton Line",
					},
					{
						RouteID:         "G2",
						Name:            "G2 - GEORGETOWN UNIV - HOWARD UNIV",
						LineDescription: "P Street-LeDroit Park Line",
					},
					{
						RouteID:         "G2v1",
						Name:            "G2 - GEORGETOWN UNIV - P+14TH ST",
						LineDescription: "P Street-LeDroit Park Line",
					},
					{
						RouteID:         "G8",
						Name:            "G8 - FARRAGUT SQUARE - AVONDALE",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G8v1",
						Name:            "G8 - RHODE ISLAND AVE STA-FARRAGUT SQ",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G8v2",
						Name:            "G8 - FARRAGUT SQ - RHODE ISLAND STA",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G8v3",
						Name:            "G8 - BROOKLAND STA - AVONDALE",
						LineDescription: "Rhode Island Avenue Line",
					},
					{
						RouteID:         "G9",
						Name:            "G9 - FRANKLIN SQ - RHODE ISLD + EAST PM",
						LineDescription: "Rhode Island Avenue Limited Line",
					},
					{
						RouteID:         "G9v1",
						Name:            "G9 - FRANKLIN SQ - RHODE ISLD + EASTERN",
						LineDescription: "Rhode Island Avenue Limited Line",
					},
					{
						RouteID:         "H1",
						Name:            "H1 - C + 17TH ST - BROOKLAND CUA STA",
						LineDescription: "Brookland-Potomac Park Line",
					},
					{
						RouteID:         "H11",
						Name:            "H11 - HEATHER HILL - NAYLOR RD STA",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H12",
						Name:            "H12 - HEATHER HILL-NAYLOR RD (MACYS)",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H12v1",
						Name:            "H12 - HEATHER HILL-NAYLOR RD STA (P&R)",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H13",
						Name:            "H13 - HEATHER HILL - NAYLOR RD STA",
						LineDescription: "Marlow Heights-Temple Hills Line",
					},
					{
						RouteID:         "H2",
						Name:            "H2 - TENLEYTOWN AU-ST - BROOKLND CUA STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H3",
						Name:            "H3 - TENLEYTOWN STA -BROOKLND CUA STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H4",
						Name:            "H4 - TENLEYTWN AU-STA -BROOKLAND CUA STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H4v1",
						Name:            "H4 - COLUMBIA RD+14TH -TENLEYTOWN-AU STA",
						LineDescription: "Crosstown Line",
					},
					{
						RouteID:         "H6",
						Name:            "H6 - BROOKLAND CUA STA - FORT LINCOLN",
						LineDescription: "Brookland-Fort Lincoln Line",
					},
					{
						RouteID:         "H6v1",
						Name:            "H6 - FORT LINCOLN - FORT LINCOLN",
						LineDescription: "Brookland-Fort Lincoln Line",
					},
					{
						RouteID:         "H8",
						Name:            "H8 - MT PLEASANT+17TH - RHODE ISLAND",
						LineDescription: "Park Road-Brookland Line",
					},
					{
						RouteID:         "H9",
						Name:            "H9 - RHODE ISLAND - FORT DR + 1ST",
						LineDescription: "Park Road-Brookland Line",
					},
					{
						RouteID:         "J1",
						Name:            "J1 - MEDCAL CTR STA - SILVR SPRNG STA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J1v1",
						Name:            "J1 - SILVR SPRNG STA - MONT MALL/BAT LN",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J12",
						Name:            "J12 - ADDISON RD STA - FORESTVILLE",
						LineDescription: "Marlboro Pike Line",
					},
					{
						RouteID:         "J12v1",
						Name:            "J12 - ADDISON RD-FORESTVILLE VIA PRES PKY",
						LineDescription: "Marlboro Pike Line",
					},
					{
						RouteID:         "J2",
						Name:            "J2 - MONTGOMRY MALL - SILVR SPRNG STA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J2v1",
						Name:            "J2 - MEDICAL CTR STA - SILVR SPRNG STA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J2v2",
						Name:            "J2 - SILVER SPRNG - MONT MALL/BATTRY LA",
						LineDescription: "Bethesda-Silver Spring Line",
					},
					{
						RouteID:         "J4",
						Name:            "J4 - BETHESDA STA - COLLEGE PARK STA",
						LineDescription: "College Park-Bethesda Limited",
					},
					{
						RouteID:         "K12",
						Name:            "K12 - BRANCH AVE ST-SUITLAND ST",
						LineDescription: "Forestville Line",
					},
					{
						RouteID:         "K12v1",
						Name:            "K12 - PENN MAR - SUITLAND ST",
						LineDescription: "Forestville Line",
					},
					{
						RouteID:         "K12v2",
						Name:            "K12 - BRANCH AVE - PENN MAR",
						LineDescription: "Forestville Line",
					},
					{
						RouteID:         "K2",
						Name:            "K2 - FT TOTTEN STA - TAKOMA STA",
						LineDescription: "Takoma-Fort Totten Line",
					},
					{
						RouteID:         "K6",
						Name:            "K6 - FORT TOTTEN STA - WHITE OAK",
						LineDescription: "New Hampshire Ave-Maryland Line",
					},
					{
						RouteID:         "K6v1",
						Name:            "K6 - TLTC - FORT TOTTEN STA",
						LineDescription: "New Hampshire Ave-Maryland Line",
					},
					{
						RouteID:         "K9",
						Name:            "K9 - FORT TOTTEN STA - FDA/FRC (PM)",
						LineDescription: "New Hampshire Ave-Maryland Limited Line",
					},
					{
						RouteID:         "K9v1",
						Name:            "K9 - FORT TOTTEN STA - FDA/FRC (AM)",
						LineDescription: "New Hampshire Ave-Maryland Limited Line",
					},
					{
						RouteID:         "L1",
						Name:            "L1 - POTOMAC PK  - CHEVY CHASE",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L2",
						Name:            "L2 - FARRAGUT SQ - CHEVY CHASE",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L2v1",
						Name:            "L2 - VAN NESS-UDC STA - CHEVY CHASE",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L2v2",
						Name:            "L2 - FARRAGUT SQ - BETHESDA",
						LineDescription: "Connecticut Ave Line",
					},
					{
						RouteID:         "L8",
						Name:            "L8 - FRIENDSHIP HTS STA - ASPEN HILL",
						LineDescription: "Connecticut Ave-Maryland Line",
					},
					{
						RouteID:         "M4",
						Name:            "M4 - SIBLEY HOSPITAL - TENLEYTOWN/AU STA",
						LineDescription: "Nebraska Ave Line",
					},
					{
						RouteID:         "M4v1",
						Name:            "M4 - TENLEYTOWN   - PINEHRST CIR",
						LineDescription: "Nebraska Ave Line",
					},
					{
						RouteID:         "M4v2",
						Name:            "M4 - PINEHRST CIR - TENLEYTOWN",
						LineDescription: "Nebraska Ave Line",
					},
					{
						RouteID:         "M6",
						Name:            "M6 - POTOMAC AVE - ALABAMA + PENN",
						LineDescription: "Fairfax Village Line",
					},
					{
						RouteID:         "M6v1",
						Name:            "M6 - ALABAMA + PENN - FAIRFAX VILLAGE",
						LineDescription: "Fairfax Village Line",
					},
					{
						RouteID:         "MW1",
						Name:            "MW1 - BRADDOCK RD - PENTAGON CITY",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "MW1v1",
						Name:            "MW1 - POTOMAC YARD - CRYSTAL CITY",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "MW1v2",
						Name:            "MW1 - CRYSTAL CITY - BRADDOCK RD",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "MW1v3",
						Name:            "MW1 - BRADDOCK RD - CRYSTAL CITY",
						LineDescription: "Metroway-Potomac Yard Line",
					},
					{
						RouteID:         "N2",
						Name:            "N2 - FRIENDSHIP HTS - FARRAGUT SQ",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "N4",
						Name:            "N4 - FRIENDSHP HTS - POTOMAC PARK",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "N4v1",
						Name:            "N4 - FRIENDSHP HTS - FARRAGUT  SQ",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "N6",
						Name:            "N6 - FRNDSHIP HTS - FARRAGUT SQ",
						LineDescription: "Massachusetts Ave Line",
					},
					{
						RouteID:         "NH1",
						Name:            "NH1 - NATIONAL HARBOR-SOUTHERN AVE STA",
						LineDescription: "National Harbor-Southern Avenue Line",
					},
					{
						RouteID:         "NH2",
						Name:            "NH2 - HUNTINGTON STA-NATIONAL HARBOR",
						LineDescription: "National Harbor-Alexandria Line",
					},
					{
						RouteID:         "P12",
						Name:            "P12 - EASTOVER - ADDISON RD STA (NO HOSP)",
						LineDescription: "Eastover-Addison Road Line",
					},
					{
						RouteID:         "P12v1",
						Name:            "P12 - IVERSON MALL - ADDISON RD STA",
						LineDescription: "Eastover-Addison Road Line",
					},
					{
						RouteID:         "P12v2",
						Name:            "P12 - SUITLAND STA - ADDISON RD STA",
						LineDescription: "Eastover-Addison Road Line",
					},
					{
						RouteID:         "P18",
						Name:            "P18 - FT WASH P+R LOT - SOUTHERN AVE",
						LineDescription: "Oxon Hill-Fort Washington Line",
					},
					{
						RouteID:         "P19",
						Name:            "P19 - FT WASH P+R LOT-SOUTHERN AVE STA",
						LineDescription: "Oxon Hill-Fort Washington Line",
					},
					{
						RouteID:         "P6",
						Name:            "P6 - ANACOSTIA STA - RHODE ISLAND STA",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v1",
						Name:            "P6 - ECKINGTON - RHODE ISLAND AVE",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v2",
						Name:            "P6 - RHODE ISLAND AVE - ECKINGTON",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v3",
						Name:            "P6 - ANACOSTIA STA - ARCHIVES",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "P6v4",
						Name:            "P6 - ARCHIVES - ANACOSTIA",
						LineDescription: "Anacostia-Eckington Line",
					},
					{
						RouteID:         "Q1",
						Name:            "Q1 - SILVR SPRNG STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q2",
						Name:            "Q2 - SILVR SPRNG STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q2v1",
						Name:            "Q2 - MONT COLLEGE - SILVR SPRNG STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q2v2",
						Name:            "Q2 - SILVR SPRNG STA - MONT COLLEGE",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q4",
						Name:            "Q4 - SILVER SPRNG STA - ROCKVILLE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q4v1",
						Name:            "Q4 - WHEATON STA - ROCKVILLE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q5",
						Name:            "Q5 - WHEATON STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q6",
						Name:            "Q6 - WHEATON STA - SHADY GRVE STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "Q6v1",
						Name:            "Q6 - ROCKVILLE STA - WHEATON STA",
						LineDescription: "Viers Mill Road Line",
					},
					{
						RouteID:         "R1",
						Name:            "R1 - FORT TOTTEN STA - ADELPHI",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R12",
						Name:            "R12 - DEANWOOD STA - GREENBELT STA",
						LineDescription: "Kenilworth Avenue Line",
					},
					{
						RouteID:         "R12v1",
						Name:            "R12 - DEANWOOD STA - GREENBELT STA",
						LineDescription: "Kenilworth Avenue Line",
					},
					{
						RouteID:         "R2",
						Name:            "R2 - FORT TOTTEN - CALVERTON",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R2v1",
						Name:            "R2 - HIGH POINT HS - FORT TOTTEN",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R2v2",
						Name:            "R2 - POWDER MILL+CHERRY HILL - CALVERTON",
						LineDescription: "Riggs Road Line",
					},
					{
						RouteID:         "R4",
						Name:            "R4 - BROOKLAND STA- HIGHVIEW",
						LineDescription: "Queens Chapel Road Line",
					},
					{
						RouteID:         "REX",
						Name:            "REX - FT BELVOIR POST - KING ST STA",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv1",
						Name:            "REX - FT BELVOIR COMM HOSP - KING ST STA",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv2",
						Name:            "REX - KING ST STA - FT BELVOIR COMM HOSP",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv3",
						Name:            "REX - KING ST STA - WOODLAWN",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "REXv4",
						Name:            "REX - WOODLAWN - KING ST STA",
						LineDescription: "Richmond Highway Express",
					},
					{
						RouteID:         "S1",
						Name:            "S1 - NORTHERN DIVISION - POTOMAC PK",
						LineDescription: "16th Street-Potomac Park Line",
					},
					{
						RouteID:         "S1v1",
						Name:            "S1 - VIRGINIA+E - COLORDO+16TH",
						LineDescription: "16th Street-Potomac Park Line",
					},
					{
						RouteID:         "S2",
						Name:            "S2 - FED TRIANGLE  - SILVER SPRNG",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S2v1",
						Name:            "S2 - 16TH & HARVARD - MCPHERSON SQ",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S35",
						Name:            "S35 - BRANCH + RANDLE CIR - FT DUPONT",
						LineDescription: "Fort Dupont Shuttle Line",
					},
					{
						RouteID:         "S4",
						Name:            "S4 - FED TRIANGLE - SILVER SPRNG",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S4v1",
						Name:            "S4 - SILVER SPRING STA - FRANKLIN SQ",
						LineDescription: "16th Street Line",
					},
					{
						RouteID:         "S41",
						Name:            "S41 - CARVER TERRACE - RHODE ISLAND AVE",
						LineDescription: "Rhode Island Ave-Carver Terrace Line",
					},
					{
						RouteID:         "S80",
						Name:            "S80 - FRANCONIA-SPRNGFLD - METRO PARK",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S80v1",
						Name:            "S80 - FRANCONIA-SPRINGFLD - HILTON",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S80v2",
						Name:            "S80 - HILTON - FRANCONIA-SPRNGFLD",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S9",
						Name:            "S9 - FRANKLIN SQ - SILVER SPRING STA",
						LineDescription: "16th Street Limited Line",
					},
					{
						RouteID:         "S9v1",
						Name:            "S9 - FRANKLIN SQ - COLORADO + 16TH",
						LineDescription: "16th Street Limited Line",
					},
					{
						RouteID:         "S91",
						Name:            "S91 - FRANCONIA SPRINGFLD STA SHUTTLE",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "S91v1",
						Name:            "S91 - FRANCONIA SPRINGFLD STA SHUTTLE",
						LineDescription: "Springfield Circulator-Metro Park Shuttle (TAGS)",
					},
					{
						RouteID:         "T14",
						Name:            "T14 - RHD ISLND AVE STA-NEW CARRLTN STA",
						LineDescription: "Rhode Island Ave-New Carrollton Line",
					},
					{
						RouteID:         "T14v1",
						Name:            "T14 - MT RAINIER - NEW CARRLTN STA",
						LineDescription: "Rhode Island Ave-New Carrollton Line",
					},
					{
						RouteID:         "T18",
						Name:            "T18 - R I AVE STA - NEW CARROLLTON STA",
						LineDescription: "Annapolis Road Line",
					},
					{
						RouteID:         "T18v1",
						Name:            "T18 - BLADENSBURG HS - NEW CARROLLTON STA",
						LineDescription: "Annapolis Road Line",
					},
					{
						RouteID:         "T2",
						Name:            "T2 - FRIENDSHIP HTS - ROCKVILLE STA",
						LineDescription: "River Road Line",
					},
					{
						RouteID:         "U4",
						Name:            "U4 - MINNESOTA AVE - SHERIFF RD",
						LineDescription: "Sheriff Road-River Terrace Line",
					},
					{
						RouteID:         "U4v1",
						Name:            "U4 - SHERIFF RD - MINNESOTA STA",
						LineDescription: "Sheriff Road-River Terrace Line",
					},
					{
						RouteID:         "U4v2",
						Name:            "U4 - RIVER TERRACE - MINNESOTA AVE",
						LineDescription: "Sheriff Road-River Terrace Line",
					},
					{
						RouteID:         "U5",
						Name:            "U5 - MINNESOTA AVE-MARSHALL HTS",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U6",
						Name:            "U6 - MINNESOTA - LINCOLN HEIGHTS",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U6v1",
						Name:            "U6 - 37TH + RIDGE - PLUMMER ES",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U6v2",
						Name:            "U6 - LINCOLN HTS - E CAP + 47TH ST NE",
						LineDescription: "Marshall Heights Line",
					},
					{
						RouteID:         "U7",
						Name:            "U7 - RIDGE + ANACOSTIA - DEANWOOD",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v1",
						Name:            "U7 - MINNESOTA STA - KENILWORTH HAYES",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v2",
						Name:            "U7 - KENILWORTH HAYES - MINNESOTA STA",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v3",
						Name:            "U7 - MINNESOTA STA - DEANWOOD",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "U7v4",
						Name:            "U7 - DEANWOOD - MINNESOTA AVE",
						LineDescription: "Deanwood-Minnesota Ave Line",
					},
					{
						RouteID:         "V1",
						Name:            "V1 - BUR OF ENGRVNG - BENNING HTS",
						LineDescription: "Benning Heights-M Street Line",
					},
					{
						RouteID:         "V12",
						Name:            "V12 - SUITLAND STA - ADDISON RD STA",
						LineDescription: "District Heights-Suitland Line",
					},
					{
						RouteID:         "V14",
						Name:            "V14 - PENN MAR - DEANWOOD STA",
						LineDescription: "District Heights - Seat Pleasant Line",
					},
					{
						RouteID:         "V14v1",
						Name:            "V14 - PENN MAR  - ADDISON RD STA",
						LineDescription: "District Heights - Seat Pleasant Line",
					},
					{
						RouteID:         "V2",
						Name:            "V2 - ANACOSTIA - CAPITOL HGTS",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V2v1",
						Name:            "V2 - MINNESOTA AVE - ANACOSTIA",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V4",
						Name:            "V4 - 1ST + K ST SE - CAPITOL HGTS",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V4v1",
						Name:            "V4 - MINN AVE STA-CAPITOL HGTS",
						LineDescription: "Capitol Heights-Minnesota Avenue Line",
					},
					{
						RouteID:         "V7",
						Name:            "V7 - CONGRESS HGTS - MINN STA",
						LineDescription: "Benning Heights-Alabama Ave Line",
					},
					{
						RouteID:         "V8",
						Name:            "V8 - BENNING HGTS - MINN STA",
						LineDescription: "Benning Heights-Alabama Ave Line",
					},
					{
						RouteID:         "W1",
						Name:            "W1 - FORT DRUM  - SOUTHERN AVE STA",
						LineDescription: "Shipley Terrace - Fort Drum Line",
					},
					{
						RouteID:         "W14",
						Name:            "W14 - FT WASHINGTON-SOUTHERN AVE STA",
						LineDescription: "Bock Road Line",
					},
					{
						RouteID:         "W14v1",
						Name:            "W14 - SOUTHERN AVE - FRIENDLY",
						LineDescription: "Bock Road Line",
					},
					{
						RouteID:         "W14v2",
						Name:            "W14 - FRIENDLY - SOUTHERN AVE",
						LineDescription: "Bock Road Line",
					},
					{
						RouteID:         "W2",
						Name:            "W2 - MALCOM X+OAKWD - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v1",
						Name:            "W2 - NYLDR+GOOD HOP- HOWARD+ANACOSTIA",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v2",
						Name:            "W2 - NAYLOR+GOODHOPE - MALCM X+OAKWOOD",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v3",
						Name:            "W2 - HOWRD+ANACOSTIA-NAYLOR+GOODHOPE",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v4",
						Name:            "W2 - ANACOSTIA STA - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v5",
						Name:            "W2 - ANACOSTIA STA - SOUTHERN AVE STA",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v6",
						Name:            "W2 - MELLN+M L KNG - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W2v7",
						Name:            "W2 - SOUTHERN AVE-WASHINGTON OVERLOOK",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W3",
						Name:            "W3 - MALCOM X+OAKWD - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W3v1",
						Name:            "W3 - MELLN+M L KNG - UNITED MEDICAL CTR",
						LineDescription: "United Medical Center-Anacostia Line",
					},
					{
						RouteID:         "W4",
						Name:            "W4 - ANACOSTIA STA - DEANWOOD STA",
						LineDescription: "Deanwood-Alabama Ave Line",
					},
					{
						RouteID:         "W4v1",
						Name:            "W4 - MALCOLM X & PORTLAND - DEANWOOD",
						LineDescription: "Deanwood-Alabama Ave Line",
					},
					{
						RouteID:         "W4v2",
						Name:            "W4 - DEANWOOD STA - MALCLM X+ ML KING",
						LineDescription: "Deanwood-Alabama Ave Line",
					},
					{
						RouteID:         "W45",
						Name:            "W45 - TENLEYTOWN STA - 16TH + SHEPHERD",
						LineDescription: "Mt Pleasant-Tenleytown Line",
					},
					{
						RouteID:         "W47",
						Name:            "W47 - TENLEYTOWN STA - COLUMBIA HTS STA",
						LineDescription: "Mt Pleasant-Tenleytown Line",
					},
					{
						RouteID:         "W5",
						Name:            "W5 - DC VILLAGE-USCG (VIA ANAC)",
						LineDescription: "Anacostia-Fort Drum Line",
					},
					{
						RouteID:         "W6",
						Name:            "W6 - ANACOSTIA - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W6v1",
						Name:            "W6 - NAYLOR+GOODHOPE - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W8",
						Name:            "W8 - ANACOSTIA - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W8v1",
						Name:            "W8 - ANACOSTIA-NAYLOR+GOODHOPE",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "W8v2",
						Name:            "W8 - NAYLOR+GOODHOPE - ANACOSTIA",
						LineDescription: "Garfield-Anacostia Loop Line",
					},
					{
						RouteID:         "X1",
						Name:            "X1 - FOGGY BOTTOM+GWU- MINNESOTA STA",
						LineDescription: "Benning Road Line",
					},
					{
						RouteID:         "X2",
						Name:            "X2 - LAFAYETTE SQ - MINNESOTA STA",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X2v1",
						Name:            "X2 - PHELPS HS - MINNESOTA STA",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X2v2",
						Name:            "X2 - FRIENDSHIP EDISON PCS-MINNESOTA STA",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X2v3",
						Name:            "X2 - PHELPS HS - LAFAYTTE SQ",
						LineDescription: "Benning Road-H Street Line",
					},
					{
						RouteID:         "X3",
						Name:            "X3 - DUKE ELLINGTON BR - MINNESOTA STN",
						LineDescription: "Benning Road Line",
					},
					{
						RouteID:         "X3v1",
						Name:            "X3 - KIPP DC - MINNESOTA AVE STN",
						LineDescription: "Benning Road Line",
					},
					{
						RouteID:         "X8",
						Name:            "X8 - UNION STATION - CARVER TERR",
						LineDescription: "Maryland Ave Line",
					},
					{
						RouteID:         "X9",
						Name:            "X9 - NY AVE & 12TH NW - CAPITOL HTS STA",
						LineDescription: "Benning Road-H Street Limited Line",
					},
					{
						RouteID:         "X9v1",
						Name:            "X9 - NY AVE & 12TH ST NW - MINNESOTA AVE",
						LineDescription: "Benning Road-H Street Limited Line",
					},
					{
						RouteID:         "X9v2",
						Name:            "X9 - MINNESOTA AVE ST - NY AVE & 12TH ST",
						LineDescription: "Benning Road-H Street Limited Line",
					},
					{
						RouteID:         "Y2",
						Name:            "Y2 - SILVER SPRING STA - MONTG MED CTR",
						LineDescription: "Georgia Ave-Maryland Line",
					},
					{
						RouteID:         "Y7",
						Name:            "Y7 - SILVER SPRING STA - ICC P&R",
						LineDescription: "Georgia Ave-Maryland Line",
					},
					{
						RouteID:         "Y8",
						Name:            "Y8 - SILVER SPR STA - MONTGOMERY MED CTR",
						LineDescription: "Georgia Ave-Maryland Line",
					},
					{
						RouteID:         "Z11",
						Name:            "Z11 - SILVR SPRING - BURTONSVILLE P&R",
						LineDescription: "Greencastle-Briggs Chaney Express Line",
					},
					{
						RouteID:         "Z11v1",
						Name:            "Z11 - GREENCASTLE - SILVR SPRING",
						LineDescription: "Greencastle-Briggs Chaney Express Line",
					},
					{
						RouteID:         "Z2",
						Name:            "Z2 - SILVR SPRING - OLNEY (NO BLAKE HS)",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z2v1",
						Name:            "Z2 - SILVER SPRING - BONIFANT & NH",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z2v2",
						Name:            "Z2 - SILVER SPRING STA-OLNEY (BLAKE HS)",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z2v3",
						Name:            "Z2 - COLESVILLE-SILVER SPRING",
						LineDescription: "Colesville - Ashton Line",
					},
					{
						RouteID:         "Z6",
						Name:            "Z6 - SILVR SPRNG STA - BURTONSVILLE",
						LineDescription: "Calverton-Westfarm Line",
					},
					{
						RouteID:         "Z6v1",
						Name:            "Z6 - CASTLE BLVD - SILVER SPRING STA",
						LineDescription: "Calverton-Westfarm Line",
					},
					{
						RouteID:         "Z6v2",
						Name:            "Z6 - SILVER SPRING -CASTLE BLVD",
						LineDescription: "Calverton-Westfarm Line",
					},
					{
						RouteID:         "Z7",
						Name:            "Z7 - SLVER SPRNG STA-S LAUREL P&R (4COR)",
						LineDescription: "Laurel-Burtonsville Express Line",
					},
					{
						RouteID:         "Z7v1",
						Name:            "Z7 - S LAUREL P&R-SILVR SPR (NO 4COR)",
						LineDescription: "Laurel-Burtonsville Express Line",
					},
					{
						RouteID:         "Z8",
						Name:            "Z8 - SILVER SRING STA - BRIGGS CHANEY",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v1",
						Name:            "Z8 - WHITE OAK - SILVER SPRING STA",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v2",
						Name:            "Z8 - SILVR SPRNG - CSTLE BLVD VERIZON",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v3",
						Name:            "Z8 - SILVER SPRING STA - CASTLE BLVD",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v4",
						Name:            "Z8 - SILVER SPRING - GRNCSTLE (VERIZON)",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v5",
						Name:            "Z8 - SILVER SPRING STA - GREENCASTLE",
						LineDescription: "Fairland Line",
					},
					{
						RouteID:         "Z8v6",
						Name:            "Z8 - SILVER SPRING STA - WHITE OAK",
						LineDescription: "Fairland Line",
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
			continue
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
			continue
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

func TestGetRoutes(t *testing.T) {
	jsonAndXmlPaths := []string{"/Bus.svc/json/jRoutes", "/Bus.svc/Routes"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetRoutes")
			continue
		}

		for _, request := range testRequests {
			response, err := testService.GetRoutes()

			if err != nil {
				t.Errorf("error calling GetRoutes, error: %s", err.Error())
				continue
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}