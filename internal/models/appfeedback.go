package models

import (
	cm "orientation-training-api/internal/common"
	"time"
)

// AppFeedback represents user feedback about the application
type AppFeedback struct {
	cm.BaseModel

	Rating   float64   `json:"rating" db:"rating"`
	Feedback string    `json:"feedback" db:"feedback"`
	SubmitAt time.Time `json:"submit_at" db:"submit_at"`
	UserID   int       `json:"user_id" db:"user_id"`
}
