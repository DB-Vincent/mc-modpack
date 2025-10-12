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
		modName := args[0]

		path := workingDirectory
		if workingDirectory == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				return
			}
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

		// Check if mod exists in config file
		modExists, index := config.HasMod(*cfg, modName)

		if modExists {
			log.Info(fmt.Sprintf("Updating %s from version %s to %s", modName, cfg.Mods[index].Version, version.ModVersion))
			cfg.Mods[index].Version = version.ModVersion
			cfg.Mods[index].VersionId = version.VersionId
		} else {
			log.Info(fmt.Sprintf("Adding %s version %s to modpack", modName, version.ModVersion))
			cfg.Mods = append(cfg.Mods, config.Mod{
				Name:      modName,
				VersionId: version.VersionId,
				Version:   version.ModVersion,
			})
		}

		if err = modrinth.DownloadFile(path, version.Files[0]); err != nil {
			log.Error(fmt.Sprintf("Failed to download mod file '%s': %v", version.Files[0].Name, err))
			return
		}
		log.Info(fmt.Sprintf("Successfully downloaded %s", version.Files[0].Name))

		err = config.Update(path, *cfg)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to save modpack configuration: %v", err))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
