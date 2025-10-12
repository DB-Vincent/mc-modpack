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

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Removes a mod from the current modpack",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		modName := args[0]

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

		// Check if mod exists in config file
		modExists, index := config.HasMod(*cfg, modName)

		if modExists {
			mod, err := modrinth.GetSpecificVersion(modName, cfg.Mods[index].VersionId)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to get information for requested mod: %s", err.Error()))
				return
			}

			cfg.Mods = append(cfg.Mods[:index], cfg.Mods[index+1:]...)
			config.Update(path, *cfg)
			log.Info(fmt.Sprintf("Removed %s from list of mods", modName))

			// Remove jar file of deleted mod
			if err = os.Remove(fmt.Sprintf("%s/%s", path, mod.Files[0].Name)); err != nil {
				log.Error(fmt.Sprintf("Failed removing the mod's jar-file: %s", err.Error()))
				return
			}
			log.Info(fmt.Sprintf("Deleted %s", mod.Files[0].Name))
		} else {
			log.Info(fmt.Sprintf("Mod with name %s does not exist in configuration file", modName))
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
