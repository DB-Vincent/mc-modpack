/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package modrinth

type Version struct {
	GameVersions []string     `json:"game_versions"`
	Loaders      []string     `json:"loaders"`
	ModVersion   string       `json:"version_number"`
	VersionId    string       `json:"id"`
	Type         string       `json:"version_type"`
	Files        []File       `json:"files"`
	Dependencies []Dependency `json:"dependencies"`
}

type File struct {
	Hashes FileHash `json:"hashes"`
	Url    string   `json:"url"`
	Name   string   `json:"filename"`
}

type FileHash struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
}

type Dependency struct {
	VersionId string `json:"version_id"`
	ProjectId string `json:"project_id"`
	Type      string `json:"dependency_type"`
}

type Project struct {
	Id   string `json:"id"`
	Slug string `json:"slug"`
}
