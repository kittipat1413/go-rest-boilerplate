package cmd

import (
	"context"
	"fmt"
	"log"

	"go-rest-boilerplate/cmd/prompts"
	"go-rest-boilerplate/config"
	"go-rest-boilerplate/data"
	"go-rest-boilerplate/internal"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [migrations-dir]",
	Short: "Runs all migration scripts in migrations-dir. Defaults to ./data/migrations",
	RunE:  runMigrateCmd,
}

func runMigrateCmd(cmd *cobra.Command, args []string) error {
	return runMigration(data.IntentMigrate, args)
}

func runMigration(intent data.Intent, args []string) (err error) {
	defer internal.WrapErr("migrate", &err)

	cfg := config.MustConfigure()
	prompt := prompts.New(cfg, args)
	dir := prompt.OptionalStr("migrations dir", "./data/migrations")

	db, err := data.Connect(cfg)
	if err != nil {
		return err
	}

	scope, err := data.NewScope(context.Background(), db)
	if err != nil {
		cfg.Fatalln("db connection error", err)
	} else {
		defer scope.End(&err)
	}

	migrator := data.NewMigrator(db, dir)
	plans, dirty, err := migrator.Plan(scope.Context(), intent)
	if err != nil {
		return err
	}

	if len(plans) == 0 {
		log.Println("no changes")
		return nil
	}

	for _, plan := range plans {
		fmt.Println(plan)
	}

	if dirty {
		log.Println("some migrations are missing or have changed content")
		if !prompt.YesNo("proceed anyway") {
			return nil
		}
	}

	log.Println(len(plans), "migrations planned")
	if !prompt.YesNo("apply changes") {
		return nil
	}

	for _, plan := range plans {
		fmt.Println(plan)
		if err := migrator.Apply(scope.Context(), plan); err != nil {
			log.Fatalln("failed to run migration", err)
		}
	}

	log.Println(len(plans), "migration(s) applied")
	return nil
}
