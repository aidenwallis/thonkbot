package common

import "time"

type Quote struct {
	Username  string
	Message   string
	CreatedAt time.Time
}
