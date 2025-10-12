/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package modrinth

import (
	"net/http"
	"time"
)

const (
	// HTTPTimeout is the timeout duration for HTTP requests
	HTTPTimeout = 10 * time.Second
)

var httpClient = &http.Client{
	Timeout: HTTPTimeout,
}

// Little function to send requests
func SendRequest(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	return httpClient.Do(req)
}
