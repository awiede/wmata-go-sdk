package wmata

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"net/http"
	"net/http/httptest"
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
	testResponses, exist := testData[req.URL.Path]

	if !exist {
		return nil, errors.New("no test data found")
	}

	for _, response := range testResponses {
		fmt.Println(req.URL.RawQuery)
		if response.rawQuery == req.URL.RawQuery {
			if response.headerField == "" || response.headerValue == req.Header.Get(response.headerField) {
				rr := httptest.NewRecorder()

				rr.WriteHeader(response.responseHttpCode)
				_, writeErr := rr.Write([]byte(response.response))

				return rr.Result(), writeErr
			}
		}
	}

	return nil, errors.New("no data found")
}

type testResponseData struct {
	rawQuery             string
	requestURL           string
	queryParams          map[string]string
	response             string
	responseHttpCode     int
	responseFormat       ResponseType
	headerField          string
	headerValue          string
	expectedError        error
	unmarshalledResponse interface{}
}

type testType struct {
	XMLName xml.Name `json:"-" xml:"http://foo.bar.test TestTypeResp"`
	Foo     string   `json:"foo" xml:"Foo"`
	Bar     string   `json:"bar" xml:"Bar"`
}

var testData = map[string][]testResponseData{
	"/Misc/Validate": {
		{
			rawQuery:         "",
			response:         "42",
			responseHttpCode: http.StatusOK,
			headerField:      APIKeyHeader,
			headerValue:      "123456789",
		},
		{
			rawQuery:         "",
			response:         "invalid API key",
			responseHttpCode: http.StatusUnauthorized,
			headerField:      APIKeyHeader,
			headerValue:      "987654321",
			expectedError:    errors.New("invalid API key"),
		},
	},
	"/test": {
		{
			rawQuery:         "Test=true",
			queryParams:      map[string]string{"Test": "true"},
			requestURL:       "http://foo.bar.test/test",
			response:         `{"foo": "hello", "bar": "world"}`,
			responseHttpCode: http.StatusOK,
			responseFormat:   JSON,
			headerField:      APIKeyHeader,
			headerValue:      "123456789",
			unmarshalledResponse: &testType{
				Foo: "hello",
				Bar: "world",
			},
		},
		{
			rawQuery:         "AnotherTest=true",
			queryParams:      map[string]string{"AnotherTest": "true"},
			requestURL:       "http://foo.bar.test/test",
			response:         `<TestTypeResp xmlns="http://foo.bar.test"><Foo>hello</Foo><Bar>world</Bar></TestTypeResp>`,
			responseHttpCode: http.StatusOK,
			responseFormat:   XML,
			headerField:      APIKeyHeader,
			headerValue:      "123456789",
			unmarshalledResponse: &testType{
				XMLName: xml.Name{
					Space: "http://foo.bar.test",
					Local: "TestTypeResp",
				},
				Foo: "hello",
				Bar: "world",
			},
		},
		{
			rawQuery:         "ErrorTest=true",
			queryParams:      map[string]string{"ErrorTest": "true"},
			requestURL:       "http://foo.bar.test/test",
			response:         `<TestTypeResp xmlns="http://foo.bar.test"><Foo>hello</Foo><Bar>world</Bar></TestTypeResp>`,
			responseHttpCode: http.StatusOK,
			responseFormat:   3,
			headerField:      APIKeyHeader,
			headerValue:      "123456789",
			expectedError:    errors.New("invalid response type"),
		},
	},
}

func TestNewWMATAClient(t *testing.T) {
	apiKey := "123456789"
	client := http.Client{
		Timeout: time.Second * 60,
	}
	wmataClient := NewWMATAClient(apiKey, client)

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

func TestValidateAPIKey(t *testing.T) {
	testRequests, exist := testData["/Misc/Validate"]

	if !exist {
		t.Errorf("no test data found for ValidateAPIKey")
		return
	}

	for _, request := range testRequests {
		wmataClient := Client{
			APIKey:     request.headerValue,
			HTTPClient: &testHttpClient{},
		}

		responseStatus, responseErr := wmataClient.ValidateAPIKey()

		if responseErr != nil && responseErr.Error() != request.expectedError.Error() {
			t.Errorf("unexpected error: %s", responseErr)
			continue
		}

		if responseStatus != request.responseHttpCode {
			t.Errorf("unexpected response status: %d", responseStatus)
			continue
		}
	}
}

func TestBuildAndSendGetRequest(t *testing.T) {
	testRequests, exist := testData["/test"]

	if !exist {
		t.Errorf("no test data found for BuildAndSendGetRequest")
		return
	}

	for _, request := range testRequests {
		wmataClient := Client{
			APIKey:     request.headerValue,
			HTTPClient: &testHttpClient{},
		}

		test := testType{}

		responseErr := wmataClient.BuildAndSendGetRequest(request.responseFormat, request.requestURL, request.queryParams, &test)

		if responseErr != nil {
			if request.expectedError == nil || responseErr.Error() != request.expectedError.Error() {
				t.Errorf("unexpected error: %s", responseErr)
			}
			continue
		}

		if !reflect.DeepEqual(&test, request.unmarshalledResponse) {
			t.Error(pretty.Diff(&test, request.unmarshalledResponse))
			continue
		}
	}

}
