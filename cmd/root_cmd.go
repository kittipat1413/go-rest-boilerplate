package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Starts the application",
}

func init() {
	rootCmd.AddCommand(
		migrateCmd,
		newMigrationCmd,
		printConfigCmd,
		rollbackCmd,
		serveCmd,
		testWorkerCmd,
		workerCmd,
		scheduleRunnerCmd,
	)
}
