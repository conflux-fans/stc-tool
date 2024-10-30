/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/conflux-fans/storage-cli/cmd"
	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/core"
)

func main() {
	config.Init()

	core.Init()
	cmd.Execute()
}
