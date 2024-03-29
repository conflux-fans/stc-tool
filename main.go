/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/conflux-fans/storage-cli/cmd"
	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/core"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Init()
	cfg := config.Get()

	logLevel, err := logrus.ParseLevel(cfg.Log)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(logLevel)

	core.Init()
	cmd.Execute()
}
