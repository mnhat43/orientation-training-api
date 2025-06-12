package requestparams

// SubmitAppFeedbackRequest represents a request to submit app feedback
type SubmitAppFeedbackRequest struct {
	Rating      int    `json:"rating" valid:"required,range(1|5)"`
	Feedback    string `json:"feedback" valid:"required"`
	SubmittedAt string `json:"submittedAt,omitempty"`
}

// DeleteAppFeedbackRequest represents a request to delete app feedback
type DeleteAppFeedbackRequest struct {
	ID int `json:"id" valid:"required"`
}

// GetAppFeedbackListRequest represents a request to get app feedback list
type GetAppFeedbackListRequest struct {
	Limit  int `json:"limit" valid:"required,range(1|100)"`
	Offset int `json:"offset" valid:"optional,min(0)"`
}
