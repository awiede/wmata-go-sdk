package wmata

import (
	"log"
	"net/http"
)

// CloseResponseBody is a helper function to close response body and log error
func CloseResponseBody(response *http.Response) {
	if closeErr := response.Body.Close(); closeErr != nil {
		log.Printf("error closing response body: %s", closeErr)
	}
}