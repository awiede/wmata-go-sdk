package incidents

import (
	"github.com/awiede/wmata-go-sdk/wmata"
	"strings"
)

const incidentsServiceBaseURL = "https://api.wmata.com/Incidents.svc"

type GetBusIncidentsResponse struct {
	BusIncidents []BusIncident `json:"BusIncidents" xml:"BusIncidents"`
}

type BusIncident struct {
	DateUpdated    string   `json:"DateUpdated" xml:"DateUpdated"`
	Description    string   `json:"Description" xml:"Description"`
	IncidentID     string   `json:"IncidentID" xml:"IncidentID"`
	IncidentType   string   `json:"IncidentType" xml:"IncidentType"`
	RoutesAffected []string `json:"RoutesAffected" xml:"RoutesAffected"`
}

type GetElevatorEscalatorOutagesResponse struct {
	ElevatorIncidents []ElevatorIncident `json:"ElevatorIncidents" xml:"ElevatorIncidents"`
}

type ElevatorIncident struct {
	DateOutOfService string `json:"DateOutOfServ" xml:"DateOutOfServ"`
	DateUpdated      string `json:"DateUpdated" xml:"DateUpdated"`
	// Deprecated: DisplayOrder response field is deprecated
	DisplayOrder             int    `json:"DisplayOrder" xml:"DisplayOrder"`
	EstimatedReturnToService string `json:"EstimatedReturnToService" xml:"EstimatedReturnToService"`
	LocationDescription      string `json:"LocationDescription" xml:"LocationDescription"`
	StationCode              string `json:"StationCode" xml:"StationCode"`
	StationName              string `json:"StationName" xml:"StationName"`
	// Deprecated: SymptomCode response field is deprecated
	SymptomCode        string `json:"SymptomCode" xml:"SymptomCode"`
	SymptomDescription string `json:"SymptomDescription" xml:"SymptomDescription"`
	// Deprecated: TimeOutOfService response field is deprecated, use time portion of DateOutOfService
	TimeOutOfService string `json:"TimeOutOfService" xml:"TimeOutOfService"`
	UnitName         string `json:"UnitName" xml:"UnitName"`
	// Deprecated: UnitStatus response field is deprecated
	UnitStatus string `json:"UnitStatus" xml:"UnitStatus"`
	UnitType   string `json:"UnitType" xml:"UnitType"`
}

type GetRailIncidentsResponse struct {
	RailIncidents []RailIncident `json:"Incidents" xml:"Incidents"`
}

type RailIncident struct {
	DateUpdated string `json:"DateUpdated" xml:"DateUpdated"`
	// Deprecated: DelaySeverity response field is deprecated
	DelaySeverity string `json:"DelaySeverity" xml:"DelaySeverity"`
	Description   string `json:"Description" xml:"Description"`
	// Deprecated: EmergencyText response field is deprecated
	EmergencyText string `json:"EmergencyText" xml:"EmergencyText"`
	// Deprecated: EndLocationFullName response field is deprecated
	EndLocationFullName string `json:"EndLocationFullName" xml:"EndLocationFullName"`
	IncidentID          string `json:"IncidentID" xml:"IncidentID"`
	IncidentType        string `json:"IncidentType" xml:"IncidentType"`
	// LinesAffected returns a semi-colon and space separated list of line codes - will be updated to an array eventually
	LinesAffected string `json:"LinesAffected" xml:"LinesAffected"`
	// Deprecated: PassengerDelay response field is deprecated
	PassengerDelay int `json:"PassengerDelay" xml:"PassengerDelay"`
	// Deprecated: StartLocationFullName response field is deprecated
	StartLocationFullName string `json:"StartLocationFullName" xml:"StartLocationFullName"`
}

type Incidents interface {
	GetBusIncidents(route string) (*GetBusIncidentsResponse, error)
	GetOutages(stationCode string) (*GetElevatorEscalatorOutagesResponse, error)
	GetRailIncidents() (*GetRailIncidentsResponse, error)
}

var _ Incidents = (*Service)(nil)

// NewService returns a new Incidents service with a reference to an existing wmata.Client
func NewService(client *wmata.Client, responseType wmata.ResponseType) *Service {
	return &Service{
		client:       client,
		responseType: responseType,
	}
}

type Service struct {
	client       *wmata.Client
	responseType wmata.ResponseType
}

func (incidentService *Service) GetBusIncidents(route string) (*GetBusIncidentsResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(incidentsServiceBaseURL)

	switch incidentService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/BusIncidents")
	case wmata.XML:
		requestUrl.WriteString("/BusIncidents")
	}

	busIncident := GetBusIncidentsResponse{}

	return &busIncident, incidentService.client.BuildAndSendGetRequest(incidentService.responseType, requestUrl.String(), map[string]string{"Route": route}, &busIncident)

}

func (incidentService *Service) GetOutages(stationCode string) (*GetElevatorEscalatorOutagesResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(incidentsServiceBaseURL)

	switch incidentService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/ElevatorIncidents")
	case wmata.XML:
		requestUrl.WriteString("/ElevatorIncidents")
	}

	outages := GetElevatorEscalatorOutagesResponse{}

	return &outages, incidentService.client.BuildAndSendGetRequest(incidentService.responseType, requestUrl.String(), map[string]string{"StationCode": stationCode}, &outages)
}

func (incidentService *Service) GetRailIncidents() (*GetRailIncidentsResponse, error) {
	var requestUrl strings.Builder
	requestUrl.WriteString(incidentsServiceBaseURL)

	switch incidentService.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/Incidents")
	case wmata.XML:
		requestUrl.WriteString("/Incidents")
	}

	railIncidents := GetRailIncidentsResponse{}

	return &railIncidents, incidentService.client.BuildAndSendGetRequest(incidentService.responseType, requestUrl.String(), nil, &railIncidents)
}
