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
				return
			}
		}

		if err := config.Exists(path); err != nil {
			log.Error(fmt.Sprintf("Failed to check if config exists: %s", err.Error()))
			return
		}

		cfg, err := config.Load(path)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to load configuration file: %s", err.Error()))
			return
		}

		if len(cfg.Mods) == 0 {
			log.Info("No mods found in modpack configuration")
			return
		}

		log.Info(fmt.Sprintf("Downloading %d mod(s) for modpack '%s'", len(cfg.Mods), cfg.Name))

		for _, mod := range cfg.Mods {
			modVersion, err := modrinth.GetSpecificVersion(mod.Name, mod.VersionId)
			if err != nil {
				log.Error(fmt.Sprintf("Unable to retrieve information for mod '%s': %v", mod.Name, err))
				continue
			}

			if err = modrinth.DownloadFile(path, modVersion.Files[0]); err != nil {
				log.Error(fmt.Sprintf("Failed to download mod '%s' (%s): %v", mod.Name, modVersion.Files[0].Name, err))
				continue
			}
			log.Info(fmt.Sprintf("Downloaded: %s", modVersion.Files[0].Name))
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
