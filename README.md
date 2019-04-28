# wmata-go-sdk
[![Build Status](https://travis-ci.com/awiede/wmata-go-sdk.svg?branch=master)](https://travis-ci.com/awiede/wmata-go-sdk) [![Coverage Status](https://coveralls.io/repos/github/awiede/wmata-go-sdk/badge.svg?branch=master)](https://coveralls.io/github/awiede/wmata-go-sdk?branch=master)


WIP - Go client for WMATA API

# Install

To install this package run:
```bash
go get -u github.com/awiede/wmata-go-sdk
```

## Register

In order to use the WMATA API, you will need to register with WMATA and get and API key. Information on how to get an API key can be found [here](https://developer.wmata.com/).

# Usage
This project is broken up into separate packages for each grouping of [WMATA APIs](https://developer.wmata.com/docs/services/) (*Please refer to these docs as the source of truth for all business rules*) 

The application is split up into the following services:
* [wmata](https://github.com/awiede/wmata-go-sdk/tree/master/wmata) - Top level package. Houses `client` configuration to make API calls. (*Note: This also houses the [Misc Method](https://developer.wmata.com/docs/services/5923434c08d33c0f201a600a/operations/5923437c031f5914d0204bcf) API for health checks*).
* [businfo](https://github.com/awiede/wmata-go-sdk/tree/master/wmata/businfo) - Service methods corresponding to [Bus Route and Stop](https://developer.wmata.com/docs/services/54763629281d83086473f231/operations/5476362a281d830c946a3d68) API.
* [buspredictions](https://github.com/awiede/wmata-go-sdk/tree/master/wmata/buspredictions) - Service methods corresponding to [Real-Time Bus Predictions](https://developer.wmata.com/docs/services/5476365e031f590f38092508/operations/5476365e031f5909e4fe331d) API.
* [incidents](https://github.com/awiede/wmata-go-sdk/tree/master/wmata/incidents) - Service methods corresponding to [Indicents](https://developer.wmata.com/docs/services/54763641281d83086473f232/operations/54763641281d830c946a3d75) API.
* [railinfo](https://github.com/awiede/wmata-go-sdk/tree/master/wmata/railinfo) - Service methods corresponding to [Rail Station Information](https://developer.wmata.com/docs/services/5476364f031f590f38092507/operations/5476364f031f5909e4fe330c) API.
* [railpredictions](https://github.com/awiede/wmata-go-sdk/tree/master/wmata/railpredictions) - Service methods corresponding to [Real-Time Rail Predictions](https://developer.wmata.com/docs/services/547636a6f9182302184cda78/operations/547636a6f918230da855363f) API.
* [trainpositions](https://github.com/awiede/wmata-go-sdk/tree/master/wmata/trainpositions) - Service methods corresponding to [Train Positions](https://developer.wmata.com/docs/services/5763fa6ff91823096cac1057/operations/5763fb35f91823096cac1058) API.

## Creating a `wmata.Client`

All requests using this SDK are routed through a `wmata.Client`. To create a client you will need an API key, and can optionally specify an `http.Client` with all associated configurations.

### Example

```go
apiKey := "123456789"
client := http.Client{
    Timeout: time.Second * 60,
}
wmataClient := wmata.NewWMATAClient(apiKey, client)

// Defaults with a 30 second timeout
defaultClient := wmata.NewWMATADefaultClient(apiKey)
```

## Creating and Using a Service

All services have a `<package-name>.NewService` function which will build a new service. To create a service a `wmata.Client` is required, as well as a choice of either communicating to WMATA via their `JSON` or `XML` endpoints.

### Example
```go
// Create rail info service using WMATA client in previous example
railInfoService := railinfo.NewService(wmataClient, wmata.JSON)

// Invoke the GetLines() method of the rail service to get active line info
activeRailLines, err := railInfoService.GetLines()
```