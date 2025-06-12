package response

import (
	"orientation-training-api/internal/models"
	"time"
)

// FeedbackWithUser represents app feedback with associated user information
type FeedbackWithUser struct {
	ID        int       `json:"id"`
	Rating    int       `json:"rating"`
	Feedback  string    `json:"feedback"`
	SubmitAt  time.Time `json:"submit_at"`
	CreatedAt time.Time `json:"created_at"`
	User      struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	} `json:"user"`
}

// CreateFeedbackWithUserFromAppFeedback converts an AppFeedback model to a FeedbackWithUser response
func CreateFeedbackWithUserFromAppFeedback(feedback *models.AppFeedback, user *models.User, userProfile *models.UserProfile, roleName string) *FeedbackWithUser {
	feedbackWithUser := &FeedbackWithUser{
		ID:        feedback.ID,
		Rating:    feedback.Rating,
		Feedback:  feedback.Feedback,
		SubmitAt:  feedback.SubmitAt,
		CreatedAt: feedback.CreatedAt,
	}

	feedbackWithUser.User.ID = user.ID
	feedbackWithUser.User.Email = user.Email
	feedbackWithUser.User.FirstName = userProfile.FirstName
	feedbackWithUser.User.LastName = userProfile.LastName
	feedbackWithUser.User.Role = roleName

	return feedbackWithUser
}
