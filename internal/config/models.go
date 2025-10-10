/*
  Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package config

type Config struct {
  McVersion string
  Name      string
  Loader    string
  Mods      []Mod
}

type Mod struct {
  Name    string
  Version string
}
