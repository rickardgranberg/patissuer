package main

import (
	"log"

	"github.com/rickardgranberg/patissuer/cmd/patissuer/cmd"
)

var (
	// version is the application version
	version string = "dev"
	// commit is the git commit
	commit string = "deadbeef"
	// buildTime is the build time
	buildTime string
)

func main() {
	if err := cmd.Execute(version, commit, buildTime); err != nil {
		log.Fatal(err)
	}
}
