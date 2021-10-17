package models

import "time"

type Reward struct {
	ID        string
	Token     string
	Amount    float64
	User      string
	Pool      string
	EarnedAt  time.Time
	ClaimedAt time.Time
}
