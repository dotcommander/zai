package main

import (
	"github.com/dotcommander/zai/cmd"
	"github.com/dotcommander/zai/internal/config"
)

func main() {
	config.SetDefaults()
	cmd.Execute()
}
