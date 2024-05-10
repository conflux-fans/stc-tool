/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new template",
	Long:  `Create a new template`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.CreateTemplate(name, keys)
		if err != nil {
			logger.Fail(err.Error())
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
