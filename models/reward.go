package models

import "time"

type Reward struct {
	ID string
	// Token is the currency being paid out.
	Token string
	// Amount is the amount of the token to be rewards.
	Amount float64
	// User is the id of the user who has earned this reward.
	User string
	// Pool is the id of the reward pool this reward belongs to.
	Pool string
	// ShadowUser indicates if the person who earned the reward was signed up on ApeSurvey at the time they earned the reward.
	ShadowUser bool
	// EarnedAt indicates when the participant earned this reward.
	EarnedAt time.Time
	// ClaimedAt indicates when the reward was claimed by the user.
	ClaimedAt time.Time
	// Reminders is the number of reminder emails have been sent to the user email.
	Reminders int
	// RemindedAt timestamps when the last reminder email notification was sent.
	RemindedAt time.Time
	// RefundedAt timestamps when the reward has refunded back to the reward pool.
	RefundedAt time.Time
	// RejectedAt timestamps when/if the participant rejected the reward.
	RejectedAt time.Time
	// ExpiresAt timestamps when/if the reward expires due to being unclaimed from a closing reward pool.
	ExpiresAt time.Time
}
