/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/DB-Vincent/mc-modpack/internal/config"
	"github.com/spf13/cobra"
)

var (
	name      string
	mcVersion string
	loader    string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a directory where you want to store your modpack",
	Long:  "Initializes a directory where you want to store your modpack",
	Run: func(cmd *cobra.Command, args []string) {
		path := workingDirectory
		if workingDirectory == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Error(fmt.Sprintf("Failed to get the current working directory: %s", err.Error()))
				return
			}
		}

		// Create directory if it doesn't exist
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, 0755)
			if err != nil {
				log.Error(fmt.Sprintf("Unable to create modpack directory '%s': %v", path, err))
				return
			}
			log.Info(fmt.Sprintf("Created modpack directory: %s", path))
		}

		// Create a config based on the inputs
		configContent := config.Config{
			McVersion: mcVersion,
			Name:      name,
			Loader:    loader,
		}

		// Create configuration file
		err := config.Update(path, configContent)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to create modpack configuration file: %v", err))
			return
		}

		log.Info(fmt.Sprintf("Successfully initialized modpack '%s' for Minecraft %s (%s)", name, mcVersion, loader))
	},
}

func init() {
	initCmd.Flags().StringVar(&name, "name", "", "Name of the modpack you want to create")
	if err := initCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
	initCmd.Flags().StringVar(&mcVersion, "mc-version", "", "Version of Minecraft to create the modpack for")
	if err := initCmd.MarkFlagRequired("mc-version"); err != nil {
		panic(err)
	}
	initCmd.Flags().StringVar(&loader, "loader", "fabric", "Modding platform used in the modpack you want to create")
	rootCmd.AddCommand(initCmd)
}
