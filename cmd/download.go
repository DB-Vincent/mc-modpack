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

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads all mods in modpack",
	Run: func(cmd *cobra.Command, args []string) {
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

    for _, mod := range cfg.Mods {
      modVersion, err := modrinth.GetSpecificVersion(mod.Name, mod.VersionId)
      if err != nil {
        log.Error(fmt.Sprintf("Failed to get information for requested mod: %s", err.Error()))
      }
      
      if err = modrinth.DownloadFile(path, modVersion.Files[0]); err != nil {
        log.Error(fmt.Sprintf("Failed to download mod file: %s", err.Error()))
      }
      log.Info(fmt.Sprintf("Downloaded file %s", modVersion.Files[0].Name))
    }
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
