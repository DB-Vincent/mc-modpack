/*
  Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package modrinth

type Version struct {
  GameVersions []string `json:"game_versions"`
  Loaders      []string `json:"loaders"`
  ModVersion   string   `json:"version_number"`
  Type         string   `json:"version_type"`
  Files        []File   `json:"files"`
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
