package models

import "time"

type RewardPool interface {
	TopUpPool(amt float64)
	ClosePool()
}

// ResponseRewardPool represents the reward pool objects where each participant is paid out for their response to the survey.
type ResponseRewardPool struct {
	ID             string
	Title          string
	Token          string
	TotalFunds     float64
	Active         bool
	TotalEarned    float64
	TotalRewarded  float64
	ResponseReward float64
	Responses      int
	SurveyID       string
	CreatedAt      time.Time
	ClosedAt       time.Time
	UserID         string
	Participants   []string
}

// RaffleRewardPool represnts the reward pool object where each participant is entered into a raffle for a crypto prize payout.
// The raffle can have multiple winners which each take a equal share of the total pool funds.
type RaffleRewardPool struct {
	ID           string
	Title        string
	Token        string
	TotalFunds   float64
	Active       bool
	Entries      int
	WinnerCount  int
	SurveyID     string
	CreatedAt    time.Time
	ClosesAt     time.Time
	UserID       string
	Participants []string
	Winners      []string
}
