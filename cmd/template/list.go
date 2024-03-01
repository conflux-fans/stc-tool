/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List templates",
	Long:  `List templates`,
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := core.ListTemplate()
		logrus.WithField("result", templates).WithError(err).Info("list templates")
		if err != nil {
			logrus.WithError(err).Error("Failed to create template")
		} else {
			fmt.Printf("Templates: %v\n", strings.Join(templates, ", "))
		}
	},
}

func init() {
	templateCmd.AddCommand(listCmd)
}
