package main

import (
	_ "embed"
	"strings"

	"github.com/wearetechnative/bmc/cmd"
)

//go:embed VERSION-bmc
var versionRaw string

func main() {
	if cmd.Version == "dev" {
		cmd.Version = strings.TrimSpace(versionRaw)
	}
	cmd.Execute()
}
