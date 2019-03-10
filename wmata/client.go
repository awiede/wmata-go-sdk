package wmata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
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

func (client *Client) BuildAndSendGetRequest(url string, queryParams map[string]string, apiResponse *interface{}) error {
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

	return json.Unmarshal(body, apiResponse)
}