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
        log.Error(fmt.Sprintf("Failed to get current working directory: %s", err.Error()))
      }
    }

    if err := config.Exists(path); err != nil {
      log.Error(fmt.Sprintf("Failed to check if config exists: %s", err.Error()))
    } 

    cfg, err := config.Load(path)
    if err != nil {
      log.Error(fmt.Sprintf("Failed to load configuration file: %s", err.Error()))
    }

    version, err := modrinth.GetLatestVersion(modName, cfg.McVersion, cfg.Loader)
    if err != nil {
      log.Error(fmt.Sprintf("Failed retrieving latest version of specified mod: %s", err.Error()))
    }

    // Check if mod exists in config file
    modExists, index := config.HasMod(*cfg, modName)
    
    if modExists {
      log.Info(fmt.Sprintf("Current version: %s, updating to new version: %s", cfg.Mods[index].Version, version.ModVersion))
      cfg.Mods[index].Version = version.ModVersion
      cfg.Mods[index].VersionId = version.VersionId
    } else {
      log.Info(fmt.Sprintf("Adding %s with version %s", modName, version.ModVersion))
      cfg.Mods = append(cfg.Mods, config.Mod{
        Name: modName,
        VersionId: version.VersionId, 
        Version: version.ModVersion,
      })
    }

    if err = modrinth.DownloadFile(path, version.Files[0]); err != nil {
      log.Error(fmt.Sprintf("Failed to download mod file: %s", err.Error()))
    }
    log.Info(fmt.Sprintf("Downloaded mod file %s", version.Files[0].Name))

    err = config.Update(path, *cfg)
    if err != nil {
      log.Error(fmt.Sprintf("Failed updating the configuration file: %s", err.Error()))
    }
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
