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
			rawQuery:     "StationCode=K08",
			stringParam1: "K08",
			jsonResponse: `{"StationsParking":[{"Code":"K08","Notes":"North Kiss & Ride - 45 short term metered spaces. South Kiss & Ride - 26 short term metered spaces.  101 spaces metered for 12-hr. max @ $1.00 per 60 mins. 17 spaces metered for 7-hr. max. @ $1.00 per 60 mins. Parking available from 8:30 AM to 2 AM.","AllDayParking":{"TotalCount":5169,"RiderCost":4.95,"NonRiderCost":4.95,"SaturdayRiderCost":0,"SaturdayNonRiderCost":0},"ShortTermParking":{"TotalCount":71,"Notes":"Parking available in section B between 8:30 AM - 3:30 PM and 7 PM - 2 AM, in section D between 10 AM - 2 PM."}}]}`,
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
			jsonResponse: `{"Path":[{"LineCode":"RD","StationCode":"A09","StationName":"Bethesda","SeqNum":1,"DistanceToPrev":0},{"LineCode":"RD","StationCode":"A08","StationName":"Friendship Heights","SeqNum":2,"DistanceToPrev":9095},{"LineCode":"RD","StationCode":"A07","StationName":"Tenleytown-AU","SeqNum":3,"DistanceToPrev":4135},{"LineCode":"RD","StationCode":"A06","StationName":"Van Ness-UDC","SeqNum":4,"DistanceToPrev":5841},{"LineCode":"RD","StationCode":"A05","StationName":"Cleveland Park","SeqNum":5,"DistanceToPrev":3320},{"LineCode":"RD","StationCode":"A04","StationName":"Woodley Park-Zoo/Adams Morgan","SeqNum":6,"DistanceToPrev":3740},{"LineCode":"RD","StationCode":"A03","StationName":"Dupont Circle","SeqNum":7,"DistanceToPrev":6304},{"LineCode":"RD","StationCode":"A02","StationName":"Farragut North","SeqNum":8,"DistanceToPrev":2711},{"LineCode":"RD","StationCode":"A01","StationName":"Metro Center","SeqNum":9,"DistanceToPrev":4178},{"LineCode":"RD","StationCode":"B01","StationName":"Gallery Pl-Chinatown","SeqNum":10,"DistanceToPrev":1505},{"LineCode":"RD","StationCode":"B02","StationName":"Judiciary Square","SeqNum":11,"DistanceToPrev":1875},{"LineCode":"RD","StationCode":"B03","StationName":"Union Station","SeqNum":12,"DistanceToPrev":3446},{"LineCode":"RD","StationCode":"B35","StationName":"NoMa-Gallaudet U","SeqNum":13,"DistanceToPrev":3553},{"LineCode":"RD","StationCode":"B04","StationName":"Rhode Island Ave-Brentwood","SeqNum":14,"DistanceToPrev":5771}]}`,
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
	},
	"/Rail.svc/json/jStationEntrances": {
		{
			rawQuery: "Lat=38.897383&Lon=-77.007262&Radius=500",
			requestType: GetStationEntrancesRequest{
				latitude:  38.897383,
				longitude: -77.007262,
				radius:    500,
			},
			jsonResponse: `{"Entrances":[{"ID":"54","Name":"SOUTH ENTRANCE (MASS AVE EXIT, NORTHEAST CORNER OF 1ST ST & MASSACHUSETTS AVE)","StationCode1":"B03","StationCode2":"","Description":"Station entrance from 1st St NE to southeast corner of the Union station building.","Lat":38.897383,"Lon":-77.007262},{"ID":"55","Name":"NORTH ENTRANCE (1ST ST EXIT, WEST SIDE OF 1ST ST BETWEEN G ST AND MASSACHUSETTS AVE)","StationCode1":"B03","StationCode2":"","Description":"Station entrance from northeast corner of Massachusetts Ave NE and 1st NE.","Lat":38.89845,"Lon":-77.007243},{"ID":"53","Name":"ENTRANCE FROM AMTRAK, MARC, VRE TRAINS","StationCode1":"B03","StationCode2":"","Description":"Escalator entrance from the passageway to  AMTRAK, MARC, VRE TRAINS","Lat":38.898541,"Lon":-77.006984}]}`,
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
	"/Rail.svc/json/jStations": {
		{
			rawQuery:     "LineCode=GR",
			stringParam1: "GR",
			jsonResponse: `{"Stations":[{"Code":"E01","Name":"Mt Vernon Sq 7th St-Convention Center","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.905604,"Lon":-77.022256,"Address":{"Street":"700 M St. NW","City":"Washington","State":"DC","Zip":"20001"}},{"Code":"E02","Name":"Shaw-Howard U","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.912919,"Lon":-77.022194,"Address":{"Street":"1701 8th St. NW","City":"Washington","State":"DC","Zip":"20001"}},{"Code":"E03","Name":"U Street/African-Amer Civil War Memorial/Cardozo","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.916489,"Lon":-77.028938,"Address":{"Street":"1240 U Street NW","City":"Washington","State":"DC","Zip":"20009"}},{"Code":"E04","Name":"Columbia Heights","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.928672,"Lon":-77.032775,"Address":{"Street":"3030 14th St. NW","City":"Washington","State":"DC","Zip":"20009"}},{"Code":"E05","Name":"Georgia Ave-Petworth","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.936077,"Lon":-77.024728,"Address":{"Street":"3700 Georgia Avenue NW","City":"Washington","State":"DC","Zip":"20010"}},{"Code":"E06","Name":"Fort Totten","StationTogether1":"B06","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.951777,"Lon":-77.002174,"Address":{"Street":"550 Galloway Street NE","City":"Washington","State":"DC","Zip":"20011"}},{"Code":"E07","Name":"West Hyattsville","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.954931,"Lon":-76.969881,"Address":{"Street":"2700 Hamilton St.","City":"Hyattsville","State":"MD","Zip":"20782"}},{"Code":"E08","Name":"Prince George's Plaza","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.965276,"Lon":-76.956182,"Address":{"Street":"3575 East West Highway","City":"Hyattsville","State":"MD","Zip":"20782"}},{"Code":"E09","Name":"College Park-U of Md","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.978523,"Lon":-76.928432,"Address":{"Street":"4931 Calvert Road","City":"College Park","State":"MD","Zip":"20740"}},{"Code":"E10","Name":"Greenbelt","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":39.011036,"Lon":-76.911362,"Address":{"Street":"5717 Greenbelt Metro Drive","City":"Greenbelt","State":"MD","Zip":"20740"}},{"Code":"F01","Name":"Gallery Pl-Chinatown","StationTogether1":"B01","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.89834,"Lon":-77.021851,"Address":{"Street":"630 H St. NW","City":"Washington","State":"DC","Zip":"20001"}},{"Code":"F02","Name":"Archives-Navy Memorial-Penn Quarter","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.893893,"Lon":-77.021902,"Address":{"Street":"701 Pennsylvania Avenue NW","City":"Washington","State":"DC","Zip":"20004"}},{"Code":"F03","Name":"L'Enfant Plaza","StationTogether1":"D03","StationTogether2":"","LineCode1":"GR","LineCode2":"YL","LineCode3":null,"LineCode4":null,"Lat":38.884775,"Lon":-77.021964,"Address":{"Street":"600 Maryland Avenue SW","City":"Washington","State":"DC","Zip":"20024"}},{"Code":"F04","Name":"Waterfront","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.876221,"Lon":-77.017491,"Address":{"Street":"399 M Street SW","City":"Washington","State":"DC","Zip":"20024"}},{"Code":"F05","Name":"Navy Yard-Ballpark","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.876588,"Lon":-77.005086,"Address":{"Street":"200 M Street SE","City":"Washington","State":"DC","Zip":"20003"}},{"Code":"F06","Name":"Anacostia","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.862072,"Lon":-76.995648,"Address":{"Street":"1101 Howard Road SE","City":"Washington","State":"DC","Zip":"20020"}},{"Code":"F07","Name":"Congress Heights","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.845334,"Lon":-76.98817,"Address":{"Street":"1290 Alabama Avenue SE","City":"Washington","State":"DC","Zip":"20020"}},{"Code":"F08","Name":"Southern Avenue","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.840974,"Lon":-76.97536,"Address":{"Street":"1411 Southern Avenue","City":"Temple Hills","State":"MD","Zip":"20748"}},{"Code":"F09","Name":"Naylor Road","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.851187,"Lon":-76.956565,"Address":{"Street":"3101 Branch Avenue","City":"Temple Hills","State":"MD","Zip":"20748"}},{"Code":"F10","Name":"Suitland","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.843891,"Lon":-76.932022,"Address":{"Street":"4500 Silver Hill Road","City":"Suitland","State":"MD","Zip":"20746"}},{"Code":"F11","Name":"Branch Ave","StationTogether1":"","StationTogether2":"","LineCode1":"GR","LineCode2":null,"LineCode3":null,"LineCode4":null,"Lat":38.826995,"Lon":-76.912134,"Address":{"Street":"4704 Old Soper Road","City":"Suitland","State":"MD","Zip":"20746"}}]}`,
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
			t.Errorf("error calling GetParkingInformation for station: %s Error: %s", request.stringParam1, err.Error())
			return
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Errorf("unexpected response. Expected: %v but got: %v", response, request.unmarshalledResponse)
		}
	}

}

func TestGetPathBetweenStations(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Rail.svc/json/jPath"]

	if !exist {
		t.Errorf("no data found for GetPathBetweenStations")
		return
	}

	for _, request := range testRequests {
		response, err := testService.GetPathBetweenStations(request.stringParam1, request.stringParam2)

		if err != nil {
			t.Errorf("error calling GetPathBetweenStations for FromStation: %s ToStation: %s error: %s", request.stringParam1, request.stringParam2, err.Error())
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Errorf("unexpected response. Expected: %v but got: %v", response, request.unmarshalledResponse)
		}
	}
}

func TestGetStationEntrances(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Rail.svc/json/jStationEntrances"]

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
			t.Errorf("unexpected response. Expected: %v but got: %v", response, request.unmarshalledResponse)
		}

	}

}

func TestGetStationList(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Rail.svc/json/jStations"]

	if !exist {
		t.Errorf("no data found for GetStationList")
		return
	}

	for _, request := range testRequests {
		response, err := testService.GetStationList(request.stringParam1)

		if err != nil {
			t.Errorf("error calling GetstationList for line code: %s error: %s", request.stringParam1, err.Error())
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Errorf("unexpected response. Expected: %v but got: %v", response, request.unmarshalledResponse)
		}

	}

}
