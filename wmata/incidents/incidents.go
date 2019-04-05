package incidents

type Incidents interface {
	GetBusIncidents(route string) (*GetBusIncidentsResponse, error)
	GetOutages(stationCode string) (*GetElevatorEscalatorOutagesResponse, error)
	GetRailIncidents()
}

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
