package entity

import "time"

type Info struct {
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Infos []Info
