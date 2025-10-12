/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package cmd

import (
	"os"

	"github.com/DB-Vincent/mc-modpack/internal/logger"
	"github.com/spf13/cobra"
)

var (
	log = logger.New()

	rootCmd = &cobra.Command{
		Use:   "mc-modpack",
		Short: "A simple CLI tool to update the mods in your modpack",
		Long:  "A simple CLI tool to update the mods in your modpack",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetVerbose(verbose)
		},
	}

	workingDirectory string
	verbose          bool
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// getWorkingDirectory returns the working directory path, using the flag value or current directory
func getWorkingDirectory() (string, error) {
	if workingDirectory != "" {
		return workingDirectory, nil
	}
	return os.Getwd()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&workingDirectory, "directory", "", "Directory in which to work")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}
