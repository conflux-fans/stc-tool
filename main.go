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
