/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-tool/cmd"
	"github.com/zero-gravity-labs/zerog-storage-tool/config"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	config.Init()
	core.Init()
	cmd.Execute()
}
