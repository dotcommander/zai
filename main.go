package main

import (
	"github.com/garyblankenship/zai/cmd"
	"github.com/garyblankenship/zai/internal/config"
)

func main() {
	config.SetDefaults()
	cmd.Execute()
}