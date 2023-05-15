package cmd

import (
	"go-rest-boilerplate/config"
	"go-rest-boilerplate/worker"

	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:     "worker",
	Short:   "Starts background worker",
	RunE:    runWorkerCmd,
	Aliases: []string{"w", "work"},
}

func runWorkerCmd(cmd *cobra.Command, args []string) error {
	cfg := config.MustConfigure()
	cfg.Println("starting worker")
	return worker.New(cfg).Start()
}
