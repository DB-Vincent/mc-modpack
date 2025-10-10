/*
  Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mc-modpack",
	Short: "A simple CLI tool to update the mods in your modpack",
	Long: "A simple CLI tool to update the mods in your modpack",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() { }
