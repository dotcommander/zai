package main

import (
	"github.com/dotcommander/zai/cmd"
	"github.com/dotcommander/zai/internal/config"
	"github.com/dotcommander/zai/internal/version"
)

func main() {
	version.Init()
	config.SetDefaults()
	cmd.Execute()
}
