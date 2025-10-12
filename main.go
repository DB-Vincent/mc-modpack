/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package main

import (
	"fmt"
	"os"

	"github.com/DB-Vincent/mc-modpack/cmd"
)

var version = "dev" // This will be overridden by ldflags during build

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version)
		return
	}
	cmd.Execute()
}
