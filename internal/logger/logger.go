/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package logger

import (
	"fmt"
  "os"
)

type CliLogger struct {
  verbose bool
}

func New() (*CliLogger) {
  return &CliLogger{
    verbose: false,
  }
}

func (logger *CliLogger) SetVerbose(verbose bool) {
  logger.verbose = verbose
}

func (logger *CliLogger) Info(msg string) {
  fmt.Printf("[INFO] %s\n", msg)
}

func (logger *CliLogger) Warn(msg string) {
  fmt.Printf("[WARN] %s\n", msg)
}

func (logger *CliLogger) Error(msg string) {
  fmt.Printf("[ERRO] %s\n", msg)
  os.Exit(1)
}

func (logger *CliLogger) Debug(msg string) {
  // Verbose mode has not been enabled
  if (!logger.verbose) {
    return
  }

  fmt.Printf("[DEBU] %s\n", msg)
}
