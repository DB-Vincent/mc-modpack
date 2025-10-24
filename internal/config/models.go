/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package config

type Config struct {
	McVersion string
	Name      string
	Loader    string

	Mods         []Mod
	Dependencies []Dependency
}

type Mod struct {
	Name      string
	Version   string
	VersionId string
}

type Dependency struct {
	Name       string
	Version    string
	VersionId  string
	RequiredBy []string
}
