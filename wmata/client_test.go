package wmata

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

// testHttpClient is a mock implementation of wmata.HTTPClient interface used for testing purposes
type testHttpClient struct {
}

// ensure testHttpClient implements wmata.HTTPClient interface
var _ HTTPClient = (*testHttpClient)(nil)

// Do stubs out an httpClient.Do request
func (httpClient *testHttpClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestNewWMATAClient(t *testing.T) {
	apiKey := "123456789"
	client := http.Client{
		Timeout: time.Second * 60,
	}
	wmataClient := NewWMATAClient(apiKey, client, JSON)

	if apiKey != wmataClient.APIKey {
		t.Errorf("invalid API key: %s", wmataClient.APIKey)
	}

	if reflect.DeepEqual(client, wmataClient.HTTPClient) {
		t.Errorf("invalid http client: %v", wmataClient.HTTPClient)
	}

	defaultClient := NewWMATADefaultClient(apiKey)

	if apiKey != defaultClient.APIKey {
		t.Errorf("invalid API key: %s", wmataClient.APIKey)
	}

	httpClient := defaultClient.HTTPClient.(*http.Client)

	if (time.Second * 30) != httpClient.Timeout {
		t.Errorf("incorrect timeout value: %d", httpClient.Timeout)
	}

}
