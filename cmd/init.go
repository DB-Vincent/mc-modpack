/*
  Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
  "errors"
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
	Long: "Initializes a directory where you want to store your modpack",
	Run: func(cmd *cobra.Command, args []string) {
    path := workingDirectory
    if workingDirectory == "" {
      var err error
      path, err = os.Getwd()
      if err != nil {
        panic(err)
      }
    }


    if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
      err := os.Mkdir(path, 0755)
      if err != nil {
        panic(err)
      }  
    }

    // Create a config based on the inputs
    configContent := config.Config{
      McVersion: mcVersion,
      Name:      name,
      Loader:    loader,
    }

    // Create configuration file
    config.Update(path, configContent)
	},
}

func init() {
  initCmd.Flags().StringVar(&name, "name", "", "Name of the modpack you want to create")
 	if err := initCmd.MarkFlagRequired("name"); err != nil { panic(err) }
  initCmd.Flags().StringVar(&mcVersion, "mc-version", "", "Version of Minecraft to create the modpack for")
 	if err := initCmd.MarkFlagRequired("mc-version"); err != nil { panic(err) }
  initCmd.Flags().StringVar(&loader, "loader", "fabric", "Modding platform used in the modpack you want to create")
	rootCmd.AddCommand(initCmd)
}
