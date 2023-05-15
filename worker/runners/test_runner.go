package runners

import (
	"context"
	"errors"
	"fmt"
	"go-rest-boilerplate/config"
	"time"
)

type TestRunner struct {
	Arg          string    `json:"arg"`
	RegisteredAt time.Time `json:"registered_at"`
	Error        string    `json:"error"`
}

var _ Interface = &TestRunner{}

func (r *TestRunner) Name() string { return "test-runner" }
func (r *TestRunner) Clone() Interface {
	clone := *r
	return &clone
}

func (r *TestRunner) Run(ctx context.Context, scheduler func(time.Time)) error {
	cfg := config.FromContext(ctx)
	defer cfg.Println("TestRunner finished")
	defer scheduler(time.Now().Add(time.Hour * 12))

	cfg.Printf(
		"TestRunner with arg `%s` registered at `%s`\n",
		r.Arg,
		r.RegisteredAt,
	)

	time.Sleep(3 * time.Second)
	if r.Error != "" {
		return fmt.Errorf("mock error: %w", errors.New(r.Error))
	} else {
		return nil
	}
}
