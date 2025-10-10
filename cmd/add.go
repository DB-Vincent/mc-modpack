/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"os"

	"github.com/DB-Vincent/mc-modpack/internal/config"
	"github.com/DB-Vincent/mc-modpack/internal/modrinth"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a mod to the current modpack",
	Run: func(cmd *cobra.Command, args []string) {
    path, err := os.Getwd()
    if err != nil {
      panic(err)
    }

    if err := config.Exists(path); err != nil {
      panic(err)
    } 

    cfg, err := config.Load(path)
    if err != nil {
      panic(err)
    }

    version, err := modrinth.GetLatestVersion("better-advancements", cfg.McVersion, cfg.Loader)
    if err != nil {
      panic(err)
    }

    cfg.Mods = append(cfg.Mods, config.Mod{
      Name: "better-advancements",
      Version: version.ModVersion,
    })
    err = config.Update(path, *cfg)
    if err != nil {
      panic(err)
    }
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
