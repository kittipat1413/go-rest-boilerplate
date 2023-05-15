package cmd

import (
	"go-rest-boilerplate/data"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [migrations-dir]",
	Short: "Revert one previously ran migration.",
	RunE:  runRollbackCmd,
}

func runRollbackCmd(cmd *cobra.Command, args []string) error {
	return runMigration(data.IntentRollback, args)
}
