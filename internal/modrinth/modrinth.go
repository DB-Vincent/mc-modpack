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
	"os"
)

// Retrieve latest version of mod for given Minecraft version
func GetLatestVersion(projectId, mcVersion, modLoader string) (*Version, error) {
	typePriority := []string{"release", "beta", "alpha"}

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

	return nil, fmt.Errorf("no versions found with the given priority (release, beta, alpha)")
}

// Get project information by project ID
func GetProject(projectId string) (*Project, error) {
	baseURL := fmt.Sprintf("https://api.modrinth.com/v2/project/%s", projectId)

	res, err := SendRequest("GET", baseURL)
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

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// Get a specific version of a mod
func GetSpecificVersion(name, id string) (Version, error) {
	baseURL := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version/%s", name, id)

	res, err := SendRequest("GET", baseURL)
	if err != nil {
		return Version{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Version{}, fmt.Errorf("API returned status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Version{}, err
	}

	var resVersion Version
	if err := json.Unmarshal(body, &resVersion); err != nil {
		return Version{}, err
	}

	return resVersion, nil
}

// Download mod file
func DownloadFile(location string, file File) error {
	// Create destination file
	out, err := os.Create(fmt.Sprintf("%s/%s", location, file.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	// Get data
	res, err := SendRequest("GET", file.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check if server responded OK
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned bad status: %s", res.Status)
	}

	// Write data to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	// Verify hash
	hash, err := getHash(fmt.Sprintf("%s/%s", location, file.Name))
	if err != nil {
		return err
	}
	if hash != file.Hashes.Sha512 {
		return fmt.Errorf("Downloaded file hash is different from expected hash")
	}

	return nil
}
