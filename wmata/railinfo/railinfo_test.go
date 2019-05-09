package railinfo

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
	stringParam1         string
	stringParam2         string
	requestType          interface{}
	response             string
	unmarshalledResponse interface{}
	expectedError        error
}

var testData = map[string][]testResponseData{
	"/Rail.svc/json/jLines": {
		testResponseData{
			rawQuery: "",
			response: `{"Lines":[{"LineCode":"BL","DisplayName":"Blue","StartStationCode":"J03","EndStationCode":"G05","InternalDestination1":"","InternalDestination2":""},{"LineCode":"GR","DisplayName":"Green","StartStationCode":"F11","EndStationCode":"E10","InternalDestination1":"","InternalDestination2":""},{"LineCode":"OR","DisplayName":"Orange","StartStationCode":"K08","EndStationCode":"D13","InternalDestination1":"","InternalDestination2":""},{"LineCode":"RD","DisplayName":"Red","StartStationCode":"A15","EndStationCode":"B11","InternalDestination1":"A11","InternalDestination2":"B08"},{"LineCode":"SV","DisplayName":"Silver","StartStationCode":"N06","EndStationCode":"G05","InternalDestination1":"","InternalDestination2":""},{"LineCode":"YL","DisplayName":"Yellow","StartStationCode":"C15","EndStationCode":"E06","InternalDestination1":"E01","InternalDestination2":""}]}`,
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
			response:     `{"StationsParking":[{"Code":"B08","Notes":"Parking is available at Montgomery County lots and garages.","AllDayParking":{"TotalCount":0,"RiderCost":null,"NonRiderCost":null,"SaturdayRiderCost":null,"SaturdayNonRiderCost":null},"ShortTermParking":{"TotalCount":0,"Notes":null}}]}`,
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
			rawQuery:     "StationCode=K08",
			stringParam1: "K08",
			response:     `{"StationsParking":[{"Code":"K08","Notes":"North Kiss & Ride - 45 short term metered spaces. South Kiss & Ride - 26 short term metered spaces.  101 spaces metered for 12-hr. max @ $1.00 per 60 mins. 17 spaces metered for 7-hr. max. @ $1.00 per 60 mins. Parking available from 8:30 AM to 2 AM.","AllDayParking":{"TotalCount":5169,"RiderCost":4.95,"NonRiderCost":4.95,"SaturdayRiderCost":0,"SaturdayNonRiderCost":0},"ShortTermParking":{"TotalCount":71,"Notes":"Parking available in section B between 8:30 AM - 3:30 PM and 7 PM - 2 AM, in section D between 10 AM - 2 PM."}}]}`,
			unmarshalledResponse: &GetParkingInformationResponse{
				ParkingInformation: []StationParking{
					{
						StationCode: "K08",
						Notes:       "North Kiss & Ride - 45 short term metered spaces. South Kiss & Ride - 26 short term metered spaces.  101 spaces metered for 12-hr. max @ $1.00 per 60 mins. 17 spaces metered for 7-hr. max. @ $1.00 per 60 mins. Parking available from 8:30 AM to 2 AM.",
						AllDay: AllDayParking{
							TotalCount:           5169,
							RiderCost:            4.95,
							NonRiderCost:         4.95,
							SaturdayRiderCost:    0,
							SaturdayNonRiderCost: 0,
						},
						ShortTerm: ShortTermParking{
							TotalCount: 71,
							Notes:      "Parking available in section B between 8:30 AM - 3:30 PM and 7 PM - 2 AM, in section D between 10 AM - 2 PM.",
						},
					},
				},
			},
		},
	},
	"/Rail.svc/json/jPath": {
		{
			rawQuery:     "FromStationCode=A09&ToStationCode=B04",
			stringParam1: "A09",
			stringParam2: "B04",
			response:     `{"Path":[{"LineCode":"RD","StationCode":"A09","StationName":"Bethesda","SeqNum":1,"DistanceToPrev":0},{"LineCode":"RD","StationCode":"A08","StationName":"Friendship Heights","SeqNum":2,"DistanceToPrev":9095},{"LineCode":"RD","StationCode":"A07","StationName":"Tenleytown-AU","SeqNum":3,"DistanceToPrev":4135},{"LineCode":"RD","StationCode":"A06","StationName":"Van Ness-UDC","SeqNum":4,"DistanceToPrev":5841},{"LineCode":"RD","StationCode":"A05","StationName":"Cleveland Park","SeqNum":5,"DistanceToPrev":3320},{"LineCode":"RD","StationCode":"A04","StationName":"Woodley Park-Zoo/Adams Morgan","SeqNum":6,"DistanceToPrev":3740},{"LineCode":"RD","StationCode":"A03","StationName":"Dupont Circle","SeqNum":7,"DistanceToPrev":6304},{"LineCode":"RD","StationCode":"A02","StationName":"Farragut North","SeqNum":8,"DistanceToPrev":2711},{"LineCode":"RD","StationCode":"A01","StationName":"Metro Center","SeqNum":9,"DistanceToPrev":4178},{"LineCode":"RD","StationCode":"B01","StationName":"Gallery Pl-Chinatown","SeqNum":10,"DistanceToPrev":1505},{"LineCode":"RD","StationCode":"B02","StationName":"Judiciary Square","SeqNum":11,"DistanceToPrev":1875},{"LineCode":"RD","StationCode":"B03","StationName":"Union Station","SeqNum":12,"DistanceToPrev":3446},{"LineCode":"RD","StationCode":"B35","StationName":"NoMa-Gallaudet U","SeqNum":13,"DistanceToPrev":3553},{"LineCode":"RD","StationCode":"B04","StationName":"Rhode Island Ave-Brentwood","SeqNum":14,"DistanceToPrev":5771}]}`,
			unmarshalledResponse: &GetPathBetweenStationsResponse{
				Path: []PathItem{
					{
						LineCode:                  "RD",
						StationCode:               "A09",
						StationName:               "Bethesda",
						SequenceNumber:            1,
						DistanceToPreviousStation: 0,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A08",
						StationName:               "Friendship Heights",
						SequenceNumber:            2,
						DistanceToPreviousStation: 9095,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A07",
						StationName:               "Tenleytown-AU",
						SequenceNumber:            3,
						DistanceToPreviousStation: 4135,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A06",
						StationName:               "Van Ness-UDC",
						SequenceNumber:            4,
						DistanceToPreviousStation: 5841,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A05",
						StationName:               "Cleveland Park",
						SequenceNumber:            5,
						DistanceToPreviousStation: 3320,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A04",
						StationName:               "Woodley Park-Zoo/Adams Morgan",
						SequenceNumber:            6,
						DistanceToPreviousStation: 3740,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A03",
						StationName:               "Dupont Circle",
						SequenceNumber:            7,
						DistanceToPreviousStation: 6304,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A02",
						StationName:               "Farragut North",
						SequenceNumber:            8,
						DistanceToPreviousStation: 2711,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A01",
						StationName:               "Metro Center",
						SequenceNumber:            9,
						DistanceToPreviousStation: 4178,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B01",
						StationName:               "Gallery Pl-Chinatown",
						SequenceNumber:            10,
						DistanceToPreviousStation: 1505,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B02",
						StationName:               "Judiciary Square",
						SequenceNumber:            11,
						DistanceToPreviousStation: 1875,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B03",
						StationName:               "Union Station",
						SequenceNumber:            12,
						DistanceToPreviousStation: 3446,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B35",
						StationName:               "NoMa-Gallaudet U",
						SequenceNumber:            13,
						DistanceToPreviousStation: 3553,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B04",
						StationName:               "Rhode Island Ave-Brentwood",
						SequenceNumber:            14,
						DistanceToPreviousStation: 5771,
					},
				},
			},
		},
		{
			expectedError: errors.New("fromStation and toStation are required parameters"),
		},
	},
	"/Rail.svc/json/jStationEntrances": {
		{
			rawQuery: "Lat=38.897383&Lon=-77.007262&Radius=500",
			requestType: GetStationEntrancesRequest{
				Latitude:  38.897383,
				Longitude: -77.007262,
				Radius:    500,
			},
			response: `{"Entrances":[{"ID":"54","Name":"SOUTH ENTRANCE (MASS AVE EXIT, NORTHEAST CORNER OF 1ST ST & MASSACHUSETTS AVE)","StationCode1":"B03","StationCode2":"","Description":"Station entrance from 1st St NE to southeast corner of the Union station building.","Lat":38.897383,"Lon":-77.007262},{"ID":"55","Name":"NORTH ENTRANCE (1ST ST EXIT, WEST SIDE OF 1ST ST BETWEEN G ST AND MASSACHUSETTS AVE)","StationCode1":"B03","StationCode2":"","Description":"Station entrance from northeast corner of Massachusetts Ave NE and 1st NE.","Lat":38.89845,"Lon":-77.007243},{"ID":"53","Name":"ENTRANCE FROM AMTRAK, MARC, VRE TRAINS","StationCode1":"B03","StationCode2":"","Description":"Escalator entrance from the passageway to  AMTRAK, MARC, VRE TRAINS","Lat":38.898541,"Lon":-77.006984}]}`,
			unmarshalledResponse: &GetStationEntrancesResponse{
				Entrances: []StationEntrance{
					{
						ID:           "54",
						Name:         "SOUTH ENTRANCE (MASS AVE EXIT, NORTHEAST CORNER OF 1ST ST & MASSACHUSETTS AVE)",
						StationCode1: "B03",
						StationCode2: "",
						Description:  "Station entrance from 1st St NE to southeast corner of the Union station building.",
						Latitude:     38.897383,
						Longitude:    -77.007262,
					},
					{
						ID:           "55",
						Name:         "NORTH ENTRANCE (1ST ST EXIT, WEST SIDE OF 1ST ST BETWEEN G ST AND MASSACHUSETTS AVE)",
						StationCode1: "B03",
						StationCode2: "",
						Description:  "Station entrance from northeast corner of Massachusetts Ave NE and 1st NE.",
						Latitude:     38.89845,
						Longitude:    -77.007243,
					},
					{
						ID:           "53",
						Name:         "ENTRANCE FROM AMTRAK, MARC, VRE TRAINS",
						StationCode1: "B03",
						StationCode2: "",
						Description:  "Escalator entrance from the passageway to  AMTRAK, MARC, VRE TRAINS",
						Latitude:     38.898541,
						Longitude:    -77.006984,
					},
				},
			},
		},
	},
	"/Rail.svc/json/jStationInfo": {
		{
			rawQuery:     "StationCode=A07",
			stringParam1: "A07",
			response:     `{"Code":"A07","Name":"Tenleytown-AU","StationTogether1":"","StationTogether2":"","LineCode1":"RD","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.947808,"Lon":-77.079615,"Address":{"Street":"4501 Wisconsin Avenue NW","City":"Washington","State":"DC","Zip":"20016"}}`,
			unmarshalledResponse: &GetStationInformationResponse{
				Address: StationAddress{
					Street: "4501 Wisconsin Avenue NW",
					City:   "Washington",
					State:  "DC",
					Zip:    "20016",
				},
				Latitude:         38.947808,
				LineCode1:        "RD",
				LineCode2:        "",
				Longitude:        -77.079615,
				Name:             "Tenleytown-AU",
				StationCode:      "A07",
				StationTogether1: "",
				StationTogether2: "",
			},
		},
		{
			rawQuery:     "StationCode=F01",
			stringParam1: "F01",
			response:     `{"Code":"F01","Name":"Gallery Pl-Chinatown","StationTogether1":"B01","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.89834,"Lon":-77.021851,"Address":{"Street":"630 H St. NW","City":"Washington","State":"DC","Zip":"20001"}}`,
			unmarshalledResponse: &GetStationInformationResponse{
				Address: StationAddress{
					Street: "630 H St. NW",
					City:   "Washington",
					State:  "DC",
					Zip:    "20001",
				},
				Latitude:         38.89834,
				LineCode1:        "GR",
				LineCode2:        "YL",
				Longitude:        -77.021851,
				Name:             "Gallery Pl-Chinatown",
				StationCode:      "F01",
				StationTogether1: "B01",
				StationTogether2: "",
			},
		},
		{
			expectedError: errors.New("stationCode is a required parameter"),
		},
	},
	"/Rail.svc/json/jStations": {
		{
			rawQuery:     "LineCode=GR",
			stringParam1: "GR",
			response:     `{"Stations":[{"Code":"E01","Name":"Mt Vernon Sq 7th St-Convention Center","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.905604,"Lon":-77.022256,"Address":{"Street":"700 M St. NW","City":"Washington","State":"DC","Zip":"20001"}},{"Code":"E02","Name":"Shaw-Howard U","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.912919,"Lon":-77.022194,"Address":{"Street":"1701 8th St. NW","City":"Washington","State":"DC","Zip":"20001"}},{"Code":"E03","Name":"U Street/African-Amer Civil War Memorial/Cardozo","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.916489,"Lon":-77.028938,"Address":{"Street":"1240 U Street NW","City":"Washington","State":"DC","Zip":"20009"}},{"Code":"E04","Name":"Columbia Heights","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.928672,"Lon":-77.032775,"Address":{"Street":"3030 14th St. NW","City":"Washington","State":"DC","Zip":"20009"}},{"Code":"E05","Name":"Georgia Ave-Petworth","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.936077,"Lon":-77.024728,"Address":{"Street":"3700 Georgia Avenue NW","City":"Washington","State":"DC","Zip":"20010"}},{"Code":"E06","Name":"Fort Totten","StationTogether1":"B06","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.951777,"Lon":-77.002174,"Address":{"Street":"550 Galloway Street NE","City":"Washington","State":"DC","Zip":"20011"}},{"Code":"E07","Name":"West Hyattsville","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.954931,"Lon":-76.969881,"Address":{"Street":"2700 Hamilton St.","City":"Hyattsville","State":"MD","Zip":"20782"}},{"Code":"E08","Name":"Prince George's Plaza","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.965276,"Lon":-76.956182,"Address":{"Street":"3575 East West Highway","City":"Hyattsville","State":"MD","Zip":"20782"}},{"Code":"E09","Name":"College Park-U of Md","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.978523,"Lon":-76.928432,"Address":{"Street":"4931 Calvert Road","City":"College Park","State":"MD","Zip":"20740"}},{"Code":"E10","Name":"Greenbelt","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":39.011036,"Lon":-76.911362,"Address":{"Street":"5717 Greenbelt Metro Drive","City":"Greenbelt","State":"MD","Zip":"20740"}},{"Code":"F01","Name":"Gallery Pl-Chinatown","StationTogether1":"B01","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.89834,"Lon":-77.021851,"Address":{"Street":"630 H St. NW","City":"Washington","State":"DC","Zip":"20001"}},{"Code":"F02","Name":"Archives-Navy Memorial-Penn Quarter","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.893893,"Lon":-77.021902,"Address":{"Street":"701 Pennsylvania Avenue NW","City":"Washington","State":"DC","Zip":"20004"}},{"Code":"F03","Name":"L'Enfant Plaza","StationTogether1":"D03","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.884775,"Lon":-77.021964,"Address":{"Street":"600 Maryland Avenue SW","City":"Washington","State":"DC","Zip":"20024"}},{"Code":"F04","Name":"Waterfront","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.876221,"Lon":-77.017491,"Address":{"Street":"399 M Street SW","City":"Washington","State":"DC","Zip":"20024"}},{"Code":"F05","Name":"Navy Yard-Ballpark","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.876588,"Lon":-77.005086,"Address":{"Street":"200 M Street SE","City":"Washington","State":"DC","Zip":"20003"}},{"Code":"F06","Name":"Anacostia","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.862072,"Lon":-76.995648,"Address":{"Street":"1101 Howard Road SE","City":"Washington","State":"DC","Zip":"20020"}},{"Code":"F07","Name":"Congress Heights","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.845334,"Lon":-76.98817,"Address":{"Street":"1290 Alabama Avenue SE","City":"Washington","State":"DC","Zip":"20020"}},{"Code":"F08","Name":"Southern Avenue","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.840974,"Lon":-76.97536,"Address":{"Street":"1411 Southern Avenue","City":"Temple Hills","State":"MD","Zip":"20748"}},{"Code":"F09","Name":"Naylor Road","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.851187,"Lon":-76.956565,"Address":{"Street":"3101 Branch Avenue","City":"Temple Hills","State":"MD","Zip":"20748"}},{"Code":"F10","Name":"Suitland","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.843891,"Lon":-76.932022,"Address":{"Street":"4500 Silver Hill Road","City":"Suitland","State":"MD","Zip":"20746"}},{"Code":"F11","Name":"Branch Ave","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.826995,"Lon":-76.912134,"Address":{"Street":"4704 Old Soper Road","City":"Suitland","State":"MD","Zip":"20746"}}]}`,
			unmarshalledResponse: &GetStationListResponse{
				Stations: []GetStationListResponseItem{
					{
						StationCode:      "E01",
						Name:             "Mt Vernon Sq 7th St-Convention Center",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.905604,
						Longitude:        -77.022256,
						Address: StationAddress{
							Street: "700 M St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20001",
						},
					},
					{
						StationCode:      "E02",
						Name:             "Shaw-Howard U",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.912919,
						Longitude:        -77.022194,
						Address: StationAddress{
							Street: "1701 8th St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20001",
						},
					},
					{
						StationCode:      "E03",
						Name:             "U Street/African-Amer Civil War Memorial/Cardozo",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.916489,
						Longitude:        -77.028938,
						Address: StationAddress{
							Street: "1240 U Street NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20009",
						},
					},
					{
						StationCode:      "E04",
						Name:             "Columbia Heights",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.928672,
						Longitude:        -77.032775,
						Address: StationAddress{
							Street: "3030 14th St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20009",
						},
					},
					{
						StationCode:      "E05",
						Name:             "Georgia Ave-Petworth",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.936077,
						Longitude:        -77.024728,
						Address: StationAddress{
							Street: "3700 Georgia Avenue NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20010",
						},
					},
					{
						StationCode:      "E06",
						Name:             "Fort Totten",
						StationTogether1: "B06",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.951777,
						Longitude:        -77.002174,
						Address: StationAddress{
							Street: "550 Galloway Street NE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20011",
						},
					},
					{
						StationCode:      "E07",
						Name:             "West Hyattsville",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.954931,
						Longitude:        -76.969881,
						Address: StationAddress{
							Street: "2700 Hamilton St.",
							City:   "Hyattsville",
							State:  "MD",
							Zip:    "20782",
						},
					},
					{
						StationCode:      "E08",
						Name:             "Prince George's Plaza",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.965276,
						Longitude:        -76.956182,
						Address: StationAddress{
							Street: "3575 East West Highway",
							City:   "Hyattsville",
							State:  "MD",
							Zip:    "20782",
						},
					},
					{
						StationCode:      "E09",
						Name:             "College Park-U of Md",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.978523,
						Longitude:        -76.928432,
						Address: StationAddress{
							Street: "4931 Calvert Road",
							City:   "College Park",
							State:  "MD",
							Zip:    "20740",
						},
					},
					{
						StationCode:      "E10",
						Name:             "Greenbelt",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         39.011036,
						Longitude:        -76.911362,
						Address: StationAddress{
							Street: "5717 Greenbelt Metro Drive",
							City:   "Greenbelt",
							State:  "MD",
							Zip:    "20740",
						},
					},
					{
						StationCode:      "F01",
						Name:             "Gallery Pl-Chinatown",
						StationTogether1: "B01",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.89834,
						Longitude:        -77.021851,
						Address: StationAddress{
							Street: "630 H St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20001",
						},
					},
					{
						StationCode:      "F02",
						Name:             "Archives-Navy Memorial-Penn Quarter",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.893893,
						Longitude:        -77.021902,
						Address: StationAddress{
							Street: "701 Pennsylvania Avenue NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20004",
						},
					},
					{
						StationCode:      "F03",
						Name:             "L'Enfant Plaza",
						StationTogether1: "D03",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.884775,
						Longitude:        -77.021964,
						Address: StationAddress{
							Street: "600 Maryland Avenue SW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20024",
						},
					},
					{
						StationCode:      "F04",
						Name:             "Waterfront",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.876221,
						Longitude:        -77.017491,
						Address: StationAddress{
							Street: "399 M Street SW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20024",
						},
					},
					{
						StationCode:      "F05",
						Name:             "Navy Yard-Ballpark",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.876588,
						Longitude:        -77.005086,
						Address: StationAddress{
							Street: "200 M Street SE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20003",
						},
					},
					{
						StationCode:      "F06",
						Name:             "Anacostia",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.862072,
						Longitude:        -76.995648,
						Address: StationAddress{
							Street: "1101 Howard Road SE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20020",
						},
					},
					{
						StationCode:      "F07",
						Name:             "Congress Heights",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.845334,
						Longitude:        -76.98817,
						Address: StationAddress{
							Street: "1290 Alabama Avenue SE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20020",
						},
					},
					{
						StationCode:      "F08",
						Name:             "Southern Avenue",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.840974,
						Longitude:        -76.97536,
						Address: StationAddress{
							Street: "1411 Southern Avenue",
							City:   "Temple Hills",
							State:  "MD",
							Zip:    "20748",
						},
					},
					{
						StationCode:      "F09",
						Name:             "Naylor Road",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.851187,
						Longitude:        -76.956565,
						Address: StationAddress{
							Street: "3101 Branch Avenue",
							City:   "Temple Hills",
							State:  "MD",
							Zip:    "20748",
						},
					},
					{
						StationCode:      "F10",
						Name:             "Suitland",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.843891,
						Longitude:        -76.932022,
						Address: StationAddress{
							Street: "4500 Silver Hill Road",
							City:   "Suitland",
							State:  "MD",
							Zip:    "20746",
						},
					},
					{
						StationCode:      "F11",
						Name:             "Branch Ave",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.826995,
						Longitude:        -76.912134,
						Address: StationAddress{
							Street: "4704 Old Soper Road",
							City:   "Suitland",
							State:  "MD",
							Zip:    "20746",
						},
					},
				},
			},
		},
	},
	"/Rail.svc/json/jStationTimes": {
		{
			rawQuery:     "StationCode=F01",
			stringParam1: "F01",
			response:     `{"StationTimes":[{"Code":"F01","StationName":"Gallery Pl-Chinatown","Monday":{"OpeningTime":"05:15","FirstTrains":[{"Time":"05:25","DestinationStation":"E10"},{"Time":"05:26","DestinationStation":"F11"},{"Time":"05:32","DestinationStation":"C15"}],"LastTrains":[{"Time":"23:19","DestinationStation":"C15"},{"Time":"23:28","DestinationStation":"F11"},{"Time":"23:48","DestinationStation":"E10"}]},"Tuesday":{"OpeningTime":"05:15","FirstTrains":[{"Time":"05:25","DestinationStation":"E10"},{"Time":"05:26","DestinationStation":"F11"},{"Time":"05:32","DestinationStation":"C15"}],"LastTrains":[{"Time":"23:19","DestinationStation":"C15"},{"Time":"23:28","DestinationStation":"F11"},{"Time":"23:48","DestinationStation":"E10"}]},"Wednesday":{"OpeningTime":"05:15","FirstTrains":[{"Time":"05:25","DestinationStation":"E10"},{"Time":"05:26","DestinationStation":"F11"},{"Time":"05:32","DestinationStation":"C15"}],"LastTrains":[{"Time":"23:19","DestinationStation":"C15"},{"Time":"23:28","DestinationStation":"F11"},{"Time":"23:48","DestinationStation":"E10"}]},"Thursday":{"OpeningTime":"05:15","FirstTrains":[{"Time":"05:25","DestinationStation":"E10"},{"Time":"05:26","DestinationStation":"F11"},{"Time":"05:32","DestinationStation":"C15"}],"LastTrains":[{"Time":"23:19","DestinationStation":"C15"},{"Time":"23:28","DestinationStation":"F11"},{"Time":"23:48","DestinationStation":"E10"}]},"Friday":{"OpeningTime":"05:15","FirstTrains":[{"Time":"05:25","DestinationStation":"E10"},{"Time":"05:26","DestinationStation":"F11"},{"Time":"05:32","DestinationStation":"C15"}],"LastTrains":[{"Time":"00:49","DestinationStation":"C15"},{"Time":"00:58","DestinationStation":"F11"},{"Time":"01:18","DestinationStation":"E10"}]},"Saturday":{"OpeningTime":"07:15","FirstTrains":[{"Time":"07:25","DestinationStation":"E10"},{"Time":"07:26","DestinationStation":"F11"},{"Time":"07:32","DestinationStation":"C15"}],"LastTrains":[{"Time":"00:49","DestinationStation":"C15"},{"Time":"00:58","DestinationStation":"F11"},{"Time":"01:18","DestinationStation":"E10"}]},"Sunday":{"OpeningTime":"08:15","FirstTrains":[{"Time":"08:25","DestinationStation":"E10"},{"Time":"08:26","DestinationStation":"F11"},{"Time":"08:32","DestinationStation":"C15"}],"LastTrains":[{"Time":"22:49","DestinationStation":"C15"},{"Time":"22:58","DestinationStation":"F11"},{"Time":"23:18","DestinationStation":"E10"}]}}]}`,
			unmarshalledResponse: &GetStationTimingsResponse{
				StationTimes: []StationTime{
					{
						StationCode: "F01",
						StationName: "Gallery Pl-Chinatown",
						Monday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Tuesday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Wednesday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Thursday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Friday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "00:49",
									DestinationStation: "C15",
								},
								{
									Time:               "00:58",
									DestinationStation: "F11",
								},
								{
									Time:               "01:18",
									DestinationStation: "E10",
								},
							},
						},
						Saturday: StationDayItem{
							OpeningTime: "07:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "07:25",
									DestinationStation: "E10",
								},
								{
									Time:               "07:26",
									DestinationStation: "F11",
								},
								{
									Time:               "07:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "00:49",
									DestinationStation: "C15",
								},
								{
									Time:               "00:58",
									DestinationStation: "F11",
								},
								{
									Time:               "01:18",
									DestinationStation: "E10",
								},
							},
						},
						Sunday: StationDayItem{
							OpeningTime: "08:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "08:25",
									DestinationStation: "E10",
								},
								{
									Time:               "08:26",
									DestinationStation: "F11",
								},
								{
									Time:               "08:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "22:49",
									DestinationStation: "C15",
								},
								{
									Time:               "22:58",
									DestinationStation: "F11",
								},
								{
									Time:               "23:18",
									DestinationStation: "E10",
								},
							},
						},
					},
				},
			},
		},
	},
	"/Rail.svc/json/jSrcStationToDstStationInfo": {
		{
			rawQuery:     "FromStationCode=F01&ToStationCode=A07",
			stringParam1: "F01",
			stringParam2: "A07",
			response:     `{"StationToStationInfos":[{"SourceStation":"F01","DestinationStation":"A07","CompositeMiles":4.93,"RailTime":15,"RailFare":{"PeakTime":2.9,"OffPeakTime":2.45,"SeniorDisabled":1.45}}]}`,
			unmarshalledResponse: &GetStationToStationInformationResponse{
				StationToStationInformation: []StationToStation{
					{
						CompositeMiles:     4.93,
						DestinationStation: "A07",
						Fare: RailFare{
							OffPeakTime:    2.45,
							PeakTime:       2.9,
							SeniorDisabled: 1.45,
						},
						Time:          15,
						SourceStation: "F01",
					},
				},
			},
		},
	},
	"/Rail.svc/Lines": {
		testResponseData{
			rawQuery: "",
			response: `<LinesResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Lines><Line><DisplayName>Blue</DisplayName><EndStationCode>G05</EndStationCode><InternalDestination1/><InternalDestination2/><LineCode>BL</LineCode><StartStationCode>J03</StartStationCode></Line><Line><DisplayName>Green</DisplayName><EndStationCode>E10</EndStationCode><InternalDestination1/><InternalDestination2/><LineCode>GR</LineCode><StartStationCode>F11</StartStationCode></Line><Line><DisplayName>Orange</DisplayName><EndStationCode>D13</EndStationCode><InternalDestination1/><InternalDestination2/><LineCode>OR</LineCode><StartStationCode>K08</StartStationCode></Line><Line><DisplayName>Red</DisplayName><EndStationCode>B11</EndStationCode><InternalDestination1>A11</InternalDestination1><InternalDestination2>B08</InternalDestination2><LineCode>RD</LineCode><StartStationCode>A15</StartStationCode></Line><Line><DisplayName>Silver</DisplayName><EndStationCode>G05</EndStationCode><InternalDestination1/><InternalDestination2/><LineCode>SV</LineCode><StartStationCode>N06</StartStationCode></Line><Line><DisplayName>Yellow</DisplayName><EndStationCode>E06</EndStationCode><InternalDestination1>E01</InternalDestination1><InternalDestination2/><LineCode>YL</LineCode><StartStationCode>C15</StartStationCode></Line></Lines></LinesResp>`,
			unmarshalledResponse: &GetLinesResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "LinesResp",
				},
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
	"/Rail.svc/StationParking": {
		{
			rawQuery:     "StationCode=B08",
			stringParam1: "B08",
			response:     `<StationParkingResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><StationsParking><StationParking><Code>B08</Code><Notes>Parking is available at Montgomery County lots and garages.</Notes><AllDayParking><TotalCount>0</TotalCount><RiderCost i:nil="true"/><NonRiderCost i:nil="true"/><SaturdayRiderCost i:nil="true"/><SaturdayNonRiderCost i:nil="true"/></AllDayParking><ShortTermParking><TotalCount>0</TotalCount><Notes i:nil="true"/></ShortTermParking></StationParking></StationsParking></StationParkingResp>`,
			unmarshalledResponse: &GetParkingInformationResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "StationParkingResp",
				},
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
			rawQuery:     "StationCode=K08",
			stringParam1: "K08",
			response:     `<StationParkingResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><StationsParking><StationParking><Code>K08</Code><Notes>North Kiss &amp; Ride - 45 short term metered spaces. South Kiss &amp; Ride - 26 short term metered spaces.  101 spaces metered for 12-hr. max @ $1.00 per 60 mins. 17 spaces metered for 7-hr. max. @ $1.00 per 60 mins. Parking available from 8:30 AM to 2 AM.</Notes><AllDayParking><TotalCount>5169</TotalCount><RiderCost>4.95</RiderCost><NonRiderCost>4.95</NonRiderCost><SaturdayRiderCost>0</SaturdayRiderCost><SaturdayNonRiderCost>0</SaturdayNonRiderCost></AllDayParking><ShortTermParking><TotalCount>71</TotalCount><Notes>Parking available in section B between 8:30 AM - 3:30 PM and 7 PM - 2 AM, in section D between 10 AM - 2 PM.</Notes></ShortTermParking></StationParking></StationsParking></StationParkingResp>`,
			unmarshalledResponse: &GetParkingInformationResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "StationParkingResp",
				},
				ParkingInformation: []StationParking{
					{
						StationCode: "K08",
						Notes:       "North Kiss & Ride - 45 short term metered spaces. South Kiss & Ride - 26 short term metered spaces.  101 spaces metered for 12-hr. max @ $1.00 per 60 mins. 17 spaces metered for 7-hr. max. @ $1.00 per 60 mins. Parking available from 8:30 AM to 2 AM.",
						AllDay: AllDayParking{
							TotalCount:           5169,
							RiderCost:            4.95,
							NonRiderCost:         4.95,
							SaturdayRiderCost:    0,
							SaturdayNonRiderCost: 0,
						},
						ShortTerm: ShortTermParking{
							TotalCount: 71,
							Notes:      "Parking available in section B between 8:30 AM - 3:30 PM and 7 PM - 2 AM, in section D between 10 AM - 2 PM.",
						},
					},
				},
			},
		},
	},
	"/Rail.svc/Path": {
		{
			rawQuery:     "FromStationCode=A09&ToStationCode=B04",
			stringParam1: "A09",
			stringParam2: "B04",
			response:     `<PathResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Path><MetroPathItem><DistanceToPrev>0</DistanceToPrev><LineCode>RD</LineCode><SeqNum>1</SeqNum><StationCode>A09</StationCode><StationName>Bethesda</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>9095</DistanceToPrev><LineCode>RD</LineCode><SeqNum>2</SeqNum><StationCode>A08</StationCode><StationName>Friendship Heights</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>4135</DistanceToPrev><LineCode>RD</LineCode><SeqNum>3</SeqNum><StationCode>A07</StationCode><StationName>Tenleytown-AU</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>5841</DistanceToPrev><LineCode>RD</LineCode><SeqNum>4</SeqNum><StationCode>A06</StationCode><StationName>Van Ness-UDC</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>3320</DistanceToPrev><LineCode>RD</LineCode><SeqNum>5</SeqNum><StationCode>A05</StationCode><StationName>Cleveland Park</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>3740</DistanceToPrev><LineCode>RD</LineCode><SeqNum>6</SeqNum><StationCode>A04</StationCode><StationName>Woodley Park-Zoo/Adams Morgan</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>6304</DistanceToPrev><LineCode>RD</LineCode><SeqNum>7</SeqNum><StationCode>A03</StationCode><StationName>Dupont Circle</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>2711</DistanceToPrev><LineCode>RD</LineCode><SeqNum>8</SeqNum><StationCode>A02</StationCode><StationName>Farragut North</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>4178</DistanceToPrev><LineCode>RD</LineCode><SeqNum>9</SeqNum><StationCode>A01</StationCode><StationName>Metro Center</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>1505</DistanceToPrev><LineCode>RD</LineCode><SeqNum>10</SeqNum><StationCode>B01</StationCode><StationName>Gallery Pl-Chinatown</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>1875</DistanceToPrev><LineCode>RD</LineCode><SeqNum>11</SeqNum><StationCode>B02</StationCode><StationName>Judiciary Square</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>3446</DistanceToPrev><LineCode>RD</LineCode><SeqNum>12</SeqNum><StationCode>B03</StationCode><StationName>Union Station</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>3553</DistanceToPrev><LineCode>RD</LineCode><SeqNum>13</SeqNum><StationCode>B35</StationCode><StationName>NoMa-Gallaudet U</StationName></MetroPathItem><MetroPathItem><DistanceToPrev>5771</DistanceToPrev><LineCode>RD</LineCode><SeqNum>14</SeqNum><StationCode>B04</StationCode><StationName>Rhode Island Ave-Brentwood</StationName></MetroPathItem></Path></PathResp>`,
			unmarshalledResponse: &GetPathBetweenStationsResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "PathResp",
				},
				Path: []PathItem{
					{
						LineCode:                  "RD",
						StationCode:               "A09",
						StationName:               "Bethesda",
						SequenceNumber:            1,
						DistanceToPreviousStation: 0,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A08",
						StationName:               "Friendship Heights",
						SequenceNumber:            2,
						DistanceToPreviousStation: 9095,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A07",
						StationName:               "Tenleytown-AU",
						SequenceNumber:            3,
						DistanceToPreviousStation: 4135,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A06",
						StationName:               "Van Ness-UDC",
						SequenceNumber:            4,
						DistanceToPreviousStation: 5841,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A05",
						StationName:               "Cleveland Park",
						SequenceNumber:            5,
						DistanceToPreviousStation: 3320,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A04",
						StationName:               "Woodley Park-Zoo/Adams Morgan",
						SequenceNumber:            6,
						DistanceToPreviousStation: 3740,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A03",
						StationName:               "Dupont Circle",
						SequenceNumber:            7,
						DistanceToPreviousStation: 6304,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A02",
						StationName:               "Farragut North",
						SequenceNumber:            8,
						DistanceToPreviousStation: 2711,
					},
					{
						LineCode:                  "RD",
						StationCode:               "A01",
						StationName:               "Metro Center",
						SequenceNumber:            9,
						DistanceToPreviousStation: 4178,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B01",
						StationName:               "Gallery Pl-Chinatown",
						SequenceNumber:            10,
						DistanceToPreviousStation: 1505,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B02",
						StationName:               "Judiciary Square",
						SequenceNumber:            11,
						DistanceToPreviousStation: 1875,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B03",
						StationName:               "Union Station",
						SequenceNumber:            12,
						DistanceToPreviousStation: 3446,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B35",
						StationName:               "NoMa-Gallaudet U",
						SequenceNumber:            13,
						DistanceToPreviousStation: 3553,
					},
					{
						LineCode:                  "RD",
						StationCode:               "B04",
						StationName:               "Rhode Island Ave-Brentwood",
						SequenceNumber:            14,
						DistanceToPreviousStation: 5771,
					},
				},
			},
		},
	},
	"/Rail.svc/StationEntrances": {
		{
			rawQuery: "Lat=38.897383&Lon=-77.007262&Radius=500",
			requestType: GetStationEntrancesRequest{
				Latitude:  38.897383,
				Longitude: -77.007262,
				Radius:    500,
			},
			response: `<StationEntrancesResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Entrances><StationEntrance><Description>Station entrance from 1st St NE to southeast corner of the Union station building.</Description><ID>54</ID><Lat>38.897383</Lat><Lon>-77.007262</Lon><Name>SOUTH ENTRANCE (MASS AVE EXIT, NORTHEAST CORNER OF 1ST ST &amp; MASSACHUSETTS AVE)</Name><StationCode1>B03</StationCode1><StationCode2/></StationEntrance><StationEntrance><Description>Station entrance from northeast corner of Massachusetts Ave NE and 1st NE.</Description><ID>55</ID><Lat>38.89845</Lat><Lon>-77.007243</Lon><Name>NORTH ENTRANCE (1ST ST EXIT, WEST SIDE OF 1ST ST BETWEEN G ST AND MASSACHUSETTS AVE)</Name><StationCode1>B03</StationCode1><StationCode2/></StationEntrance><StationEntrance><Description>Escalator entrance from the passageway to  AMTRAK, MARC, VRE TRAINS</Description><ID>53</ID><Lat>38.898541</Lat><Lon>-77.006984</Lon><Name>ENTRANCE FROM AMTRAK, MARC, VRE TRAINS</Name><StationCode1>B03</StationCode1><StationCode2/></StationEntrance></Entrances></StationEntrancesResp>`,
			unmarshalledResponse: &GetStationEntrancesResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "StationEntrancesResp",
				},
				Entrances: []StationEntrance{
					{
						ID:           "54",
						Name:         "SOUTH ENTRANCE (MASS AVE EXIT, NORTHEAST CORNER OF 1ST ST & MASSACHUSETTS AVE)",
						StationCode1: "B03",
						StationCode2: "",
						Description:  "Station entrance from 1st St NE to southeast corner of the Union station building.",
						Latitude:     38.897383,
						Longitude:    -77.007262,
					},
					{
						ID:           "55",
						Name:         "NORTH ENTRANCE (1ST ST EXIT, WEST SIDE OF 1ST ST BETWEEN G ST AND MASSACHUSETTS AVE)",
						StationCode1: "B03",
						StationCode2: "",
						Description:  "Station entrance from northeast corner of Massachusetts Ave NE and 1st NE.",
						Latitude:     38.89845,
						Longitude:    -77.007243,
					},
					{
						ID:           "53",
						Name:         "ENTRANCE FROM AMTRAK, MARC, VRE TRAINS",
						StationCode1: "B03",
						StationCode2: "",
						Description:  "Escalator entrance from the passageway to  AMTRAK, MARC, VRE TRAINS",
						Latitude:     38.898541,
						Longitude:    -77.006984,
					},
				},
			},
		},
	},
	"/Rail.svc/StationInfo": {
		{
			rawQuery:     "StationCode=A07",
			stringParam1: "A07",
			response:     `<Station xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Address><City>Washington</City><State>DC</State><Street>4501 Wisconsin Avenue NW</Street><Zip>20016</Zip></Address><Code>A07</Code><Lat>38.947808</Lat><LineCode1>RD</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.079615</Lon><Name>Tenleytown-AU</Name><StationTogether1/><StationTogether2/></Station>`,
			unmarshalledResponse: &GetStationInformationResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "Station",
				},
				Address: StationAddress{
					Street: "4501 Wisconsin Avenue NW",
					City:   "Washington",
					State:  "DC",
					Zip:    "20016",
				},
				Latitude:         38.947808,
				LineCode1:        "RD",
				LineCode2:        "",
				Longitude:        -77.079615,
				Name:             "Tenleytown-AU",
				StationCode:      "A07",
				StationTogether1: "",
				StationTogether2: "",
			},
		},
		{
			rawQuery:     "StationCode=F01",
			stringParam1: "F01",
			response:     `<Station xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Address><City>Washington</City><State>DC</State><Street>630 H St. NW</Street><Zip>20001</Zip></Address><Code>F01</Code><Lat>38.89834</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.021851</Lon><Name>Gallery Pl-Chinatown</Name><StationTogether1>B01</StationTogether1><StationTogether2/></Station>`,
			unmarshalledResponse: &GetStationInformationResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "Station",
				},
				Address: StationAddress{
					Street: "630 H St. NW",
					City:   "Washington",
					State:  "DC",
					Zip:    "20001",
				},
				Latitude:         38.89834,
				LineCode1:        "GR",
				LineCode2:        "YL",
				Longitude:        -77.021851,
				Name:             "Gallery Pl-Chinatown",
				StationCode:      "F01",
				StationTogether1: "B01",
				StationTogether2: "",
			},
		},
	},
	"/Rail.svc/Stations": {
		{
			rawQuery:     "LineCode=GR",
			stringParam1: "GR",
			response:     `<StationsResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><Stations><Station><Address><City>Washington</City><State>DC</State><Street>700 M St. NW</Street><Zip>20001</Zip></Address><Code>E01</Code><Lat>38.905604</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.022256</Lon><Name>Mt Vernon Sq 7th St-Convention Center</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>1701 8th St. NW</Street><Zip>20001</Zip></Address><Code>E02</Code><Lat>38.912919</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.022194</Lon><Name>Shaw-Howard U</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>1240 U Street NW</Street><Zip>20009</Zip></Address><Code>E03</Code><Lat>38.916489</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.028938</Lon><Name>U Street/African-Amer Civil War Memorial/Cardozo</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>3030 14th St. NW</Street><Zip>20009</Zip></Address><Code>E04</Code><Lat>38.928672</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.032775</Lon><Name>Columbia Heights</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>3700 Georgia Avenue NW</Street><Zip>20010</Zip></Address><Code>E05</Code><Lat>38.936077</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.024728</Lon><Name>Georgia Ave-Petworth</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>550 Galloway Street NE</Street><Zip>20011</Zip></Address><Code>E06</Code><Lat>38.951777</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.002174</Lon><Name>Fort Totten</Name><StationTogether1>B06</StationTogether1><StationTogether2/></Station><Station><Address><City>Hyattsville</City><State>MD</State><Street>2700 Hamilton St.</Street><Zip>20782</Zip></Address><Code>E07</Code><Lat>38.954931</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.969881</Lon><Name>West Hyattsville</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Hyattsville</City><State>MD</State><Street>3575 East West Highway</Street><Zip>20782</Zip></Address><Code>E08</Code><Lat>38.965276</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.956182</Lon><Name>Prince George's Plaza</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>College Park</City><State>MD</State><Street>4931 Calvert Road</Street><Zip>20740</Zip></Address><Code>E09</Code><Lat>38.978523</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.928432</Lon><Name>College Park-U of Md</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Greenbelt</City><State>MD</State><Street>5717 Greenbelt Metro Drive</Street><Zip>20740</Zip></Address><Code>E10</Code><Lat>39.011036</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.911362</Lon><Name>Greenbelt</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>630 H St. NW</Street><Zip>20001</Zip></Address><Code>F01</Code><Lat>38.89834</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.021851</Lon><Name>Gallery Pl-Chinatown</Name><StationTogether1>B01</StationTogether1><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>701 Pennsylvania Avenue NW</Street><Zip>20004</Zip></Address><Code>F02</Code><Lat>38.893893</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.021902</Lon><Name>Archives-Navy Memorial-Penn Quarter</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>600 Maryland Avenue SW</Street><Zip>20024</Zip></Address><Code>F03</Code><Lat>38.884775</Lat><LineCode1>GR</LineCode1><LineCode2>YL</LineCode2><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.021964</Lon><Name>L'Enfant Plaza</Name><StationTogether1>D03</StationTogether1><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>399 M Street SW</Street><Zip>20024</Zip></Address><Code>F04</Code><Lat>38.876221</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.017491</Lon><Name>Waterfront</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>200 M Street SE</Street><Zip>20003</Zip></Address><Code>F05</Code><Lat>38.876588</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-77.005086</Lon><Name>Navy Yard-Ballpark</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>1101 Howard Road SE</Street><Zip>20020</Zip></Address><Code>F06</Code><Lat>38.862072</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.995648</Lon><Name>Anacostia</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Washington</City><State>DC</State><Street>1290 Alabama Avenue SE</Street><Zip>20020</Zip></Address><Code>F07</Code><Lat>38.845334</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.98817</Lon><Name>Congress Heights</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Temple Hills</City><State>MD</State><Street>1411 Southern Avenue</Street><Zip>20748</Zip></Address><Code>F08</Code><Lat>38.840974</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.97536</Lon><Name>Southern Avenue</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Temple Hills</City><State>MD</State><Street>3101 Branch Avenue</Street><Zip>20748</Zip></Address><Code>F09</Code><Lat>38.851187</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.956565</Lon><Name>Naylor Road</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Suitland</City><State>MD</State><Street>4500 Silver Hill Road</Street><Zip>20746</Zip></Address><Code>F10</Code><Lat>38.843891</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.932022</Lon><Name>Suitland</Name><StationTogether1/><StationTogether2/></Station><Station><Address><City>Suitland</City><State>MD</State><Street>4704 Old Soper Road</Street><Zip>20746</Zip></Address><Code>F11</Code><Lat>38.826995</Lat><LineCode1>GR</LineCode1><LineCode2 i:nil="true"/><LineCode3 i:nil="true"/><LineCode4 i:nil="true"/><Lon>-76.912134</Lon><Name>Branch Ave</Name><StationTogether1/><StationTogether2/></Station></Stations></StationsResp>`,
			unmarshalledResponse: &GetStationListResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "StationsResp",
				},
				Stations: []GetStationListResponseItem{
					{
						StationCode:      "E01",
						Name:             "Mt Vernon Sq 7th St-Convention Center",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.905604,
						Longitude:        -77.022256,
						Address: StationAddress{
							Street: "700 M St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20001",
						},
					},
					{
						StationCode:      "E02",
						Name:             "Shaw-Howard U",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.912919,
						Longitude:        -77.022194,
						Address: StationAddress{
							Street: "1701 8th St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20001",
						},
					},
					{
						StationCode:      "E03",
						Name:             "U Street/African-Amer Civil War Memorial/Cardozo",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.916489,
						Longitude:        -77.028938,
						Address: StationAddress{
							Street: "1240 U Street NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20009",
						},
					},
					{
						StationCode:      "E04",
						Name:             "Columbia Heights",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.928672,
						Longitude:        -77.032775,
						Address: StationAddress{
							Street: "3030 14th St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20009",
						},
					},
					{
						StationCode:      "E05",
						Name:             "Georgia Ave-Petworth",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.936077,
						Longitude:        -77.024728,
						Address: StationAddress{
							Street: "3700 Georgia Avenue NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20010",
						},
					},
					{
						StationCode:      "E06",
						Name:             "Fort Totten",
						StationTogether1: "B06",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.951777,
						Longitude:        -77.002174,
						Address: StationAddress{
							Street: "550 Galloway Street NE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20011",
						},
					},
					{
						StationCode:      "E07",
						Name:             "West Hyattsville",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.954931,
						Longitude:        -76.969881,
						Address: StationAddress{
							Street: "2700 Hamilton St.",
							City:   "Hyattsville",
							State:  "MD",
							Zip:    "20782",
						},
					},
					{
						StationCode:      "E08",
						Name:             "Prince George's Plaza",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.965276,
						Longitude:        -76.956182,
						Address: StationAddress{
							Street: "3575 East West Highway",
							City:   "Hyattsville",
							State:  "MD",
							Zip:    "20782",
						},
					},
					{
						StationCode:      "E09",
						Name:             "College Park-U of Md",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.978523,
						Longitude:        -76.928432,
						Address: StationAddress{
							Street: "4931 Calvert Road",
							City:   "College Park",
							State:  "MD",
							Zip:    "20740",
						},
					},
					{
						StationCode:      "E10",
						Name:             "Greenbelt",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         39.011036,
						Longitude:        -76.911362,
						Address: StationAddress{
							Street: "5717 Greenbelt Metro Drive",
							City:   "Greenbelt",
							State:  "MD",
							Zip:    "20740",
						},
					},
					{
						StationCode:      "F01",
						Name:             "Gallery Pl-Chinatown",
						StationTogether1: "B01",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.89834,
						Longitude:        -77.021851,
						Address: StationAddress{
							Street: "630 H St. NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20001",
						},
					},
					{
						StationCode:      "F02",
						Name:             "Archives-Navy Memorial-Penn Quarter",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.893893,
						Longitude:        -77.021902,
						Address: StationAddress{
							Street: "701 Pennsylvania Avenue NW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20004",
						},
					},
					{
						StationCode:      "F03",
						Name:             "L'Enfant Plaza",
						StationTogether1: "D03",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "YL",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.884775,
						Longitude:        -77.021964,
						Address: StationAddress{
							Street: "600 Maryland Avenue SW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20024",
						},
					},
					{
						StationCode:      "F04",
						Name:             "Waterfront",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.876221,
						Longitude:        -77.017491,
						Address: StationAddress{
							Street: "399 M Street SW",
							City:   "Washington",
							State:  "DC",
							Zip:    "20024",
						},
					},
					{
						StationCode:      "F05",
						Name:             "Navy Yard-Ballpark",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.876588,
						Longitude:        -77.005086,
						Address: StationAddress{
							Street: "200 M Street SE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20003",
						},
					},
					{
						StationCode:      "F06",
						Name:             "Anacostia",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.862072,
						Longitude:        -76.995648,
						Address: StationAddress{
							Street: "1101 Howard Road SE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20020",
						},
					},
					{
						StationCode:      "F07",
						Name:             "Congress Heights",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.845334,
						Longitude:        -76.98817,
						Address: StationAddress{
							Street: "1290 Alabama Avenue SE",
							City:   "Washington",
							State:  "DC",
							Zip:    "20020",
						},
					},
					{
						StationCode:      "F08",
						Name:             "Southern Avenue",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.840974,
						Longitude:        -76.97536,
						Address: StationAddress{
							Street: "1411 Southern Avenue",
							City:   "Temple Hills",
							State:  "MD",
							Zip:    "20748",
						},
					},
					{
						StationCode:      "F09",
						Name:             "Naylor Road",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.851187,
						Longitude:        -76.956565,
						Address: StationAddress{
							Street: "3101 Branch Avenue",
							City:   "Temple Hills",
							State:  "MD",
							Zip:    "20748",
						},
					},
					{
						StationCode:      "F10",
						Name:             "Suitland",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.843891,
						Longitude:        -76.932022,
						Address: StationAddress{
							Street: "4500 Silver Hill Road",
							City:   "Suitland",
							State:  "MD",
							Zip:    "20746",
						},
					},
					{
						StationCode:      "F11",
						Name:             "Branch Ave",
						StationTogether1: "",
						StationTogether2: "",
						LineCode1:        "GR",
						LineCode2:        "",
						LineCode3:        "",
						LineCode4:        "",
						Latitude:         38.826995,
						Longitude:        -76.912134,
						Address: StationAddress{
							Street: "4704 Old Soper Road",
							City:   "Suitland",
							State:  "MD",
							Zip:    "20746",
						},
					},
				},
			},
		},
	},
	"/Rail.svc/StationTimes": {
		{
			rawQuery:     "StationCode=F01",
			stringParam1: "F01",
			response:     `<StationTimeResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><StationTimes><StationTime><Code>F01</Code><StationName>Gallery Pl-Chinatown</StationName><Monday><OpeningTime>05:15</OpeningTime><FirstTrains><Train><Time>05:25</Time><DestinationStation>E10</DestinationStation></Train><Train><Time>05:26</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>05:32</Time><DestinationStation>C15</DestinationStation></Train></FirstTrains><LastTrains><Train><Time>23:19</Time><DestinationStation>C15</DestinationStation></Train><Train><Time>23:28</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>23:48</Time><DestinationStation>E10</DestinationStation></Train></LastTrains></Monday><Tuesday><OpeningTime>05:15</OpeningTime><FirstTrains><Train><Time>05:25</Time><DestinationStation>E10</DestinationStation></Train><Train><Time>05:26</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>05:32</Time><DestinationStation>C15</DestinationStation></Train></FirstTrains><LastTrains><Train><Time>23:19</Time><DestinationStation>C15</DestinationStation></Train><Train><Time>23:28</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>23:48</Time><DestinationStation>E10</DestinationStation></Train></LastTrains></Tuesday><Wednesday><OpeningTime>05:15</OpeningTime><FirstTrains><Train><Time>05:25</Time><DestinationStation>E10</DestinationStation></Train><Train><Time>05:26</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>05:32</Time><DestinationStation>C15</DestinationStation></Train></FirstTrains><LastTrains><Train><Time>23:19</Time><DestinationStation>C15</DestinationStation></Train><Train><Time>23:28</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>23:48</Time><DestinationStation>E10</DestinationStation></Train></LastTrains></Wednesday><Thursday><OpeningTime>05:15</OpeningTime><FirstTrains><Train><Time>05:25</Time><DestinationStation>E10</DestinationStation></Train><Train><Time>05:26</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>05:32</Time><DestinationStation>C15</DestinationStation></Train></FirstTrains><LastTrains><Train><Time>23:19</Time><DestinationStation>C15</DestinationStation></Train><Train><Time>23:28</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>23:48</Time><DestinationStation>E10</DestinationStation></Train></LastTrains></Thursday><Friday><OpeningTime>05:15</OpeningTime><FirstTrains><Train><Time>05:25</Time><DestinationStation>E10</DestinationStation></Train><Train><Time>05:26</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>05:32</Time><DestinationStation>C15</DestinationStation></Train></FirstTrains><LastTrains><Train><Time>00:49</Time><DestinationStation>C15</DestinationStation></Train><Train><Time>00:58</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>01:18</Time><DestinationStation>E10</DestinationStation></Train></LastTrains></Friday><Saturday><OpeningTime>07:15</OpeningTime><FirstTrains><Train><Time>07:25</Time><DestinationStation>E10</DestinationStation></Train><Train><Time>07:26</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>07:32</Time><DestinationStation>C15</DestinationStation></Train></FirstTrains><LastTrains><Train><Time>00:49</Time><DestinationStation>C15</DestinationStation></Train><Train><Time>00:58</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>01:18</Time><DestinationStation>E10</DestinationStation></Train></LastTrains></Saturday><Sunday><OpeningTime>08:15</OpeningTime><FirstTrains><Train><Time>08:25</Time><DestinationStation>E10</DestinationStation></Train><Train><Time>08:26</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>08:32</Time><DestinationStation>C15</DestinationStation></Train></FirstTrains><LastTrains><Train><Time>22:49</Time><DestinationStation>C15</DestinationStation></Train><Train><Time>22:58</Time><DestinationStation>F11</DestinationStation></Train><Train><Time>23:18</Time><DestinationStation>E10</DestinationStation></Train></LastTrains></Sunday></StationTime></StationTimes></StationTimeResp>`,
			unmarshalledResponse: &GetStationTimingsResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "StationTimeResp",
				},
				StationTimes: []StationTime{
					{
						StationCode: "F01",
						StationName: "Gallery Pl-Chinatown",
						Monday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Tuesday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Wednesday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Thursday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "23:19",
									DestinationStation: "C15",
								},
								{
									Time:               "23:28",
									DestinationStation: "F11",
								},
								{
									Time:               "23:48",
									DestinationStation: "E10",
								},
							},
						},
						Friday: StationDayItem{
							OpeningTime: "05:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "05:25",
									DestinationStation: "E10",
								},
								{
									Time:               "05:26",
									DestinationStation: "F11",
								},
								{
									Time:               "05:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "00:49",
									DestinationStation: "C15",
								},
								{
									Time:               "00:58",
									DestinationStation: "F11",
								},
								{
									Time:               "01:18",
									DestinationStation: "E10",
								},
							},
						},
						Saturday: StationDayItem{
							OpeningTime: "07:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "07:25",
									DestinationStation: "E10",
								},
								{
									Time:               "07:26",
									DestinationStation: "F11",
								},
								{
									Time:               "07:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "00:49",
									DestinationStation: "C15",
								},
								{
									Time:               "00:58",
									DestinationStation: "F11",
								},
								{
									Time:               "01:18",
									DestinationStation: "E10",
								},
							},
						},
						Sunday: StationDayItem{
							OpeningTime: "08:15",
							FirstTrains: []StationTrainInformation{
								{
									Time:               "08:25",
									DestinationStation: "E10",
								},
								{
									Time:               "08:26",
									DestinationStation: "F11",
								},
								{
									Time:               "08:32",
									DestinationStation: "C15",
								},
							},
							LastTrains: []StationTrainInformation{
								{
									Time:               "22:49",
									DestinationStation: "C15",
								},
								{
									Time:               "22:58",
									DestinationStation: "F11",
								},
								{
									Time:               "23:18",
									DestinationStation: "E10",
								},
							},
						},
					},
				},
			},
		},
	},
	"/Rail.svc/SrcStationToDstStationInfo": {
		{
			rawQuery:     "FromStationCode=F01&ToStationCode=A07",
			stringParam1: "F01",
			stringParam2: "A07",
			response:     `<StationToStationInfoResp xmlns="http://www.wmata.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><StationToStationInfos><StationToStationInfo><CompositeMiles>4.93</CompositeMiles><DestinationStation>A07</DestinationStation><RailFare><OffPeakTime>2.45</OffPeakTime><PeakTime>2.90</PeakTime><SeniorDisabled>1.45</SeniorDisabled></RailFare><RailTime>15</RailTime><SourceStation>F01</SourceStation></StationToStationInfo></StationToStationInfos></StationToStationInfoResp>`,
			unmarshalledResponse: &GetStationToStationInformationResponse{
				XMLName: xml.Name{
					Space: "http://www.wmata.com",
					Local: "StationToStationInfoResp",
				},
				StationToStationInformation: []StationToStation{
					{
						CompositeMiles:     4.93,
						DestinationStation: "A07",
						Fare: RailFare{
							OffPeakTime:    2.45,
							PeakTime:       2.9,
							SeniorDisabled: 1.45,
						},
						Time:          15,
						SourceStation: "F01",
					},
				},
			},
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

func TestGetLines(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jLines", "/Rail.svc/Lines"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testRequests, exist := testData[path]
		testService := setupTestService(responseFormats[i])

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
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}

func TestGetParkingInformation(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jStationParking", "/Rail.svc/StationParking"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetLines")
			return
		}

		for _, request := range testRequests {
			response, err := testService.GetParkingInformation(request.stringParam1)

			if err != nil {
				t.Errorf("error calling GetParkingInformation for station: %s Error: %s", request.stringParam1, err.Error())
				return
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}

func TestGetPathBetweenStations(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jPath", "/Rail.svc/Path"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetPathBetweenStations")
			return
		}

		for _, request := range testRequests {
			response, err := testService.GetPathBetweenStations(request.stringParam1, request.stringParam2)

			if err != nil {
				if request.expectedError != nil && request.expectedError.Error() == err.Error() {
					continue
				}

				t.Errorf("error calling GetPathBetweenStations for FromStation: %s ToStation: %s error: %s", request.stringParam1, request.stringParam2, err.Error())
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}

func TestGetStationEntrances(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jStationEntrances", "/Rail.svc/StationEntrances"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetStationEntrances")
			return
		}

		for _, request := range testRequests {
			getStationRequest := request.requestType.(GetStationEntrancesRequest)
			response, err := testService.GetStationEntrances(&getStationRequest)

			if err != nil {
				t.Errorf("error calling GetStationEntrances for request: %v error: %s", request, err.Error())
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}

func TestGetStationInformation(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jStationInfo", "/Rail.svc/StationInfo"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequest, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetStationInformation")
			return
		}

		for _, request := range testRequest {
			response, err := testService.GetStationInformation(request.stringParam1)

			if err != nil {
				if request.expectedError != nil && request.expectedError.Error() == err.Error() {
					continue
				}

				t.Errorf("error calling GetStationInfromation for station code: %s error: %s", request.stringParam1, err.Error())
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}

func TestGetStationList(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jStations", "/Rail.svc/Stations"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetStationList")
			return
		}

		for _, request := range testRequests {
			response, err := testService.GetStationList(request.stringParam1)

			if err != nil {
				t.Errorf("error calling GetStationList for line code: %s error: %s", request.stringParam1, err.Error())
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}

func TestGetStationTimings(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jStationTimes", "/Rail.svc/StationTimes"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetStationTimings")
			return
		}

		for _, request := range testRequests {
			response, err := testService.GetStationTimings(request.stringParam1)

			if err != nil {
				t.Errorf("error calling GetStationTimings for station: %s error: %s", request.stringParam1, err.Error())
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}

func TestGetStationToStationInformation(t *testing.T) {
	jsonAndXmlPaths := []string{"/Rail.svc/json/jSrcStationToDstStationInfo", "/Rail.svc/SrcStationToDstStationInfo"}
	responseFormats := []wmata.ResponseType{wmata.JSON, wmata.XML}

	for i, path := range jsonAndXmlPaths {
		testService := setupTestService(responseFormats[i])
		testRequests, exist := testData[path]

		if !exist {
			t.Errorf("no data found for GetStationToStationInformation")
			return
		}

		for _, request := range testRequests {
			response, err := testService.GetStationToStationInformation(request.stringParam1, request.stringParam2)

			if err != nil {
				t.Errorf("error calling GetStationToStationInformation, FromStation: %s ToStation: %s error: %s", request.stringParam1, request.stringParam2, err.Error())
			}

			if !reflect.DeepEqual(response, request.unmarshalledResponse) {
				t.Error(pretty.Diff(response, request.unmarshalledResponse))
			}
		}
	}
}
