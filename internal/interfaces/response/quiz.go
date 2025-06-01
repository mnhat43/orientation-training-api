package response

// PendingReviewResponse represents the response for GetQuizPendingReview
type PendingReviewResponse struct {
	UserID     int                 `json:"user_id"`
	Fullname   string              `json:"fullname"`
	Department string              `json:"department"`
	Avatar     string              `json:"avatar"`
	Reviews    []PendingReviewItem `json:"reviews"`
}

// PendingReviewItem represents a single review item in the pending reviews
type PendingReviewItem struct {
	SubmissionID int     `json:"submission_id"`
	CourseTitle  string  `json:"course_title"`
	QuestionText string  `json:"question_text"`
	AnswerText   string  `json:"answer_text"`
	SubmittedAt  string  `json:"submitted_at"`
	MaxScore     float64 `json:"maxScore"`
}
