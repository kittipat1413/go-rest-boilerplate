package runners

import (
	"context"
	"time"
)

type (
	Interface interface {
		Name() string
		Run(context.Context, func(time.Time)) error
	}

	// Cloner should be implemented by all runners that requires different arguments
	// for each run and must be serialized/deserialized to/from the jobs table.
	Cloner interface {
		// Clone should returns a clone of itself and initialized to receive
		// a new set of payload from the database (via json.Unmarshal) and run.
		Clone() Interface
	}
)

var allRunners = []Interface{
	&TestRunner{},
}

func New(name string) Interface {
	for _, runner := range allRunners {
		if runner.Name() == name {
			if cloner, ok := runner.(Cloner); ok {
				return cloner.Clone()
			} else {
				return runner
			}
		}
	}
	return nil
}
