package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "root",
		Short: "SumoLogic Ecosystem",
		Long:  `SumoLogic Interview Take-Home Server`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
