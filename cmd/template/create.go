/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-tool/core"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new template",
	Long:  `create a new template`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.CreateTemplate(name, keys)
		if err != nil {
			logrus.WithError(err).Error("Failed to create template")
		}
	},
}

func init() {
	templateCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVar(&name, "name", "", "template name")
	createCmd.PersistentFlags().StringSliceVar(&keys, "keys", nil, "keys array")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("keys")
}
