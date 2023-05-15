package cmd

import (
	"context"
	"strconv"
	"time"

	"go-rest-boilerplate/cmd/prompts"
	"go-rest-boilerplate/config"
	"go-rest-boilerplate/data"
	"go-rest-boilerplate/internal"
	"go-rest-boilerplate/worker"
	"go-rest-boilerplate/worker/runners"

	"github.com/spf13/cobra"
)

var testWorkerCmd = &cobra.Command{
	Use:   "test-worker",
	Short: "Enqueue a test runner to test/debug the background worker",
	RunE:  runTestWorkerCmd,
}

func runTestWorkerCmd(cmd *cobra.Command, args []string) (err error) {
	defer internal.WrapErr("test-worker", &err)

	cfg := config.MustConfigure()
	prompt := prompts.New(cfg, args)
	arg, error :=
		prompt.Str("arg to test runner"),
		prompt.OptionalStr("runner error", "")

	db, err := data.Connect(cfg)
	if err != nil {
		return err
	}

	ctx := data.NewContext(config.NewContext(context.Background(), cfg), db)
	runner := &runners.TestRunner{Arg: arg, Error: error, RegisteredAt: time.Now()}
	if jobID, err := worker.ScheduleNow(ctx, runner); err != nil {
		return err
	} else {
		cfg.Println("scheduled job " + strconv.FormatInt(jobID, 10))
		return nil
	}
}
