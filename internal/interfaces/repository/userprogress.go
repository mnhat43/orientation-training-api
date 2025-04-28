package repository

import (
	m "orientation-training-api/internal/models"
)

// UserProgressRepository defines methods for accessing user progress data
type UserProgressRepository interface {
	GetUserProgress(userID int, courseID int) (m.UserProgress, error)
	SaveUserProgress(userProgress *m.UserProgress) error
}
