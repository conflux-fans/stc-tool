/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/zero-gravity-labs/zerog-storage-tool/cmd"
	"github.com/zero-gravity-labs/zerog-storage-tool/config"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

func main() {
	config.Init()
	core.Init()
	cmd.Execute()
}
