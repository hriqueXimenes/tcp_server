package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	awaitCmd = &cobra.Command{
		Use:   "await",
		Short: "Command to wait X milliseconds",
		Long:  `This command waits for a specified number of milliseconds.`,
		Run:   awaitCommandExecute,
	}
)

func init() {
	awaitCmd.Flags().IntP("time", "t", 1000, "Time to wait in milliseconds")
	rootCmd.AddCommand(awaitCmd)
}

func awaitCommandExecute(cmd *cobra.Command, args []string) {
	t, err := cmd.Flags().GetInt("time")
	if err != nil {
		fmt.Println("Error getting time:", err)
		return
	}

	time.Sleep(time.Duration(t) * time.Millisecond)

	fmt.Println("Await finished.")
}
