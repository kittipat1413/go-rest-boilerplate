package internal

import (
	"fmt"
)

func WrapErr(name string, errptr *error) {
	if *errptr != nil {
		*errptr = fmt.Errorf(name+": %w", *errptr)
	}
}
