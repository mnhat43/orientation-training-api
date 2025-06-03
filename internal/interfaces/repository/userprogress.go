package repository

import (
	m "orientation-training-api/internal/models"
)

// UserProgressRepository defines methods for accessing user progress data
type UserProgressRepository interface {
	GetSingleUserProgress(userID int, courseID int) (m.UserProgress, error)
	SaveUserProgress(userProgress *m.UserProgress) error
	GetUserProgressByCourseID(courseID int) ([]m.UserProgress, error)
	GetAllUserProgressByUserID(userID int) ([]m.UserProgress, error)
	ReviewUserProgress(userID int, courseID int, performanceRating float64, performanceComment string, reviewedBy int) error
}
