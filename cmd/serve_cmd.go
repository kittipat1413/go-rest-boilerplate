package cmd

import (
	"go-rest-boilerplate/config"
	"go-rest-boilerplate/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Starts application server",
	RunE:    runServeCmd,
	Aliases: []string{"s", "server"},
}

func runServeCmd(cmd *cobra.Command, args []string) error {
	cfg := config.MustConfigure()
	cfg.Println("starting web server on: " + cfg.ListenAddr())
	return server.New(cfg).Start()
}
