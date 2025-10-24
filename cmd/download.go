/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"fmt"

	"github.com/DB-Vincent/mc-modpack/internal/config"
	"github.com/DB-Vincent/mc-modpack/internal/modrinth"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads all mods in modpack",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getWorkingDirectory()
		if err != nil {
			log.Error(fmt.Sprintf("Failed to get current working directory: %s", err.Error()))
			return
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

		if len(cfg.Mods) == 0 && len(cfg.Dependencies) == 0 {
			log.Info("No mods or dependencies found in modpack configuration")
			return
		}

		totalItems := len(cfg.Mods) + len(cfg.Dependencies)
		log.Info(fmt.Sprintf("Downloading %d item(s) for modpack '%s' (%d mods, %d dependencies)", totalItems, cfg.Name, len(cfg.Mods), len(cfg.Dependencies)))

		// Download mods
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
			log.Info(fmt.Sprintf("Downloaded mod: %s", modVersion.Files[0].Name))
		}

		// Download dependencies
		for _, dependency := range cfg.Dependencies {
			dependencyVersion, err := modrinth.GetSpecificVersion(dependency.Name, dependency.VersionId)
			if err != nil {
				log.Error(fmt.Sprintf("Unable to retrieve information for dependency '%s': %v", dependency.Name, err))
				continue
			}

			if err = modrinth.DownloadFile(path, dependencyVersion.Files[0]); err != nil {
				log.Error(fmt.Sprintf("Failed to download dependency '%s' (%s): %v", dependency.Name, dependencyVersion.Files[0].Name, err))
				continue
			}
			log.Info(fmt.Sprintf("Downloaded dependency: %s", dependencyVersion.Files[0].Name))
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
