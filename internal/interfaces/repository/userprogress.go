package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

// UserProgressRepository defines methods for accessing user progress data
type UserProgressRepository interface {
	// GetUserProgress returns the current progress for a user in a specific course
	GetUserProgress(getUserProgressParams *param.GetUserProgressParams) (m.UserProgress, error)
	
	// SaveUserProgress creates or updates user progress
	SaveUserProgress(userProgress *m.UserProgress) error
}
