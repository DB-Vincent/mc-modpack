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
	Use:   "del [mod-name]",
	Short: "Removes a mod or dependency from the current modpack",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		itemName := args[0]
		isDependency, _ := cmd.Flags().GetBool("dependency")

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

		if isDependency {
			// Handle dependency removal
			dependencyExists, index := config.HasDependency(*cfg, itemName)

			if dependencyExists {
				// Check if other mods still require this dependency
				if len(cfg.Dependencies[index].RequiredBy) > 0 {
					log.Info(fmt.Sprintf("Cannot remove dependency %s - still required by: %v", itemName, cfg.Dependencies[index].RequiredBy))
					return
				}

				dependency, err := modrinth.GetSpecificVersion(itemName, cfg.Dependencies[index].VersionId)
				if err != nil {
					log.Error(fmt.Sprintf("Failed to get information for requested dependency: %s", err.Error()))
					return
				}

				cfg.Dependencies = append(cfg.Dependencies[:index], cfg.Dependencies[index+1:]...)
				config.Update(path, *cfg)
				log.Info(fmt.Sprintf("Removed %s from list of dependencies", itemName))

				// Remove jar file of deleted dependency
				if err = os.Remove(fmt.Sprintf("%s/%s", path, dependency.Files[0].Name)); err != nil {
					log.Error(fmt.Sprintf("Failed removing the dependency's jar-file: %s", err.Error()))
					return
				}
				log.Info(fmt.Sprintf("Deleted %s", dependency.Files[0].Name))
			} else {
				log.Info(fmt.Sprintf("Dependency with name %s does not exist in configuration file", itemName))
			}
		} else {
			// Handle mod removal
			modExists, index := config.HasMod(*cfg, itemName)

			if modExists {
				mod, err := modrinth.GetSpecificVersion(itemName, cfg.Mods[index].VersionId)
				if err != nil {
					log.Error(fmt.Sprintf("Failed to get information for requested mod: %s", err.Error()))
					return
				}

				cfg.Mods = append(cfg.Mods[:index], cfg.Mods[index+1:]...)

				// Remove the mod from RequiredBy lists of dependencies
				removeModFromDependencies(cfg, itemName)

				config.Update(path, *cfg)
				log.Info(fmt.Sprintf("Removed %s from list of mods", itemName))

				// Remove jar file of deleted mod
				if err = os.Remove(fmt.Sprintf("%s/%s", path, mod.Files[0].Name)); err != nil {
					log.Error(fmt.Sprintf("Failed removing the mod's jar-file: %s", err.Error()))
					return
				}
				log.Info(fmt.Sprintf("Deleted %s", mod.Files[0].Name))
			} else {
				log.Info(fmt.Sprintf("Mod with name %s does not exist in configuration file", itemName))
			}
		}
	},
}

func init() {
	delCmd.Flags().BoolP("dependency", "d", false, "Remove a dependency instead of a mod")
	rootCmd.AddCommand(delCmd)
}

// Helper function to remove a mod from RequiredBy lists of dependencies
func removeModFromDependencies(cfg *config.Config, modName string) {
	for i := range cfg.Dependencies {
		// Remove the mod from RequiredBy list
		for j, requiredBy := range cfg.Dependencies[i].RequiredBy {
			if requiredBy == modName {
				cfg.Dependencies[i].RequiredBy = append(cfg.Dependencies[i].RequiredBy[:j], cfg.Dependencies[i].RequiredBy[j+1:]...)
				break
			}
		}
	}
}
