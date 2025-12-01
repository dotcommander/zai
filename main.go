package main

import (
	"zai/cmd"
	"zai/internal/config"
)

func main() {
	config.SetDefaults()
	cmd.Execute()
}