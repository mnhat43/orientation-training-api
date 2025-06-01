package models

import (
	cm "orientation-training-api/internal/common"
)

type Quiz struct {
	cm.BaseModel

	Title      string  `json:"title" pg:"title,notnull"`
	Difficulty int     `json:"difficulty" pg:"difficulty,notnull"`
	TotalScore float64 `json:"total_score" pg:"total_score,notnull"`
	TimeLimit  int     `json:"time_limit" pg:"time_limit,notnull"`
}

type QuizQuestion struct {
	cm.BaseModel

	QuizID            int     `json:"quiz_id" pg:"quiz_id,notnull"`
	QuestionType      int     `json:"question_type" pg:"question_type,notnull"`
	QuestionText      string  `json:"question_text" pg:"question_text,notnull"`
	Explanation       string  `json:"explanation" pg:"explanation"`
	Weight            float64 `json:"weight" pg:"weight,notnull"`
	IsMultipleCorrect bool    `json:"is_multiple_correct" pg:"is_multiple_correct,notnull"`

	// Relationships
	Answers []QuizAnswer `json:"answers" pg:"rel:has-many"`
	Quiz    Quiz         `json:"quiz" pg:"rel:belongs-to"`
}

type QuizAnswer struct {
	cm.BaseModel

	QuizQuestionID int    `json:"quiz_question_id" pg:"quiz_question_id,notnull"`
	AnswerText     string `json:"answer_text" pg:"answer_text,notnull"`
	IsCorrect      bool   `json:"is_correct" pg:"is_correct,notnull"`
}

type QuizSubmission struct {
	cm.BaseModel

	UserID            int     `json:"user_id" pg:"user_id,notnull"`
	QuizID            int     `json:"quiz_id" pg:"quiz_id,notnull"`
	QuizQuestionID    int     `json:"quiz_question_id" pg:"quiz_question_id,notnull"`
	AnswerText        string  `json:"answer_text" pg:"answer_text"`
	SelectedAnswerIds []int   `json:"selected_answer_ids" pg:"selected_answer_ids,array"`
	Score             float64 `json:"score" pg:"score"`
	Attempt           int     `json:"attempt" pg:"attempt"`
	Reviewed          bool    `json:"reviewed" pg:"reviewed,notnull"`
	Feedback          string  `json:"feedback" pg:"feedback"`
	SubmittedAt       string  `json:"submitted_at" pg:"submitted_at,notnull"`

	// Relationships
	User         User         `json:"user" pg:"rel:belongs-to"`
	Quiz         Quiz         `json:"quiz" pg:"rel:belongs-to"`
	QuizQuestion QuizQuestion `json:"quiz_question" pg:"rel:belongs-to"`
}
