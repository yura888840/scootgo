package domain

import "time"

type Timer struct {
	ID        string
	UserID    string
	ScooterID string
	StartedAt time.Time
}
