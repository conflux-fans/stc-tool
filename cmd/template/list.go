/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package template

import (
	"github.com/conflux-fans/storage-cli/core"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List templates",
	Long:  `List templates`,
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := core.ListTemplate()
		logger.Get().WithField("result", templates).WithError(err).Info("list templates")
		if err != nil {
			logger.Fail(err.Error())
			return
		}
		// params := map[string]string{"Templates: %v\n": strings.Join(templates, ", ")}
		logger.SuccessfWithList(templates, "Get templates")
	},
}

func init() {
	templateCmd.AddCommand(listCmd)
}
