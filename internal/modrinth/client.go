/*
  Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package modrinth

import (
  "net/http"
  "time"
)

var httpClient = &http.Client{
  Timeout: 10 * time.Second, // Timeout of 10 seconds on requests
}

// Little function to send requests
func SendRequest(method, url string) (*http.Response, error) {
  req, err := http.NewRequest(method, url, nil)
  if err != nil {
    return nil, err
  }

  return httpClient.Do(req)
}
