package incidents

import (
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
			_, writeErr := rr.Write([]byte(response.jsonResponse))

			return rr.Result(), writeErr
		}
	}

	return nil, errors.New("no data found")
}

type testResponseData struct {
	rawQuery             string
	param                string
	jsonResponse         string
	unmarshalledResponse interface{}
}

var testData = map[string][]testResponseData{
	"/Incidents.svc/json/BusIncidents": {
		{
			rawQuery:     "Route=R2",
			param:        "R2",
			jsonResponse: `{"BusIncidents":[{"IncidentID":"5F61BD99-1DF7-48DD-95C1-F83DBDDE713D","IncidentType":"Delay","RoutesAffected":["R2"],"Description":"Due to traffic congestion on Riggs Rd btwn Adelphi Rd and Knowllwood Dr, buses are experiencing delays.","DateUpdated":"2019-04-05T16:04:27"}]}`,
			unmarshalledResponse: &GetBusIncidentsResponse{
				BusIncidents: []BusIncident{
					{
						IncidentID:   "5F61BD99-1DF7-48DD-95C1-F83DBDDE713D",
						IncidentType: "Delay",
						RoutesAffected: []string{
							"R2",
						},
						Description: "Due to traffic congestion on Riggs Rd btwn Adelphi Rd and Knowllwood Dr, buses are experiencing delays.",
						DateUpdated: "2019-04-05T16:04:27",
					},
				},
			},
		},
	},
	"/Incidents.svc/json/ElevatorIncidents": {
		{
			rawQuery:     "StationCode=",
			param:        "",
			jsonResponse: `{"ElevatorIncidents":[{"UnitName":"A01E06","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"A01","StationName":"Metro Center, G and 11th St Entrance","LocationDescription":"Escalator between mezzanine and platform to Shady Grove","SymptomCode":null,"TimeOutOfService":"1607","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T16:07:00","DateUpdated":"2019-04-05T16:08:59","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"A07X01","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"A07","StationName":"Tenleytown-AU","LocationDescription":"Elevator between street and platform, east side of Wisconsin Avenue","SymptomCode":null,"TimeOutOfService":"1517","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T15:17:00","DateUpdated":"2019-04-05T16:11:19","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"A07X05","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"A07","StationName":"Tenleytown-AU","LocationDescription":"Escalator between middle landing and mezzanine","SymptomCode":null,"TimeOutOfService":"0712","SymptomDescription":"Major Repair","DisplayOrder":0,"DateOutOfServ":"2019-03-26T07:12:00","DateUpdated":"2019-04-05T10:39:15","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"A11X01","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"A11","StationName":"Grosvenor-Strathmore","LocationDescription":"Escalator between street and mezzanine (under Rockville Pike)","SymptomCode":null,"TimeOutOfService":"0714","SymptomDescription":"Preventive Maintenance Inspection","DisplayOrder":0,"DateOutOfServ":"2019-04-04T07:14:00","DateUpdated":"2019-04-04T14:37:16","EstimatedReturnToService":"2019-04-06T23:59:59"},{"UnitName":"B02N01","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"B02","StationName":"Judiciary Square, Building Museum Entrance","LocationDescription":"Escalator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"0434","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2018-11-19T04:34:00","DateUpdated":"2019-03-25T15:33:24","EstimatedReturnToService":"2019-04-17T23:59:59"},{"UnitName":"B02N03","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"B02","StationName":"Judiciary Square, Building Museum Entrance","LocationDescription":"Escalator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"0437","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2018-11-19T04:37:00","DateUpdated":"2019-02-09T08:12:58","EstimatedReturnToService":"2019-04-17T23:59:59"},{"UnitName":"B04X02","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"B04","StationName":"Rhode Island Ave-Brentwood","LocationDescription":"Garage elevator","SymptomCode":null,"TimeOutOfService":"0814","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T08:14:00","DateUpdated":"2019-04-05T16:20:28","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"B08N05","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"B08","StationName":"Silver Spring, North Side/Colesville Road Entrance","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"1255","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T12:55:00","DateUpdated":"2019-04-05T14:48:58","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"B08S01","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"B08","StationName":"Silver Spring, South Side/Bus Terminal Entrance","LocationDescription":"Elevator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"1615","SymptomDescription":"Customer Incident","DisplayOrder":0,"DateOutOfServ":"2019-04-05T16:15:00","DateUpdated":"2019-04-05T16:18:14","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"B35N02","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"B35","StationName":"NoMa-Gallaudet, Florida Avenue Entrance","LocationDescription":"Elevator between bike trail and mezzanine","SymptomCode":null,"TimeOutOfService":"0502","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2019-03-22T05:02:00","DateUpdated":"2019-03-22T05:04:58","EstimatedReturnToService":"2019-06-22T23:59:59"},{"UnitName":"C03W04","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"C03","StationName":"Farragut West, 18th Street Entrance","LocationDescription":"Escalator between mezzanine and platform to Vienna/Franconia-Springfield","SymptomCode":null,"TimeOutOfService":"1547","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T15:47:00","DateUpdated":"2019-04-05T15:48:31","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"C06X01","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"C06","StationName":"Arlington Cemetery","LocationDescription":"Escalator between street and mezzanine north side of Memorial Drive","SymptomCode":null,"TimeOutOfService":"0625","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T06:25:00","DateUpdated":"2019-04-05T13:14:05","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"C09X06","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"C09","StationName":"Crystal City","LocationDescription":"Escalator between middle landing and mezzanine","SymptomCode":null,"TimeOutOfService":"0944","SymptomDescription":"Minor Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-05T09:44:00","DateUpdated":"2019-04-05T13:39:53","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"C14X01","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"C14","StationName":"Eisenhower Avenue","LocationDescription":"Escalator between mezzanine and platform to Mt. Vernon Sq","SymptomCode":null,"TimeOutOfService":"1355","SymptomDescription":"Minor Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-04T13:55:00","DateUpdated":"2019-04-04T17:31:27","EstimatedReturnToService":"2019-04-06T23:59:59"},{"UnitName":"C15N05","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"C15","StationName":"Huntington, Huntington Ave. Entrance","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"0857","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2019-03-01T08:57:00","DateUpdated":"2019-03-11T04:43:54","EstimatedReturnToService":"2019-04-24T23:59:59"},{"UnitName":"C15N06","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"C15","StationName":"Huntington, Huntington Ave. Entrance","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"0415","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2019-01-07T04:15:00","DateUpdated":"2019-02-09T08:12:58","EstimatedReturnToService":"2019-04-17T23:59:59"},{"UnitName":"D02S02","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"D02","StationName":"Smithsonian, Independence Avenue Entrance","LocationDescription":"Escalator between street mezzanine","SymptomCode":null,"TimeOutOfService":"1512","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T15:12:00","DateUpdated":"2019-04-05T15:13:17","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"D04X01","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"D04","StationName":"Federal Center SW","LocationDescription":"Elevator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"1449","SymptomDescription":"Power Outage","DisplayOrder":0,"DateOutOfServ":"2019-04-02T14:49:00","DateUpdated":"2019-04-05T02:33:45","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"D04X02","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"D04","StationName":"Federal Center SW","LocationDescription":"Escalator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"1338","SymptomDescription":"Major Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-02T13:38:00","DateUpdated":"2019-04-04T23:49:31","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"D07X03","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"D07","StationName":"Potomac Ave","LocationDescription":"Escalator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"1504","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T15:04:00","DateUpdated":"2019-04-05T15:05:58","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"D09X04","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"D09","StationName":"Minnesota Ave","LocationDescription":"Garage Elevator","SymptomCode":null,"TimeOutOfService":"1616","SymptomDescription":"Fire Alarm","DisplayOrder":0,"DateOutOfServ":"2019-04-01T16:16:00","DateUpdated":"2019-04-05T15:27:49","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"D10X02","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"D10","StationName":"Deanwood","LocationDescription":"Elevator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"0737","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T07:37:00","DateUpdated":"2019-04-05T09:49:54","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"D10X03","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"D10","StationName":"Deanwood","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"0508","SymptomDescription":"Minor Repair","DisplayOrder":0,"DateOutOfServ":"2019-03-29T05:08:00","DateUpdated":"2019-04-03T14:27:56","EstimatedReturnToService":"2019-04-08T23:59:59"},{"UnitName":"D13X04","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"D13","StationName":"New Carrollton","LocationDescription":"Garage elevator","SymptomCode":null,"TimeOutOfService":"2031","SymptomDescription":"Minor Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-03T20:31:00","DateUpdated":"2019-04-04T14:55:19","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"D13X05","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"D13","StationName":"New Carrollton","LocationDescription":"Garage elevator","SymptomCode":null,"TimeOutOfService":"2032","SymptomDescription":"Minor Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-03T20:32:00","DateUpdated":"2019-04-04T19:52:56","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"E07X04","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"E07","StationName":"West Hyattsville","LocationDescription":"Escalator between mezzanine and platform to Branch Ave","SymptomCode":null,"TimeOutOfService":"1859","SymptomDescription":"Minor Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-04T18:59:00","DateUpdated":"2019-04-05T16:19:08","EstimatedReturnToService":"2019-04-06T23:59:59"},{"UnitName":"F02X02","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"F02","StationName":"Archives-Navy Memorial-Penn Quarter","LocationDescription":"Escalator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"0358","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2019-02-20T03:58:00","DateUpdated":"2019-02-20T03:59:27","EstimatedReturnToService":"2019-05-20T23:59:59"},{"UnitName":"F05W03","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"F05","StationName":"Navy Yard-Ballpark, Half Street Entrance","LocationDescription":"Escalator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"0246","SymptomDescription":"Major Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-05T02:46:00","DateUpdated":"2019-04-05T15:18:06","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"F06S02","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"F06","StationName":"Anacostia, Howard Road Entrance","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"1206","SymptomDescription":"Minor Repair","DisplayOrder":0,"DateOutOfServ":"2019-03-28T12:06:00","DateUpdated":"2019-04-02T06:25:08","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"F07X04","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"F07","StationName":"Congress Heights","LocationDescription":"Escalator between mezzanine and Alabama Avenue","SymptomCode":null,"TimeOutOfService":"0456","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2018-10-22T04:56:00","DateUpdated":"2019-04-01T06:49:31","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"F07X05","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"F07","StationName":"Congress Heights","LocationDescription":"Escalator between mezzanine and Alabama Avenue","SymptomCode":null,"TimeOutOfService":"0459","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2018-10-22T04:59:00","DateUpdated":"2019-04-01T06:49:47","EstimatedReturnToService":"2019-04-05T23:59:59"},{"UnitName":"J02X01","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"J02","StationName":"Van Dorn Street","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"0425","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2019-01-15T04:25:00","DateUpdated":"2019-02-09T08:12:58","EstimatedReturnToService":"2019-05-24T23:59:59"},{"UnitName":"J02X02","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"J02","StationName":"Van Dorn Street","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"0255","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2019-02-05T02:55:00","DateUpdated":"2019-02-09T08:12:58","EstimatedReturnToService":"2019-05-24T23:59:59"},{"UnitName":"K01X01","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"K01","StationName":"Court House","LocationDescription":"Elevator between street and mezzanine","SymptomCode":null,"TimeOutOfService":"2040","SymptomDescription":"Major Repair","DisplayOrder":0,"DateOutOfServ":"2019-04-04T20:40:00","DateUpdated":"2019-04-05T09:48:01","EstimatedReturnToService":"2019-04-06T23:59:59"},{"UnitName":"K01X02","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"K01","StationName":"Court House","LocationDescription":"Escalator between street (Wilson Blvd) and middle landing/tunnel","SymptomCode":null,"TimeOutOfService":"0541","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2018-09-14T05:41:00","DateUpdated":"2019-03-25T02:47:13","EstimatedReturnToService":"2019-04-06T23:59:59"},{"UnitName":"K01X05","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"K01","StationName":"Court House","LocationDescription":"Escalator between middle landing and mezzanine","SymptomCode":null,"TimeOutOfService":"0402","SymptomDescription":"Modernization","DisplayOrder":0,"DateOutOfServ":"2019-03-15T04:02:00","DateUpdated":"2019-03-25T02:47:44","EstimatedReturnToService":"2019-04-06T23:59:59"},{"UnitName":"K02X01","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"K02","StationName":"Clarendon","LocationDescription":"Escalator between street and middle landing/tunnel","SymptomCode":null,"TimeOutOfService":"1010","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T10:10:00","DateUpdated":"2019-04-05T13:14:43","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"K08X02","UnitType":"ESCALATOR","UnitStatus":null,"StationCode":"K08","StationName":"Vienna/Fairfax-GMU","LocationDescription":"Escalator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"1322","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T13:22:00","DateUpdated":"2019-04-05T13:24:30","EstimatedReturnToService":"2019-04-07T23:59:59"},{"UnitName":"N06X03","UnitType":"ELEVATOR","UnitStatus":null,"StationCode":"N06","StationName":"Wiehle-Reston East","LocationDescription":"Elevator between mezzanine and platform","SymptomCode":null,"TimeOutOfService":"1325","SymptomDescription":"Service Call","DisplayOrder":0,"DateOutOfServ":"2019-04-05T13:25:00","DateUpdated":"2019-04-05T13:27:02","EstimatedReturnToService":"2019-04-07T23:59:59"}]}`,
			unmarshalledResponse: &GetElevatorEscalatorOutagesResponse{
				ElevatorIncidents: []ElevatorIncident{
					{
						UnitName:                 "A01E06",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "A01",
						StationName:              "Metro Center, G and 11th St Entrance",
						LocationDescription:      "Escalator between mezzanine and platform to Shady Grove",
						SymptomCode:              "",
						TimeOutOfService:         "1607",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T16:07:00",
						DateUpdated:              "2019-04-05T16:08:59",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "A07X01",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "A07",
						StationName:              "Tenleytown-AU",
						LocationDescription:      "Elevator between street and platform, east side of Wisconsin Avenue",
						SymptomCode:              "",
						TimeOutOfService:         "1517",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T15:17:00",
						DateUpdated:              "2019-04-05T16:11:19",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "A07X05",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "A07",
						StationName:              "Tenleytown-AU",
						LocationDescription:      "Escalator between middle landing and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0712",
						SymptomDescription:       "Major Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-03-26T07:12:00",
						DateUpdated:              "2019-04-05T10:39:15",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "A11X01",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "A11",
						StationName:              "Grosvenor-Strathmore",
						LocationDescription:      "Escalator between street and mezzanine (under Rockville Pike)",
						SymptomCode:              "",
						TimeOutOfService:         "0714",
						SymptomDescription:       "Preventive Maintenance Inspection",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-04T07:14:00",
						DateUpdated:              "2019-04-04T14:37:16",
						EstimatedReturnToService: "2019-04-06T23:59:59",
					},
					{
						UnitName:                 "B02N01",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "B02",
						StationName:              "Judiciary Square, Building Museum Entrance",
						LocationDescription:      "Escalator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0434",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2018-11-19T04:34:00",
						DateUpdated:              "2019-03-25T15:33:24",
						EstimatedReturnToService: "2019-04-17T23:59:59",
					},
					{
						UnitName:                 "B02N03",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "B02",
						StationName:              "Judiciary Square, Building Museum Entrance",
						LocationDescription:      "Escalator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0437",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2018-11-19T04:37:00",
						DateUpdated:              "2019-02-09T08:12:58",
						EstimatedReturnToService: "2019-04-17T23:59:59",
					},
					{
						UnitName:                 "B04X02",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "B04",
						StationName:              "Rhode Island Ave-Brentwood",
						LocationDescription:      "Garage elevator",
						SymptomCode:              "",
						TimeOutOfService:         "0814",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T08:14:00",
						DateUpdated:              "2019-04-05T16:20:28",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "B08N05",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "B08",
						StationName:              "Silver Spring, North Side/Colesville Road Entrance",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "1255",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T12:55:00",
						DateUpdated:              "2019-04-05T14:48:58",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "B08S01",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "B08",
						StationName:              "Silver Spring, South Side/Bus Terminal Entrance",
						LocationDescription:      "Elevator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "1615",
						SymptomDescription:       "Customer Incident",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T16:15:00",
						DateUpdated:              "2019-04-05T16:18:14",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "B35N02",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "B35",
						StationName:              "NoMa-Gallaudet, Florida Avenue Entrance",
						LocationDescription:      "Elevator between bike trail and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0502",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2019-03-22T05:02:00",
						DateUpdated:              "2019-03-22T05:04:58",
						EstimatedReturnToService: "2019-06-22T23:59:59",
					},
					{
						UnitName:                 "C03W04",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "C03",
						StationName:              "Farragut West, 18th Street Entrance",
						LocationDescription:      "Escalator between mezzanine and platform to Vienna/Franconia-Springfield",
						SymptomCode:              "",
						TimeOutOfService:         "1547",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T15:47:00",
						DateUpdated:              "2019-04-05T15:48:31",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "C06X01",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "C06",
						StationName:              "Arlington Cemetery",
						LocationDescription:      "Escalator between street and mezzanine north side of Memorial Drive",
						SymptomCode:              "",
						TimeOutOfService:         "0625",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T06:25:00",
						DateUpdated:              "2019-04-05T13:14:05",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "C09X06",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "C09",
						StationName:              "Crystal City",
						LocationDescription:      "Escalator between middle landing and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0944",
						SymptomDescription:       "Minor Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T09:44:00",
						DateUpdated:              "2019-04-05T13:39:53",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "C14X01",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "C14",
						StationName:              "Eisenhower Avenue",
						LocationDescription:      "Escalator between mezzanine and platform to Mt. Vernon Sq",
						SymptomCode:              "",
						TimeOutOfService:         "1355",
						SymptomDescription:       "Minor Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-04T13:55:00",
						DateUpdated:              "2019-04-04T17:31:27",
						EstimatedReturnToService: "2019-04-06T23:59:59",
					},
					{
						UnitName:                 "C15N05",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "C15",
						StationName:              "Huntington, Huntington Ave. Entrance",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "0857",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2019-03-01T08:57:00",
						DateUpdated:              "2019-03-11T04:43:54",
						EstimatedReturnToService: "2019-04-24T23:59:59",
					},
					{
						UnitName:                 "C15N06",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "C15",
						StationName:              "Huntington, Huntington Ave. Entrance",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "0415",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2019-01-07T04:15:00",
						DateUpdated:              "2019-02-09T08:12:58",
						EstimatedReturnToService: "2019-04-17T23:59:59",
					},
					{
						UnitName:                 "D02S02",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "D02",
						StationName:              "Smithsonian, Independence Avenue Entrance",
						LocationDescription:      "Escalator between street mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "1512",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T15:12:00",
						DateUpdated:              "2019-04-05T15:13:17",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "D04X01",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "D04",
						StationName:              "Federal Center SW",
						LocationDescription:      "Elevator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "1449",
						SymptomDescription:       "Power Outage",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-02T14:49:00",
						DateUpdated:              "2019-04-05T02:33:45",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "D04X02",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "D04",
						StationName:              "Federal Center SW",
						LocationDescription:      "Escalator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "1338",
						SymptomDescription:       "Major Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-02T13:38:00",
						DateUpdated:              "2019-04-04T23:49:31",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "D07X03",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "D07",
						StationName:              "Potomac Ave",
						LocationDescription:      "Escalator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "1504",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T15:04:00",
						DateUpdated:              "2019-04-05T15:05:58",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "D09X04",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "D09",
						StationName:              "Minnesota Ave",
						LocationDescription:      "Garage Elevator",
						SymptomCode:              "",
						TimeOutOfService:         "1616",
						SymptomDescription:       "Fire Alarm",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-01T16:16:00",
						DateUpdated:              "2019-04-05T15:27:49",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "D10X02",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "D10",
						StationName:              "Deanwood",
						LocationDescription:      "Elevator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "0737",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T07:37:00",
						DateUpdated:              "2019-04-05T09:49:54",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "D10X03",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "D10",
						StationName:              "Deanwood",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "0508",
						SymptomDescription:       "Minor Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-03-29T05:08:00",
						DateUpdated:              "2019-04-03T14:27:56",
						EstimatedReturnToService: "2019-04-08T23:59:59",
					},
					{
						UnitName:                 "D13X04",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "D13",
						StationName:              "New Carrollton",
						LocationDescription:      "Garage elevator",
						SymptomCode:              "",
						TimeOutOfService:         "2031",
						SymptomDescription:       "Minor Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-03T20:31:00",
						DateUpdated:              "2019-04-04T14:55:19",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "D13X05",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "D13",
						StationName:              "New Carrollton",
						LocationDescription:      "Garage elevator",
						SymptomCode:              "",
						TimeOutOfService:         "2032",
						SymptomDescription:       "Minor Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-03T20:32:00",
						DateUpdated:              "2019-04-04T19:52:56",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "E07X04",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "E07",
						StationName:              "West Hyattsville",
						LocationDescription:      "Escalator between mezzanine and platform to Branch Ave",
						SymptomCode:              "",
						TimeOutOfService:         "1859",
						SymptomDescription:       "Minor Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-04T18:59:00",
						DateUpdated:              "2019-04-05T16:19:08",
						EstimatedReturnToService: "2019-04-06T23:59:59",
					},
					{
						UnitName:                 "F02X02",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "F02",
						StationName:              "Archives-Navy Memorial-Penn Quarter",
						LocationDescription:      "Escalator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0358",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2019-02-20T03:58:00",
						DateUpdated:              "2019-02-20T03:59:27",
						EstimatedReturnToService: "2019-05-20T23:59:59",
					},
					{
						UnitName:                 "F05W03",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "F05",
						StationName:              "Navy Yard-Ballpark, Half Street Entrance",
						LocationDescription:      "Escalator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0246",
						SymptomDescription:       "Major Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T02:46:00",
						DateUpdated:              "2019-04-05T15:18:06",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "F06S02",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "F06",
						StationName:              "Anacostia, Howard Road Entrance",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "1206",
						SymptomDescription:       "Minor Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-03-28T12:06:00",
						DateUpdated:              "2019-04-02T06:25:08",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "F07X04",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "F07",
						StationName:              "Congress Heights",
						LocationDescription:      "Escalator between mezzanine and Alabama Avenue",
						SymptomCode:              "",
						TimeOutOfService:         "0456",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2018-10-22T04:56:00",
						DateUpdated:              "2019-04-01T06:49:31",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "F07X05",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "F07",
						StationName:              "Congress Heights",
						LocationDescription:      "Escalator between mezzanine and Alabama Avenue",
						SymptomCode:              "",
						TimeOutOfService:         "0459",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2018-10-22T04:59:00",
						DateUpdated:              "2019-04-01T06:49:47",
						EstimatedReturnToService: "2019-04-05T23:59:59",
					},
					{
						UnitName:                 "J02X01",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "J02",
						StationName:              "Van Dorn Street",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "0425",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2019-01-15T04:25:00",
						DateUpdated:              "2019-02-09T08:12:58",
						EstimatedReturnToService: "2019-05-24T23:59:59",
					},
					{
						UnitName:                 "J02X02",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "J02",
						StationName:              "Van Dorn Street",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "0255",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2019-02-05T02:55:00",
						DateUpdated:              "2019-02-09T08:12:58",
						EstimatedReturnToService: "2019-05-24T23:59:59",
					},
					{
						UnitName:                 "K01X01",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "K01",
						StationName:              "Court House",
						LocationDescription:      "Elevator between street and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "2040",
						SymptomDescription:       "Major Repair",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-04T20:40:00",
						DateUpdated:              "2019-04-05T09:48:01",
						EstimatedReturnToService: "2019-04-06T23:59:59",
					},
					{
						UnitName:                 "K01X02",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "K01",
						StationName:              "Court House",
						LocationDescription:      "Escalator between street (Wilson Blvd) and middle landing/tunnel",
						SymptomCode:              "",
						TimeOutOfService:         "0541",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2018-09-14T05:41:00",
						DateUpdated:              "2019-03-25T02:47:13",
						EstimatedReturnToService: "2019-04-06T23:59:59",
					},
					{
						UnitName:                 "K01X05",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "K01",
						StationName:              "Court House",
						LocationDescription:      "Escalator between middle landing and mezzanine",
						SymptomCode:              "",
						TimeOutOfService:         "0402",
						SymptomDescription:       "Modernization",
						DisplayOrder:             0,
						DateOutOfService:         "2019-03-15T04:02:00",
						DateUpdated:              "2019-03-25T02:47:44",
						EstimatedReturnToService: "2019-04-06T23:59:59",
					},
					{
						UnitName:                 "K02X01",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "K02",
						StationName:              "Clarendon",
						LocationDescription:      "Escalator between street and middle landing/tunnel",
						SymptomCode:              "",
						TimeOutOfService:         "1010",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T10:10:00",
						DateUpdated:              "2019-04-05T13:14:43",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "K08X02",
						UnitType:                 "ESCALATOR",
						UnitStatus:               "",
						StationCode:              "K08",
						StationName:              "Vienna/Fairfax-GMU",
						LocationDescription:      "Escalator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "1322",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T13:22:00",
						DateUpdated:              "2019-04-05T13:24:30",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
					{
						UnitName:                 "N06X03",
						UnitType:                 "ELEVATOR",
						UnitStatus:               "",
						StationCode:              "N06",
						StationName:              "Wiehle-Reston East",
						LocationDescription:      "Elevator between mezzanine and platform",
						SymptomCode:              "",
						TimeOutOfService:         "1325",
						SymptomDescription:       "Service Call",
						DisplayOrder:             0,
						DateOutOfService:         "2019-04-05T13:25:00",
						DateUpdated:              "2019-04-05T13:27:02",
						EstimatedReturnToService: "2019-04-07T23:59:59",
					},
				},
			},
		},
	},
	"/Incidents.svc/json/Incidents": {
		{
			rawQuery:     "",
			param:        "",
			jsonResponse: `{"Incidents":[{"DateUpdated":"2010-07-29T14:21:28","DelaySeverity":null,"Description":"Red Line: Expect residual delays to Glenmont due to an earlier signal problem outside Forest Glen.","EmergencyText":null,"EndLocationFullName":null,"IncidentID":"3754F8B2-A0A6-494E-A4B5-82C9E72DFA74","IncidentType":"Delay","LinesAffected":"RD;","PassengerDelay":0,"StartLocationFullName":null}]}`,
			unmarshalledResponse: &GetRailIncidentsResponse{
				RailIncidents: []RailIncident{
					{
						DateUpdated:           "2010-07-29T14:21:28",
						DelaySeverity:         "",
						Description:           "Red Line: Expect residual delays to Glenmont due to an earlier signal problem outside Forest Glen.",
						EmergencyText:         "",
						EndLocationFullName:   "",
						IncidentID:            "3754F8B2-A0A6-494E-A4B5-82C9E72DFA74",
						IncidentType:          "Delay",
						LinesAffected:         "RD;",
						PassengerDelay:        0,
						StartLocationFullName: "",
					},
				},
			},
		},
	},
}

// setupTestService creates a service struct with a mock http client
func setupTestService() *Service {
	return &Service{
		client: &wmata.Client{
			HTTPClient: &testClient{},
		},
	}
}

func TestGetBusIncidents(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Incidents.svc/json/BusIncidents"]

	if !exist {
		t.Errorf("no data found for GetBusIncidents")
		return
	}

	for _, request := range testRequests {
		response, err := testService.GetBusIncidents(request.param)

		if err != nil {
			t.Errorf("error calling GetBusIncidents, Route: %s error: %s", request.param, err.Error())
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Errorf("unexpected response. Expected: %v but got: %v", response, request.unmarshalledResponse)
		}
	}
}

func TestGetElevatorEscalatorOutages(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Incidents.svc/json/ElevatorIncidents"]

	if !exist {
		t.Errorf("no data found for GetElevatorEscalatorOutages")
		return
	}

	for _, request := range testRequests {
		response, err := testService.GetOutages(request.param)

		if err != nil {
			t.Errorf("error calling GetOutages for station: %s error: %s", request.param, err.Error())
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Error(pretty.Diff(response, request.unmarshalledResponse))
		}
	}
}

func TestGetRailIncidents(t *testing.T) {
	testService := setupTestService()

	testRequests, exist := testData["/Incidents.svc/json/Incidents"]

	if !exist {
		t.Errorf("no data found for GetRailIncidents")
		return
	}

	for _, request := range testRequests {
		response, err := testService.GetRailIncidents()

		if err != nil {
			t.Errorf("error calling GetRailIncidents: %s", err.Error())
		}

		if !reflect.DeepEqual(response, request.unmarshalledResponse) {
			t.Error(pretty.Diff(response, request.unmarshalledResponse))
		}
	}
}
