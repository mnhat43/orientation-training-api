package requestparams

// CreateQuizParams defines parameters for creating a new quiz
type CreateQuizParams struct {
	Title      string  `json:"title" valid:"required"`
	Difficulty int     `json:"difficulty" valid:"required"`
	TotalScore float64 `json:"total_score" valid:"required"`
	TimeLimit  int     `json:"time_limit" valid:"required"` // in minutes
}

// UpdateQuizParams defines parameters for updating an existing quiz
type UpdateQuizParams struct {
	ID         int     `json:"id" valid:"required"`
	Title      string  `json:"title" valid:"required"`
	Difficulty int     `json:"difficulty" valid:"required"`
	TotalScore float64 `json:"total_score" valid:"required"`
	TimeLimit  int     `json:"time_limit" valid:"required"` // in minutes
}

// QuizListParams defines parameters for fetching a list of quizzes
type QuizListParams struct {
	Title       string `json:"title"`
	CurrentPage int    `json:"current_page"`
	RowPerPage  int    `json:"row_per_page"`
}

// GetQuizParams defines parameters for fetching a single quiz
type GetQuizParams struct {
	QuizID int `json:"quiz_id" valid:"required"`
}

// DeleteQuizParams defines parameters for deleting a quiz
type DeleteQuizParams struct {
	QuizID int `json:"quiz_id" valid:"required"`
}

// QuizAnswerParam defines parameters for a quiz answer
type QuizAnswerParam struct {
	AnswerText string `json:"answer_text" valid:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

// CreateQuizQuestionParams defines parameters for creating a quiz question
type CreateQuizQuestionParams struct {
	QuizID            int               `json:"quiz_id" valid:"required"`
	QuestionType      int               `json:"question_type" valid:"required"` // 1 for multiple choice, 2 for text
	QuestionText      string            `json:"question_text" valid:"required"`
	Weight            float64           `json:"weight" valid:"required"`
	IsMultipleCorrect bool              `json:"is_multiple_correct"`
	Answers           []QuizAnswerParam `json:"answers" valid:"required"`
}

// SubmitQuizAnswersParams defines parameters for submitting answers to a quiz
type SubmitQuizAnswersParams struct {
	QuizID            int    `json:"quiz_id" valid:"required"`
	QuestionID        int    `json:"question_id" valid:"required"`
	AnswerText        string `json:"answer_text"`
	SelectedAnswerIds []int  `json:"selected_answer_ids"`
}

// QuizAnswer represents a single answer to a quiz question
type QuizAnswer struct {
	QuestionID        int    `json:"question_id" valid:"required"`
	AnswerText        string `json:"answer_text"`
	SelectedAnswerIds []int  `json:"selected_answer_ids"`
}

// SubmitFullQuizParams defines parameters for submitting a complete quiz with multiple answers
type SubmitFullQuizParams struct {
	QuizID      int          `json:"quiz_id" valid:"required"`
	Answers     []QuizAnswer `json:"answers" valid:"required"`
	SubmittedAt string       `json:"submitted_at" valid:"required"`
}

// GetQuizResultsParams defines parameters for fetching quiz results
type GetQuizResultsParams struct {
	QuizID int `json:"quiz_id" valid:"required"`
	UserID int `json:"user_id"` // Only used by managers
}
