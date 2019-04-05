package incidents

import (
	"encoding/json"
	"github.com/awiede/wmata-go-sdk/wmata"
	"io/ioutil"
	"net/http"
)

type GetBusIncidentsResponse struct {
	BusIncidents []BusIncident `json:"BusIncidents"`
}

type BusIncident struct {
	DateUpdated    string   `json:"DateUpdated"`
	Description    string   `json:"Description"`
	IncidentID     string   `json:"IncidentID"`
	IncidentType   string   `json:"IncidentType"`
	RoutesAffected []string `json:"RoutesAffected"`
}

type GetElevatorEscalatorOutagesResponse struct {
	ElevatorIncidents []ElevatorIncident `json:"ElevatorIncidents"`
}

type ElevatorIncident struct {
	DateOutOfService string `json:"DateOutOfServ"`
	DateUpdated      string `json:"DateUpdated"`
	// Deprecated: DisplayOrder response field is deprecated
	DisplayOrder             int    `json:"DisplayOrder"`
	EstimatedReturnToService string `json:"EstimatedReturnToService"`
	LocationDescription      string `json:"LocationDescription"`
	StationCode              string `json:"StationCode"`
	// Deprecated: SymptomCode response field is deprecated
	SymptomCode        string `json:"SymptomCode"`
	SymptomDescription string `json:"SymptomDescription"`
	// Deprecated: TimeOutOfService response field is deprecated, use time portion of DateOutOfServ
	TimeOutOfService string `json:"TimeOutOfService"`
	UnitName         string `json:"UnitName"`
	// Deprecated: UnitStatus response field is deprecated
	UnitStatus string `json:"UnitStatus"`
	UnitType   string `json:"UnitType"`
}

type GetRailIncidentsResponse struct {
	RailIncidents []RailIncident `json:"Incidents"`
}

type RailIncident struct {
	DateUpdated string `json:"DateUpdated"`
	// Deprecated: DelaySeverity response field is deprecated
	DelaySeverity string `json:"DelaySeverity"`
	Description   string `json:"Description"`
	// Deprecated: EmergencyText response field is deprecated
	EmergencyText string `json:"EmergencyText"`
	// Deprecated: EndLocationFullName response field is deprecated
	EndLocationFullName string `json:"EndLocationFullName"`
	IncidentID          string `json:"IncidentID"`
	IncidentType        string `json:"IncidentType"`
	// LinesAffected returns a semi-colon and space separated list of line codes - will be updated to an array eventually
	LinesAffected string `json:"LinesAffected"`
	// Deprecated: PassengerDelay response field is deprecated
	PassengerDelay int `json:"PassengerDelay"`
	// Deprecated: StartLocationFullName response field is deprecated
	StartLocationFullName string `json:"StartLocationFullName"`
}

type Incidents interface {
	GetBusIncidents(route string) (*GetBusIncidentsResponse, error)
	GetOutages(stationCode string) (*GetElevatorEscalatorOutagesResponse, error)
	GetRailIncidents() (*GetRailIncidentsResponse, error)
}

var _ Incidents = (*Service)(nil)

type Service struct {
	client *wmata.Client
}

func (incidentService *Service) GetBusIncidents(route string) (*GetBusIncidentsResponse, error) {
	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Incidents.svc/json/BusIncidents", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, incidentService.client.APIKey)

	query := request.URL.Query()
	query.Add("Route", route)

	request.URL.RawQuery = query.Encode()

	response, responseErr := incidentService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	busIncident := GetBusIncidentsResponse{}

	unmarshalErr := json.Unmarshal(body, &busIncident)

	return &busIncident, unmarshalErr

}

func (incidentService *Service) GetOutages(stationCode string) (*GetElevatorEscalatorOutagesResponse, error) {
	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Incidents.svc/json/ElevatorIncidents", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, incidentService.client.APIKey)

	query := request.URL.Query()
	query.Add("StationCode", stationCode)

	request.URL.RawQuery = query.Encode()

	response, responseErr := incidentService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	outages := GetElevatorEscalatorOutagesResponse{}

	unmarshalErr := json.Unmarshal(body, &outages)

	return &outages, unmarshalErr
}

func (incidentService *Service) GetRailIncidents() (*GetRailIncidentsResponse, error) {
	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Incidents.svc/json/Incidents", nil)

	if requestErr != nil {
		return nil, requestErr
	}

	request.Header.Set(wmata.APIKeyHeader, incidentService.client.APIKey)

	response, responseErr := incidentService.client.HTTPClient.Do(request)

	if responseErr != nil {
		return nil, responseErr
	}

	defer wmata.CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	railIncidents := GetRailIncidentsResponse{}

	unmarshalErr := json.Unmarshal(body, &railIncidents)

	return &railIncidents, unmarshalErr
}
