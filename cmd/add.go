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
				if err = addDependency(cfg, dependency.ProjectId, modName, path); err != nil {
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

func addDependency(cfg *config.Config, dependencyProjectId string, requiredByMod string, path string) error {
	// Get project information to get the actual name
	project, err := modrinth.GetProject(dependencyProjectId)
	if err != nil {
		return fmt.Errorf("failed to get project information for dependency '%s': %v", dependencyProjectId, err)
	}

	version, err := modrinth.GetLatestVersion(dependencyProjectId, cfg.McVersion, cfg.Loader)
	if err != nil {
		return fmt.Errorf("failed to find latest version for dependency '%s' (MC %s, %s): %v", project.Slug, cfg.McVersion, cfg.Loader, err)
	}

	// Check if dependency exists in config file (using the project slug as the name)
	dependencyExists, index := config.HasDependency(*cfg, project.Slug)
	dependencyDownload := true

	if dependencyExists {
		// Add the requiring mod to the RequiredBy list if not already present
		if !contains(cfg.Dependencies[index].RequiredBy, requiredByMod) {
			cfg.Dependencies[index].RequiredBy = append(cfg.Dependencies[index].RequiredBy, requiredByMod)
		}

		if cfg.Dependencies[index].Version != version.ModVersion {
			log.Info(fmt.Sprintf("Updating dependency %s from version %s to %s", version.Files[0].Name, cfg.Dependencies[index].Version, version.ModVersion))
			cfg.Dependencies[index].Version = version.ModVersion
			cfg.Dependencies[index].VersionId = version.VersionId
		} else {
			log.Info(fmt.Sprintf("Already have dependency %s with version %s downloaded", version.Files[0].Name, cfg.Dependencies[index].Version))
			dependencyDownload = false
		}
	} else {
		log.Info(fmt.Sprintf("Adding dependency %s version %s to modpack (required by %s)", version.Files[0].Name, version.ModVersion, requiredByMod))
		cfg.Dependencies = append(cfg.Dependencies, config.Dependency{
			Name:       project.Slug, // Use the actual project title instead of project ID
			VersionId:  version.VersionId,
			Version:    version.ModVersion,
			RequiredBy: []string{requiredByMod},
		})
	}

	if dependencyDownload {
		if err = modrinth.DownloadFile(path, version.Files[0]); err != nil {
			return fmt.Errorf("failed to download dependency '%s': %v", version.Files[0].Name, err)
		}
		log.Info(fmt.Sprintf("Successfully downloaded dependency %s", version.Files[0].Name))
	}

	return nil
}

// Helper function to check if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(addCmd)
}
