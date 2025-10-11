/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
  "fmt"
	"os"

	"github.com/DB-Vincent/mc-modpack/internal/config"
	"github.com/DB-Vincent/mc-modpack/internal/modrinth"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [mod-name]",
	Short: "Adds a mod to the current modpack",
  Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
    if (len(args) != 1) {
      cmd.Usage()
    }
    modName := args[0]

    path := workingDirectory
    if workingDirectory == "" {
      var err error
      path, err = os.Getwd()
      if err != nil {
        panic(err)
      }
    }

    if err := config.Exists(path); err != nil {
      panic(err)
    } 

    cfg, err := config.Load(path)
    if err != nil {
      panic(err)
    }

    version, err := modrinth.GetLatestVersion(modName, cfg.McVersion, cfg.Loader)
    if err != nil {
      panic(err)
    }

    // Check if the mod already exists in the config file, update if it does
    modExists := false
    if (len(cfg.Mods) > 0) {
      for i, mod := range cfg.Mods {
        if (mod.Name == modName) {
          modExists = true
          fmt.Printf("Current version: %s, new version: %s\n", mod.Version, version.ModVersion)
          cfg.Mods[i].Version = version.ModVersion
          fmt.Printf("Updated to %s\n", cfg.Mods[i].Version)
        }
      }
    }

    // Mod doesn't exist in config file, so we need to add it
    if !modExists {
      fmt.Printf("Adding %s with version %s", modName, version.ModVersion)
      cfg.Mods = append(cfg.Mods, config.Mod{
        Name: modName,
        Version: version.ModVersion,
      })
    }

    err = config.Update(path, *cfg)
    if err != nil {
      panic(err)
    }
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
