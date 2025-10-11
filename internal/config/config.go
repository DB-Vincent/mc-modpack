/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package config

import (
  "errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

// Update a configuration file (gets created if it doesn't exist yet)
func Update(location string, config Config) (error) {
  // Create file
  file, err := os.Create(fmt.Sprintf("%s/modpack.toml", location))
  if err != nil {
    return err
  }
  defer file.Close()

  // Encode new config to file
  encoder := toml.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
	  return err
  }

  return nil
}

// Retrieve configuration from a configuration file
func Load(location string) (*Config, error) {
  var config Config

  // Open file
  file, err := os.ReadFile(fmt.Sprintf("%s/modpack.toml", location))
  if err != nil {
    return nil, err
  }

  // Marshall file content into Config struct
  err = toml.Unmarshal(file, &config)
  if err != nil {
    return nil, err
  }

  return &config, nil
}

// Check if the configuration file exists
func Exists(location string) (error) {
  if _, err := os.Stat(fmt.Sprintf("%s/modpack.toml", location)); errors.Is(err, os.ErrNotExist) {
    return fmt.Errorf("Config file does not exist. First init using the init command.")
  }

  return nil
}

// Check if configuration file has a given mod 
func HasMod(cfg Config, modName string) (bool, int) {
  if (len(cfg.Mods) > 0) {
    for i, mod := range cfg.Mods {
      if (mod.Name == modName) {
        // Mod was found
        return true, i
      }
    } 
  }

  // Mod was not found
  return false, -1
}
