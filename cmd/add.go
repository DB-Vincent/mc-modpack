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

var addCmd = &cobra.Command{
	Use:   "add [mod-name]",
	Short: "Adds a mod to the current modpack",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		modName := args[0]

		path, err := getWorkingDirectory()
		if err != nil {
			log.Error(fmt.Sprintf("Failed to get current working directory: %s", err.Error()))
			return
		}

		if err := config.Exists(path); err != nil {
			return
		}

		cfg, err := config.Load(path)
		if err != nil {
			return
		}

		version, err := modrinth.GetLatestVersion(modName, cfg.McVersion, cfg.Loader)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to find latest version for mod '%s' (MC %s, %s): %v", modName, cfg.McVersion, cfg.Loader, err))
			return
		}

		if err = addMod(cfg, modName, path); err != nil {
			log.Error(fmt.Sprintf("Failed to add mod '%s': %v", modName, err))
			return
		}

		if len(version.Dependencies) > 0 {
			log.Info(fmt.Sprintf("Found %d dependencies, downloading..", len(version.Dependencies)))
			for _, dependency := range version.Dependencies {
				if err = addMod(cfg, dependency.ProjectId, path); err != nil {
					log.Error(fmt.Sprintf("Failed to add dependency: %v", err))
					return
				}
			}
		}

		err = config.Update(path, *cfg)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to save modpack configuration: %v", err))
			return
		}
	},
}

func addMod(cfg *config.Config, modName string, path string) error {
	version, err := modrinth.GetLatestVersion(modName, cfg.McVersion, cfg.Loader)
	if err != nil {
		return fmt.Errorf("failed to find latest version for mod '%s' (MC %s, %s): %v", modName, cfg.McVersion, cfg.Loader, err)
	}

	// Check if mod exists in config file
	modExists, index := config.HasMod(*cfg, modName)
	modDownload := true

	if modExists {
		if cfg.Mods[index].Version != version.ModVersion {
			log.Info(fmt.Sprintf("Updating %s from version %s to %s", version.Files[0].Name, cfg.Mods[index].Version, version.ModVersion))
			cfg.Mods[index].Version = version.ModVersion
			cfg.Mods[index].VersionId = version.VersionId
		} else {
			log.Info(fmt.Sprintf("Already have %s with version %s downloaded", version.Files[0].Name, cfg.Mods[index].Version))
			modDownload = false
		}
	} else {
		log.Info(fmt.Sprintf("Adding %s version %s to modpack", version.Files[0].Name, version.ModVersion))
		cfg.Mods = append(cfg.Mods, config.Mod{
			Name:      modName,
			VersionId: version.VersionId,
			Version:   version.ModVersion,
		})
	}

	if modDownload {
		if err = modrinth.DownloadFile(path, version.Files[0]); err != nil {
			return fmt.Errorf("failed to download mod '%s': %v", version.Files[0].Name, err)
		}
		log.Info(fmt.Sprintf("Successfully downloaded %s", version.Files[0].Name))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
