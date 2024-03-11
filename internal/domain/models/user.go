package models

import "time"

type User struct {
	ID        int64
	Email     string
	PassHash  []byte
	CreatedAt time.Time
	VisitedAt time.Time
}
