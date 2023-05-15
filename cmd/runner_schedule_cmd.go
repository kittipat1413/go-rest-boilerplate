package cmd

import (
	"context"
	"errors"
	"strconv"

	"go-rest-boilerplate/cmd/prompts"
	"go-rest-boilerplate/config"
	"go-rest-boilerplate/data"
	"go-rest-boilerplate/internal"
	"go-rest-boilerplate/worker"
	"go-rest-boilerplate/worker/runners"

	"github.com/spf13/cobra"
)

var scheduleRunnerCmd = &cobra.Command{
	Use:   "runner-schedule [runner...]",
	Short: "Schedule runner",
	RunE:  runRunnerCmd,
}

func runRunnerCmd(cmd *cobra.Command, args []string) (err error) {
	defer internal.WrapErr("runner-schedule", &err)

	cfg := config.MustConfigure()
	prompt := prompts.New(cfg, args)

	if prompt.Len() != 0 {
		for prompt.Len() > 0 {
			runnerName := prompt.Str("")
			runner := runners.New(runnerName)
			if runner == nil {
				cfg.Println(runnerName + " runner does not exists")
				continue
			}
			err := scheduleJob(cfg, runner)
			if err != nil {
				return err
			}
		}
	} else {
		runnerName := prompt.Str("runner's name")
		runner := runners.New(runnerName)
		if runner == nil {
			err = errors.New("runner does not exists")
			return
		}
		err := scheduleJob(cfg, runner)
		if err != nil {
			return err
		}
	}

	return nil
}

func scheduleJob(cfg *config.Config, runner runners.Interface) (err error) {
	db, err := data.Connect(cfg)
	if err != nil {
		return err
	}
	ctx := data.NewContext(config.NewContext(context.Background(), cfg), db)
	if jobID, err := worker.ScheduleNowIfNotExists(ctx, runner); err != nil {
		return err
	} else if jobID == 0 {
		cfg.Println(runner.Name() + " job has already existed")
		return nil
	} else {
		cfg.Printf("scheduled %s job %s", runner.Name(), strconv.FormatInt(jobID, 10))
		return nil
	}

}
