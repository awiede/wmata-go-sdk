package wmata

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	LineCodeAll    = ""
	LineCodeBlue   = "BL"
	LineCodeGreen  = "GR"
	LineCodeOrange = "OR"
	LineCodeRed    = "RD"
	LineCodeSilver = "SV"
	LineCodeYellow = "YL"

	APIKeyHeader = "api_key"
)

type ResponseType int

const (
	JSON ResponseType = iota
	XML
)

// CloseResponseBody is a helper function to close response body and log error
func CloseResponseBody(response *http.Response) {
	if closeErr := response.Body.Close(); closeErr != nil {
		log.Printf("error closing response body: %s", closeErr)
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a wmata specific http client that includes authentication information
type Client struct {
	APIKey     string
	HTTPClient HTTPClient
}

// NewWMATADefaultClient returns a new client to make requests to the WMATA API
// This creates a default http.Client with a 30 second timeout
func NewWMATADefaultClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// NewWMATAClient returns a new client to make requests to the WMATA API
func NewWMATAClient(apiKey string, httpClient http.Client) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: &httpClient,
	}
}

// ValidateAPIKey sends a validation request to the WMATA API to verify the given API works and that WMATA is available.
// Returns 200, nil if able to connect and receive a success (200) response from WMATA - otherwise returns status code and error message
func (client *Client) ValidateAPIKey() (int, error) {
	request, requestErr := http.NewRequest(http.MethodGet, "https://api.wmata.com/Misc/Validate", nil)

	if requestErr != nil {
		return http.StatusInternalServerError, requestErr
	}

	request.Header.Add(APIKeyHeader, client.APIKey)

	response, responseErr := client.HTTPClient.Do(request)

	if responseErr != nil {
		return response.StatusCode, responseErr
	}

	defer CloseResponseBody(response)

	if response.StatusCode != http.StatusOK {
		body, readErr := ioutil.ReadAll(response.Body)

		if readErr != nil {
			return response.StatusCode, readErr
		}

		return response.StatusCode, errors.New(string(body))

	}

	return response.StatusCode, nil

}

// BuildAndSendGetRequest constructs and sends a generic HTTP GET request against the WMATA API
func (client *Client) BuildAndSendGetRequest(responseFormat ResponseType, url string, queryParams map[string]string, apiResponse interface{}) error {
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)

	if requestErr != nil {
		return requestErr
	}

	request.Header.Add(APIKeyHeader, client.APIKey)

	if queryParams != nil {
		query := request.URL.Query()

		for key, value := range queryParams {
			query.Add(key, value)
		}

		request.URL.RawQuery = query.Encode()
	}

	response, responseErr := client.HTTPClient.Do(request)

	if responseErr != nil {
		return responseErr
	}

	defer CloseResponseBody(response)

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return readErr
	}

	switch responseFormat {
	case JSON:
		return json.Unmarshal(body, &apiResponse)
	case XML:
		return xml.Unmarshal(body, &apiResponse)
	default:
		return errors.New("invalid response type")
	}

}
