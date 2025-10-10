/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package modrinth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
  "net/url"
)


// Retrieve latest version of mod for given Minecraft version
func GetLatestVersion(projectId, mcVersion, modLoader string) (*Version, error) {
  typePriority := []string{"stable", "beta", "alpha"}

	baseURL := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", projectId)
	params := url.Values{}
	params.Add("game_versions", fmt.Sprintf("[\"%s\"]", mcVersion))
	params.Add("loaders", fmt.Sprintf("[\"%s\"]", modLoader))

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}
	u.RawQuery = params.Encode()

	res, err := SendRequest("GET", u.String())
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  if res.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("API returned status %d", res.StatusCode)
  }

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var versions []Version
	if err := json.Unmarshal(body, &versions); err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found")
	}

  for _, priority := range typePriority {
    for _, version := range versions {
      if version.Type == priority {
        return &version, nil
      }
    }
  }

	return nil, nil
}
